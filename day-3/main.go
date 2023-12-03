package main

import "fmt"

func isAdj(field [][]byte, x int, y int) (bool, int, int, byte) {
	maxy, maxx := sizeNumField(field)
	for _, step := range STEPS8 {
		nx, ny := x+step[0], y+step[1]
		if nx < 0 || ny < 0 || nx >= maxx || ny >= maxy {
			continue
		}
		if ch := field[ny][nx]; ch == '.' || isNumber(ch) {
			continue
		} else {
			return true, nx, ny, ch
		}
	}
	return false, -1, -1, 0
}

func main() {
	lines := input()
	debugf("file data: %+v", lines)

	field := make([][]byte, 0, len(lines))
	for _, line := range lines {
		field = append(field, []byte(line))
	}

	nums := make([]int, 0, 1)
	var num = new(int)
	var adj bool
	adjmx := make(map[Point3][]*int)
	for y, line := range field {
		for x, ch := range line {
			if isNumber(ch) {
				numv := *num*10 + int(ch-'0')
				*num = numv
				debugf("ch: %c, num upd: %d", ch, *num)
				found, adjx, adjy, adjch := isAdj(field, x, y)
				adj = adj || found
				if found {
					debugf("found adjacency")
					ix := Point3{adjx, adjy, int(adjch)}
					if _, ok := adjmx[ix]; !ok {
						adjmx[ix] = make([]*int, 0, 1)
					}
					if len(adjmx[ix]) == 0 || adjmx[ix][len(adjmx[ix])-1] != num {
						adjmx[ix] = append(adjmx[ix], num)
						debugf("adj: %+v", adjmx[ix])
					}
				}
			} else if *num > 0 {
				debugf("found number: %d, isAdj: %t", *num, adj)
				if adj {
					nums = append(nums, *num)
				}
				num = new(int)
				*num = 0
				adj = false
			}
		}
		if *num > 0 {
			debugf("found number: %d, isAdj: %t", num, adj)
			if adj {
				nums = append(nums, *num)
			}
			num = new(int)
			*num = 0
			adj = false
		}
	}

	sum1 := 0
	for _, num := range nums {
		sum1 += num
	}
	printf("sum1: %d", sum1)

	gear := uint64(0)
	for p3, nn := range adjmx {
		if p3.z != int('*') || len(nn) != 2 {
			debugf("found gear candidate: %c(%s)", byte(p3.z), printIntPtrs(nn))
			continue
		}
		debugf("found gear: %d * %d", *nn[0], *nn[1])
		gear += uint64(*nn[0] * *nn[1])
	}

	printf("gear: %d", gear)
}

func printIntPtrs(nn []*int) string {
	ns := make([]int, 0, len(nn))
	for _, n := range nn {
		ns = append(ns, *n)
	}
	return fmt.Sprintf("%+v", ns)
}
