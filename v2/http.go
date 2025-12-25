package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Request is a typed HTTP request builder.
type Request struct {
	req    *http.Request
	client *http.Client
}

// Req is a pipeline wrapper around Option[Request].
type Req struct{ Option[Request] }

// Resp is a pipeline wrapper around Option[*http.Response].
type Resp struct{ Option[*http.Response] }

// Client constructs a request builder with a custom client.
func Client(c *http.Client) Req {
	if c == nil {
		c = http.DefaultClient
	}
	return Req{Ok(Request{client: c})}
}

// Get creates a GET request builder using http.DefaultClient.
func Get(url string) Req { return Client(http.DefaultClient).Get(url) }

// Post creates a POST request builder using http.DefaultClient.
func Post(url string) Req { return Client(http.DefaultClient).Post(url) }

// Put creates a PUT request builder using http.DefaultClient.
func Put(url string) Req { return Client(http.DefaultClient).Put(url) }

// Patch creates a PATCH request builder using http.DefaultClient.
func Patch(url string) Req { return Client(http.DefaultClient).Patch(url) }

// Delete creates a DELETE request builder using http.DefaultClient.
func Delete(url string) Req { return Client(http.DefaultClient).Delete(url) }

// Head creates a HEAD request builder using http.DefaultClient.
func Head(url string) Req { return Client(http.DefaultClient).Head(url) }

// Options creates an OPTIONS request builder using http.DefaultClient.
func Options(url string) Req { return Client(http.DefaultClient).Options(url) }

// Trace creates a TRACE request builder using http.DefaultClient.
func Trace(url string) Req { return Client(http.DefaultClient).Trace(url) }

// Connect creates a CONNECT request builder using http.DefaultClient.
func Connect(url string) Req { return Client(http.DefaultClient).Connect(url) }

func (r Req) withMethod(method, url string, body io.Reader) Req {
	if r.err != nil {
		return r
	}
	client := r.v.client
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequest(method, url, body)
	return Req{Wrap(Request{req: req, client: client}, err)}
}

func (r Req) Get(url string) Req     { return r.withMethod(http.MethodGet, url, nil) }
func (r Req) Post(url string) Req    { return r.withMethod(http.MethodPost, url, nil) }
func (r Req) Put(url string) Req     { return r.withMethod(http.MethodPut, url, nil) }
func (r Req) Patch(url string) Req   { return r.withMethod(http.MethodPatch, url, nil) }
func (r Req) Delete(url string) Req  { return r.withMethod(http.MethodDelete, url, nil) }
func (r Req) Head(url string) Req    { return r.withMethod(http.MethodHead, url, nil) }
func (r Req) Options(url string) Req { return r.withMethod(http.MethodOptions, url, nil) }
func (r Req) Trace(url string) Req   { return r.withMethod(http.MethodTrace, url, nil) }
func (r Req) Connect(url string) Req { return r.withMethod(http.MethodConnect, url, nil) }

// WithClient sets the HTTP client.
func (r Req) WithClient(c *http.Client) Req {
	if r.err != nil {
		return r
	}
	if c == nil {
		c = http.DefaultClient
	}
	req := r.v.req
	return Req{Ok(Request{req: req, client: c})}
}

// WithHeader adds a header.
func (r Req) WithHeader(k, v string) Req {
	if r.err != nil {
		return r
	}
	if r.v.req == nil {
		return Req{Err[Request](errors.New("lambda/v2: request is not initialized"))}
	}
	rr := r.v
	rr.req.Header.Add(k, v)
	return Req{Ok(rr)}
}

// SetHeader sets a header (replaces existing values).
func (r Req) SetHeader(k, v string) Req {
	if r.err != nil {
		return r
	}
	if r.v.req == nil {
		return Req{Err[Request](errors.New("lambda/v2: request is not initialized"))}
	}
	rr := r.v
	rr.req.Header.Set(k, v)
	return Req{Ok(rr)}
}

// BasicAuth sets HTTP basic auth.
func (r Req) BasicAuth(user, password string) Req {
	if r.err != nil {
		return r
	}
	if r.v.req == nil {
		return Req{Err[Request](errors.New("lambda/v2: request is not initialized"))}
	}
	rr := r.v
	rr.req.SetBasicAuth(user, password)
	return Req{Ok(rr)}
}

// WithBody sets the request body (no copy). It also sets ContentLength when possible.
func (r Req) WithBody(body io.Reader) Req {
	if r.err != nil {
		return r
	}
	if r.v.req == nil {
		return Req{Err[Request](errors.New("lambda/v2: request is not initialized"))}
	}
	rr := r.v
	rr.req.Body = io.NopCloser(body)
	return Req{Ok(rr)}
}

// WithJSONBody marshals v as JSON and sets it as the request body, also setting Content-Type.
//
// Note: this takes `any` because Go does not allow methods with type parameters.
func (r Req) WithJSONBody(v any) Req {
	if r.err != nil {
		return r
	}
	if r.v.req == nil {
		return Req{Err[Request](errors.New("lambda/v2: request is not initialized"))}
	}
	b, err := json.Marshal(v)
	if err != nil {
		return Req{Err[Request](err)}
	}
	rr := r.v
	rr.req.Body = io.NopCloser(bytes.NewReader(b))
	rr.req.ContentLength = int64(len(b))
	rr.req.Header.Set("Content-Type", "application/json")
	return Req{Ok(rr)}
}

// Do executes the request using the configured client.
func (r Req) Do(ctx context.Context) Resp {
	if r.err != nil {
		return Resp{Err[*http.Response](r.err)}
	}
	if r.v.req == nil {
		return Resp{Err[*http.Response](errors.New("lambda/v2: request is not initialized"))}
	}
	client := r.v.client
	if client == nil {
		client = http.DefaultClient
	}
	req := r.v.req
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	res, err := client.Do(req)
	return Resp{Wrap(res, err)}
}

// Slurp reads the full response body and closes it.
func (r Resp) Slurp() Bytes {
	if r.err != nil {
		return Bytes{Err[[]byte](r.err)}
	}
	if r.v == nil {
		return Bytes{Err[[]byte](errors.New("lambda/v2: nil response"))}
	}
	if r.v.Body == nil {
		return Bytes{Err[[]byte](errors.New("lambda/v2: nil response body"))}
	}
	return Slurp(r.v.Body)
}

// StatusCode returns the response status code.
func (r Resp) StatusCode() Option[int] {
	if r.err != nil {
		return Err[int](r.err)
	}
	if r.v == nil {
		return Err[int](errors.New("lambda/v2: nil response"))
	}
	return Ok(r.v.StatusCode)
}

// Ensure bodies created by WithBody are non-nil.
var _ = errors.Join


