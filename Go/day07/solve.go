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

var cardValuesPart1 = map[byte]int{
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

var cardValuesPart2 = map[byte]int{
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

func (hand Hand) getStrength(replaceJokers bool) int {
	hist := utils.BuildHistogram(hand.Hand)

	item, maxCount := utils.MaxFromMap(hist)

	// The best thing to do with a joker (if we're using them) is simply to add it to whatever the
	// largest number of cards is; more of a card is always better than anything you can do with
	// an equal number of a given card
	if replaceJokers && maxCount != 5 {
		if item == 'J' {
			// If jokers was the max, we'll dump them into whatever is the second highest
			hist[item] = 0
			nextItem, _ := utils.MaxFromMap(hist)
			hist[nextItem] += maxCount
		} else {
			// Otherwise, just jump the jokers into the highest card
			hist[item] += hist['J']
		}

		// Re-calculate the max with jokers (and remove the jokers from the histogram)
		hist['J'] = 0
		item, maxCount = utils.MaxFromMap(hist)
	}

	// The max count of cards will tell the strength
	switch maxCount {
	case 5:
		return FiveOfAKind
	case 4:
		return FourOfAKind
	case 3:
		// Check next highest; if there is 3 + 2, we have full house
		hist[item] = 0
		_, newMax := utils.MaxFromMap(hist)
		if newMax == 2 {
			return FullHouse
		}

		// Otherwise, 3 of a kind
		return ThreeOfAKind
	case 2:
		// Check next highest; if there is 2 + 2, we have two pair
		hist[item] = 0
		_, secondMax := utils.MaxFromMap(hist)
		if secondMax == 2 {
			return TwoPair
		}

		// Otherwise, one pair
		return OnePair
	case 1: // Just a high card
		return HighCard
	default:
		panic("Bad hand.")
	}
}

func getHandComparer(strengthFunc func(Hand) int, cardVals map[byte]int) func(Hand, Hand) int {
	return func(h1, h2 Hand) int {
		h1Str := strengthFunc(h1)
		h2Str := strengthFunc(h2)
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
}

func calculateTotalWinnings(hands []Hand, replaceJokers bool, cardVals map[byte]int) int {
	slices.SortFunc(hands, getHandComparer(func(h Hand) int { return h.getStrength(replaceJokers) }, cardVals))

	winnings := 0
	for i, h := range hands {
		winnings += (i + 1) * h.Bid
	}

	return winnings
}

func PartA(path string) int {
	hands, err := parseInput(path)
	utils.CheckError(err)

	return calculateTotalWinnings(hands, false, cardValuesPart1)
}

func PartB(path string) int {
	hands, err := parseInput(path)
	utils.CheckError(err)

	return calculateTotalWinnings(hands, true, cardValuesPart2)
}
