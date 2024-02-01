package quadkey

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/captaincoordinates/tile-id-api/handler"
	"github.com/captaincoordinates/tile-id-api/handler/common"
)

type QuadkeyTileHandler struct{}

func (self QuadkeyTileHandler) Identifier() string {
	return common.QuadkeyIdentifier
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
		case common.QuadkeyIdentifier:
			return quadkey
		case common.ZxyIdentifier:
			return fmt.Sprintf("%d/%d/%d", z, x, y)
		case common.TmsIdentifier:
			return fmt.Sprintf("%d/%d/%d", z, x, common.FlipY(z, y))
		default:
			panic(fmt.Sprintf("Unknown identifier %s", identifier))
		}
	}, handler.NoReturnableError
}
