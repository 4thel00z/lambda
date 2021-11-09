package v1

import (
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
)

type RequestWrapper struct {
	Request *http.Request
	Client  *http.Client
}

func (o Option) Get(url string) Option {
	req, err := http.NewRequest("GET", url, nil)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  o.value.(RequestWrapper).Client,
		}, err: err,
	}
}

func (o Option) Options(url string) Option {
	req, err := http.NewRequest("OPTIONS", url, nil)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  o.value.(RequestWrapper).Client,
		}, err: err,
	}
}

func (o Option) Head(url string) Option {
	req, err := http.NewRequest("HEAD", url, nil)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  o.value.(RequestWrapper).Client,
		}, err: err,
	}
}

func (o Option) Post(url string, body io.Reader) Option {
	req, err := http.NewRequest("POST", url, body)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  o.value.(RequestWrapper).Client,
		}, err: err,
	}
}

func (o Option) Connect(url string) Option {
	req, err := http.NewRequest("CONNECT", url, nil)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  o.value.(RequestWrapper).Client,
		}, err: err,
	}
}

func (o Option) Trace(url string) Option {
	req, err := http.NewRequest("TRACE", url, nil)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  o.value.(RequestWrapper).Client,
		}, err: err,
	}
}

func (o Option) Patch(url string, body io.Reader) Option {
	req, err := http.NewRequest("PATCH", url, body)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  o.value.(RequestWrapper).Client,
		}, err: err,
	}
}

func (o Option) Put(url string, body io.Reader) Option {
	req, err := http.NewRequest("PUT", url, body)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  o.value.(RequestWrapper).Client,
		}, err: err,
	}
}
func (o Option) Delete(url string, body io.Reader) Option {
	req, err := http.NewRequest("DELETE", url, body)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  o.value.(RequestWrapper).Client,
		}, err: err,
	}
}
func Get(url string) Option {
	req, err := http.NewRequest("GET", url, nil)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  http.DefaultClient,
		}, err: err,
	}
}

func Options(url string) Option {
	req, err := http.NewRequest("OPTIONS", url, nil)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  http.DefaultClient,
		}, err: err,
	}
}

func Head(url string) Option {
	req, err := http.NewRequest("HEAD", url, nil)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  http.DefaultClient,
		}, err: err,
	}
}

func Post(url string, body io.Reader) Option {
	req, err := http.NewRequest("POST", url, body)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  http.DefaultClient,
		}, err: err,
	}
}

func Connect(url string) Option {
	req, err := http.NewRequest("CONNECT", url, nil)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  http.DefaultClient,
		}, err: err,
	}
}

func Trace(url string) Option {
	req, err := http.NewRequest("TRACE", url, nil)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  http.DefaultClient,
		}, err: err,
	}
}

func Patch(url string, body io.Reader) Option {
	req, err := http.NewRequest("PATCH", url, body)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  http.DefaultClient,
		}, err: err,
	}
}

func Put(url string, body io.Reader) Option {
	req, err := http.NewRequest("PUT", url, body)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  http.DefaultClient,
		}, err: err,
	}
}
func Delete(url string, body io.Reader) Option {
	req, err := http.NewRequest("DELETE", url, body)
	return Option{
		value: RequestWrapper{
			Request: req,
			Client:  http.DefaultClient,
		}, err: err,
	}
}

func (o Option) AddHeader(k, v string) Option {
	if o.err != nil {
		return o
	}
	r := o.value.(RequestWrapper)
	// This is fine, as long this RequestWrapper never leaves this option
	r.Request.Header.Add(k, v)
	o.value = r
	return o
}

func (o Option) BasicAuth(user, password string) Option {
	if o.err != nil {
		return o
	}
	r := o.value.(RequestWrapper)
	auth := []byte(user)
	auth = append(auth, []byte(":")...)
	auth = append(auth, []byte(password)...)
	// This is fine, as long this RequestWrapper never leaves this option
	r.Request.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString(auth))
	o.value = r
	return o
}

func (o Option) SetHeader(k, v string) Option {
	if o.err != nil {
		return o
	}
	r := o.value.(RequestWrapper)
	// This is fine, as long this RequestWrapper never leaves this option
	r.Request.Header.Set(k, v)
	o.value = r
	return o
}

func (o Option) DeleteHeader(k string) Option {
	if o.err != nil {
		return o
	}

	r := o.value.(RequestWrapper)
	// This is fine, as long this RequestWrapper never leaves this option
	r.Request.Header.Del(k)
	o.value = r
	return o
}

func (o Option) Header(k string) Option {
	if o.err != nil {
		return o
	}
	return WrapValue(o.value.(RequestWrapper).Request.Header.Get(k))
}
func (o Option) Headers() Option {
	if o.err != nil {
		return o
	}
	return WrapValue(o.value.(RequestWrapper).Request.Header.Clone)
}

func (o Option) Client(c *http.Client) Option {
	if o.err != nil {
		return o
	}
	r := o.value.(RequestWrapper)
	return WrapValue(RequestWrapper{
		Request: r.Request,
		Client:  c,
	})

}

func Client(c *http.Client) Option {
	return WrapValue(RequestWrapper{
		Request: nil,
		Client:  c,
	})
}

func (o Option) Do() Option {
	if o.err != nil {
		return o
	}
	r := o.value.(RequestWrapper)
	if r.Request == nil {
		return WrapError(errors.New("RequestWrapper is empty"))
	}
	res, err := r.Client.Do(r.Request)
	return Wrap(res, err)
}

func (o Option) UnwrapRequestWrapper() RequestWrapper {
	if o.err != nil {
		log.Fatalln(o.err)
	}
	return o.value.(RequestWrapper)
}
