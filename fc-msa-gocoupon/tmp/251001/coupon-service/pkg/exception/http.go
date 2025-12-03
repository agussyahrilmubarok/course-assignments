package exception

import "net/http"

type Http struct {
	Code    int
	Message string
	Err     error
}

func (e *Http) Error() string {
	return e.Message
}

func NewBadRequest(msg string, err error) *Http {
	return &Http{Code: http.StatusBadRequest, Message: msg, Err: err}
}

func NewUnauthorized(msg string, err error) *Http {
	return &Http{Code: http.StatusUnauthorized, Message: msg, Err: err}
}

func NewForbidden(msg string, err error) *Http {
	return &Http{Code: http.StatusForbidden, Message: msg, Err: err}
}

func NewNotFound(msg string, err error) *Http {
	return &Http{Code: http.StatusNotFound, Message: msg, Err: err}
}

func NewMethodNotAllowed(msg string, err error) *Http {
	return &Http{Code: http.StatusMethodNotAllowed, Message: msg, Err: err}
}

func NewConflict(msg string, err error) *Http {
	return &Http{Code: http.StatusConflict, Message: msg, Err: err}
}

func NewUnprocessableEntity(msg string, err error) *Http {
	return &Http{Code: http.StatusUnprocessableEntity, Message: msg, Err: err}
}

func NewTooManyRequests(msg string, err error) *Http {
	return &Http{Code: http.StatusTooManyRequests, Message: msg, Err: err}
}

func NewRequestTimeout(msg string, err error) *Http {
	return &Http{Code: http.StatusRequestTimeout, Message: msg, Err: err}
}

func NewInternal(msg string, err error) *Http {
	return &Http{Code: http.StatusInternalServerError, Message: msg, Err: err}
}

func NewBadGateway(msg string, err error) *Http {
	return &Http{Code: http.StatusBadGateway, Message: msg, Err: err}
}

func NewServiceUnavailable(msg string, err error) *Http {
	return &Http{Code: http.StatusServiceUnavailable, Message: msg, Err: err}
}

func NewGatewayTimeout(msg string, err error) *Http {
	return &Http{Code: http.StatusGatewayTimeout, Message: msg, Err: err}
}
