package main

const (
	RIGHT int = 1 << iota
	DOWN
	LEFT
	UP
)

var STEPS = map[int]Point2{
	RIGHT: {x: 1, y: 0},
	DOWN:  {x: 0, y: 1},
	LEFT:  {x: -1, y: 0},
	UP:    {x: 0, y: -1},
}

// x: 0..256       ( 8 bits)
// y: 0..256       ( 8 bits)
// dir: 4b         ( 4 bits)
// budget: 0..maxb ( 4 bits)
// cost: 0..1024   (10 bits)
// ---
// total:

const (
	MSK_X      int = 0xFF
	MSK_Y          = 0xFF
	MSK_D          = 0xF
	MSK_B          = 0xF
	MSK_C          = 0x3FF
	MSK_NOCOST     = 0xFFFFFF
)

const (
	OFF_X int = 0
	OFF_Y     = 8
	OFF_D     = 16
	OFF_B     = 20
	OFF_C     = 24
)

func traverse(F [][]int, start, end Point2, minb, maxb int) int {
	V := make(map[int]int)
	H := NewBinHeap(func(a, b int) bool {
		c1, c2 := (a>>OFF_C)&MSK_C, (b>>OFF_C)&MSK_C
		if c1 == c2 {

		}
		return c1 < c2
	})

	H.Push((start.x << OFF_X) | (start.y << OFF_Y) | (RIGHT << OFF_D) | (0 << OFF_B) | (0 << OFF_C))
	H.Push((start.x << OFF_X) | (start.y << OFF_Y) | (DOWN << OFF_D) | (0 << OFF_B) | (0 << OFF_C))

	minT := ALOT
	var head int
	var x, y, nx, ny, dir, budget, nbudget, cost, ncost, k int
	var ND [3]int
	for H.Size() > 0 {
		head = H.Pop()
		x = (head >> OFF_X) & MSK_X
		y = (head >> OFF_Y) & MSK_Y
		dir = (head >> OFF_D) & MSK_D
		budget = (head >> OFF_B) & MSK_B
		cost = (head >> OFF_C) & MSK_C
		k = head & MSK_NOCOST

		if x == end.x && y == end.y {
			return cost
		}

		V[k] = cost

		ND = [3]int{dir, (dir>>1 | dir<<3) & MSK_D, (dir<<1 | dir>>3) & MSK_D}
		for _, ndir := range ND {
			nx, ny = x+STEPS[ndir].x, y+STEPS[ndir].y
			nbudget = budget
			if dir == ndir {
				nbudget += 1
			} else {
				if nbudget < minb {
					continue
				}
				nbudget = 0
			}
			if (nbudget >= maxb) || nx < 0 || ny < 0 || nx >= len(F[0]) || ny >= len(F) {
				continue
			}
			ncost = cost + F[ny][nx]
			nk := (nx << OFF_X) | (ny << OFF_Y) | (ndir << OFF_D) | (nbudget << OFF_B)
			if _, ok := V[nk]; ok {
				continue
			}
			V[nk] = cost

			H.Push(nk | (ncost << OFF_C))
		}
	}

	return minT
}

func parseField(ss []string) [][]int {
	F := makeIntField(len(ss), len(ss[0]))
	for y, s := range ss {
		for x, ch := range s {
			F[y][x] = int(ch - '0')
		}
	}

	return F
}

func main() {
	lines := input()
	F := parseField(lines)

	tt := traverse(F, Point2{x: 0, y: 0}, Point2{x: len(F[0]) - 1, y: len(F) - 1}, 0, 3)

	printf("min traverse: %d", tt)

	tt2 := traverse(F, Point2{x: 0, y: 0}, Point2{x: len(F[0]) - 1, y: len(F) - 1}, 3, 10)
	printf("min traverse2: %d", tt2)
}
