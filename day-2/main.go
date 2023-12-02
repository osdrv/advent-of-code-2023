package main

import (
	"fmt"
	"strings"
)

type Game struct {
	id    int
	plays [][3]int
}

var COLORS = map[string]int{
	"red":   0,
	"green": 1,
	"blue":  2,
}

func parseGame(in string) *Game {
	g := &Game{}
	pp := strings.Split(in, ": ")
	var id int
	_, err := fmt.Sscanf(pp[0], "Game %d", &id)
	noerr(err)
	body := pp[1]

	g.id = id
	g.plays = make([][3]int, 0, 1)

	for _, p := range strings.Split(body, "; ") {
		var play [3]int
		for _, c := range strings.Split(p, ", ") {
			var num int
			var color string
			_, err := fmt.Sscanf(c, "%d %s", &num, &color)
			noerr(err)
			play[COLORS[color]] = num
		}
		g.plays = append(g.plays, play)
	}

	debugf("game: %+v", g)

	return g
}

var CAP_CUBES = [3]int{12, 13, 14}

func main() {
	lines := input()

	printf("file data: %+v", lines)

	games := make([]*Game, 0, len(lines))
	for _, line := range lines {
		game := parseGame(line)
		games = append(games, game)
	}

	sum1 := 0
GAME:
	for _, game := range games {
		printf("game: %+v", game)
		for _, play := range game.plays {
			if !(play[0] <= CAP_CUBES[0] && play[1] <= CAP_CUBES[1] && play[2] <= CAP_CUBES[2]) {
				printf("impossible")
				continue GAME
			}
		}
		sum1 += game.id
	}

	printf("sum1: %d", sum1)

	sum2 := 0
	for _, game := range games {
		var mc [3]int
		for _, play := range game.plays {
			mc[0] = max(mc[0], play[0])
			mc[1] = max(mc[1], play[1])
			mc[2] = max(mc[2], play[2])
		}
		sum2 += mc[0] * mc[1] * mc[2]
	}

	printf("sum2: %d", sum2)
}
