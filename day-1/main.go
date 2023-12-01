package main

import (
	"os"
	"strings"
)

var digs = []string{
	"one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
}

func main() {
	fi, err := os.Open("INPUT")
	noerr(err)
	defer fi.Close()

	lines := readLines(fi)

	sum1, sum2 := 0, 0

	for _, line := range lines {
		f, l := -1, -1
		fs, ls := ALOT, -ALOT
		for ix, c := range []byte(line) {
			if isNumber(c) {
				if fs == ALOT {
					f = int(c - '0')
					fs = ix
				}
				l = int(c - '0')
				ls = ix
			}
		}
		sum1 += f*10 + l
		debugf("fs: %d, ls: %d", fs, ls)
		for v_1, dig := range digs {
			if ix := strings.Index(line, dig); ix != -1 {
				if ix < fs {
					fs = ix
					f = v_1 + 1
				}
				ix = strings.LastIndex(line, dig)
				if ix > ls {
					ls = ix
					l = v_1 + 1
				}

			}
		}
		debugf("line: %s, f: %d, l: %d", line, f, l)
		sum2 += f*10 + l
	}

	printf("sum1: %d", sum1)
	printf("sum2: %d", sum2)
}
