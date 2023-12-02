package main

var digs = []string{
	"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
}

var digsmap = make(map[string]int)

func init() {
	for ix, dig := range digs {
		digsmap[dig] = ix
	}
}

func parseIntWithDigs(dig string) int {
	if ix, ok := digsmap[dig]; ok {
		return ix
	}
	return parseInt(dig)
}

func main() {
	lines := input()
	sum1, sum2 := 0, 0

	for _, line := range lines {
		mf, fsalt := FirstIndexOfAny(line, "0", "1", "2", "3", "4", "5", "6", "7", "8", "9")
		ml, lsalt := LastIndexOfAny(line, "0", "1", "2", "3", "4", "5", "6", "7", "8", "9")
		assert(fsalt != -1, "pattern found")
		assert(lsalt != -1, "pattern found")

		sum1 += parseInt(mf)*10 + parseInt(ml)
	}

	for _, line := range lines {
		mf, fsalt := FirstIndexOfAny(line, "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine")
		ml, lsalt := LastIndexOfAny(line, "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine")

		debugf("line: %s, mf: %s(%d), ml: %s(%d)", line, mf, fsalt, ml, lsalt)

		assert(fsalt != -1, "pattern found")
		assert(lsalt != -1, "pattern found")

		sum2 += parseIntWithDigs(mf)*10 + parseIntWithDigs(ml)
	}

	printf("sum1: %d", sum1)
	printf("sum2: %d", sum2)
}
