package day07

import (
	"aoc/utils"
	"bufio"
	"os"
	"slices"
	"strconv"
)

type Hand struct {
	Hand []byte
	Bid  int
}

const (
	HighCard     = 1
	OnePair      = 2
	TwoPair      = 3
	ThreeOfAKind = 4
	FullHouse    = 5
	FourOfAKind  = 6
	FiveOfAKind  = 7
)

var cardVals = map[byte]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 11,
	'Q': 12,
	'K': 13,
	'A': 14,
}

var cardVals2 = map[byte]int{
	'J': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'Q': 11,
	'K': 12,
	'A': 13,
}

func parseHand(line string) (Hand, error) {
	spaceScan := utils.NewStringDelimiterScanner(line, " ")

	hand, err := utils.ReadStringFromScanner(spaceScan)
	if err != nil {
		return Hand{}, err
	}

	bid, err := utils.ReadItemFromScanner(spaceScan, strconv.Atoi)
	if err != nil {
		return Hand{}, err
	}

	return Hand{[]byte(hand), bid}, nil
}

func parseInput(path string) ([]Hand, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	hands, err := utils.ReadItems(scanner, parseHand, false)
	if err != nil {
		return nil, err
	}

	return hands, nil
}

func (hand Hand) getStrength() int {
	hist := utils.BuildHistogram(hand.Hand)

	item, maxCount := utils.MaxFromMap(hist)

	switch maxCount {
	case 5:
		return FiveOfAKind
	case 4:
		return FourOfAKind
	case 3:
		hist[item] = 0

		_, newMax := utils.MaxFromMap(hist)
		if newMax == 2 {
			return FullHouse
		}

		// 3 of a kind
		return ThreeOfAKind
	case 2:
		hist[item] = 0
		_, secondMax := utils.MaxFromMap(hist)

		// Two pair
		if secondMax == 2 {
			return TwoPair
		}

		// One pair
		return OnePair
	case 1:
		return HighCard
	default:
		panic("Bad hand.")
	}
}

func (hand Hand) getStrength2() int {
	hist := utils.BuildHistogram(hand.Hand)

	item, maxCount := utils.MaxFromMap(hist)
	if maxCount != 5 {
		if item == 'J' {
			hist[item] = 0
			nextItem, _ := utils.MaxFromMap(hist)
			hist[nextItem] += maxCount
		} else {
			hist[item] += hist['J']
		}

		hist['J'] = 0
		item, maxCount = utils.MaxFromMap(hist)
	}

	switch maxCount {
	case 5:
		return FiveOfAKind
	case 4:
		return FourOfAKind
	case 3:
		hist[item] = 0

		_, newMax := utils.MaxFromMap(hist)
		if newMax == 2 {
			return FullHouse
		}

		// 3 of a kind
		return ThreeOfAKind
	case 2:
		hist[item] = 0
		_, secondMax := utils.MaxFromMap(hist)

		// Two pair
		if secondMax == 2 {
			return TwoPair
		}

		// One pair
		return OnePair
	case 1:
		return HighCard
	default:
		panic("Bad hand.")
	}
}

func compareHands(h1, h2 Hand) int {
	h1Str := h1.getStrength()
	h2Str := h2.getStrength()
	diff := h1Str - h2Str

	if diff != 0 {
		return diff
	}

	for i := 0; i < len(h1.Hand); i++ {
		i1 := cardVals[h1.Hand[i]]
		i2 := cardVals[h2.Hand[i]]
		if i1 > i2 {
			return 1
		} else if i1 < i2 {
			return -1
		}
	}

	return 0
}

func compareHands2(h1, h2 Hand) int {
	h1Str := h1.getStrength2()
	h2Str := h2.getStrength2()
	diff := h1Str - h2Str

	if diff != 0 {
		return diff
	}

	for i := 0; i < len(h1.Hand); i++ {
		i1 := cardVals2[h1.Hand[i]]
		i2 := cardVals2[h2.Hand[i]]
		if i1 > i2 {
			return 1
		} else if i1 < i2 {
			return -1
		}
	}

	return 0
}

func PartA(path string) int {
	hands, err := parseInput(path)
	utils.CheckError(err)

	slices.SortFunc(hands, compareHands)

	winnings := 0
	for i, h := range hands {
		winnings += (i + 1) * h.Bid
	}

	return winnings
}

func PartB(path string) int {
	hands, err := parseInput(path)
	utils.CheckError(err)

	slices.SortFunc(hands, compareHands2)

	winnings := 0
	for i, h := range hands {
		winnings += (i + 1) * h.Bid
	}

	return winnings
}
