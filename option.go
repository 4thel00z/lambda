package lambda

type Option struct {
	value interface{}
	err   error
}

type Producer func(o Option) interface{}
type ErrorHandler func(err error)

func (o Option) Or(i interface{}) interface{} {
	if o.err != nil {
		return i
	}
	return nil
}

func (o Option) Unwrap() interface{} {
	if o.err != nil {
		panic(o.err)
	}
	return o.value
}

func (o Option) Catch(e ErrorHandler) interface{} {
	if o.err != nil {
		e(o.err)
		return nil
	}
	return o.value
}
