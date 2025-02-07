package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/lallenfrancisl/greenlight-api/internal/data"
	"github.com/lallenfrancisl/greenlight-api/internal/validator"
	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.config.Limiter.Enabled {
			next.ServeHTTP(w, r)

			return
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(
					rate.Limit(app.config.Limiter.RPS), app.config.Limiter.Burst,
				),
				lastSeen: time.Now(),
			}
		}

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)

			return
		}

		// We are not using defer here that would only get called after
		// next.ServeHTTP() returns. Which means every middleware
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authrizationHeader := r.Header.Get("Authorization")

		if authrizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)

			return
		}

		headerParts := strings.Split(authrizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidCredentialsResponse(w, r)

			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidCredentialsResponse(w, r)

			return
		}

		user, err := app.repo.Users.GetByToken(data.ScopeAuthentication, token)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.invalidCredentialsResponse(w, r)
			} else {
				app.serverErrorResponse(w, r, err)
			}

			return
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if !user.Activated {
			app.userNotActivatedResponse(w, r)

			return
		}

		next.ServeHTTP(w, r)
	})

	return app.requireAuthenticatedUser(fn)
}
