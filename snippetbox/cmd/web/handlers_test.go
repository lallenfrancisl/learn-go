package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lallenfrancisl/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ping(rr, r)

	rs := rr.Result()

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}

func TestPingE2E(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	statusCode, _, body := ts.get(t, "/ping")

	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}

func TestSnippetView(t *testing.T) {
    app := newTestApplication(t)

    ts := newTestServer(t, app.routes())
    defer ts.Close()

    tests := []struct {
        name     string
        urlPath  string
        wantCode int
        wantBody string
    }{
        {
            name:     "Valid ID",
            urlPath:  "/snippets/1",
            wantCode: http.StatusOK,
            wantBody: "An old silent pond...",
        },
        {
            name:     "Non-existent ID",
            urlPath:  "/snippets/2",
            wantCode: http.StatusNotFound,
        },
        {
            name:     "Negative ID",
            urlPath:  "/snippets/-1",
            wantCode: http.StatusNotFound,
        },
        {
            name:     "Decimal ID",
            urlPath:  "/snippets/1.23",
            wantCode: http.StatusNotFound,
        },
        {
            name:     "String ID",
            urlPath:  "/snippets/foo",
            wantCode: http.StatusNotFound,
        },
        {
            name:     "Empty ID",
            urlPath:  "/snippets/",
            wantCode: http.StatusNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            code, _, body := ts.get(t, tt.urlPath)

            assert.Equal(t, code, tt.wantCode)

            if tt.wantBody != "" {
                assert.StringContains(t, body, tt.wantBody)
            }
        })
    }
}
