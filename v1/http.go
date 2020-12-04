package v1

import (
	"io"
	"net/http"
)

type request struct {
	req    *http.Request
	client *http.Client
}

func Get(url string) Option {
	req, err := http.NewRequest("GET", url, nil)
	return Option{
		value: request{
			req:    req,
			client: http.DefaultClient,
		}, err: err,
	}
}

func Options(url string) Option {
	req, err := http.NewRequest("OPTIONS", url, nil)
	return Option{
		value: request{
			req:    req,
			client: http.DefaultClient,
		}, err: err,
	}
}

func Head(url string) Option {
	req, err := http.NewRequest("HEAD", url, nil)
	return Option{
		value: request{
			req:    req,
			client: http.DefaultClient,
		}, err: err,
	}
}

func Post(url string, body io.Reader) Option {
	req, err := http.NewRequest("POST", url, body)
	return Option{
		value: request{
			req:    req,
			client: http.DefaultClient,
		}, err: err,
	}
}

func Connect(url string) Option {
	req, err := http.NewRequest("CONNECT", url, nil)
	return Option{
		value: request{
			req:    req,
			client: http.DefaultClient,
		}, err: err,
	}
}

func Trace(url string) Option {
	req, err := http.NewRequest("TRACE", url, nil)
	return Option{
		value: request{
			req:    req,
			client: http.DefaultClient,
		}, err: err,
	}
}

func Patch(url string, body io.Reader) Option {
	req, err := http.NewRequest("PATCH", url, body)
	return Option{
		value: request{
			req:    req,
			client: http.DefaultClient,
		}, err: err,
	}
}

func Put(url string, body io.Reader) Option {
	req, err := http.NewRequest("PUT", url, body)
	return Option{
		value: request{
			req:    req,
			client: http.DefaultClient,
		}, err: err,
	}
}
func Delete(url string, body io.Reader) Option {
	req, err := http.NewRequest("DELETE", url, body)
	return Option{
		value: request{
			req:    req,
			client: http.DefaultClient,
		}, err: err,
	}
}

func (o Option) AddHeader(k, v string) Option {
	if o.err != nil {
		return o
	}
	r := o.value.(request)
	// This is fine, as long this request never leaves this option
	r.req.Header.Add(k, v)
	o.value = r
	return o
}
func (o Option) SetHeader(k, v string) Option {
	if o.err != nil {
		return o
	}
	r := o.value.(request)
	// This is fine, as long this request never leaves this option
	r.req.Header.Set(k, v)
	o.value = r
	return o
}

func (o Option) DeleteHeader(k string) Option {
	if o.err != nil {
		return o
	}

	r := o.value.(request)
	// This is fine, as long this request never leaves this option
	r.req.Header.Del(k)
	o.value = r
	return o
}

func (o Option) Header(k string) Option {
	if o.err != nil {
		return o
	}
	return WrapValue(o.value.(request).req.Header.Get(k))
}
func (o Option) Headers() Option {
	if o.err != nil {
		return o
	}
	return WrapValue(o.value.(request).req.Header.Clone)
}

func (o Option) Client(c *http.Client) Option {
	if o.err != nil {
		return o
	}
	r := o.value.(request)
	return WrapValue(request{
		req:    r.req,
		client: c,
	})

}

func (o Option) Do() Option {
	if o.err != nil {
		return o
	}
	r := o.value.(request)
	res, err := r.client.Do(r.req)
	return Wrap(res, err)
}
