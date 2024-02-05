package handler

import (
	"net/http"
)

type TileHandler interface {
	Identifier() string
	PathPattern() string
	Keys(request *http.Request) (map[string]string, ReturnableError)
	AsZXY(request *http.Request) ([3]int, error)
}

func NewReturnableError(statusCode int, errorMessage string) ReturnableError {
	return ReturnableError{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	}
}

type ReturnableError struct {
	StatusCode   int
	ErrorMessage string
}

var NoReturnableError ReturnableError = ReturnableError{}
