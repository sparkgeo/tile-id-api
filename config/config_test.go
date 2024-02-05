package config

import (
	"fmt"
	"testing"
)

func TestGetListenPortValid(t *testing.T) {
	requestedPort := uint(12345)
	envGetter := func(key string) string {
		return fmt.Sprint(requestedPort)
	}
	returnedPort := (ConfigUtil{
		envGetter: envGetter,
	}).GetListenPort()
	if returnedPort != requestedPort {
		t.Errorf("Incorrect port returned, expected '%d', got '%d'", requestedPort, returnedPort)
	}
}
