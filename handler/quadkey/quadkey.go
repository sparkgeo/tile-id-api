package quadkey

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/captaincoordinates/tile-id-api/config"
	"github.com/captaincoordinates/tile-id-api/constants"
	"github.com/captaincoordinates/tile-id-api/handler"
	"github.com/captaincoordinates/tile-id-api/handler/common"
	"github.com/captaincoordinates/tile-id-api/log"
)

var logger = log.NewLogger(config.DefaultEnvGetter)

type QuadkeyTileHandler struct {
	flipYProvider         common.FlipYProvider
	quadkeyToZxyProvider  common.QuadkeyToZxyProvider
	pathParamsMapProvider func(*http.Request) map[string]string
}

func NewQuadkeyTileHandler() *QuadkeyTileHandler {
	return &QuadkeyTileHandler{
		flipYProvider:         common.FlipY,
		quadkeyToZxyProvider:  common.QuadkeyToZxy,
		pathParamsMapProvider: mux.Vars,
	}
}

func (self QuadkeyTileHandler) Identifier() string {
	return constants.QuadkeyIdentifier
}

func (self QuadkeyTileHandler) PathPattern() string {
	return "{quadkey:[0-3]{0,25}}"
}

func (self QuadkeyTileHandler) Keys(request *http.Request) (map[string]string, handler.ReturnableError) {
	quadkey := self.pathParamsMapProvider(request)["quadkey"]
	zxy, err := self.quadkeyToZxyProvider(quadkey)
	if err != nil {
		logger.Warn(fmt.Sprintf("Validated quadkey that could not be converted: '%s'", quadkey))
		return nil, handler.NewReturnableError(
			422,
			fmt.Sprintf("Could not convert '%s' to required format", quadkey),
		)
	}
	z, x, y := zxy[0], zxy[1], zxy[2]
	return map[string]string{
		constants.QuadkeyIdentifier: quadkey,
		constants.ZxyIdentifier:     fmt.Sprintf("%d/%d/%d", z, x, y),
		constants.TmsIdentifier:     fmt.Sprintf("%d/%d/%d", z, x, self.flipYProvider(z, y)),
	}, handler.NoReturnableError
}

func (self QuadkeyTileHandler) AsZXY(request *http.Request) ([3]int, error) {
	quadkey := self.pathParamsMapProvider(request)["quadkey"]
	zxy, err := self.quadkeyToZxyProvider(quadkey)
	if err != nil {
		return [3]int{}, errors.New(err.Error())
	}
	return [3]int{zxy[0], zxy[1], zxy[2]}, nil
}
