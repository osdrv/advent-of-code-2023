package main

import (
	"fmt"
	"strings"
)

func computeSegments(W map[string]map[string]bool, exclude map[[2]string]bool) [][]string {
	vis := make(map[string]bool)
	for from, tos := range W {
		vis[from] = false
		for to := range tos {
			vis[to] = false
		}
	}

	var visit func(v string) []string
	visit = func(v string) []string {
		res := []string{v}
		vis[v] = true
		for to := range W[v] {
			if vis[to] {
				continue
			}
			if exclude[[2]string{v, to}] || exclude[[2]string{to, v}] {
				continue
			}
			res = append(res, visit(to)...)
		}
		return res
	}

	segms := make([][]string, 0, 1)
	for v, didVis := range vis {
		if didVis {
			continue
		}
		segms = append(segms, visit(v))
	}

	return segms
}

func graphviz(W map[string][]string) string {
	var b strings.Builder
	b.WriteString("graph G {\n")
	for from, tos := range W {
		for _, to := range tos {
			b.WriteString(fmt.Sprintf("  %s -- %s\n", from, to))
		}
	}
	b.WriteString("}")

	return b.String()
}

func main() {
	lines := input()

	W := make(map[string]map[string]bool)
	ES := make([][2]string, 0, 1)

	for _, line := range lines {
		ss := strings.SplitN(line, ": ", 2)
		from := ss[0]
		tos := strings.Split(ss[1], " ")
		if _, ok := W[from]; !ok {
			W[from] = make(map[string]bool)
		}
		for _, to := range tos {
			ES = append(ES, [2]string{from, to})
			W[from][to] = true
			if _, ok := W[to]; !ok {
				W[to] = make(map[string]bool)
			}
			W[to][from] = true
		}
	}

	var isReachableWithin func(from, to string, rang int, path map[string]bool) bool
	isReachableWithin = func(from, to string, rang int, path map[string]bool) bool {
		path[from] = true
		defer func() {
			delete(path, from)
		}()
		for next := range W[from] {
			if path[next] {
				continue
			}
			if W[next][to] {
				return true
			}
			if rang > 1 {
				if isReachableWithin(next, to, rang-1, path) {
					return true
				}
			}
		}
		return false
	}

	cand := make([][2]string, 0, 3)
Edge:
	for _, e := range ES {
		from, to := e[0], e[1]
		if isReachableWithin(from, to, 2, map[string]bool{}) {
			continue Edge
		}
		cand = append(cand, e)
	}

	debugf("len(ES)=%d", len(ES))

	debugf("candidates (%d): %+v", len(cand), cand)

	ss := computeSegments(W, map[[2]string]bool{
		{"tvj", "cvx"}: true,
		{"fsv", "spx"}: true,
		{"kdk", "nct"}: true,
	})

	assert(len(ss) == 2, "2 partitions")
	printf("res= %d", len(ss[0])*len(ss[1]))

	// Cand:
	//
	//	for i := 0; i < len(cand); i++ {
	//		for j := i + 1; j < len(cand); j++ {
	//			debugf("i: %d, j: %d, k: *", i, j)
	//			for k := j + 1; k < len(cand); k++ {
	//				segments := computeSegments(W, map[[2]string]bool{
	//					cand[i]: true, cand[j]: true, cand[k]: true,
	//				})
	//				if len(segments) == 2 {
	//					res := len(segments[0]) * len(segments[1])
	//					printf("res: %d", res)
	//					break Cand
	//				}
	//			}
	//		}
	//	}
}
