package tms

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/sparkgeo/tile-id-api/tile-id-api/internal/constants"
	"github.com/sparkgeo/tile-id-api/tile-id-api/internal/handler/common"
	"github.com/sparkgeo/tile-id-api/tile-id-api/internal/params"
)

type TmsTileHandler struct {
	intPathParamsProvider params.IntPathParamsProvider
	flipYProvider         common.FlipYProvider
	zxyToQuadkeyProvider  common.ZxyToQuadkeyProvider
}

func NewTmsTileHandler(logger logrus.FieldLogger) *TmsTileHandler {
	paramsUtil := params.NewParamsUtil(logger)
	return &TmsTileHandler{
		intPathParamsProvider: paramsUtil.IntPathParams,
		flipYProvider:         common.FlipY,
		zxyToQuadkeyProvider:  common.ZxyToQuadkey,
	}
}

func (tth TmsTileHandler) Identifier() string {
	return constants.TmsIdentifier
}

func (tth TmsTileHandler) PathPattern() string {
	return common.ZxyPathPattern
}

func (tth TmsTileHandler) Keys(request *http.Request) (map[string]string, error) {
	pathParams, err := tth.intPathParamsProvider(request, "z", "x", "y")
	if err != nil {
		return nil, err
	}
	z, x, y := pathParams[0], pathParams[1], pathParams[2]
	return map[string]string{
		constants.TmsIdentifier:     fmt.Sprintf("%d/%d/%d", z, x, y),
		constants.ZxyIdentifier:     fmt.Sprintf("%d/%d/%d", z, x, tth.flipYProvider(z, y)),
		constants.QuadkeyIdentifier: tth.zxyToQuadkeyProvider(z, x, tth.flipYProvider(z, y)),
	}, nil
}

func (tth TmsTileHandler) AsZXY(request *http.Request) ([3]int, error) {
	params, err := tth.intPathParamsProvider(request, "z", "x", "y")
	if err != nil {
		return [3]int{}, err
	}
	return [3]int{params[0], params[1], tth.flipYProvider(params[0], params[2])}, nil
}
