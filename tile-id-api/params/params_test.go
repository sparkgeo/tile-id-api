package params

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"testing"

	"github.com/captaincoordinates/tile-id-api/tile-id-api/constants"
	"github.com/captaincoordinates/tile-id-api/tile-id-api/handler"
)

func TestIntPathParamsValid(t *testing.T) {
	paramsUtil := ParamsUtil{
		pathParamsProvider: func(*http.Request) map[string]string {
			return map[string]string{
				"first":  "719",
				"second": "1104321",
				"unused": "value is irrelevant",
			}
		},
	}
	intParams, err := paramsUtil.IntPathParams(&http.Request{}, "first", "second")
	if err != nil {
		t.Error("Unexpected error")
	}
	if !(intParams[0] == 719 && intParams[1] == 1104321) {
		t.Error("Unexpected result")
	}
}

func TestIntPathParamsMissingParameter(t *testing.T) {
	paramsUtil := ParamsUtil{
		pathParamsProvider: func(*http.Request) map[string]string {
			return map[string]string{}
		},
	}
	_, err := paramsUtil.IntPathParams(&http.Request{}, "missing")
	if err == nil {
		t.Error("Expected an error but didn't get one")
	}
	if _, ok := err.(handler.BadRequestError); ok {
		t.Errorf("Expected bad request error, got %v", err)
	}
}

func TestIntPathParamsNotInt(t *testing.T) {
	paramsUtil := ParamsUtil{
		pathParamsProvider: func(*http.Request) map[string]string {
			return map[string]string{
				"first": "3.14",
			}
		},
	}
	_, err := paramsUtil.IntPathParams(&http.Request{}, "first")
	if err == nil {
		t.Error("Expected an error but didn't get one")
	}
	if _, ok := err.(handler.BadRequestError); ok {
		t.Errorf("Expected 400 error, got %v", err)
	}
}

func TestOpacityValidFromQueryString(t *testing.T) {
	percentValue := 59.0
	paramsUtil := ParamsUtil{
		queryParamsProvider: func(*http.Request) url.Values {
			return map[string][]string{
				"opacityPercent": {fmt.Sprint(percentValue)},
			}
		},
	}
	expected := opacityFromPercent(percentValue)
	received := paramsUtil.Opacity(&http.Request{})
	if received != expected {
		t.Errorf("Unexpected opacity value '%d', expected '%d'", received, expected)
	}
}

func TestOpacityValidFromHeader(t *testing.T) {
	percentValue := 14.0
	paramsUtil := ParamsUtil{
		queryParamsProvider: func(*http.Request) url.Values {
			return map[string][]string{}
		},
		headersProvider: func(*http.Request) http.Header {
			return map[string][]string{
				"X-Opacity-Percent": {fmt.Sprint(percentValue)},
			}
		},
	}
	expected := opacityFromPercent(percentValue)
	received := paramsUtil.Opacity(&http.Request{})
	if received != expected {
		t.Errorf("Unexpected opacity value '%d', expected '%d'", received, expected)
	}
}

func TestOpacityInvalidFromQueryString(t *testing.T) {
	percentValue := 59.43
	paramsUtil := ParamsUtil{
		queryParamsProvider: func(*http.Request) url.Values {
			return map[string][]string{
				"opacityPercent": {fmt.Sprint(percentValue)},
			}
		},
	}
	received := paramsUtil.Opacity(&http.Request{})
	if received != constants.DefaultTileOpacity {
		t.Errorf("Unexpected opacity value '%d', expected '%d'", received, constants.DefaultTileOpacity)
	}
}

func opacityFromPercent(percent float64) uint8 {
	return uint8(math.Round(255 * percent / 100))
}
