package log

import (
	"testing"

	"github.com/captaincoordinates/tile-id-api/tile-id-api/constants"
	"github.com/sirupsen/logrus"
)

func TestNewLoggerValidLevel(t *testing.T) {
	inputsOutputs := map[string]logrus.Level{
		"FaTAL": logrus.FatalLevel,
		"trace": logrus.TraceLevel,
	}
	for strLevel, levelValue := range inputsOutputs {
		envGetter := func(_ string) (string, bool) {
			return strLevel, true
		}
		logger := NewLogger(envGetter)
		if logger.Level != levelValue {
			t.Errorf("Unexpected log level returned for '%s': '%v'", strLevel, logger.Level)
		}
	}
}

func TestNewLoggerInvalidLevel(t *testing.T) {
	envGetter := func(_ string) (string, bool) {
		return "invalid value", true
	}
	logger := NewLogger(envGetter)
	if logger.Level != constants.DefaultLogLevel {
		t.Errorf("Unexpected log level returned for invalid value: '%v'", logger.Level)
	}
}

func TestNewLoggerMissingLevel(t *testing.T) {
	envGetter := func(_ string) (string, bool) {
		return "", false
	}
	logger := NewLogger(envGetter)
	if logger.Level != constants.DefaultLogLevel {
		t.Errorf("Unexpected log level returned for missing value: '%v'", logger.Level)
	}
}
