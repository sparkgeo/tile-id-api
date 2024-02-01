package handler

import (
	"net/http"
)

type TileHandler interface {
	Identifier() string
	PathPattern() string
	GetKeyProvider(request *http.Request) (TileHandlerKeyProvider, ReturnableError)
}

type TileHandlerKeyProvider func(
	identifier string,
) (key string)

type ReturnableError struct {
	StatusCode   int
	ErrorMessage string
}

var NoReturnableError ReturnableError = ReturnableError{}
