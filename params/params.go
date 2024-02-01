package params

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"

	"github.com/captaincoordinates/tile-id-api/constants"
	"github.com/captaincoordinates/tile-id-api/handler"
	"github.com/gorilla/mux"
)

func FetchIntPathParams(
	request *http.Request,
	paramNames ...string,
) ([]int, handler.ReturnableError) {
	vars := mux.Vars(request)
	parsedInts := make([]int, len(paramNames))
	for i, paramName := range paramNames {
		paramStr, ok := vars[paramName]
		if !ok {
			return make([]int, 0), handler.ReturnableError{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("Missing expected parameter %s", paramName),
			}
		}
		paramInt, err := strconv.ParseInt(paramStr, 10, 64)
		if err != nil {
			return make([]int, 0), handler.ReturnableError{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("Parameter %s cannot be parsed to int (value %s)", paramName, paramStr),
			}
		}
		parsedInts[i] = int(paramInt)
	}
	return parsedInts, handler.NoReturnableError
}

func Opacity(request *http.Request) uint8 {
	opacityStr := request.URL.Query().Get("opacityPercent") // case-sensitive
	if opacityStr == "" {
		opacityStr = request.Header.Get("X-Opacity-Percent") // case-insensitive
		if opacityStr == "" {
			return constants.DefaultTileOpacity
		}
	}
	opacityRegex := regexp.MustCompile("^(100|[1-9]?[0-9])$")
	if opacityRegex.MatchString(opacityStr) {
		opacityPercent, err := strconv.ParseUint(opacityStr, 10, 64)
		if err != nil {
			fmt.Println(fmt.Sprintf("Unable to parse requested opacity '%s'", opacityStr))
			return constants.DefaultTileOpacity
		}
		return uint8(math.Round(float64(opacityPercent) * float64(255) / float64(100)))
	} else {
		return constants.DefaultTileOpacity
	}
}
