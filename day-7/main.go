package main

import (
	"sort"
	"strings"
)

var CARDS = []byte{'A', 'K', 'Q', 'J', 'T', '9', '8', '7', '6', '5', '4', '3', '2'}

var CARD_STRENGTH map[byte]int
var CARD_STRENGTH_WJ map[byte]int

func init() {
	CARD_STRENGTH = make(map[byte]int, len(CARDS))
	CARD_STRENGTH_WJ = make(map[byte]int, len(CARDS))
	for ix, card := range CARDS {
		CARD_STRENGTH[card] = len(CARDS) - ix
		if card != 'J' {
			CARD_STRENGTH_WJ[card] = len(CARDS) - ix
		}
	}
	CARD_STRENGTH_WJ['J'] = 0
}

const (
	HighCardMask  = uint16(0b_000_00_111)
	OnePairMask   = uint16(0b_000_01_000)
	TwoPairMask   = uint16(0b_000_10_000)
	ThreeMask     = uint16(0b_001_00_000)
	FullHouseMask = uint16(0b_001_01_000)
	FourMask      = uint16(0b_010_00_000)
	FiveMask      = uint16(0b_100_00_000)
)

const JOKER_SHIFT = 8

var CNT_SHIFT = [6]int{0, 0, 3, 5, 6, 7}

func getHandMask(s string, MOD uint16) uint16 {
	hasJ := (MOD >> JOKER_SHIFT) > 0
	cnts := make(map[byte]int, len(s))
	for _, b := range []byte(s) {
		cnts[b]++
	}
	msk := uint16(0)
	for b, cnt := range cnts {
		if hasJ && b == 'J' {
			msk += uint16(cnt) << JOKER_SHIFT
		} else {
			msk += 1 << (CNT_SHIFT[cnt])
		}
	}
	if hasJ {
		msk = applyJokerMask(msk)
	}
	debugf("s: %s, msk: %011b", s, msk)
	return msk
}

const (
	NO_JOKER   = uint16(0b_1111_1111)
	WITH_JOKER = uint16(0b_111_1111_1111)
)

func cmpHighestCard(A, B string, MOD uint16) bool {
	hasJ := MOD>>JOKER_SHIFT > 0
	aa, bb := []byte(A), []byte(B)
	for i := 0; i < len(A); i++ {
		if aa[i] == bb[i] {
			continue
		}
		// TODO: this was supposed to be handled by the MOD mask rather than branch
		if hasJ {
			return CARD_STRENGTH_WJ[aa[i]] < CARD_STRENGTH_WJ[bb[i]]
		}
		return CARD_STRENGTH[aa[i]] < CARD_STRENGTH[bb[i]]
	}
	panic("cards are equal???")
}

// detects joker bits and adds them up to the highest-ranked group
func applyJokerMask(m uint16) uint16 {
	j := m >> JOKER_SHIFT

	if j > 0 {
		m &= NO_JOKER
		if m&FiveMask > 0 {
			return m
		} else if m&FourMask > 0 {
			m = m&(^FourMask) | FiveMask
		} else if m&ThreeMask > 0 {
			m = m&(^ThreeMask) | (ThreeMask << j)
		} else if m&TwoPairMask > 0 || m&OnePairMask > 0 {
			pairs := (((m & (OnePairMask | TwoPairMask) >> 3) & 0b11) - 1) << 3
			m = m&(^(OnePairMask | TwoPairMask)) | pairs
			m |= (TwoPairMask << j)
		} else {
			// WTF, oleg?
			delta := uint16(0)
			if m > 0 {
				m -= 1
				delta = 1
			}
			if j == 1 {
				m |= OnePairMask
			} else {
				m |= (TwoPairMask << (j - 2 + delta))
			}
		}
	}
	return m
}

func cmpCardMask(A, B string, MOD uint16) bool {
	MA, MB := getHandMask(A, MOD), getHandMask(B, MOD)
	if MA != MB {
		return MA < MB
	}
	return cmpHighestCard(A, B, MOD)
}

func getWinnings(cards []string, ranks map[string]int) int {
	win := 0
	for rank, card := range cards {
		win += (rank + 1) * ranks[card]
	}
	return win
}

func main() {
	lines := input()
	debugf("file data: %+v", lines)

	cards := make([]string, 0, len(lines))
	ranks := make(map[string]int)

	for _, line := range lines {
		ss := strings.SplitN(line, " ", 2)
		card, rank := ss[0], parseInt(ss[1])
		cards = append(cards, card)
		ranks[card] = rank
	}

	debugf("cards: %+v, ranks: %+v", cards, ranks)

	sort.Slice(cards, func(i, j int) bool {
		return cmpCardMask(cards[i], cards[j], NO_JOKER)
	})
	debugf("sorted cards: %+v", cards)

	win1 := getWinnings(cards, ranks)
	printf("total winnings 1: %d", win1)

	sort.Slice(cards, func(i, j int) bool {
		return cmpCardMask(cards[i], cards[j], WITH_JOKER)
	})
	debugf("sorted cards: %+v", cards)

	win2 := getWinnings(cards, ranks)
	printf("total winnings 2: %d", win2)

}
