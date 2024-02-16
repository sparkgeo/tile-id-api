package zxy

import (
	"fmt"
	"net/http"

	"github.com/captaincoordinates/tile-id-api/tile-id-api/internal/constants"
	"github.com/captaincoordinates/tile-id-api/tile-id-api/internal/handler/common"
	"github.com/captaincoordinates/tile-id-api/tile-id-api/internal/params"
	"github.com/sirupsen/logrus"
)

type ZxyTileHandler struct {
	intPathParamsProvider params.IntPathParamsProvider
	flipYProvider         common.FlipYProvider
	zxyToQuadkeyProvider  common.ZxyToQuadkeyProvider
}

func NewZxyTileHandler(logger logrus.FieldLogger) *ZxyTileHandler {
	paramsUtil := params.NewParamsUtil(logger)
	return &ZxyTileHandler{
		intPathParamsProvider: paramsUtil.IntPathParams,
		flipYProvider:         common.FlipY,
		zxyToQuadkeyProvider:  common.ZxyToQuadkey,
	}
}

func (zth ZxyTileHandler) Identifier() string {
	return constants.ZxyIdentifier
}

func (zth ZxyTileHandler) PathPattern() string {
	return common.ZxyPathPattern
}

func (zth ZxyTileHandler) Keys(request *http.Request) (map[string]string, error) {
	pathParams, err := zth.intPathParamsProvider(request, "z", "x", "y")
	if err != nil {
		return nil, err
	}
	z, x, y := pathParams[0], pathParams[1], pathParams[2]
	return map[string]string{
		constants.ZxyIdentifier:     fmt.Sprintf("%d/%d/%d", z, x, y),
		constants.TmsIdentifier:     fmt.Sprintf("%d/%d/%d", z, x, zth.flipYProvider(z, y)),
		constants.QuadkeyIdentifier: zth.zxyToQuadkeyProvider(z, x, y),
	}, nil
}

func (zth ZxyTileHandler) AsZXY(request *http.Request) ([3]int, error) {
	params, err := zth.intPathParamsProvider(request, "z", "x", "y")
	if err != nil {
		return [3]int{}, err
	}
	return [3]int{params[0], params[1], params[2]}, nil
}
