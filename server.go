package main

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/captaincoordinates/tile-id-api/config"
	"github.com/captaincoordinates/tile-id-api/geo"
	"github.com/captaincoordinates/tile-id-api/handler"

	"github.com/captaincoordinates/tile-id-api/handler/quadkey"
	"github.com/captaincoordinates/tile-id-api/handler/tms"
	"github.com/captaincoordinates/tile-id-api/handler/zxy"
	"github.com/captaincoordinates/tile-id-api/params"
	"github.com/gorilla/mux"
)

var tileUtil handler.TileUtil = handler.NewTileUtil()
var configUtil config.ConfigUtil = config.NewConfigUtil()
var paramsUtil params.ParamsUtil = params.NewParamsUtil()

func main() {
	handlers := []handler.TileHandler{
		tms.NewTmsTileHandler(),
		zxy.NewZxyTileHandler(),
		quadkey.NewQuadkeyTileHandler(),
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
	listenPort := configUtil.GetListenPort()
	fmt.Println(fmt.Sprintf("Server port %d", listenPort))
	listenAddress := fmt.Sprintf(":%d", listenPort)
	err := http.ListenAndServe(listenAddress, router)
	if err != nil {
		panic(err)
	}
}

func createHandlerClosure(thisHandler handler.TileHandler, allIdentifiers []string) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		encoder, supportsOpacity := tileUtil.GetEncoder(request)
		var opacity uint8 = 255
		if supportsOpacity {
			opacity = paramsUtil.Opacity(request)
		}
		writer.Header().Set("X-tile-opacity", fmt.Sprintf("%d/255", opacity))
		tileKeysByIdentifier, tileKeysErr := thisHandler.Keys(request)
		if tileKeysErr != handler.NoReturnableError {
			http.Error(
				writer,
				tileKeysErr.ErrorMessage,
				tileKeysErr.StatusCode,
			)
			return
		}
		tileKeys := make([]string, len(allIdentifiers))
		for i, identifier := range sortIdentifiers(allIdentifiers, thisHandler.Identifier()) {
			tileKey, keyExists := tileKeysByIdentifier[identifier]
			if !keyExists {
				fmt.Println(fmt.Sprintf(
					"'%s' handler does not support identifier '%s'",
					thisHandler.Identifier(),
					identifier,
				))
				continue
			}
			writer.Header().Set(fmt.Sprintf("X-tile-id-%s", identifier), tileKey)
			tileKeys[i] = fmt.Sprintf("%s: %s", identifier, tileKey)
		}
		tileZxy, zxyError := thisHandler.AsZXY(request)
		if zxyError != nil {
			fmt.Println(zxyError.Error())
		} else {
			writer.Header().Set("X-tile-bounds-ll84", geo.GetTileBounds(tileZxy[0], tileZxy[1], tileZxy[2]).ToString())
		}
		img := tileUtil.GenerateTile(
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
