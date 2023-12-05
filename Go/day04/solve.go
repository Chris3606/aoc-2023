package day04

import (
	"aoc/utils"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	Id             int
	WinningNumbers map[int]bool
	ChosenNumbers  map[int]bool
	CopiesOwned    int
}

func NewCard(id int, winningNumbers map[int]bool, chosenNumbers map[int]bool) Card {
	return Card{id, winningNumbers, chosenNumbers, 1}
}

func parseCardNumbers(s string) (map[int]bool, error) {
	numbersIt := bufio.NewScanner(strings.NewReader(s))
	numbersIt.Split(utils.ScanDelimiterFunc(" "))

	return utils.ReadItemsToMap(numbersIt, strconv.Atoi, true)
}

func parseInput(path string) ([]Card, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var cards []Card
	for scanner.Scan() {
		partsIt := utils.NewStringDelimiterScanner(scanner.Text(), ": ")

		// Get card ID
		id, err := utils.ReadItemFromScanner(partsIt, func(s string) (int, error) {
			id := 0
			_, err = fmt.Sscanf(s, "Card %d", &id)
			return id, err
		})
		if err != nil {
			return cards, err
		}

		// Get card data
		cardData, err := utils.ReadStringFromScanner(partsIt)
		if err != nil {
			return cards, err
		}

		// Parse winning and chosen numbers
		partsIt = utils.NewStringDelimiterScanner(cardData, " | ")
		winningNumbers, err := utils.ReadItemFromScanner(partsIt, parseCardNumbers)
		if err != nil {
			return cards, err
		}

		chosenNumbers, err := utils.ReadItemFromScanner(partsIt, parseCardNumbers)
		if err != nil {
			return cards, err
		}

		// Add card
		cards = append(cards, NewCard(id, winningNumbers, chosenNumbers))
	}

	return cards, nil
}

func (card Card) countWinningNumbers() int {
	sum := 0
	for k := range card.WinningNumbers {
		if card.ChosenNumbers[k] {
			sum++
		}
	}

	return sum
}

func PartA(path string) int {
	card, err := parseInput(path)
	utils.CheckError(err)

	score := 0
	for _, card := range card {
		numbers := card.countWinningNumbers()
		if numbers == 0 {
			continue
		}
		score += (1 << (numbers - 1))
	}

	return score
}

func PartB(path string) int {
	cards, err := parseInput(path)
	utils.CheckError(err)

	// Propagate cards
	for i := range cards {
		winningNumbers := cards[i].countWinningNumbers()

		for j := 1; j <= winningNumbers; j++ {
			cards[i+j].CopiesOwned += cards[i].CopiesOwned
		}
	}

	// count cards
	cardCount := 0
	for _, card := range cards {
		cardCount += card.CopiesOwned
	}

	return cardCount
}
