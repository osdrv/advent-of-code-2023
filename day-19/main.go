package main

import (
	"strings"
)

const (
	OP_LT byte = '<'
	OP_GT      = '>'
)

const (
	ACCEPT = "A"
	REJECT = "R"
)

var ATTRIX = map[string]int{
	"x": 0,
	"m": 1,
	"a": 2,
	"s": 3,
}

type Matcher struct {
	attr string
	op   byte
	val  int
	dest string
}

type WF struct {
	name     string
	matchers []Matcher
}

func parseWorkflow(s string) *WF {
	var wf WF
	name, ptr := readWord(s, 0)
	wf.name = name
	matchers := make([]Matcher, 0, 1)
	for _, part := range strings.Split(s[ptr+1:len(s)-1], ",") {
		var m Matcher
		if !strings.ContainsRune(part, ':') {
			m.dest = part
			matchers = append(matchers, m)
			continue
		}
		Scanf(part, "{string}{byte}{int}:{string}", &m.attr, &m.op, &m.val, &m.dest)
		matchers = append(matchers, m)
	}
	wf.matchers = matchers
	return &wf
}

type Part struct {
	attrs map[string]int
}

func (p Part) Rating() int {
	r := 0
	for _, v := range p.attrs {
		r += v
	}
	return r
}

func parsePart(s string) *Part {
	attrs := make(map[string]int)
	for _, pp := range strings.Split(s[1:len(s)-1], ",") {
		chs := strings.SplitN(pp, "=", 2)
		attrs[chs[0]] = parseInt(chs[1])
	}
	return &Part{attrs: attrs}
}

func multRanges(ranges [4][2]int) uint64 {
	res := uint64(1)
	for _, rng := range ranges {
		assert(rng[1] >= rng[0], "wtf is this")
		res *= uint64(rng[1] - rng[0] + 1)
	}
	return res
}

func acceptRanges(wfm map[string]*WF, ptr string, rngs [4][2]int) uint64 {
	debugf("ptr: %s, rngs: %+v", ptr, rngs)
	if ptr == ACCEPT {
		return multRanges(rngs)
	} else if ptr == REJECT {
		return 0
	}
	res := uint64(0)
	for _, mtch := range wfm[ptr].matchers {
		if mtch.op == 0 {
			res += acceptRanges(wfm, mtch.dest, rngs)
			continue
		}
		if mtch.op == OP_LT {
			nrngs := rngs
			attrix := ATTRIX[mtch.attr]
			nrngs[attrix][1] = min(mtch.val-1, rngs[attrix][1]+1)
			if nrngs[attrix][0] < nrngs[attrix][1] {
				res += acceptRanges(wfm, mtch.dest, nrngs)
			}
			if rngs[attrix][1] < mtch.val {
				break
			}
			rngs[attrix][0] = max(mtch.val, rngs[attrix][0])
		} else if mtch.op == OP_GT {
			nrngs := rngs
			attrix := ATTRIX[mtch.attr]
			nrngs[attrix][0] = max(mtch.val+1, rngs[attrix][0]-1)
			if nrngs[attrix][0] < nrngs[attrix][1] {
				res += acceptRanges(wfm, mtch.dest, nrngs)
			}
			if rngs[attrix][0] > mtch.val {
				break
			}
			rngs[attrix][1] = min(mtch.val, rngs[attrix][1])
		} else {
			panic("wtf2")
		}
	}

	return res
}

func main() {
	lines := input()
	ptr := 0
	wfm := make(map[string]*WF)
	for ptr < len(lines) {
		if len(lines[ptr]) == 0 {
			break
		}
		wf := parseWorkflow(lines[ptr])
		wfm[wf.name] = wf
		ptr++
	}
	ptr++

	tot := 0
	for ptr < len(lines) {
		part := parsePart(lines[ptr])
		if cnt := acceptRanges(wfm, "in", [4][2]int{
			{part.attrs["x"], part.attrs["x"]},
			{part.attrs["m"], part.attrs["m"]},
			{part.attrs["a"], part.attrs["a"]},
			{part.attrs["s"], part.attrs["s"]},
		}); cnt > 0 {
			debugf("cnt=%d", cnt)
			debugf("part %v is accepted", part)
			tot += part.Rating()
		}
		ptr++
	}
	printf("tot: %d", tot)

	tot2 := acceptRanges(wfm, "in", [4][2]int{{1, 4000}, {1, 4000}, {1, 4000}, {1, 4000}})
	printf("tot2: %d", tot2)
}
