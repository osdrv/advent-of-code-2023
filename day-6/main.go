package main

import (
	"math"
	"strings"
)

func solveInInts(ti, di int) (int, int) {
	t, d := float64(ti), float64(di)
	v1 := (t - math.Sqrt(t*t-4*d)) / 2
	v2 := (t + math.Sqrt(t*t-4*d)) / 2
	return int(math.Ceil(v1)), int(math.Floor(v2))
}

func main() {
	lines := input()
	times := parseInts(strings.SplitN(lines[0], ":", 2)[1])
	dists := parseInts(strings.SplitN(lines[1], ":", 2)[1])

	debugf("times: %+v, dists: %+v", times, dists)

	mult := 1
	for i := 0; i < len(times); i++ {
		v1, v2 := solveInInts(times[i], dists[i])
		debugf("v1: %d, v2: %d, tot: %d", v1, v2, v2-v1)
		mult *= (v2 - v1 + 1)
	}
	printf("mult: %d", mult)

	T := glueNum(times)
	D := glueNum(dists)

	v1, v2 := solveInInts(T, D)
	tot := v2 - v1 + 1
	printf("tot: %d", tot)

}
