package main

type node struct {
	L, R string
}

func parseNode(s string) (string, *node) {
	var lbl, L, R string
	Scanf(s, "{string} = ({string}, {string})", &lbl, &L, &R)
	return lbl, &node{
		L: L,
		R: R,
	}
}

func traverseAll(NS map[string]*node, MV []byte, STTS, ENDS []string) int {
	ptrs := STTS
	steps := 0
	mvptr := 0

	cycles := make(map[string]int)

	for {
		if isAllEnd(ptrs) {
			return steps
		}
		newptrs := make([]string, 0, len(ptrs))
		for ix, ptr := range ptrs {
			if isEnd(ptr) {
				cycles[STTS[ix]] = steps
				if len(cycles) == len(STTS) {
					debugf("cycles: %+v", cycles)
					lcm := cycles[STTS[0]]
					for i := 1; i < len(STTS); i++ {
						lcm = LCM(lcm, cycles[STTS[i]])
					}
					return lcm
				}
			}
			if MV[mvptr] == 'L' {
				ptr = NS[ptr].L
			} else {
				ptr = NS[ptr].R
			}
			newptrs = append(newptrs, ptr)
		}
		mvptr++
		steps++
		mvptr %= len(MV)
		ptrs = newptrs
	}
}

func isAllEnd(ss []string) bool {
	for _, s := range ss {
		if !isEnd(s) {
			return false
		}
	}
	return true
}

func isStart(s string) bool {
	return s[2] == 'A'
}
func isEnd(s string) bool {
	return s[2] == 'Z'
}

func main() {
	lines := input()
	MV := []byte(lines[0])

	nodes := make(map[string]*node)
	for _, line := range lines[2:] {
		label, node := parseNode(line)
		nodes[label] = node
	}

	steps := traverseAll(nodes, MV, []string{"AAA"}, []string{"ZZZ"})

	printf("steps: %d", steps)

	starts := make([]string, 0, 1)
	ends := make([]string, 0, 1)
	for label := range nodes {
		if isStart(label) {
			starts = append(starts, label)
		} else if isEnd(label) {
			ends = append(ends, label)
		}
	}

	steps2 := traverseAll(nodes, MV, starts, ends)
	printf("steps2: %d", steps2)
}
