package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/osdrv/go-z3"
)

func parseInt64(s string) int64 {
	num, err := strconv.ParseInt(s, 10, 64)
	noerr(err)
	return num
}

func parseInts64(s string) []int64 {
	chs := strings.FieldsFunc(trim(s), func(r rune) bool {
		return r == ' ' || r == ',' || r == '\t'
	})
	nums := make([]int64, 0, len(chs))
	for i := 0; i < len(chs); i++ {
		nums = append(nums, parseInt64(chs[i]))
	}
	return nums
}

func intersectXY(p1, v1, p2, v2 [3]int64) ([3]float64, bool) {
	if v1 == v2 {
		return [3]float64{0, 0, 0}, false
	}

	var k1, k2, b1, b2 float64

	k1 = float64(v1[1]) / float64(v1[0])
	k2 = float64(v2[1]) / float64(v2[0])

	b1 = float64(p1[1]) - k1*float64(p1[0])
	b2 = float64(p2[1]) - k2*float64(p2[0])

	x := (b2 - b1) / (k1 - k2)
	y := k1*x + b1

	return [3]float64{x, y, 0}, true
}

func sign(v float64) float64 {
	if v < 0 {
		return -1
	}
	return 1
}

func isInFuture(h [2][3]int64, p [3]float64) bool {
	sx := sign(p[0] - float64(h[0][0]))
	svx := sign(float64(h[1][0]))
	sy := sign(p[1] - float64(h[0][1]))
	svy := sign(float64(h[1][1]))
	return sx == svx && sy == svy
}

func main() {
	lines := input()

	H := make([][2][3]int64, 0, len(lines))

	for _, line := range lines {
		ss := strings.SplitN(line, " @ ", 2)
		p := parseInts64(ss[0])
		v := parseInts64(ss[1])
		H = append(H, [2][3]int64{
			{p[0], p[1], p[2]},
			{v[0], v[1], v[2]},
		})
	}

	//var MIN, MAX float64 = 7, 27
	var MIN, MAX float64 = 200000000000000, 400000000000000

	cnt := 0
	for i := 0; i < len(H); i++ {
		for j := i + 1; j < len(H); j++ {
			cross, ok := intersectXY(H[i][0], H[i][1], H[j][0], H[j][1])
			if !ok {
				continue
			}
			debugf("intersect: %v %v: %v", H[i], H[j], cross)
			if cross[0] >= MIN && cross[0] <= MAX && cross[1] >= MIN && cross[1] <= MAX {
				fi := isInFuture(H[i], cross)
				fj := isInFuture(H[j], cross)
				if fi && fj {
					cnt++
				}
			}
		}
	}

	printf("intersect: %d", cnt)

	cfg := z3.NewConfig()
	ctx := z3.NewContext(cfg)
	cfg.Close()
	defer ctx.Close()

	ctx.SetErrorHandler(func(ctx *z3.Context, err z3.ErrorCode) {
		printf("z3 error : %+v", err)
	})

	s := ctx.NewSolver()
	defer s.Close()

	x := ctx.Const(ctx.Symbol("x"), ctx.IntSort())
	y := ctx.Const(ctx.Symbol("y"), ctx.IntSort())
	z := ctx.Const(ctx.Symbol("z"), ctx.IntSort())
	vx := ctx.Const(ctx.Symbol("vx"), ctx.IntSort())
	vy := ctx.Const(ctx.Symbol("vy"), ctx.IntSort())
	vz := ctx.Const(ctx.Symbol("vz"), ctx.IntSort())

	for i := 0; i < 3; i++ {
		t := ctx.Const(ctx.Symbol(fmt.Sprintf("t%d", i)), ctx.IntSort())
		a := ctx.Int64(H[i][0][0], ctx.IntSort())
		b := ctx.Int64(H[i][0][1], ctx.IntSort())
		c := ctx.Int64(H[i][0][2], ctx.IntSort())
		va := ctx.Int64(H[i][1][0], ctx.IntSort())
		vb := ctx.Int64(H[i][1][1], ctx.IntSort())
		vc := ctx.Int64(H[i][1][2], ctx.IntSort())

		s.Assert(t.Gt(ctx.Int(0, ctx.IntSort())))
		s.Assert(x.Add(vx.Mul(t)).Eq(a.Add(va.Mul(t))))
		s.Assert(y.Add(vy.Mul(t)).Eq(b.Add(vb.Mul(t))))
		s.Assert(z.Add(vz.Mul(t)).Eq(c.Add(vc.Mul(t))))
	}

	if v := s.Check(); v != z3.True {
		panic("could not solve")
	}

	m := s.Model()
	asg := m.Assignments()
	sum := m.Eval(x.Add(y.Add(z)))
	m.Close()

	debugf("assignments: %+v", asg)
	printf("sum: %d", sum.Int64())
}

func GCD64(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

/*

r = P + t*V

x(t) = P.x + t * V.x
y(t) = P.y + t * V.y
z(t) = P.z + t * V.z

P0, V0 - ?

P1 + t1*V1 = P0 + t1 * V0
P2 + t2 * V2 = P0 + t2 * V0

*/
