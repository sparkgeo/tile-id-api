package log

import (
	"testing"

	"github.com/captaincoordinates/tile-id-api/tile-id-api/internal/constants"
)

func TestAllLogLevels(t *testing.T) {
	if len(AllLogLevels()) == 0 {
		t.Error("Did not receive a list of log levels as expected")
	}
}

func TestLogLevelFromStringValid(t *testing.T) {
	validLogLevelStrs := []string{"trace", "panic", "error"}
	for _, levelStr := range validLogLevelStrs {
		level := logLevelFromString(levelStr)
		if level.String() != levelStr {
			t.Errorf("got unexpected level '%s' from '%s'", level.String(), levelStr)
		}
	}
}

func TestLogLevelFromStringInvalid(t *testing.T) {
	invalidLogLevelStrs := []string{"deebg", "unknown", "432$$"}
	for _, levelStr := range invalidLogLevelStrs {
		level := logLevelFromString(levelStr)
		if level != constants.DefaultLogLevel {
			t.Errorf("expected default level '%s' from '%s'", constants.DefaultLogLevel.String(), levelStr)
		}
	}
}
