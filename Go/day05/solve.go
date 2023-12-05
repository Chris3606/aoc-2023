package day05

import (
	"aoc/utils"
	"bufio"
	"math"
	"os"
	"strconv"
)

type RangeMap struct {
	SourceRange utils.Range
	Destination int
}

func parseRangeMap(mapData string) (RangeMap, error) {
	numIt := utils.NewStringDelimiterScanner(mapData, " ")

	destination, err := utils.ReadItemFromScanner(numIt, strconv.Atoi)
	if err != nil {
		return RangeMap{}, err
	}

	srcStart, err := utils.ReadItemFromScanner(numIt, strconv.Atoi)
	if err != nil {
		return RangeMap{}, err
	}

	length, err := utils.ReadItemFromScanner(numIt, strconv.Atoi)
	if err != nil {
		return RangeMap{}, err
	}

	return RangeMap{SourceRange: utils.NewRange(srcStart, length), Destination: destination}, nil
}

func parseInput(path string) ([]int, [][]RangeMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(utils.ScanDelimiterFunc("\r\n\r\n"))

	seeds, err := utils.ReadItemFromScanner(scanner, func(s string) ([]int, error) {
		return utils.ReadItems(utils.NewStringDelimiterScanner(s[len("seeds: "):], " "), strconv.Atoi, false)
	})
	if err != nil {
		return seeds, nil, err
	}

	var maps [][]RangeMap
	for scanner.Scan() {
		// Parse groups
		groupIt := utils.NewStringDelimiterScanner(scanner.Text(), "\r\n")

		// Ignore title
		_, err = utils.ReadStringFromScanner(groupIt)
		if err != nil {
			return seeds, maps, err
		}

		// Parse range maps
		parsedMap, err := utils.ReadItems(groupIt, parseRangeMap, false)
		if err != nil {
			return seeds, maps, err
		}

		maps = append(maps, parsedMap)
	}

	return seeds, maps, nil
}

func findLocation(seed int, maps [][]RangeMap) int {
	curVal := seed
	for _, mapValues := range maps {
		for _, r := range mapValues {
			if r.SourceRange.ContainsNum(curVal) {
				curVal = r.Destination + (curVal - r.SourceRange.Start)
				break
			}
		}
	}

	return curVal
}

func PartA(path string) int {
	seeds, maps, err := parseInput(path)
	utils.CheckError(err)

	minVal := math.MaxInt
	for _, seed := range seeds {
		minVal = min(minVal, findLocation(seed, maps))

	}

	return minVal
}

func PartB(path string) string {
	_, _, err := parseInput(path)
	utils.CheckError(err)

	return "Not implemented"
}
