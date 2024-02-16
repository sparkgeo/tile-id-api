package common

import "testing"

func TestFlipY(t *testing.T) {
	inputsOutputs := []struct{ z, y, result int }{
		{z: 12, y: 3245, result: 850},
		{z: 4, y: 12, result: 3},
	}
	for _, inputOutput := range inputsOutputs {
		flipped := FlipY(inputOutput.z, inputOutput.y)
		if flipped != inputOutput.result {
			t.Errorf("Expected '%d', got '%d'", inputOutput.result, flipped)
		}
	}
}

var zxyQuadkeySets = []struct {
	z, x, y int
	quadkey string
}{
	{z: 14, x: 8187, y: 5448, quadkey: "03131313113011"},
	{z: 9, x: 123, y: 242, quadkey: "023331031"},
	{z: 0, x: 0, y: 0, quadkey: ""},
}

func TestZxyToQuadkey(t *testing.T) {
	for _, set := range zxyQuadkeySets {
		providedQuadkey := ZxyToQuadkey(set.z, set.x, set.y)
		if providedQuadkey != set.quadkey {
			t.Errorf("Expected '%s', got '%s'", set.quadkey, providedQuadkey)
		}
	}
}

func TestQuadkeyToZxy(t *testing.T) {
	for _, set := range zxyQuadkeySets {
		providedZxy, _ := QuadkeyToZxy(set.quadkey)
		z, x, y := providedZxy[0], providedZxy[1], providedZxy[2]
		if !(z == set.z && x == set.x && y == set.y) {
			t.Errorf("Expected '%d/%d/%d', got '%d/%d/%d'", set.z, set.x, set.y, z, x, y)
		}
	}
}
