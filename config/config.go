package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/captaincoordinates/tile-id-api/constants"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type ConfigUtil struct {
	envGetter func(string) string
}

func NewConfigUtil() ConfigUtil {
	return ConfigUtil{
		envGetter: os.Getenv,
	}
}

func (self ConfigUtil) GetListenPort() uint {
	configuredPortStr := self.envGetter("TILE_ID_LISTEN_PORT")
	portRegex := regexp.MustCompile("^\\d{4,5}$")
	logFail := func() {
		log.Debug(fmt.Sprintf(
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
