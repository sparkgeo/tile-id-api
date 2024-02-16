package log

import (
	"fmt"
	"strings"
)

type LogLevelError struct {
	requestedLogLevel  string
	availableLogLevels []string
}

func (err *LogLevelError) Error() string {
	return fmt.Sprintf(
		"Requested invalid log level '%s', available log levels: '%s'",
		err.requestedLogLevel,
		strings.Join(err.availableLogLevels, "', '"),
	)
}

func NewLogLevelError(requestedLogLevel string, availableLogLevels []string) *LogLevelError {
	return &LogLevelError{
		requestedLogLevel,
		availableLogLevels,
	}
}
