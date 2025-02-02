package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/lallenfrancisl/greenlight-api/internal/validator"
)

func (app *application) readIDParams(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	var js []byte
	var err error

	if app.config.env == "development" {
		js, err = json.MarshalIndent(data, "", "\t")
	} else {
		js, err = json.Marshal(data)
	}

	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) readJSON(
	w http.ResponseWriter, r *http.Request, dst interface{},
) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		if errors.As(err, &syntaxError) {
			return fmt.Errorf("body contains badly formed JSON (at character %d)", syntaxError.Offset)
		} else if errors.Is(err, io.ErrUnexpectedEOF) {
			return errors.New("body contains badly formed JSON")
		} else if errors.As(err, &unmarshalTypeError) {
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}

			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		} else if errors.Is(err, io.EOF) {
			return errors.New("body must not be empty")
		} else if errors.As(err, &invalidUnmarshalError) {
			panic(err)
		} else {
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// Returns a string value from the query string or the provided defaultValue
// if no matching key could be found
func (app *application) readString(
	qs url.Values, key string, defaultValue string,
) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

// Returns an integer value from the query or the provided defaultValue if not matching
// key could be found. If value couldn't be converted to an integer, the error is
// recorded in the validator provided
func (app *application) readInt(
	qs url.Values, key string, defaultValue int, v *validator.Validator,
) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer")

		return defaultValue
	}

	return i
}

// Reads a string value from the query string and then splits it
// into a slice on the comma character. If no matching key could be found, it returns
// the provided default value.
func (app *application) readCSV(
	qs url.Values, key string, defaultValue []string,
) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	values := strings.Split(csv, ",")

	for i, v := range values {
		values[i] = strings.Trim(v, " \n\t")
	}

	return values
}
