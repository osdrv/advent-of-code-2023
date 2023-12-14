package main

const (
	EMPTY  int = 0
	SQUARE     = 1
	ROUND      = 2
)

func tiltUp(F [][]int) [][]int {
	TF := copyNumField(F)

	for y := 0; y < len(TF); y++ {
		for x := 0; x < len(TF[y]); x++ {
			if TF[y][x] != ROUND {
				continue
			}
			for y1 := y - 1; y1 >= 0; y1-- {
				if TF[y1][x] == EMPTY {
					TF[y1][x], TF[y1+1][x] = TF[y1+1][x], TF[y1][x]
				} else {
					break
				}
			}
		}
	}

	return TF
}

func calcLoad(TF [][]int) int {
	load := 0
	for y := 0; y < len(TF); y++ {
		for x := 0; x < len(TF[y]); x++ {
			if TF[y][x] != ROUND {
				continue
			}
			load += len(TF) - y
		}
	}
	return load
}

func main() {
	lines := input()
	F := make([][]int, 0, len(lines))
	for _, line := range lines {
		row := make([]int, 0, len(line))
		for _, ch := range line {
			v := -1
			switch ch {
			case '.':
				v = EMPTY
			case '#':
				v = SQUARE
			case 'O':
				v = ROUND
			}
			row = append(row, v)
		}
		F = append(F, row)
	}

	debugf("\n%s", printNumFieldWithSubs(F, "", map[int]string{
		EMPTY: ".", SQUARE: "#", ROUND: "O",
	}))

	FF := copyNumField(F)

	TF := tiltUp(F)

	debugf("Tilted field:\n%s", printNumFieldWithSubs(TF, "", map[int]string{
		EMPTY: ".", SQUARE: "#", ROUND: "O",
	}))

	load := calcLoad(TF)

	printf("Load: %d", load)

	memo := make(map[string]int)
	cycle := 0
	cmax := 1000000000
	for cycle < cmax {
		debugf("Cycle %d", cycle)
		// NORTH
		FF = tiltUp(FF)
		// WEST
		FF = rotateMatRight(FF)
		FF = tiltUp(FF)
		// SOUTH
		FF = rotateMatRight(FF)
		FF = tiltUp(FF)
		// EAST
		FF = rotateMatRight(FF)
		FF = tiltUp(FF)
		// return back
		FF = rotateMatRight(FF)

		ss := printNumField(FF, "")

		debugf("Cycle %d:\n%s", cycle, ss)

		if prev, ok := memo[ss]; ok {
			debugf("cycle started at %d and came at %d", prev, cycle)
			cmax = (cmax-prev)%(cycle-prev) - 1
			cycle = 0
			memo = make(map[string]int)
			continue
		}

		memo[ss] = cycle
		cycle++
	}

	printf("Final load: %d", calcLoad(FF))
}
