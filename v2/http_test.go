package v2

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTP_GetDoSlurp(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET")
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("header mismatch: %q", got)
		}
		w.WriteHeader(200)
		_, _ = io.WriteString(w, "ok")
	}))
	defer srv.Close()

	out := Get(srv.URL).
		WithHeader("X-Test", "1").
		Do(context.Background()).
		Slurp().
		String().
		Must()

	if out != "ok" {
		t.Fatalf("got %q", out)
	}
}

func TestHTTP_BasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "Basic " + base64.StdEncoding.EncodeToString([]byte("u:p"))
		if got := r.Header.Get("Authorization"); got != want {
			t.Fatalf("auth mismatch: %q", got)
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	if err := Get(srv.URL).BasicAuth("u", "p").Do(context.Background()).Err(); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestHTTP_JSONBody(t *testing.T) {
	type payload struct {
		A string `json:"a"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Fatalf("content-type mismatch: %q", ct)
		}
		b, _ := io.ReadAll(r.Body)
		_ = r.Body.Close()
		if !strings.Contains(string(b), "\"a\":\"x\"") {
			t.Fatalf("unexpected body: %q", string(b))
		}
		w.WriteHeader(201)
	}))
	defer srv.Close()

	status := Post(srv.URL).
		WithJSONBody(payload{A: "x"}).
		Do(context.Background()).
		StatusCode().
		Must()

	if status != 201 {
		t.Fatalf("got %d", status)
	}
}


