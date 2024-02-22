package internal

import (
	"errors"
	"net/http"
	"testing"

	"github.com/sparkgeo/tile-id-api/tile-id-api/internal/handler"
)

func TestStatusCodeFromErrorKnown(t *testing.T) {
	knownErrorsAndStatuses := []struct {
		err        error
		statusCode int
	}{
		{err: handler.NewBadRequestError("BadRequestError"), statusCode: http.StatusBadRequest},
		{err: handler.NewUnprocessableEntityError("UnprocessableEntityError"), statusCode: http.StatusUnprocessableEntity},
	}
	for _, pair := range knownErrorsAndStatuses {
		if statusCode := statusCodeFromError(pair.err); statusCode != pair.statusCode {
			t.Errorf("Unexpected status '%d' from '%v'", statusCode, pair.err)
		}
	}
}

func TestStatusCodeFromErrorUnknown(t *testing.T) {
	if statusCode := statusCodeFromError(errors.New("default error type")); statusCode != http.StatusInternalServerError {
		t.Errorf("Unexpected status '%d' from default error", statusCode)
	}
}
