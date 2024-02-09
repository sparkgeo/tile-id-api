package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/captaincoordinates/tile-id-api/common"
	"github.com/captaincoordinates/tile-id-api/constants"
	"github.com/captaincoordinates/tile-id-api/log"
)

var DefaultEnvGetter = os.LookupEnv

var logger = log.NewLogger(DefaultEnvGetter)

type ConfigUtil struct {
	envGetter common.EnvGetter
}

func NewConfigUtil() ConfigUtil {
	return ConfigUtil{
		envGetter: os.LookupEnv,
	}
}

func (self ConfigUtil) GetListenPort() uint {
	configuredPortStr, _ := self.envGetter("TILE_ID_SERVER_PORT")
	portRegex := regexp.MustCompile("^\\d{2,5}$")
	logFail := func() {
		logger.Debug(fmt.Sprintf(
			"Unable to parse configured port '%s', returning default %d",
			configuredPortStr,
			constants.DefaultPort,
		))
	}
	if portRegex.MatchString(configuredPortStr) {
		parsedPort, err := strconv.ParseUint(configuredPortStr, 10, 64)
		if err != nil {
			logFail()
			return constants.DefaultPort
		}
		return uint(parsedPort)
	} else {
		logFail()
		return constants.DefaultPort
	}
}

func (self ConfigUtil) GetEnvVar(key string) (string, bool) {
	return self.envGetter(key)
}
