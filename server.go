package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/zxy/{z}/{x}/{y}", tileZXYHandler)
	listenPort := getListenPort()
	fmt.Println(fmt.Sprintf("Listening on port %d", listenPort))
	listenAddress := fmt.Sprintf(":%d", listenPort)
	err := http.ListenAndServe(listenAddress, router)
	if err != nil {
		panic(err)
	}
}

func tileZXYHandler(writer http.ResponseWriter, request *http.Request) {
	values, err := fetchIntPathParams(writer, request, "z", "x", "y")
	if err != nil {
		return
	}
	z, x, y := values[0], values[1], values[2]
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "Requested z: %d, x: %d, y: %d", z, x, y)
}

func fetchIntPathParams(
	writer http.ResponseWriter,
	request *http.Request,
	paramNames ...string,
) ([]int, error) {
	vars := mux.Vars(request)
	parsedInts := make([]int, len(paramNames))
	for i, paramName := range paramNames {
		paramStr, ok := vars[paramName]
		if !ok {
			http.Error(
				writer,
				fmt.Sprintf("Missing expected parameter %s", paramName),
				http.StatusBadRequest,
			)
			return make([]int, 0), errors.New("Missing parameter(s)")
		}
		paramInt, err := strconv.ParseInt(paramStr, 10, 64)
		if err != nil {
			http.Error(
				writer,
				fmt.Sprintf("Parameter %s cannot be parsed to int (value %s)", paramName, paramStr),
				http.StatusBadRequest,
			)
			return make([]int, 0), errors.New("Parameter parse error")
		}
		parsedInts[i] = int(paramInt)
	}
	return parsedInts, nil
}

func getListenPort() uint {
	const defaultPort uint = 8080
	configuredPortStr := os.Getenv("TILE_ID_LISTEN_PORT")
	portRegex, err := regexp.Compile("^\\d{4,5}$")
	if err != nil {
		panic(err)
	}
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
