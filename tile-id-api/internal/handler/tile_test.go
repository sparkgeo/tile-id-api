package handler

import (
	"bytes"
	"image"
	"image/color"
	"net/http"
	"testing"
)

func TestGenerateTile(t *testing.T) {
	// simply test that it does not panic, manual testing will be required to validate image content
	TileUtil{
		pathParamsMapProvider: func(request *http.Request) map[string]string {
			return map[string]string{}
		},
	}.GenerateTile(125, "label 1", "label 2")
}

func TestGetEncoder(t *testing.T) {
	getPathParamsMapProvider := func(extension string) func(*http.Request) map[string]string {
		return func(request *http.Request) map[string]string {
			return map[string]string{
				"extension": extension,
			}
		}
	}
	inputsOutputs := []struct {
		extension string
		encoding  string
		opacity   bool
	}{
		{extension: "jpg", encoding: "jpeg", opacity: false},
		{extension: "jpeg", encoding: "jpeg", opacity: false},
		{extension: "png", encoding: "png", opacity: true},
		{extension: "unsupported", encoding: "png", opacity: true},
		{extension: "", encoding: "png", opacity: true},
	}
	for _, inputOutput := range inputsOutputs {
		tileLogic := TileUtil{
			pathParamsMapProvider: getPathParamsMapProvider(inputOutput.extension),
		}
		img := image.NewRGBA(image.Rect(0, 0, 1, 1))
		img.Set(0, 0, color.NRGBA{255, 255, 255, 255})
		buf := new(bytes.Buffer)
		encoder, supportsOpacity := tileLogic.GetEncoder(&http.Request{})
		if supportsOpacity != inputOutput.opacity {
			t.Errorf("Unexpected opacity result for '%s'", inputOutput.extension)
		}
		encoder(buf, img)
		_, format, _ := image.Decode(bytes.NewReader(buf.Bytes()))
		if format != inputOutput.encoding {
			t.Errorf("Unexpected image format from '%s'. Expected '%s' and got '%s'", inputOutput.extension, inputOutput.encoding, format)
		}
	}
}
