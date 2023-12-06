package day06

import (
	"aoc/utils"
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Race struct {
	Time           int
	RecordDistance int
}

func parseInput(path string) ([]Race, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	// Get the strings for the time and distance lines
	timeData, err := utils.ReadStringFromScanner(scanner)
	if err != nil {
		return nil, err
	}

	distData, err := utils.ReadStringFromScanner(scanner)
	if err != nil {
		return nil, err
	}

	// Read time and distance as lists by splitting the header off and passing the rest to our parse function
	times, err := utils.ReadItems(utils.NewStringDelimiterScanner(timeData[len("Time:"):], " "), strconv.Atoi, true)
	if err != nil {
		return nil, err
	}

	recordDistances, err := utils.ReadItems(utils.NewStringDelimiterScanner(distData[len("Distance:"):], " "), strconv.Atoi, true)
	if err != nil {
		return nil, err
	}

	// Sanity check
	if len(times) != len(recordDistances) {
		return nil, errors.New("times and distances did not match in length")
	}

	// Take our arrays and translate them to our race structure
	var races []Race
	for i := 0; i < len(times); i++ {
		races = append(races, Race{Time: times[i], RecordDistance: recordDistances[i]})
	}

	return races, nil
}

func parseInput2(path string) (Race, error) {
	f, err := os.Open(path)
	if err != nil {
		return Race{}, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	// Get the strings for the time and distance lines.  Split the headers off so we have just
	// the numbers
	timeData, err := utils.ReadStringFromScanner(scanner)
	if err != nil {
		return Race{}, err
	}
	timeData = timeData[len("Time:"):]

	distData, err := utils.ReadStringFromScanner(scanner)
	if err != nil {
		return Race{}, err
	}
	distData = distData[len("Distance:"):]

	// Remove spaces from the strings to undo the elves bad formatting
	timeData = strings.ReplaceAll(timeData, " ", "")
	distData = strings.ReplaceAll(distData, " ", "")

	// Interpret the digits as numbers
	time, err := strconv.Atoi(timeData)
	if err != nil {
		return Race{}, err
	}

	dist, err := strconv.Atoi(distData)
	if err != nil {
		return Race{}, err
	}

	return Race{Time: time, RecordDistance: dist}, nil
}

func getWaysToWin(race Race) int {
	waysToWin := 0
	for timeIdx := 0; timeIdx < race.Time; timeIdx++ {
		// Figure out how far we would get if we started movement at this time unit.
		// We use this closed-form solution so we don't have to simulate it (we want O(n) rather
		// than O(n^2) or worse)
		dist := timeIdx * (race.Time - timeIdx)

		// If it's a winning distance, count it
		if dist > race.RecordDistance {
			waysToWin++
		}
	}

	return waysToWin
}

func PartA(path string) int {
	races, err := parseInput(path)
	utils.CheckError(err)

	prod := 1
	for _, race := range races {
		prod *= getWaysToWin(race)
	}

	return prod
}

func PartB(path string) int {
	race, err := parseInput2(path)
	utils.CheckError(err)

	return getWaysToWin(race)
}
