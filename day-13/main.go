package main

var CHS = map[rune]int{
	'#': 1,
	'.': 0,
}

func readField(lines []string) [][]int {
	f := make([][]int, 0, len(lines))
	for _, line := range lines {
		row := make([]int, 0, len(line))
		for _, ch := range line {
			row = append(row, CHS[ch])
		}
		f = append(f, row)
	}
	return f
}

var NONE = [2]int{-1, -1}

func findSymmetryAll(F [][]int) [][2]int {
	res := make([][2]int, 0, 1)
	for x0 := 0; x0 < len(F[0]); x0++ {
		x1 := 0
	X1:
		for x0-x1 >= 0 && x0+1+x1 < len(F[0]) {
			for y := 0; y < len(F); y++ {
				if F[y][x0-x1] != F[y][x0+1+x1] {
					break X1
				}
			}
			x1++
		}
		if x1 > 0 && (x0-x1 == -1 || (x0+1+x1 == len(F[0]))) {
			debugf("symmetry in x0: %d with depth: %d", x0, x1)
			res = append(res, [2]int{x0, x1})
		}
	}

	return res
}

func sumRefl(vsym, hsym [2]int) int {
	sum := 0
	if hsym != NONE {
		sum += hsym[0] + 1
	}
	if vsym != NONE {
		sum += (vsym[0] + 1) * 100
	}
	return sum
}

func findDiff(B, L [][2]int) [][2]int {
	M := make(map[[2]int]bool)
	for _, b := range B {
		M[b] = true
	}
	res := make([][2]int, 0, 1)
	for _, l := range L {
		if _, ok := M[l]; !ok {
			res = append(res, l)
		}
	}
	return res
}

func main() {
	lines := input()
	fields := make([][][]int, 0, 1)
	from := 0
	for to := from + 1; to < len(lines); to++ {
		if len(lines[to]) == 0 {
			fields = append(fields, readField(lines[from:to]))
			from = to + 1

		}
	}
	fields = append(fields, readField(lines[from:]))

	tot := 0
	tot2 := 0
	for _, f := range fields {
		debugf("\n%s", printNumFieldWithSubs(f, "", map[int]string{
			1: "#", 0: ".",
		}))
		hsym := findSymmetryAll(f)
		debugf("transpose mat")
		vsym := findSymmetryAll(transposeMat(f))

		sum := 0
		for _, s := range hsym {
			sum += s[0] + 1
		}
		for _, s := range vsym {
			sum += (s[0] + 1) * 100
		}
		debugf("sum: %d", sum)
		tot += sum

		found := false

	FindAnotherSym:

		for y := 0; y < len(f); y++ {
			for x := 0; x < len(f[y]); x++ {
				prev := f[y][x]
				f[y][x] = (prev + 1) & 0b1

				hsym2 := findSymmetryAll(f)
				vsym2 := findSymmetryAll(transposeMat(f))

				debugf("hsym2: %+v, vsym2: %+v", hsym2, vsym2)

				if len(hsym2) != 0 || len(vsym2) != 0 {
					hdiff := findDiff(hsym, hsym2)
					vdiff := findDiff(vsym, vsym2)
					if len(hdiff) != 0 || len(vdiff) != 0 {
						debugf("Another symmetry: hsym2: %+v, vsym2: %+v", hsym2, vsym2)
						sum2 := 0
						if len(hdiff) != 0 {
							sum2 += hdiff[0][0] + 1
						}
						if len(vdiff) != 0 {
							sum2 += (vdiff[0][0] + 1) * 100
						}
						tot2 += sum2
						found = true
					}
				}

				f[y][x] = prev

				if found {
					break FindAnotherSym
				}
			}
		}
		assert(found, "there must be an alternative line")
	}
	printf("tot: %d", tot)
	printf("tot2: %d", tot2)
}
