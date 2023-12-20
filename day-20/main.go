package main

import "strings"

const (
	FLIP = uint16('%')
	CONJ = uint16('&')
)

const (
	LOW  uint16 = 0
	HIGH        = 1

	FROM_OFF    = 1
	FROM_MASK   = 0x3F
	TO_OFF      = 7
	TO_MASK     = 0x3F
	SIGNAL_OFF  = 0
	SIGNAL_MASK = 1
)

const (
	BUTTON      = "button"
	BROADCASTER = "broadcaster"
)

func main() {
	lines := input()

	// pre-alloc 64 slots, it should be enough
	CONNS := make([][]uint16, 64)
	TYPES := make([]uint16, 64)
	CONJFROM := make([]map[uint16]uint16, 64)
	CONJMASK := make([]uint16, 64)

	IDS := make(map[string]uint16)
	NAMES := make([]string, 64)

	lastid := uint16(0)
	getId := func(name string) uint16 {
		if id, ok := IDS[name]; ok {
			return id
		}
		// start with 1
		lastid++
		IDS[name] = lastid
		return lastid
	}

	for _, line := range lines {
		ss := strings.SplitN(line, " -> ", 2)
		name := ss[0]
		outs := strings.Split(ss[1], ", ")
		outIDs := make([]uint16, 0, len(outs))
		for _, out := range outs {
			outIDs = append(outIDs, getId(out))
		}
		typ := uint16(0)
		if name[0] == '%' || name[0] == '&' {
			typ = uint16(name[0])
			name = name[1:]
		}
		id := getId(name)
		IDS[name] = id
		NAMES[id] = name
		TYPES[id] = typ
		CONNS[id] = outIDs
	}

	buttonId := getId(BUTTON)
	NAMES[buttonId] = BUTTON
	IDS[BUTTON] = buttonId

	CONNS[buttonId] = []uint16{IDS[BROADCASTER]}

	var rxfrom uint16
	for from, tos := range CONNS {
		for _, to := range tos {
			if TYPES[to] == CONJ {
				if CONJFROM[to] == nil {
					CONJFROM[to] = make(map[uint16]uint16)
				}
				newix := len(CONJFROM[to])
				CONJFROM[to][uint16(from)] = uint16(newix)
				CONJMASK[to] |= 1 << newix
			}
			if to == getId("rx") {
				rxfrom = uint16(from)
			}
		}
	}

	STATE := make([]uint16, 64)

	debugf("ids: %+v", IDS)
	debugf("conns: %+v", CONNS)
	debugf("types: %+v", TYPES)
	debugf("state: %+v", STATE)
	debugf("conjfrom: %+v", CONJFROM)
	debugf("conjmask: %+v", CONJMASK)

	low, high := 0, 0

	msk := CONJMASK[rxfrom]
	cycles := make([]int, 0, 4)

	SEND := make([]uint16, 0, 1)

	i := 0
BUTTON_PRESS:
	for {
		i++
		SEND = append(SEND, (IDS[BUTTON]<<TO_OFF)|(LOW<<SIGNAL_OFF))

		var msg uint16
		for len(SEND) > 0 {
			msg, SEND = SEND[0], SEND[1:]
			sign := (msg >> SIGNAL_OFF) & SIGNAL_MASK
			from := (msg >> FROM_OFF) & FROM_MASK
			to := (msg >> TO_OFF) & TO_MASK

			if to == IDS["rx"] {
				state := STATE[from]
				if state > 0 {
					for k, ix := range CONJFROM[from] {
						if (state&(1<<ix) > 0) && (msk&(1<<ix) > 0) {
							debugf("%s is set on button press: %d", NAMES[k], i)
							msk &= ^(1 << ix)
							cycles = append(cycles, i)
							if msk == 0 {
								lcm := cycles[0]
								for k := 1; k < len(cycles); k++ {
									lcm = LCM(lcm, cycles[k])
								}
								debugf("cycles: %+v", cycles)
								printf("rx would receive LOW after %d clicks", lcm)
								break BUTTON_PRESS
							}
						}

					}
				}
			}
			// this would never happen, but still
			if to == IDS["rx"] && sign == LOW {
				printf("Min button press to reach rx: %d", i)
			}

			switch TYPES[to] {
			case FLIP:
				if sign == LOW {
					s := STATE[to]
					s = (s + 1) & 1
					STATE[to] = s
					sign = s
				} else {
					continue
				}
			case CONJ:
				state := STATE[to]
				fromix := CONJFROM[to][from]
				if sign == HIGH {
					state |= (1 << fromix)
					state &= CONJMASK[to]
				} else {
					state &= ^(1 << fromix)
					state &= CONJMASK[to]
				}
				if state == CONJMASK[to] {
					sign = LOW
				} else {
					sign = HIGH
				}
				STATE[to] = state
			default:

			}
			for _, conn := range CONNS[to] {
				SEND = append(SEND, (to<<FROM_OFF)|(conn<<TO_OFF)|(sign<<SIGNAL_OFF))
				if sign == LOW {
					low++
				} else {
					high++
				}
			}
		}

		if i == 1000 {
			printf("low: %d, high: %d, mult: %d", low, high, low*high)
		}
	}

}
