package zxy

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sparkgeo/tile-id-api/tile-id-api/internal/constants"
	"github.com/sparkgeo/tile-id-api/tile-id-api/internal/handler/common"
)

func TestIdentifier(t *testing.T) {
	logger := logrus.New()
	if (NewZxyTileHandler(logger)).Identifier() != constants.ZxyIdentifier {
		t.Error("Unexpected identifier string")
	}
}

func TestPathPattern(t *testing.T) {
	logger := logrus.New()
	if (NewZxyTileHandler(logger)).PathPattern() != common.ZxyPathPattern {
		t.Error("Unexpected path pattern string")
	}
}

func TestKeysValid(t *testing.T) {
	z, x, y := 1, 2, 3
	tileHandler := ZxyTileHandler{
		intPathParamsProvider: getIntPathParamsProvider(z, x, y),
		flipYProvider:         flipYProvider,
		zxyToQuadkeyProvider:  zxyToQuadkeyProvider,
	}
	keysByIdentifier, _ := tileHandler.Keys(&http.Request{})
	for identifier, expectedKey := range map[string]string{
		constants.ZxyIdentifier:     "1/2/3",
		constants.QuadkeyIdentifier: "1_2_3",
		constants.TmsIdentifier:     "1/2/9",
	} {
		providedKey := keysByIdentifier[identifier]
		if providedKey != expectedKey {
			t.Errorf("'%s' produced incorrect key '%s'. Expected '%s'", identifier, providedKey, expectedKey)
		}
	}
}

func TestKeysInvalidRequest(t *testing.T) {
	intPathParamsProvider := func(*http.Request, ...string) ([]int, error) {
		return nil, errors.New("Spurious error")
	}
	tileHandler := ZxyTileHandler{
		intPathParamsProvider: intPathParamsProvider,
	}
	_, keysError := tileHandler.Keys(&http.Request{})
	if keysError == nil {
		t.Error("Dependency raised an error but this was not propagated to caller")
	}
}

func TestAsZXYValid(t *testing.T) {
	z, x, y := 1, 2, 3
	tileHandler := ZxyTileHandler{
		intPathParamsProvider: getIntPathParamsProvider(z, x, y),
	}
	providedZxy, _ := tileHandler.AsZXY(&http.Request{})
	expectedZxy := [3]int{1, 2, 3}
	if !reflect.DeepEqual(providedZxy, expectedZxy) {
		t.Errorf("unexpected values for z, x, y: %d, %d, %d", providedZxy[0], providedZxy[1], providedZxy[2])
	}
}

func TestAsZXYInvalidRequest(t *testing.T) {
	intPathParamsProvider := func(*http.Request, ...string) ([]int, error) {
		return nil, errors.New("Spurious error")
	}
	tileHandler := ZxyTileHandler{
		intPathParamsProvider: intPathParamsProvider,
	}
	_, err := tileHandler.AsZXY(&http.Request{})
	if err == nil {
		t.Error("Dependency raised an error but this was not propagated to caller")
	}
}

func getIntPathParamsProvider(z, x, y int) func(*http.Request, ...string) ([]int, error) {
	return func(
		request *http.Request,
		paramNames ...string,
	) ([]int, error) {
		return []int{z, x, y}, nil
	}
}

func flipYProvider(z, y int) int {
	return int(math.Round(math.Pow(float64(y), 2)))
}

func zxyToQuadkeyProvider(z, x, y int) string {
	return fmt.Sprintf("%d_%d_%d", z, x, y)
}
