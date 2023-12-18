package main

import (
	"strconv"
)

const (
	RIGHT int64 = 0
	DOWN        = 1
	LEFT        = 2
	UP          = 3
)

var STEPS = map[int64]Point2{
	RIGHT: {x: 1, y: 0},
	DOWN:  {x: 0, y: 1},
	LEFT:  {x: -1, y: 0},
	UP:    {x: 0, y: -1},
}

var DIR = map[byte]int64{
	'R': RIGHT,
	'D': DOWN,
	'L': LEFT,
	'U': UP,
}

func parseInstr(s string) (int64, int, int64) {
	var dir string
	var off int
	var clrs string

	Scanf(s, "{string} {int} (#{string})", &dir, &off, &clrs)

	clr, err := strconv.ParseInt("0x"+clrs, 0, 64)
	noerr(err)

	return DIR[dir[0]], off, clr
}

func parseInstr2(s string) (int64, int64) {
	var dir string
	var y int
	var clrs string
	Scanf(s, "{string} {int} (#{string})", &dir, &y, &clrs)
	off, err := strconv.ParseInt("0x"+clrs[:5], 0, 64)
	noerr(err)
	return int64(clrs[5] - '0'), off
}

// https://en.wikipedia.org/wiki/Shoelace_formula
func findArea(moves [][2]int64) int64 {
	S := int64(0)
	x, y := int64(0), int64(0)
	length := int64(0)
	for _, move := range moves {
		nx := x + move[1]*int64(STEPS[move[0]].x)
		ny := y + move[1]*int64(STEPS[move[0]].y)
		length += move[1]

		S += x*ny - y*nx

		x, y = nx, ny
	}

	return S/2 + length/2 + 1
}

func main() {
	lines := input()
	moves := make([][2]int64, 0, len(lines))
	moves2 := make([][2]int64, 0, len(lines))
	for _, line := range lines {
		dir, off, _ := parseInstr(line)
		moves = append(moves, [2]int64{int64(dir), int64(off)})
		dir2, off2 := parseInstr2(line)
		moves2 = append(moves2, [2]int64{dir2, off2})
	}

	debugf("moves: %+v", moves)
	debugf("moves2: %+v", moves2)

	area := findArea(moves)
	printf("area: %d", area)

	area2 := findArea(moves2)
	printf("area2: %d", area2)
}
