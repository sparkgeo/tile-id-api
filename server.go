package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"

	"github.com/captaincoordinates/tile-id-api/geo"
	"github.com/captaincoordinates/tile-id-api/handler"
	"github.com/captaincoordinates/tile-id-api/handler/quadkey"
	"github.com/captaincoordinates/tile-id-api/handler/tms"
	"github.com/captaincoordinates/tile-id-api/handler/zxy"
	"github.com/captaincoordinates/tile-id-api/params"
	"github.com/gorilla/mux"
)

func main() {
	handlers := []handler.TileHandler{
		tms.TmsTileHandler{},
		zxy.ZxyTileHandler{},
		quadkey.QuadkeyTileHandler{},
	}
	allIdentifiers := make([]string, len(handlers))
	for i, eachHandler := range handlers {
		allIdentifiers[i] = eachHandler.Identifier()
	}
	router := mux.NewRouter()
	for _, eachHandler := range handlers {
		router.HandleFunc(
			fmt.Sprintf(
				"/%s/%s{extension:(?:\\.(?:jpg|jpeg|png))?}",
				eachHandler.Identifier(),
				eachHandler.PathPattern(),
			),
			createHandlerClosure(eachHandler, allIdentifiers),
		)
	}
	// write some unit tests
	// README with installation instructions
	// demo and discuss hosting
	// proper logging with debug, warn, info etc?
	listenPort := getListenPort()
	fmt.Println(fmt.Sprintf("Listening on port %d", listenPort))
	listenAddress := fmt.Sprintf(":%d", listenPort)
	err := http.ListenAndServe(listenAddress, router)
	if err != nil {
		panic(err)
	}
}

func createHandlerClosure(thisHandler handler.TileHandler, allIdentifiers []string) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		encoder, supportsOpacity := handler.GetEncoder(request)
		var opacity uint8 = 255
		if supportsOpacity {
			opacity = params.Opacity(request)
		}
		writer.Header().Set("X-tile-opacity", fmt.Sprintf("%d/255", opacity))
		tileKeyProvider, keyProviderErr := thisHandler.GetKeyProvider(request)
		if keyProviderErr != handler.NoReturnableError {
			http.Error(
				writer,
				keyProviderErr.ErrorMessage,
				keyProviderErr.StatusCode,
			)
			return
		}
		tileKeys := make([]string, len(allIdentifiers))
		for i, identifier := range sortIdentifiers(allIdentifiers, thisHandler.Identifier()) {
			tileKey := tileKeyProvider(identifier)
			writer.Header().Set(fmt.Sprintf("X-tile-id-%s", identifier), tileKey)
			tileKeys[i] = fmt.Sprintf("%s: %s", identifier, tileKey)
		}
		tileZxy, zxyError := thisHandler.AsZXY(request)
		if zxyError != nil {
			fmt.Println(zxyError.Error())
		} else {
			writer.Header().Set("X-tile-bounds-ll84", geo.GetTileBounds(tileZxy[0], tileZxy[1], tileZxy[2]).ToString())
		}
		img := handler.GenerateTile(
			opacity,
			tileKeys...,
		)
		encodeErr := encoder(writer, img)
		if encodeErr != nil {
			http.Error(
				writer,
				fmt.Sprintf("Unable to encode image: %v", encodeErr),
				http.StatusInternalServerError,
			)
			return
		}
	}
}

func sortIdentifiers(allIdentifiers []string, firstValue string) []string {
	ownAllIdentifiers := allIdentifiers[:]
	slices.Sort(ownAllIdentifiers)
	split := -1
	for i, identifier := range ownAllIdentifiers {
		if identifier == firstValue {
			split = i
		}
	}
	if split == -1 {
		return ownAllIdentifiers
	}
	if split == 0 {
		return ownAllIdentifiers
	}
	return append(append(append([]string{}, firstValue), ownAllIdentifiers[0:split]...), ownAllIdentifiers[split+1:]...)
}

func getListenPort() uint {
	const defaultPort uint = 8080
	configuredPortStr := os.Getenv("TILE_ID_LISTEN_PORT")
	portRegex := regexp.MustCompile("^\\d{4,5}$")
	logFail := func() {
		fmt.Println(fmt.Sprintf("Unable to parse configured port '%s', returning default %d", configuredPortStr, defaultPort))
	}
	if portRegex.MatchString(configuredPortStr) {
		parsedPort, err := strconv.ParseUint(configuredPortStr, 10, 64)
		if err != nil {
			logFail()
			return defaultPort
		}
		return uint(parsedPort)
	} else {
		logFail()
		return defaultPort
	}
}
