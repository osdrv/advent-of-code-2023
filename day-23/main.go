package main

var NEXTSTEP = map[int][2]int{
	'>': {1, 0},
	'<': {-1, 0},
	'v': {0, 1},
	'^': {0, -1},
}

var (
	X_OFF  uint64 = 0
	Y_OFF         = 32
	X_MASK        = 0xFFFFFFFF
	Y_MASK        = 0xFFFFFFFF
)

type QE struct {
	cur  Point2
	prev Point2
	dist int
	vis  map[Point2]bool
	dir  bool
}

func cpMap[K comparable, V any](m map[K]V) map[K]V {
	cp := make(map[K]V)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

func findLongestPath(F [][]int, start Point2, end Point2, dryMode bool) int {
	VS := make(map[Point2]map[Point2]int)
	VS[start] = make(map[Point2]int)
	VS[end] = make(map[Point2]int)

	q := make([]QE, 0, 1)
	q = append(q, QE{start, start, 0, make(map[Point2]bool), false})

	var head QE
	for len(q) > 0 {
		head, q = q[0], q[1:]
		p := Point2{head.cur.x, head.cur.y}
		head.vis[p] = true

		steps := STEPS4
		if !dryMode && F[p.y][p.x] != '.' {
			steps = [][2]int{NEXTSTEP[F[p.y][p.x]]}
			head.dir = true
		}

		if _, ok := VS[p]; ok && p != head.prev {
			// we've been to this vertex before so just terminate
			VS[head.prev][p] = head.dist
			if !head.dir {
				VS[p][head.prev] = head.dist
			}
			continue
		}

		next := make([]Point2, 0, 1)
		for _, s := range steps {
			nx, ny := p.x+s[0], p.y+s[1]
			if nx < 0 || nx >= len(F[0]) || ny < 0 || ny >= len(F) {
				continue
			}
			if F[ny][nx] == '#' {
				continue
			}
			np := Point2{x: nx, y: ny}
			if head.vis[np] {
				continue
			}
			next = append(next, np)
		}
		if len(next) > 1 {
			if _, ok := VS[p]; !ok {
				VS[p] = make(map[Point2]int)
			}
			VS[head.prev][p] = max(VS[head.prev][p], head.dist)
			if !head.dir {
				VS[p][head.prev] = max(VS[p][head.prev], head.dist)
			}
			for _, np := range next {
				if head.vis[np] {
					continue
				}
				nvis := cpMap(head.vis)
				nvis[np] = true
				q = append(q, QE{
					prev: head.cur,
					cur:  np,
					dist: 1,
					vis:  nvis,
					dir:  false,
				})
			}
		} else {
			for _, np := range next {
				if head.vis[np] {
					continue
				}
				nvis := head.vis
				nvis[np] = true
				q = append(q, QE{
					prev: head.prev,
					cur:  np,
					dist: head.dist + 1,
					vis:  nvis,
					dir:  head.dir,
				})
			}
		}
	}

	res, _ := findLongestPathInGraph(VS, map[Point2]bool{start: true}, start, end)

	return res
}

func findLongestPathInGraph(VS map[Point2]map[Point2]int, vis map[Point2]bool, curr, end Point2) (int, bool) {
	if curr == end {
		return 0, true
	}
	maxdist := 0
	found := false
	for next, dist := range VS[curr] {
		if vis[next] {
			continue
		}
		vis[next] = true
		path, ok := findLongestPathInGraph(VS, vis, next, end)
		vis[next] = false
		if !ok {
			continue
		}
		found = true
		maxdist = max(maxdist, dist+path)
	}

	return maxdist, found
}

func main() {
	lines := input()

	F := makeIntField(len(lines), len(lines[0]))
	for y, line := range lines {
		for x, ch := range line {
			F[y][x] = int(ch)
		}
	}

	var start, end Point2
	for x := 0; x < len(F[0]); x++ {
		if F[0][x] == '.' {
			start = Point2{x, 0}
		}
		if F[len(F)-1][x] == '.' {
			end = Point2{x, len(F) - 1}
		}
	}

	debugf("===== part 1 =====")
	path1 := findLongestPath(F, start, end, false)
	printf("max path1: %d", path1)

	debugf("===== part 2 =====")
	path2 := findLongestPath(F, start, end, true)
	printf("max path2: %d", path2)
}
