package handler

import (
	"net/http"
)

type TileHandler interface {
	Identifier() string
	PathPattern() string
	GetKeyProvider(request *http.Request) (TileHandlerKeyProvider, ReturnableError)
	AsZXY(request *http.Request) ([3]int, error)
}

type TileHandlerKeyProvider func(
	identifier string,
) (key string)

type ReturnableError struct {
	StatusCode   int
	ErrorMessage string
}

var NoReturnableError ReturnableError = ReturnableError{}
