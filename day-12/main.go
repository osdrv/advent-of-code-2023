package main

import (
	"strconv"
	"strings"
)

const (
	OFF = 0
	ON  = 1
	UNK = 2
)

func countArrangements(f []int, n []int) int {
	flushMemo()
	return doCnt(f, n, 0, 0)
}

func eatWsps(f []int, fptr int) int {
	for fptr < len(f) && f[fptr] == ON {
		fptr++
	}
	return fptr
}

func readGroup(f []int, fptr int, n int) (int, bool) {
	// there is no ambiguity in the group
	for fptr < len(f) && n > 0 && (f[fptr] == OFF || f[fptr] == UNK) {
		n--
		fptr++
	}
	// ensure empty space
	if n == 0 {
		if fptr >= len(f) || f[fptr] == ON || f[fptr] == UNK {
			// move one more symbol
			return fptr + 1, true
		}
	}
	return fptr, false
}

var MEMO map[string]int

func flushMemo() {
	MEMO = make(map[string]int)
}

func mkKey(f []int, fptr, nptr int) string {
	var b strings.Builder
	for i := fptr; i < len(f); i++ {
		b.WriteByte('0' + byte(f[i]))
	}
	b.WriteByte(',')
	b.WriteString(strconv.Itoa(fptr))
	b.WriteByte(',')
	b.WriteString(strconv.Itoa(nptr))
	return b.String()
}

func doCnt(f []int, n []int, fptr, nptr int) int {
	k := mkKey(f, fptr, nptr)
	if prev, ok := MEMO[k]; ok {
		return prev
	}
	fptr = eatWsps(f, fptr)
	if fptr >= len(f) {
		if nptr >= len(n) {
			// we made it to the end
			return 1
		}
		// we have not reached the end of groups
		return 0
	}
	if f[fptr] == OFF {
		// we're out of the group budget
		if nptr >= len(n) {
			return 0
		}
		nfptr, ok := readGroup(f, fptr, n[nptr])
		if !ok {
			// we could not read the group
			return 0
		}
		// progress on to the next group
		res := doCnt(f, n, nfptr, nptr+1)
		MEMO[k] = res
		return res
	}
	if f[fptr] == UNK {
		tot := 0
		f[fptr] = ON
		tot += doCnt(f, n, fptr, nptr)
		f[fptr] = OFF
		tot += doCnt(f, n, fptr, nptr)
		f[fptr] = UNK
		MEMO[k] = tot
		return tot
	}
	debugf("fptr: %d, nptr: %d, f[fptr]: %d", fptr, nptr, f[fptr])
	panic("how did I end up here?")
}

func parseField(s string) []int {
	res := make([]int, len(s))
	for i, ch := range s {
		switch ch {
		case '?':
			res[i] = UNK
		case '.':
			res[i] = ON
		case '#':
			res[i] = OFF
		default:
			panic("unexpected symbol")
		}
	}
	return res
}

func unfold(s string, n []int, mult int) (string, []int) {
	ss := make([]string, 0, mult)
	nn := make([]int, 0, len(n)*mult)
	for i := 0; i < mult; i++ {
		ss = append(ss, s)
		nn = append(nn, n...)
	}
	return strings.Join(ss, "?"), nn
}

func main() {
	lines := input()
	F := make([][]int, 0, len(lines))
	N := make([][]int, 0, len(lines))
	for _, line := range lines {
		ss := strings.SplitN(line, " ", 2)
		F = append(F, parseField(ss[0]))
		N = append(N, parseInts(ss[1]))
	}

	debugf("F:%+v", F)
	debugf("N: %+v", N)

	tot := 0
	for i := 0; i < len(F); i++ {
		cnt := countArrangements(F[i], N[i])
		debugf("line: %s, cnt: %d", lines[i], cnt)
		tot += cnt
	}

	printf("Total arrangements: %d", tot)

	tot2 := 0
	for i := 0; i < len(F); i++ {
		ss := strings.SplitN(lines[i], " ", 2)
		su, nu := unfold(ss[0], N[i], 5)
		debugf("unfold: %s, %+v", su, nu)
		fu := parseField(su)
		cnt := countArrangements(fu, nu)
		debugf("line: %s, cnt: %d", lines[i], cnt)
		tot2 += cnt
	}

	printf("Total arrangements 2: %d", tot2)
}
