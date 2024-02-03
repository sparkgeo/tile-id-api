package tms

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/captaincoordinates/tile-id-api/handler"
	"github.com/captaincoordinates/tile-id-api/handler/common"
	"github.com/captaincoordinates/tile-id-api/params"
)

type TmsTileHandler struct{}

func (self TmsTileHandler) Identifier() string {
	return common.TmsIdentifier
}

func (self TmsTileHandler) PathPattern() string {
	return common.ZxyPathPattern
}

func (self TmsTileHandler) GetKeyProvider(request *http.Request) (handler.TileHandlerKeyProvider, handler.ReturnableError) {
	params, err := params.FetchIntPathParams(request, "z", "x", "y")
	if err != handler.NoReturnableError {
		return nil, err
	}
	z, x, y := params[0], params[1], params[2]
	return func(
		identifier string,
	) (key string) {
		switch identifier {
		case common.TmsIdentifier:
			return fmt.Sprintf("%d/%d/%d", z, x, y)
		case common.ZxyIdentifier:
			return fmt.Sprintf("%d/%d/%d", z, x, common.FlipY(z, y))
		case common.QuadkeyIdentifier:
			return common.ZxyToQuadkey(z, x, common.FlipY(z, y))
		default:
			panic(fmt.Sprintf("Unknown identifier %s", identifier))
		}
	}, handler.NoReturnableError
}

func (self TmsTileHandler) AsZXY(request *http.Request) ([3]int, error) {
	params, err := params.FetchIntPathParams(request, "z", "x", "y")
	if err != handler.NoReturnableError {
		return [3]int{}, errors.New(err.ErrorMessage)
	}
	return [3]int{params[0], params[1], common.FlipY(params[0], params[2])}, nil
}
