package day02

import (
	"aoc/utils"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// "Data" from the game; maps colors to number of that item
type GameData map[string]int

// An ID and a series of game data structures indicating what was drawn
type Game struct {
	id   int
	data []GameData
}

// Parses game data in form [num] [color1], [num] [color2], ...
func parseGameData(gameData string) (GameData, error) {
	scanner := bufio.NewScanner(strings.NewReader(gameData))
	scanner.Split(utils.ScanDelimiterFunc(", "))

	data := GameData{}
	for scanner.Scan() {
		var num int
		var color string
		_, err := fmt.Sscanf(scanner.Text(), "%d %s", &num, &color)
		if err != nil {
			return data, err
		}

		data[color] = num
	}

	return data, nil
}

// Parse a series of games
func parseInput(path string) ([]Game, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var games []Game
	for scanner.Scan() {

		// Break into Game ID and game data
		partScanner := bufio.NewScanner(strings.NewReader(scanner.Text()))
		partScanner.Split(utils.ScanDelimiterFunc(": "))

		if !partScanner.Scan() {
			return nil, errors.New("no ID data found")
		}
		idData := partScanner.Text()

		if !partScanner.Scan() {
			return nil, errors.New("no game data found")
		}
		data := partScanner.Text()

		// Parse Game {int}
		var id int
		_, err = fmt.Sscanf(idData, "Game %d", &id)
		if err != nil {
			return nil, err
		}

		// Scan game data (separated by "; ")
		game := Game{id, nil}
		dataScanner := bufio.NewScanner(strings.NewReader(data))
		dataScanner.Split(utils.ScanDelimiterFunc("; "))
		for dataScanner.Scan() {
			gameData, err := parseGameData(dataScanner.Text())
			if err != nil {
				return nil, err
			}

			game.data = append(game.data, gameData)
		}

		games = append(games, game)
	}

	return games, nil
}

// Get the sum of all game IDs that are possible with the given contents
func getPossibleGameScore(games []Game, contents GameData) int {
	sum := 0

	for _, game := range games {
		possible := true
		for _, data := range game.data {
			for k, v := range data {
				if contents[k] < v {
					possible = false
					break
				}
			}
		}

		if possible {
			sum += game.id
		}
	}

	return sum
}

// Get the "power set" of a bag as defined by the problem
func getPowerSet(contents GameData) int {
	prod := 1
	for _, v := range contents {
		prod *= v
	}

	return prod
}

// Gets the minimum bag contents needed to support the game given
func getMinimalBagContents(game Game) GameData {
	var contents = GameData{}

	for _, data := range game.data {
		for k, v := range data {
			contents[k] = max(contents[k], v)
		}
	}

	return contents
}

func PartA(path string) int {
	games, err := parseInput(path)
	utils.CheckError(err)

	bagContents := map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}

	return getPossibleGameScore(games, bagContents)
}

func PartB(path string) int {
	games, err := parseInput(path)
	utils.CheckError(err)

	sum := 0
	for _, game := range games {
		sum += getPowerSet(getMinimalBagContents(game))
	}

	return sum
}
