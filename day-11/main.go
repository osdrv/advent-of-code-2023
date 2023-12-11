package main

func main() {
	lines := input()

	F := make([][]int, 0, len(lines))
	for _, line := range lines {
		row := make([]int, 0, len(line))
		for i := 0; i < len(lines[0]); i++ {
			if line[i] == '#' {
				row = append(row, 1)
			} else {
				row = append(row, 0)
			}
		}
		F = append(F, row)
	}

	PP := make([]Point2, 0, 1)
	CX := make([]int, len(F[0]))
	CY := make([]int, len(F))
	for y := 0; y < len(F); y++ {
		for x := 0; x < len(F[0]); x++ {
			if F[y][x] == 1 {
				CX[x]++
				CY[y]++
				PP = append(PP, Point2{x: x, y: y})
			}
		}
	}
	debugf("PP: %+v", PP)

	//println(printNumFieldWithSubs(F, " ", map[int]string{
	//	0: ".", 1: "#",
	//}))

	printf("dist1: %d", totDist(F, CX, CY, PP, uint64(2)))
	printf("dist2: %d", totDist(F, CX, CY, PP, uint64(1000000)))
}

func totDist(F [][]int, CX, CY []int, PP []Point2, cost uint64) uint64 {
	dist := uint64(0)
	for i := 0; i < len(PP); i++ {
		for j := i + 1; j < len(PP); j++ {
			dist += getDist(F, CX, CY, PP[i], PP[j], cost)
		}
	}
	return dist
}

func getDist(F [][]int, CX, CY []int, p1, p2 Point2, cost uint64) uint64 {
	d := uint64(abs(p1.x-p2.x) + abs(p1.y-p2.y))
	for y := min(p1.y, p2.y) + 1; y < max(p1.y, p2.y); y++ {
		if CY[y] == 0 {
			d += cost - 1
		}
	}
	for x := min(p1.x, p2.x) + 1; x < max(p1.x, p2.x); x++ {
		if CX[x] == 0 {
			d += cost - 1
		}
	}
	return d
}
