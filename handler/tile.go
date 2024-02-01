package handler

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const TILE_WIDTH int = 256
const TILE_HEIGHT int = 256

func GenerateTile(opacity uint8, labels ...string) image.Image {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{TILE_WIDTH, TILE_HEIGHT}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	fill := color.NRGBA{255, 255, 255, opacity}
	draw.Draw(img, img.Bounds(), &image.Uniform{fill}, image.Point{}, draw.Src)
	border := color.RGBA{0, 0, 0, 255}
	for x := 0; x < TILE_WIDTH; x++ {
		img.Set(x, 0, border)
		img.Set(x, TILE_HEIGHT-1, border)
	}
	for y := 0; y < TILE_HEIGHT; y++ {
		img.Set(0, y, border)
		img.Set(TILE_WIDTH-1, y, border)
	}
	for i, label := range labels {
		col := color.Black
		point := fixed.Point26_6{X: fixed.I(20), Y: fixed.I(20 + 20*i)}

		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(col),
			Face: basicfont.Face7x13,
			Dot:  point,
		}
		d.DrawString(label)
	}
	return img
}

func GetEncoder(request *http.Request) (encoder func(io.Writer, image.Image) error, supportsOpacity bool) {
	extension := regexp.MustCompile("^\\.").ReplaceAllString(
		mux.Vars(request)["extension"],
		"",
	)
	switch extension {
	case "jpg":
		fallthrough
	case "jpeg":
		return func(writer io.Writer, img image.Image) error {
			return jpeg.Encode(writer, img, &jpeg.Options{Quality: 100})
		}, false
	case "png":
		fallthrough
	default:
		return png.Encode, true
	}
}
