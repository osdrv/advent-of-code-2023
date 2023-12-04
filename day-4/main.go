package main

import "strings"

func readCard(s string) ([]int, []int) {
	ss := strings.SplitN(s, ": ", 2)
	sss := strings.SplitN(ss[1], " | ", 2)
	return parseInts(sss[0]), parseInts(sss[1])
}

func getNumWins(wn, on []int) int {
	wnm := make(map[int]bool)
	for _, n := range wn {
		wnm[n] = true
	}
	v := 0
	for _, n := range on {
		if wnm[n] {
			v += 1
		}
	}
	return v
}

func main() {
	lines := input()
	debugf("file data: %+v", lines)

	mult := make([]int, len(lines))
	for i := 0; i < len(mult); i++ {
		mult[i] = 1
	}

	cnt1 := 0
	cnt2 := 0
	for ix, line := range lines {
		wn, on := readCard(line)
		v := getNumWins(wn, on)

		if v > 0 {
			cnt1 += 1 << (v - 1)
		}
		cnt2 += mult[ix]

		for i := 1; i <= v; i++ {
			mult[ix+i] += mult[ix]
		}
	}

	debugf("%+v", mult)

	printf("total worth: %d", cnt1)
	printf("total num cards: %d", cnt2)

}
