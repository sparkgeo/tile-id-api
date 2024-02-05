package geo

import (
	"encoding/json"
	"math"
)

type TileBounds struct {
	MinLon float64 `json:"llX"`
	MinLat float64 `json:"llY"`
	MaxLon float64 `json:"urX"`
	MaxLat float64 `json:"urY"`
}

func NewTileBounds(
	minLon, minLat, maxLon, maxLat float64,
) TileBounds {
	return TileBounds{
		MinLon: minLon,
		MinLat: minLat,
		MaxLon: maxLon,
		MaxLat: maxLat,
	}
}

func (self TileBounds) ToString() string {
	bytes, err := json.Marshal(self)
	if err != nil {
		return "error serializing bounds"
	}
	return string(bytes)
}

func GetTileBounds(z int, x int, y int) TileBounds {
	minLat, maxLat := yToLatEdges(y, z)
	minLon, maxLon := xToLonEdges(x, z)
	return NewTileBounds(
		minLon,
		minLat,
		maxLon,
		maxLat,
	)
}

func mercatorYToLat(mercatorY float64) float64 {
	return math.Atan(
		math.Sinh(
			float64(mercatorY),
		),
	) * 180 / math.Pi
}

func tileCount(z int) float64 {
	return math.Pow(2, float64(z))
}

func yToLatEdges(y int, z int) (minLat float64, maxLat float64) {
	unit := 1 / tileCount(z)
	relativeY1 := float64(y) * unit
	relativeY2 := relativeY1 + unit
	mercatorY1 := math.Pi * (1 - 2*relativeY1)
	mercatorY2 := math.Pi * (1 - 2*relativeY2)
	maxLat = mercatorYToLat(mercatorY1)
	minLat = mercatorYToLat(mercatorY2)
	return minLat, maxLat
}

func xToLonEdges(x int, z int) (minLon float64, maxLon float64) {
	unit := 360 / tileCount(z)
	minLon = -180 + float64(x)*unit
	maxLon = minLon + unit
	return minLon, maxLon
}
