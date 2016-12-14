package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	i := 0
	curZ := 1.0
	preZ := 0.0
	for math.Abs(curZ - preZ) > 0.000000000000001 {
		i += 1
		preZ = curZ
		curZ = curZ - ((curZ * curZ) - x) / (2 * curZ)
	}
	fmt.Println(i)
	return curZ
}
