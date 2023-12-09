package main

var (
	MOD_IDENTITY = func(arr []int) []int {
		return arr
	}

	MOD_REVERSE = func(arr []int) []int {
		return reverseNumArr(arr)
	}
)

func getNextVal(nums []int, mod func([]int) []int) int {
	hist := make([][]int, 0, 1)
	cp := make([]int, len(nums))
	copy(cp, nums)
	hist = append(hist, mod(cp))

	hptr := 0
	for {
		all0 := true
		next := make([]int, 0, len(hist[hptr])-1)
		for i := 0; i < len(hist[hptr])-1; i++ {
			nv := hist[hptr][i+1] - hist[hptr][i]
			if nv != 0 {
				all0 = false
			}
			next = append(next, nv)
		}
		hist = append(hist, next)
		hptr++
		if all0 {
			hist[hptr] = append(hist[hptr], 0)
			break
		}
	}

	for i := len(hist) - 1; i > 0; i-- {
		nv := hist[i][len(hist[i])-1] + hist[i-1][len(hist[i-1])-1]
		hist[i-1] = append(hist[i-1], nv)
	}

	return hist[0][len(hist[0])-1]
}

func main() {
	lines := input()

	vals := make([][]int, 0, len(lines))
	for _, line := range lines {
		vals = append(vals, parseInts(line))
	}

	sumNext := 0
	for _, seq := range vals {
		nv := getNextVal(seq, MOD_IDENTITY)
		debugf("seq: %+v, nv: %d", seq, nv)
		sumNext += nv
	}

	printf("sumNext: %d", sumNext)

	sumPrev := 0
	for _, seq := range vals {
		pv := getNextVal(seq, MOD_REVERSE)
		debugf("seq: %+v, pv: %d", seq, pv)
		sumPrev += pv
	}

	printf("sumPrev: %d", sumPrev)
}
