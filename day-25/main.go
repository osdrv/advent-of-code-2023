package main

import (
	"strings"

	_ "net/http/pprof"
)

type Graph struct {
	numVx int
	VN    map[string]int
	VS    []int
	ES    [][]int
}

func NewGraph(vx map[string]struct{}) *Graph {
	numVx := len(vx)
	vn := make(map[string]int, numVx)
	vs := make([]int, numVx)
	es := makeIntField(numVx, numVx)

	for v := range vx {
		vix := len(vn)
		vn[v] = vix
		vs[vix] = 1
	}

	return &Graph{
		numVx: numVx,
		VN:    vn,
		VS:    vs,
		ES:    es,
	}
}

func (g *Graph) AddEdge(u, v string) {
	uix, vix := g.VN[u], g.VN[v]
	g.ES[uix][vix] = 1
	g.ES[vix][uix] = 1
}

func (g *Graph) contract(uix, vix int) {
	for wix := 0; wix < len(g.ES); wix++ {
		if g.VS[wix] == 0 {
			continue
		}
		g.ES[uix][wix] += g.ES[vix][wix]
		g.ES[wix][uix] += g.ES[wix][vix]
		g.ES[wix][vix] = 0
		g.ES[vix][wix] = 0
	}
	g.ES[uix][vix] = 0
	g.VS[vix] = 0
	g.numVx--
}

func (g *Graph) minCutPhase(six int) (int, int, int) {
	var W_1, V_1, V_2 int

	N := len(g.VS)
	A := make([]int, N)
	A[six] = 1
	cand := make([]int, N)
	copy(cand, g.VS)
	cand[six] = 0
	nCand := g.numVx - 1
	V_2 = six

	var maxW, maxVix, w int

	for nCand > 0 {
		maxW, maxVix = -ALOT, -1
		for vix := 0; vix < len(cand); vix++ {
			if cand[vix] == 0 {
				continue
			}
			w = 0
			for aix := 0; aix < len(A); aix++ {
				if A[aix] > 0 {
					w += g.ES[vix][aix]
				}
			}
			if w > maxW {
				maxW = w
				maxVix = vix
			}
		}
		A[maxVix] = 1
		cand[maxVix] = 0
		nCand--
		V_2 = V_1
		V_1 = maxVix
		W_1 = maxW
	}

	return V_2, V_1, W_1
}

func (g *Graph) MinCut(start string) (int, []int) {
	six := g.VN[start]
	AL := make(map[int][]int)
	debugf("g VS len: %d", len(g.VS))
	minWeight := ALOT
	var cutPartition []int
	for g.numVx > 1 {
		u, v, cutWeight := g.minCutPhase(six)
		if cutWeight < minWeight {
			minWeight = cutWeight
			cutPartition = expand(AL, v)
		}
		g.contract(u, v)
		if _, ok := AL[u]; !ok {
			AL[u] = make([]int, 0, 1)
		}
		AL[u] = append(AL[u], v)
	}

	return minWeight, cutPartition
}

func expand(M map[int][]int, v int) []int {
	res := []int{v}
	if _, ok := M[v]; !ok {
		return res
	}
	for _, vv := range M[v] {
		res = append(res, expand(M, vv)...)
	}
	return res
}

func cpMap[K comparable, V any](src map[K]V) map[K]V {
	cp := make(map[K]V, len(src))
	for k, v := range src {
		cp[k] = v
	}
	return cp
}

func main() {
	lines := input()

	var start string

	VS := make(map[string]struct{})
	ES := make(map[string][]string)

	for _, line := range lines {
		ss := strings.SplitN(line, ": ", 2)
		from := ss[0]
		VS[from] = struct{}{}
		tos := strings.Split(ss[1], " ")
		for _, to := range tos {
			VS[to] = struct{}{}
		}
		ES[from] = tos
		// we do not care about the choice of the active vertex
		if len(start) == 0 {
			start = from
		}
	}
	G := NewGraph(VS)
	for from, tos := range ES {
		for _, to := range tos {
			G.AddEdge(from, to)
		}
	}

	numV := len(G.VS)

	minCut, part := G.MinCut(start)
	printf("min cut: %d", minCut)
	printf("the result: %d", len(part)*(numV-len(part)))
}
