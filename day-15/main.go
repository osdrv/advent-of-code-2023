package main

import (
	"sort"
	"strings"
)

type Op byte

const (
	EQL   Op = '='
	MINUS    = '-'
)

func hash(s string) int {
	h := 0
	for _, ch := range s {
		h = ((h + int(ch)) * 17) % 256
	}
	return h
}

func parseStep(s string) (string, Op, int) {
	ptr := 0
	var n int
	var l string
	l, ptr = readWord(s, ptr)
	var op Op
	if match(s, ptr, byte(EQL)) {
		ptr = consume(s, ptr, byte(EQL))
		op = EQL
		n, ptr = readInt(s, ptr)
	} else if match(s, ptr, byte(MINUS)) {
		op = MINUS
	}

	return l, op, n
}

type Lens struct {
	H int
	L string
	P int
	S int
}

func main() {
	lines := input()
	debugf("file data: %+v", lines)

	for _, line := range lines {
		steps := strings.Split(line, ",")
		sum := 0
		lens := make(map[string]Lens)
		index := 0
		FP := 0
		for _, step := range steps {
			h := hash(step)
			debugf("step: %s, h: %d", step, h)
			sum += h

			label, op, s := parseStep(step)

			if op == MINUS {
				if _, ok := lens[label]; ok {
					delete(lens, label)
				}
			} else if op == EQL {
				if l, ok := lens[label]; ok {
					l.S = s
					lens[label] = l
				} else {
					lens[label] = Lens{
						L: label,
						P: index,
						S: s,
						H: hash(label),
					}
					index += 1
				}
			} else {
				panic("wtf")
			}

			debugf("after %s: %+v", step, lens)
		}

		LL := make([]Lens, 0, len(lens))
		for _, l := range lens {
			LL = append(LL, l)
		}

		sort.Slice(LL, func(a, b int) bool {
			if LL[a].H != LL[b].H {
				return LL[a].H < LL[b].H
			}
			return LL[a].P < LL[b].P
		})

		box := 1
		pos := 1
		LL[0].P = pos
		for i := 1; i < len(LL); i++ {
			if LL[i].H != LL[i-1].H {
				box += 1
				pos = 1
			} else {
				pos += 1
			}
			LL[i].P = pos
		}

		debugf("LL=%+v", LL)

		for _, l := range LL {
			v := (1 + l.H) * l.P * l.S
			debugf("%s=%d", l.L, v)
			FP += v
		}

		printf("sum=%d", sum)
		printf("FP=%d", FP)
	}
}
