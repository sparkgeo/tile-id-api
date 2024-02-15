package zxy

import (
	"fmt"
	"net/http"

	"github.com/captaincoordinates/tile-id-api/tile-id-api/constants"
	"github.com/captaincoordinates/tile-id-api/tile-id-api/handler/common"
	"github.com/captaincoordinates/tile-id-api/tile-id-api/params"
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

func (self ZxyTileHandler) Identifier() string {
	return constants.ZxyIdentifier
}

func (self ZxyTileHandler) PathPattern() string {
	return common.ZxyPathPattern
}

func (self ZxyTileHandler) Keys(request *http.Request) (map[string]string, error) {
	pathParams, err := self.intPathParamsProvider(request, "z", "x", "y")
	if err != nil {
		return nil, err
	}
	z, x, y := pathParams[0], pathParams[1], pathParams[2]
	return map[string]string{
		constants.ZxyIdentifier:     fmt.Sprintf("%d/%d/%d", z, x, y),
		constants.TmsIdentifier:     fmt.Sprintf("%d/%d/%d", z, x, self.flipYProvider(z, y)),
		constants.QuadkeyIdentifier: self.zxyToQuadkeyProvider(z, x, y),
	}, nil
}

func (self ZxyTileHandler) AsZXY(request *http.Request) ([3]int, error) {
	params, err := self.intPathParamsProvider(request, "z", "x", "y")
	if err != nil {
		return [3]int{}, err
	}
	return [3]int{params[0], params[1], params[2]}, nil
}
