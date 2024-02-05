package main

import (
	"strings"
	"testing"
)

func TestSortIdentifiers(t *testing.T) {
	sorted := sortIdentifiers([]string{
		"c",
		"b",
		"f",
		"a",
		"z",
	}, "z")
	sortedString := strings.Join(sorted, "")
	expectedString := "zabcf"
	if sortedString != expectedString {
		t.Errorf("Unexpected sort string. Expected '%s', got '%s'", expectedString, sortedString)
	}
}
