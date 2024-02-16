package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/captaincoordinates/tile-id-api/tile-id-api/internal"
	"github.com/captaincoordinates/tile-id-api/tile-id-api/internal/log"
)

func main() {
	listenPort := flag.Int("server-port", 8080, "Port the server listens on")
	logLevelStr := flag.String("log-level", "info", strings.Join(log.AllLogLevels(), " | "))
	flag.Parse()
	logger := log.NewLogger(*logLevelStr)
	logger.Debug(fmt.Sprintf("Server port %d", *listenPort))
	listenAddress := fmt.Sprintf(":%d", *listenPort)
	err := http.ListenAndServe(listenAddress, internal.ConfigureRouter(logger))
	if err != nil {
		panic(err)
	}
}
