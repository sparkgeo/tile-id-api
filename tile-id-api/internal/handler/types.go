package handler

import (
	"net/http"
)

type TileHandler interface {
	Identifier() string
	PathPattern() string
	Keys(request *http.Request) (map[string]string, error)
	AsZXY(request *http.Request) ([3]int, error)
}
