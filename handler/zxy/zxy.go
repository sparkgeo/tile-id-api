package zxy

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/captaincoordinates/tile-id-api/constants"
	"github.com/captaincoordinates/tile-id-api/handler"
	"github.com/captaincoordinates/tile-id-api/handler/common"
	"github.com/captaincoordinates/tile-id-api/params"
)

type ZxyTileHandler struct{}

func (self ZxyTileHandler) Identifier() string {
	return constants.ZxyIdentifier
}

func (self ZxyTileHandler) PathPattern() string {
	return common.ZxyPathPattern
}

func (self ZxyTileHandler) GetKeyProvider(request *http.Request) (handler.TileHandlerKeyProvider, handler.ReturnableError) {
	params, err := params.FetchIntPathParams(request, "z", "x", "y")
	if err != handler.NoReturnableError {
		return nil, err
	}
	z, x, y := params[0], params[1], params[2]
	return func(
		identifier string,
	) (key string) {
		switch identifier {
		case constants.ZxyIdentifier:
			return fmt.Sprintf("%d/%d/%d", z, x, y)
		case constants.TmsIdentifier:
			return fmt.Sprintf("%d/%d/%d", z, x, common.FlipY(z, y))
		case constants.QuadkeyIdentifier:
			return common.ZxyToQuadkey(z, x, y)
		default:
			panic(fmt.Sprintf("Unknown identifier %s", identifier))
		}
	}, handler.NoReturnableError
}

func (self ZxyTileHandler) AsZXY(request *http.Request) ([3]int, error) {
	params, err := params.FetchIntPathParams(request, "z", "x", "y")
	if err != handler.NoReturnableError {
		return [3]int{}, errors.New(err.ErrorMessage)
	}
	return [3]int{params[0], params[1], params[2]}, nil
}
