package quadkey

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/captaincoordinates/tile-id-api/constants"
	"github.com/captaincoordinates/tile-id-api/handler"
	"github.com/captaincoordinates/tile-id-api/handler/common"
)

type QuadkeyTileHandler struct{}

func (self QuadkeyTileHandler) Identifier() string {
	return constants.QuadkeyIdentifier
}

func (self QuadkeyTileHandler) PathPattern() string {
	return "{quadkey:[0-3]{0,25}}"
}

func (self QuadkeyTileHandler) GetKeyProvider(request *http.Request) (handler.TileHandlerKeyProvider, handler.ReturnableError) {
	quadkey := mux.Vars(request)["quadkey"]
	return func(
		identifier string,
	) (key string) {
		zxy, err := common.QuadkeyToZxy(quadkey)
		if err != nil {
			fmt.Println(fmt.Sprintf("Validated quadkey that could not be converted: '%s'", quadkey))
			return "unknown"
		}
		z, x, y := zxy[0], zxy[1], zxy[2]
		switch identifier {
		case constants.QuadkeyIdentifier:
			return quadkey
		case constants.ZxyIdentifier:
			return fmt.Sprintf("%d/%d/%d", z, x, y)
		case constants.TmsIdentifier:
			return fmt.Sprintf("%d/%d/%d", z, x, common.FlipY(z, y))
		default:
			panic(fmt.Sprintf("Unknown identifier %s", identifier))
		}
	}, handler.NoReturnableError
}

func (self QuadkeyTileHandler) AsZXY(request *http.Request) ([3]int, error) {
	quadkey := mux.Vars(request)["quadkey"]
	zxy, err := common.QuadkeyToZxy(quadkey)
	if err != nil {
		return [3]int{}, errors.New(err.Error())
	}
	return [3]int{zxy[0], zxy[1], zxy[2]}, nil
}
