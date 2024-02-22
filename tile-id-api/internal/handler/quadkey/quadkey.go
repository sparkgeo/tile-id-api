package quadkey

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/sparkgeo/tile-id-api/tile-id-api/internal/constants"
	"github.com/sparkgeo/tile-id-api/tile-id-api/internal/handler"
	"github.com/sparkgeo/tile-id-api/tile-id-api/internal/handler/common"
)

type QuadkeyTileHandler struct {
	flipYProvider         common.FlipYProvider
	quadkeyToZxyProvider  common.QuadkeyToZxyProvider
	pathParamsMapProvider func(*http.Request) map[string]string
	logger                logrus.FieldLogger
}

func NewQuadkeyTileHandler(logger logrus.FieldLogger) *QuadkeyTileHandler {
	return &QuadkeyTileHandler{
		flipYProvider:         common.FlipY,
		quadkeyToZxyProvider:  common.QuadkeyToZxy,
		pathParamsMapProvider: mux.Vars,
		logger:                logger,
	}
}

func (qth QuadkeyTileHandler) Identifier() string {
	return constants.QuadkeyIdentifier
}

func (qth QuadkeyTileHandler) PathPattern() string {
	return "{quadkey:[0-3]{0,25}}"
}

func (qth QuadkeyTileHandler) Keys(request *http.Request) (map[string]string, error) {
	quadkey := qth.pathParamsMapProvider(request)["quadkey"]
	zxy, err := qth.quadkeyToZxyProvider(quadkey)
	if err != nil {
		qth.logger.Warn(fmt.Sprintf("Validated quadkey that could not be converted: '%s'", quadkey))
		return nil, handler.NewUnprocessableEntityError(
			fmt.Sprintf("Could not convert '%s' to required format", quadkey),
		)
	}
	z, x, y := zxy[0], zxy[1], zxy[2]
	return map[string]string{
		constants.QuadkeyIdentifier: quadkey,
		constants.ZxyIdentifier:     fmt.Sprintf("%d/%d/%d", z, x, y),
		constants.TmsIdentifier:     fmt.Sprintf("%d/%d/%d", z, x, qth.flipYProvider(z, y)),
	}, nil
}

func (qth QuadkeyTileHandler) AsZXY(request *http.Request) ([3]int, error) {
	quadkey := qth.pathParamsMapProvider(request)["quadkey"]
	zxy, err := qth.quadkeyToZxyProvider(quadkey)
	if err != nil {
		return [3]int{}, errors.New(err.Error())
	}
	return [3]int{zxy[0], zxy[1], zxy[2]}, nil
}
