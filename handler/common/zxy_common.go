package common

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

const ZxyPathPattern string = "{z:[0-9]|1[0-9]|2[0-5]}/{x:[0-9]{1,8}}/{y:[0-9]{1,8}}"

type FlipYProvider = func(z, y int) int

func FlipY(z, y int) int {
	return int(math.Pow(2, float64(z))) - y - 1
}

type ZxyToQuadkeyProvider = func(z, x, y int) string

func ZxyToQuadkey(z, x, y int) string {
	quadkeyParts := []string{}
	for i := z; i > 0; i-- {
		b := 0
		mask := 1 << (i - 1)
		if (x & mask) != 0 {
			b++
		}
		if (y & mask) != 0 {
			b += 2
		}
		quadkeyParts = append(quadkeyParts, fmt.Sprint(b))
	}
	return strings.Join(quadkeyParts, "")
}

type QuadkeyToZxyProvider = func(quadkey string) ([]int, error)

func QuadkeyToZxy(quadkey string) ([]int, error) {
	zxy := []int{len(quadkey), 0, 0}
	mask := 1 << len(quadkey)
	for i := 0; i < len(quadkey); i++ {
		mask >>= 1
		switch string(quadkey[i]) { // assumes utf-8, which is safe with input validation in the quadkey handler
		case "0":
			// default case, already handled in initialisation
			break
		case "1":
			zxy[1] |= mask
			break
		case "2":
			zxy[2] |= mask
			break
		case "3":
			zxy[1] |= mask
			zxy[2] |= mask
			break
		default:
			return zxy, errors.New(fmt.Sprintf("Invalid quadkey '%s'", quadkey))
		}
	}
	return zxy, nil
}
