package main

import (
	"sort"
	"strconv"
	"strings"
)

type Mapping [3]uint64

var (
	mappings = []string{
		"seed-to-soil",
		"soil-to-fertilizer",
		"fertilizer-to-water",
		"water-to-light",
		"light-to-temperature",
		"temperature-to-humidity",
		"humidity-to-location",
	}
)

func parseSeeds(s string) []uint64 {
	seeds := make([]uint64, 0, 1)
	for _, v := range strings.Split(strings.SplitN(s, ": ", 2)[1], " ") {
		seeds = append(seeds, parseUInt64(v))
	}
	return seeds
}

func parseUInt64(s string) uint64 {
	n, err := strconv.ParseUint(s, 10, 64)
	noerr(err)
	return n
}

func parseCategory(ss []string, ix int) (string, []Mapping, int) {
	cname := strings.SplitN(ss[ix], " ", 2)[0]
	debugf("parse category %s", cname)
	ix++
	cmps := make([]Mapping, 0, 1)
	for ix < len(ss) {
		if len(ss[ix]) == 0 {
			// category def end
			ix++
			break
		}
		vs := strings.SplitN(ss[ix], " ", 3)
		cmps = append(cmps, Mapping{
			parseUInt64(vs[0]),
			parseUInt64(vs[1]),
			parseUInt64(vs[2]),
		})
		ix++
	}

	sort.Slice(cmps, func(i, j int) bool {
		return cmps[i][1] < cmps[j][1]
	})

	return cname, cmps, ix
}

func getNextRanges(cmps []Mapping, A, B uint64) [][4]uint64 {
	res := make([][4]uint64, 0, 1)
	ix := 0
	a, b := A, B
	for a < b && ix < len(cmps) {
		base, vmin, vmax := cmps[ix][0], cmps[ix][1], cmps[ix][1]+cmps[ix][2]-1
		if b < vmin {
			// other ranges exceed the upper bound
			break
		}
		if a > vmax {
			// the range is left to the current one, we can skip it
			ix++
			continue
		}
		if a >= vmin && b <= vmax {
			res = append(res, [4]uint64{a, b, base + (a - vmin), base + (b - vmin)})
			a = b
		} else if a < vmin {
			res = append(res, [4]uint64{a, vmin - 1, a, vmin - 1})
			a = vmin
		} else if b > vmax {
			res = append(res, [4]uint64{a, vmax, base + (a - vmin), base + (vmax - vmin)})
			a = vmax + 1
		} else {
			panic("how did I even get here?")
		}
		ix++
	}
	if a < b {
		res = append(res, [4]uint64{a, b, a, b})
	}
	return res
}

func lookup(cmps []Mapping, val uint64) uint64 {
	for _, mp := range cmps {
		if val >= mp[1] && val < mp[1]+mp[2] {
			return mp[0] + val - mp[1]
		}
	}
	return val
}

func main() {
	lines := input()

	seeds := parseSeeds(lines[0])
	debugf("seeds: %+v", seeds)

	mps := make(map[string][]Mapping)
	ix := 2
	var cname string
	var cmps []Mapping
	for ix < len(lines) {
		cname, cmps, ix = parseCategory(lines, ix)
		mps[cname] = cmps
	}

	debugf("mps: %+v", mps)

	minloc := ALOT64u

	for _, seed := range seeds {
		v := seed
		for _, m := range mappings {
			v = lookup(mps[m], v)
		}
		minloc = min(minloc, v)
	}

	printf("min loc: %d", minloc)

	minlocall := ALOT64u

	for i := 0; i < len(seeds); i += 2 {
		seedmin, seedmax := seeds[i], seeds[i]+seeds[i+1]-1
		debugf("inspecting seed range: [%d, %d]", seedmin, seedmax)
		ranges := [][4]uint64{{0, 1, seedmin, seedmax}}

		for _, m := range mappings {
			debugf("performing mapping: %s", m)
			newranges := make([][4]uint64, 0, 1)
			for _, rng := range ranges {
				nrgs := getNextRanges(mps[m], rng[2], rng[3])
				newranges = append(newranges, nrgs...)
			}
			sort.Slice(newranges, func(a, b int) bool {
				return newranges[a][2] < newranges[b][2]
			})
			ranges = newranges
		}

		for _, rng := range ranges {
			minlocall = min(minlocall, rng[2])
		}
	}

	printf("min loc all: %d", minlocall)

}
