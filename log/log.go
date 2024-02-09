package log

import (
	"fmt"
	"strings"

	"github.com/captaincoordinates/tile-id-api/common"
	"github.com/captaincoordinates/tile-id-api/constants"
	"github.com/sirupsen/logrus"
)

var logLevelByString map[string]logrus.Level = map[string]logrus.Level{
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
	"error": logrus.ErrorLevel,
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
	"trace": logrus.TraceLevel,
}

func NewLogger(envGetter common.EnvGetter) *logrus.Logger {
	logLevelString, exists := envGetter("TILE_ID_LOG_LEVEL")
	var logLevel logrus.Level
	if exists {
		logLevel, exists = logLevelByString[strings.ToLower(logLevelString)]
		if !exists {
			var allKeys []string
			for key := range logLevelByString {
				allKeys = append(allKeys, key)
			}
			fmt.Println(fmt.Sprintf(
				"requested log level '%s' does not exist, using default. Permitted values: [%s]",
				logLevelString,
				strings.Join(allKeys, ", "),
			))
			logLevel = constants.DefaultLogLevel
		} else {
			fmt.Println(fmt.Sprintf(
				"successfully set log level with '%s'",
				logLevelString,
			))
		}
	} else {
		logLevel = constants.DefaultLogLevel
	}
	logger := logrus.New()
	logger.SetLevel(logLevel)
	return logger
}
