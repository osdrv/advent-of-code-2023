package main

import (
	"fmt"
	"strings"
)

func findSymbol(F [][]byte, ch byte) Point2 {
	for y := 0; y < len(F); y++ {
		for x := 0; x < len(F[y]); x++ {
			if F[y][x] == ch {
				return Point2{x: x, y: y}
			}
		}
	}
	panic(fmt.Sprintf("Symbol %c not found", ch))
}

var (
	UP    = Point2{x: 0, y: -1}
	DOWN  = Point2{x: 0, y: 1}
	LEFT  = Point2{x: -1, y: 0}
	RIGHT = Point2{x: 1, y: 0}
)

func getNextStep(F [][]byte, P Point2, D Point2) (Point2, Point2, bool) {
	NP := Point2{x: P.x + D.x, y: P.y + D.y}
	if NP.x >= len(F[0]) || NP.x < 0 || NP.y >= len(F) || NP.y < 0 {
		return NP, Point2{0, 0}, false
	}

	ND := Point2{x: 0, y: 0}
	nch := F[NP.y][NP.x]
	if D == UP {
		if nch == '7' {
			ND = LEFT
		} else if nch == 'F' {
			ND = RIGHT
		} else if nch == '|' {
			ND = UP
		}
	} else if D == DOWN {
		if nch == 'J' {
			ND = LEFT
		} else if nch == 'L' {
			ND = RIGHT
		} else if nch == '|' {
			ND = DOWN
		}
	} else if D == RIGHT {
		if nch == '7' {
			ND = DOWN
		} else if nch == 'J' {
			ND = UP
		} else if nch == '-' {
			ND = RIGHT
		}
	} else if D == LEFT {
		if nch == 'L' {
			ND = UP
		} else if nch == 'F' {
			ND = DOWN
		} else if nch == '-' {
			ND = LEFT
		}
	}

	return NP, ND, ND.x != 0 || ND.y != 0
}

var ESCAPE_MEMO = make(map[Point2]bool)

func canEscape(F [][]rune, p Point2) bool {
	if res, ok := ESCAPE_MEMO[p]; ok {
		return res
	}
	V := make(map[Point2]bool)
	V[p] = true
	ps := make([]Point2, 0, 1)
	ps = append(ps, p)
	var cp Point2
	H := make([]Point2, 0, 1)
	res := false
	for len(ps) > 0 {
		cp, ps = ps[0], ps[1:]
		if prev, ok := ESCAPE_MEMO[cp]; ok {
			res = prev
			goto ReturnRes
		}
		H = append(H, cp)
		for _, step := range STEPS4 {
			np := Point2{x: cp.x + step[0], y: cp.y + step[1]}
			if V[np] {
				continue
			}
			if np.x >= len(F[0]) || np.x < 0 || np.y >= len(F) || np.y < 0 {
				res = true
				goto ReturnRes
			}
			if F[np.y][np.x] == '.' {
				ps = append(ps, np)
				V[np] = true
			}
		}
	}
ReturnRes:
	for _, p := range H {
		ESCAPE_MEMO[p] = res
	}
	return res
}

func isPixel(F [][]rune, p Point2) bool {
	y, x := p.y, p.x
	return F[y][x] == '.' && F[y+1][x] == '.' && F[y+1][x+1] == '.'
}

func scaleRoute2x(F [][]byte, R map[Point2]bool, FD, LD Point2) [][]rune {
	NF := make([][]rune, 0, len(F)*2)
	for i := 0; i < len(F)*2; i++ {
		NF = append(NF, make([]rune, len(F[0])*2))
	}
	for p := range R {
		switch F[p.y][p.x] {
		case '7':
			NF[p.y*2][p.x*2] = printRune(F[p.y][p.x])
			NF[p.y*2+1][p.x*2] = printRune('|')
		case 'L':
			NF[p.y*2][p.x*2] = printRune(F[p.y][p.x])
			NF[p.y*2][p.x*2+1] = printRune('-')
		case 'F':
			NF[p.y*2][p.x*2] = printRune(F[p.y][p.x])
			NF[p.y*2][p.x*2+1] = printRune('-')
			NF[p.y*2+1][p.x*2] = printRune('|')
		case 'S':
			NF[p.y*2][p.x*2] = printRune(F[p.y][p.x])
			if FD == RIGHT {
				NF[p.y*2][p.x*2+1] = printRune('-')
			} else if FD == DOWN {
				NF[p.y*2+1][p.x*2] = printRune('|')
			}
			if LD == LEFT {
				NF[p.y*2][p.x*2+1] = printRune('-')
			} else if LD == UP {
				NF[p.y*2+1][p.x*2] = printRune('|')
			}
		case 'J':
			NF[p.y*2][p.x*2] = printRune(F[p.y][p.x])
		case '-':
			NF[p.y*2][p.x*2] = '-'
			NF[p.y*2][p.x*2+1] = '-'
		case '|':
			NF[p.y*2][p.x*2] = '|'
			NF[p.y*2+1][p.x*2] = '|'
		default:
			NF[p.y*2][p.x*2] = '.'
			NF[p.y*2+1][p.x*2] = '.'
			NF[p.y*2][p.x*2+1] = '.'
			NF[p.y*2+1][p.x*2+1] = '.'
		}
	}
	for y := 0; y < len(NF); y++ {
		for x := 0; x < len(NF[0]); x++ {
			if NF[y][x] == 0 {
				NF[y][x] = '.'
			}
		}
	}
	return NF
}

func printRune(ch byte) rune {
	switch ch {
	case '7':
		return '┐'
	case 'F':
		return '┌'
	case 'L':
		return '└'
	case 'J':
		return '┘'
	case '-':
		return '-'
	case '|':
		return '|'
	case 'S':
		return 'S'
	default:
		return '?'
	}
}

func main() {
	lines := input()
	F := make([][]byte, 0, len(lines))
	for _, line := range lines {
		F = append(F, []byte(line))
	}

	S := findSymbol(F, 'S')
	P := S

	debugf("Start: %v", S)

	var D Point2
	var FD, LD Point2
	for _, dc := range []Point2{UP, DOWN, LEFT, RIGHT} {
		if _, _, ok := getNextStep(F, P, dc); ok {
			D = dc
			FD = D
			break
		}
	}

	debugf("starting dir: %v", D)

	ROUTE := make([]Point2, 0, 1)
	V := make(map[Point2]bool)
	for {
		NP := Point2{x: P.x + D.x, y: P.y + D.y}
		if NP.x >= len(F[0]) || NP.x < 0 || NP.y >= len(F) || NP.y < 0 {
			panic("out of field bounds")
		}
		ROUTE = append(ROUTE, NP)
		V[NP] = true
		if F[NP.y][NP.x] == 'S' {
			LD = D
			break
		}
		np, nd, ok := getNextStep(F, P, D)
		if !ok {
			printf("P: %+v, D: %+v", P, D)
		}
		assert(ok, "the next step is defined")
		P = np
		D = nd
	}

	debugf("route: %d", ROUTE)
	printf("farthest steps: %d", len(ROUTE)/2)

	var buf strings.Builder
	for y := 0; y < len(F); y++ {
		for x := 0; x < len(F[0]); x++ {
			p := Point2{x: x, y: y}
			if _, ok := V[p]; ok {
				buf.WriteRune(printRune(F[y][x]))
			} else {
				buf.WriteRune('.')
			}
		}
		buf.WriteByte('\n')
	}
	debugf("\n%s\n", buf.String())

	F2 := scaleRoute2x(F, V, FD, LD)
	var buf2 strings.Builder
	for y := 0; y < len(F2); y++ {
		for x := 0; x < len(F2[0]); x++ {
			buf2.WriteRune(F2[y][x])
		}
		buf2.WriteByte('\n')
	}

	debugf("\n%s\n", buf2.String())

	pixels := 0
	for y := 0; y < len(F2); y += 2 {
		for x := 0; x < len(F2[0]); x += 2 {
			p := Point2{x: x, y: y}
			if !isPixel(F2, p) || canEscape(F2, p) {
				continue
			}
			pixels++
		}
	}
	printf("pixels: %d", pixels)
}
