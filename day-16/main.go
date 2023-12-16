package main

const (
	RIGHT int = 1 << iota
	DOWN
	LEFT
	UP
)

const PAD = 8

var NEXT = map[int]Point2{
	RIGHT: {x: 1, y: 0},
	DOWN:  {x: 0, y: 1},
	LEFT:  {x: -1, y: 0},
	UP:    {x: 0, y: -1},
}

var RULES = map[int]int{
	RIGHT<<PAD | '.': RIGHT,
	LEFT<<PAD | '.':  LEFT,
	DOWN<<PAD | '.':  DOWN,
	UP<<PAD | '.':    UP,

	RIGHT<<PAD | '/': UP,
	LEFT<<PAD | '/':  DOWN,
	DOWN<<PAD | '/':  LEFT,
	UP<<PAD | '/':    RIGHT,

	RIGHT<<PAD | '\\': DOWN,
	LEFT<<PAD | '\\':  UP,
	DOWN<<PAD | '\\':  RIGHT,
	UP<<PAD | '\\':    LEFT,

	RIGHT<<PAD | '-': RIGHT,
	LEFT<<PAD | '-':  LEFT,
	DOWN<<PAD | '-':  LEFT | RIGHT,
	UP<<PAD | '-':    LEFT | RIGHT,

	RIGHT<<PAD | '|': UP | DOWN,
	LEFT<<PAD | '|':  UP | DOWN,
	DOWN<<PAD | '|':  DOWN,
	UP<<PAD | '|':    UP,
}

func getNextDir(F [][]byte, p Point3) []Point3 {
	next := make([]Point3, 0, 2)
	np := Point3{x: p.x + NEXT[p.z].x, y: p.y + NEXT[p.z].y}

	if np.x < 0 || np.x >= len(F[0]) || np.y < 0 || np.y >= len(F) {
		return next
	}

	CH := F[np.y][np.x]

	ND := RULES[p.z<<PAD|int(CH)]
	for _, d := range [4]int{RIGHT, LEFT, DOWN, UP} {
		if ND&d > 0 {
			nnp := np
			nnp.z = d
			next = append(next, nnp)
		}
	}

	return next
}

func getEnergized(F [][]byte, start Point3) int {
	heads := make([]Point3, 0, 1)
	heads = append(heads, start)

	M := make(map[int]bool)
	H := make(map[Point3]bool)
	for len(heads) > 0 {
		newheads := make([]Point3, 0, 1)
		for _, head := range heads {
			newnewheads := getNextDir(F, head)
			for _, newhead := range newnewheads {
				if H[newhead] {
					continue
				}
				H[newhead] = true
				M[newhead.x<<PAD|newhead.y] = true
				newheads = append(newheads, newhead)
			}
		}
		heads = newheads
		debugf("newheads: %+v", newheads)
	}

	return len(M)
}

func main() {
	lines := input()
	F := make([][]byte, 0, len(lines))
	for _, line := range lines {
		F = append(F, []byte(line))
	}

	en := getEnergized(F, Point3{x: -1, y: 0, z: RIGHT})
	printf("energized: %d", en)

	maxen := en
	for y := 0; y < len(F); y++ {
		maxen = max(maxen, getEnergized(F, Point3{x: -1, y: y, z: RIGHT}))
		maxen = max(maxen, getEnergized(F, Point3{x: len(F[0]), y: y, z: LEFT}))
	}
	for x := 0; x < len(F[0]); x++ {
		maxen = max(maxen, getEnergized(F, Point3{x: x, y: -1, z: DOWN}))
		maxen = max(maxen, getEnergized(F, Point3{x: x, y: len(F), z: UP}))
	}

	printf("max energized: %d", maxen)
}
