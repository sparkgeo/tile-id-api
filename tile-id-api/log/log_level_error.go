package log

import (
	"fmt"
	"strings"
)

type LogLevelError struct {
	requestedLogLevel  string
	availableLogLevels []string
}

func (self *LogLevelError) Error() string {
	return fmt.Sprintf(
		"Requested invalid log level '%s', available log levels: '%s'",
		self.requestedLogLevel,
		strings.Join(self.availableLogLevels, "', '"),
	)
}

func NewLogLevelError(requestedLogLevel string, availableLogLevels []string) *LogLevelError {
	return &LogLevelError{
		requestedLogLevel,
		availableLogLevels,
	}
}
