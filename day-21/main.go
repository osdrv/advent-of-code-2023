package main

import (
	"github.com/openacid/slimarray/polyfit"
)

const (
	STONE  = 1
	GARDEN = 0
)

func countPoints(F [][]int, start Point2, N int) int {
	pp := make(map[Point2]bool)
	pp[start] = true
	var nx, ny, fy, fx int
	var np Point2
	eq := make([]float64, 0, 3)
	for i := 0; i < N; i++ {
		newpp := make(map[Point2]bool)
		for p := range pp {
			for _, step := range STEPS4 {
				nx, ny = p.x+step[0], p.y+step[1]
				np = Point2{x: nx, y: ny}
				fy, fx = ny%len(F), nx%len(F[0])
				if fy < 0 {
					fy += len(F)
				}
				if fx < 0 {
					fx += len(F[0])
				}
				if F[fy][fx] == STONE {
					continue
				}
				newpp[np] = true
			}
		}
		pp = newpp
		if (i+1)%len(F) == N%len(F) {
			debugf("steps: %d, points: %d", i+1, len(pp))
			eq = append(eq, float64(len(pp)))
			if len(eq) == 3 {
				fit := polyfit.NewFit([]float64{0, 1, 2}, eq, 2).Solve()
				x := N / len(F)
				res := int(fit[2])*x*x + int(fit[1])*x + int(fit[0])
				return res
			}
		}
	}

	debugf("\n%s", print2DMapWithSubs(pp, map[bool]string{true: "O"}))
	return len(pp)
}

func main() {
	lines := input()
	F := makeIntField(len(lines), len(lines[0]))
	var start Point2
	for y, line := range lines {
		for x, ch := range line {
			switch ch {
			case '#':
				F[y][x] = 1
			case '.':
				F[y][x] = 0
			case 'S':
				F[y][x] = 0
				start = Point2{x: x, y: y}
			}
		}
	}

	points1 := countPoints(F, start, 64)
	printf("points1: %d", points1)

	points2 := countPoints(F, start, 26501365)
	printf("points2: %d", points2)
}
