package main

import (
	"fmt"
	"sort"
	"strings"
)

type Brick struct {
	p1, p2 Point3
	sup    []*Brick
	supBy  []*Brick
}

func NewBrick(p1, p2 Point3) *Brick {
	return &Brick{p1: p1, p2: p2}
}

func parseBrick(s string) *Brick {
	ss := strings.SplitN(s, "~", 2)
	nn1 := parseInts(ss[0])
	nn2 := parseInts(ss[1])

	p1 := Point3{x: nn1[0], y: nn1[1], z: nn1[2]}
	p2 := Point3{x: nn2[0], y: nn2[1], z: nn2[2]}
	if p2.z < p1.z {
		p1, p2 = p2, p1
	}

	return NewBrick(p1, p2)
}

func (b *Brick) SupportedBy(other *Brick) {
	b.supBy = append(b.supBy, other)
}

func (b *Brick) Support(other *Brick) {
	b.sup = append(b.sup, other)
}

func (b *Brick) GetExclusivelySupported() []*Brick {
	sup := make(map[*Brick]struct{})
	for _, s := range b.sup {
		if len(s.supBy) > 1 {
			continue
		}
		sup[s] = struct{}{}
		for _, ss := range s.GetExclusivelySupported() {
			sup[ss] = struct{}{}
		}
	}
	res := make([]*Brick, 0, len(sup))
	for s := range sup {
		res = append(res, s)
	}
	return res
}

func (b *Brick) Move(d Point3) *Brick {
	p1, p2 := b.p1, b.p2
	p1.x += d.x
	p1.y += d.y
	p1.z += d.z
	p2.x += d.x
	p2.y += d.y
	p2.z += d.z
	return NewBrick(p1, p2)
}

func (b *Brick) String() string {
	return fmt.Sprintf("[%v %v]", b.p1, b.p2)
}

func (b *Brick) Overlaps(other *Brick) bool {
	return overlaps3D(b.p1, b.p2, other.p1, other.p2)
}

func overlaps1D(x11, x12, x21, x22 int) bool {
	return (x11 <= x21 && x12 >= x21) ||
		(x11 <= x22 && x12 >= x22) ||
		(x21 <= x11 && x22 >= x11) ||
		(x21 <= x12 && x22 >= x12)
}

func overlaps2D(p11, p12, p21, p22 Point2) bool {
	return overlaps1D(p11.x, p12.x, p21.x, p22.x) && overlaps1D(p11.y, p12.y, p21.y, p22.y)
}

func overlaps3D(p11, p12, p21, p22 Point3) bool {
	return overlaps2D(Point2{p11.x, p11.y}, Point2{p12.x, p12.y}, Point2{p21.x, p21.y}, Point2{p22.x, p22.y}) &&
		overlaps2D(Point2{p11.x, p11.z}, Point2{p12.x, p12.z}, Point2{p21.x, p21.z}, Point2{p22.x, p22.z}) &&
		overlaps2D(Point2{p11.y, p11.z}, Point2{p12.y, p12.z}, Point2{p21.y, p21.z}, Point2{p22.y, p22.z})
}

func computeChain(b *Brick, fail map[*Brick]bool) {
	fail[b] = true
	chain := make([]*Brick, 0, 1)
Sup:
	for _, s := range b.sup {
		for _, sb := range s.supBy {
			if !fail[sb] {
				continue Sup
			}
		}
		chain = append(chain, s)
	}
	for _, b := range chain {
		computeChain(b, fail)
	}
}

func main() {
	lines := input()
	bricks := make([]*Brick, 0, len(lines))
	for _, line := range lines {
		b := parseBrick(line)
		bricks = append(bricks, b)
	}

	sort.Slice(bricks, func(a, b int) bool {
		return min(bricks[a].p1.z, bricks[a].p2.z) < min(bricks[b].p1.z, bricks[b].p2.z)
	})

	debugf("bricks: %v", bricks)

	settled := make([]*Brick, 0, len(bricks))

	for _, brick := range bricks {
		cb := brick.Move(Point3{0, 0, 0})
		for min(cb.p1.z, cb.p2.z) > 1 {
			nb := cb.Move(Point3{0, 0, -1})
			stopFall := false
			for _, b := range settled {
				if nb.Overlaps(b) {
					b.Support(cb)
					cb.SupportedBy(b)
					stopFall = true
				}
			}
			if stopFall {
				break
			}
			cb = nb
		}
		settled = append(settled, cb)
	}

	debugf("settled: %v", settled)

	sort.Slice(settled, func(a, b int) bool {
		return min(settled[a].p1.z, settled[a].p2.z) < min(settled[b].p1.z, settled[b].p2.z)
	})

	must := make([]*Brick, 0, 1)
	free := make(map[*Brick]struct{})
	for _, b := range settled {
		wouldFail := false
		for _, s := range b.sup {
			if len(s.supBy) == 1 {
				wouldFail = true
			}
		}
		if wouldFail {
			must = append(must, b)
		} else {
			free[b] = struct{}{}
		}
	}

	printf("could be desintegrated: %d", len(free))

	willfail := 0
	for _, b := range must {
		f := make(map[*Brick]bool)
		computeChain(b, f)
		willfail += len(f) - 1
	}
	printf("would fail: %d", willfail)
}
