package params

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/captaincoordinates/tile-id-api/config"
	"github.com/captaincoordinates/tile-id-api/constants"
	"github.com/captaincoordinates/tile-id-api/handler"
	"github.com/captaincoordinates/tile-id-api/log"
	"github.com/gorilla/mux"
)

var logger = log.NewLogger(config.DefaultEnvGetter)

type IntPathParamsProvider = func(request *http.Request, paramNames ...string) ([]int, handler.ReturnableError)

type ParamsUtil struct {
	pathParamsProvider  func(*http.Request) map[string]string
	queryParamsProvider func(*http.Request) url.Values
	headersProvider     func(*http.Request) http.Header
}

func NewParamsUtil() ParamsUtil {
	return ParamsUtil{
		pathParamsProvider: mux.Vars,
		queryParamsProvider: func(request *http.Request) url.Values {
			return request.URL.Query()
		},
		headersProvider: func(request *http.Request) http.Header {
			return request.Header
		},
	}
}

func (self ParamsUtil) IntPathParams(
	request *http.Request,
	paramNames ...string,
) ([]int, handler.ReturnableError) {
	vars := self.pathParamsProvider(request)
	parsedInts := make([]int, len(paramNames))
	for i, paramName := range paramNames {
		paramStr, ok := vars[paramName]
		if !ok {
			return make([]int, 0), handler.NewReturnableError(
				http.StatusBadRequest,
				fmt.Sprintf("Missing expected parameter %s", paramName),
			)
		}
		paramInt, err := strconv.ParseInt(paramStr, 10, 64)
		if err != nil {
			return make([]int, 0), handler.NewReturnableError(
				http.StatusBadRequest,
				fmt.Sprintf("Parameter %s cannot be parsed to int (value %s)", paramName, paramStr),
			)
		}
		parsedInts[i] = int(paramInt)
	}
	return parsedInts, handler.NoReturnableError
}

func (self ParamsUtil) Opacity(request *http.Request) uint8 {
	opacityStr := self.queryParamsProvider(request).Get("opacityPercent") // case-sensitive
	if opacityStr == "" {
		opacityStr = self.headersProvider(request).Get("X-Opacity-Percent") // case-insensitive
		if opacityStr == "" {
			return constants.DefaultTileOpacity
		}
	}
	opacityRegex := regexp.MustCompile("^(100|[1-9]?[0-9])$")
	if opacityRegex.MatchString(opacityStr) {
		opacityPercent, err := strconv.ParseUint(opacityStr, 10, 64)
		if err != nil {
			logger.Debug(fmt.Sprintf("Unable to parse requested opacity '%s'", opacityStr))
			return constants.DefaultTileOpacity
		}
		return uint8(math.Round(float64(opacityPercent) * float64(255) / float64(100)))
	} else {
		return constants.DefaultTileOpacity
	}
}
