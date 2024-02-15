package config

import (
	"os"

	"github.com/captaincoordinates/tile-id-api/tile-id-api/common"
	"github.com/captaincoordinates/tile-id-api/tile-id-api/log"
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

func (self ConfigUtil) GetEnvVar(key string) (string, bool) {
	return self.envGetter(key)
}
