package geo

import "testing"

func TestToString(t *testing.T) {
	generated := NewTileBounds(-127.5, 54, -126.5, 55).ToString()
	if generated != "{\"llX\":-127.5,\"llY\":54,\"urX\":-126.5,\"urY\":55}" {
		t.Errorf("Unexpected JSON output: '%s'", generated)
	}
}

func TestGetTileBounds(t *testing.T) {
	inputsOutputs := []struct {
		z, x, y            int
		llx, lly, urx, ury float64
	}{
		{z: 8, x: 37, y: 81, llx: -127.968749, lly: 54.162436, urx: -126.562500, ury: 54.977616},
		{z: 0, x: 0, y: 0, llx: -180, lly: -85.051129, urx: 180, ury: 85.051129},
	}
	equivalenceThreshold := 0.000001
	for _, inputOutput := range inputsOutputs {
		tileBounds := GetTileBounds(inputOutput.z, inputOutput.x, inputOutput.y)
		if !(equivalentFloat64(tileBounds.MinLon, inputOutput.llx, equivalenceThreshold) &&
			equivalentFloat64(tileBounds.MinLat, inputOutput.lly, equivalenceThreshold) &&
			equivalentFloat64(tileBounds.MaxLon, inputOutput.urx, equivalenceThreshold) &&
			equivalentFloat64(tileBounds.MaxLat, inputOutput.ury, equivalenceThreshold)) {
			t.Errorf(
				"Unexpected discrepancy in calculated tile bounds. z: %d, x: %d, y: %d, calculated bounds: '%s', expected bounds llX: %f, llY: %f, urX: %f, urY: %f",
				inputOutput.z, inputOutput.x, inputOutput.y,
				tileBounds.ToString(),
				inputOutput.llx, inputOutput.lly, inputOutput.urx, inputOutput.ury,
			)
		}
	}
}

func equivalentFloat64(a, b, threshold float64) bool {
	return (a - b) <= threshold
}
