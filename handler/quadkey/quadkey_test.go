package quadkey

import (
	"errors"
	"math"
	"net/http"
	"reflect"
	"testing"

	"github.com/captaincoordinates/tile-id-api/constants"
	"github.com/captaincoordinates/tile-id-api/handler"
)

func TestIdentifier(t *testing.T) {
	if (NewQuadkeyTileHandler()).Identifier() != constants.QuadkeyIdentifier {
		t.Error("Unexpected identifier string")
	}
}

func TestKeysValid(t *testing.T) {
	quadkey := "0123"
	z, x, y := 9, 10, 11
	tileHandler := QuadkeyTileHandler{
		flipYProvider:         flipYProvider,
		quadkeyToZxyProvider:  getQuadkeyToZxyProvider(z, x, y),
		pathParamsMapProvider: getPathParamsMapProvider(quadkey),
	}
	keysByIdentifier, _ := tileHandler.Keys(&http.Request{})
	for identifier, expectedKey := range map[string]string{
		constants.QuadkeyIdentifier: quadkey,
		constants.TmsIdentifier:     "9/10/121",
		constants.ZxyIdentifier:     "9/10/11",
	} {
		providedKey := keysByIdentifier[identifier]
		if providedKey != expectedKey {
			t.Errorf("'%s' produced incorrect key '%s'. Expected '%s'", identifier, providedKey, expectedKey)
		}
	}
}

func TestKeysInvalidRequest(t *testing.T) {
	tileHandler := QuadkeyTileHandler{
		pathParamsMapProvider: getPathParamsMapProvider("0123"),
		quadkeyToZxyProvider:  quadkeyToZxyProviderError,
	}
	_, keysError := tileHandler.Keys(&http.Request{})
	if keysError == handler.NoReturnableError {
		t.Error("Dependency raised an error but this was not propagated to caller")
	}
}

func TestAsZXYValid(t *testing.T) {
	z, x, y := 9, 10, 11
	tileHandler := QuadkeyTileHandler{
		pathParamsMapProvider: getPathParamsMapProvider("0123"),
		quadkeyToZxyProvider:  getQuadkeyToZxyProvider(z, x, y),
	}
	providedZxy, _ := tileHandler.AsZXY(&http.Request{})
	expectedZxy := [3]int{z, x, y}
	if !reflect.DeepEqual(providedZxy, expectedZxy) {
		t.Errorf("unexpected values for z, x, y: %d, %d, %d", providedZxy[0], providedZxy[1], providedZxy[2])
	}
}

func TestAsZXYInvalidRequest(t *testing.T) {
	tileHandler := QuadkeyTileHandler{
		pathParamsMapProvider: getPathParamsMapProvider("0123"),
		quadkeyToZxyProvider:  quadkeyToZxyProviderError,
	}
	_, err := tileHandler.AsZXY(&http.Request{})
	if err == nil {
		t.Error("Dependency raised an error but this was not propagated to caller")
	}
}

func getPathParamsMapProvider(quadkeyValue string) func(*http.Request) map[string]string {
	return func(*http.Request) map[string]string {
		return map[string]string{
			"quadkey": quadkeyValue,
		}
	}
}

func getQuadkeyToZxyProvider(z, x, y int) func(string) ([]int, error) {
	return func(quadkey string) ([]int, error) {
		return []int{z, x, y}, nil
	}
}

func quadkeyToZxyProviderError(quadkey string) ([]int, error) {
	return nil, errors.New("Spurious error")
}

func flipYProvider(z, y int) int {
	return int(math.Round(math.Pow(float64(y), 2)))
}
