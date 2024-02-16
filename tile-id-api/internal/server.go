package internal

import (
	"fmt"
	"net/http"

	"github.com/captaincoordinates/tile-id-api/tile-id-api/internal/geo"
	"github.com/captaincoordinates/tile-id-api/tile-id-api/internal/handler"
	"github.com/captaincoordinates/tile-id-api/tile-id-api/internal/params"
	"github.com/sirupsen/logrus"

	"github.com/captaincoordinates/tile-id-api/tile-id-api/internal/handler/quadkey"
	"github.com/captaincoordinates/tile-id-api/tile-id-api/internal/handler/tms"
	"github.com/captaincoordinates/tile-id-api/tile-id-api/internal/handler/zxy"
	"github.com/gorilla/mux"
)

func ConfigureRouter(logger logrus.FieldLogger) *mux.Router {
	tileUtil := handler.NewTileUtil()
	paramsUtil := params.NewParamsUtil(logger)
	handlers := []handler.TileHandler{
		tms.NewTmsTileHandler(logger),
		zxy.NewZxyTileHandler(logger),
		quadkey.NewQuadkeyTileHandler(logger),
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
			createHandlerClosure(eachHandler, logger, paramsUtil, tileUtil, allIdentifiers),
		)
	}
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./swagger/"))))
	router.HandleFunc("/openapi.yml", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "openapi.yml")
	})
	router.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/docs/", http.StatusMovedPermanently)
	})
	return router
}

func createHandlerClosure(
	thisHandler handler.TileHandler,
	logger logrus.FieldLogger,
	paramsUtil params.ParamsUtil,
	tileUtil handler.TileUtil,
	allIdentifiers []string,
) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		encoder, supportsOpacity := tileUtil.GetEncoder(request)
		var opacity uint8 = 255
		if supportsOpacity {
			opacity = paramsUtil.Opacity(request)
		}
		writer.Header().Set("X-tile-opacity", fmt.Sprintf("%d/255", opacity))
		tileKeysByIdentifier, tileKeysErr := thisHandler.Keys(request)
		if tileKeysErr != nil {
			http.Error(
				writer,
				tileKeysErr.Error(),
				statusCodeFromError(tileKeysErr),
			)
			return
		}
		tileKeys := make([]string, len(allIdentifiers))
		for i, identifier := range handler.SortIdentifiers(allIdentifiers, thisHandler.Identifier()) {
			tileKey, keyExists := tileKeysByIdentifier[identifier]
			if !keyExists {
				logger.Warn(fmt.Sprintf(
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
			logger.Warn(zxyError.Error())
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

func statusCodeFromError(err error) int {
	if _, ok := err.(handler.BadRequestError); ok {
		return http.StatusBadRequest
	}
	if _, ok := err.(handler.UnprocessableEntityError); ok {
		return http.StatusUnprocessableEntity
	}
	return http.StatusInternalServerError
}
