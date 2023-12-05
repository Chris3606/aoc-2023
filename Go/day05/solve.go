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

func seedsToRanges(seedValues []int) []utils.Range {
	var ranges []utils.Range
	for i := 0; i < len(seedValues); i += 2 {
		ranges = append(ranges, utils.NewRange(seedValues[i], seedValues[i+1]))
	}

	return ranges
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

func evaluateMap(seedRange []utils.Range, mapping []RangeMap) []utils.Range {

	var shiftDeltas []int
	for i := 0; i < len(seedRange); i++ {
		shiftDeltas = append(shiftDeltas, 0)
	}

	for _, mapRange := range mapping {
		l := len(seedRange)
		for i := 0; i < l; i++ {
			seed := seedRange[i]
			if seed.Overlaps(mapRange.SourceRange) {
				// Split off the portion that overlaps
				overlapping := utils.Range{Start: max(seed.Start, mapRange.SourceRange.Start), End: min(seed.End, mapRange.SourceRange.End)}

				// Add any portions that don't overlap back to the list (won't be processed this range since we recorded the length prior to modification)
				if overlapping.End < seed.End {
					seedRange = append(seedRange, utils.Range{Start: overlapping.End + 1, End: seed.End})
					shiftDeltas = append(shiftDeltas, 0)
				}
				if overlapping.Start > seed.Start {
					seedRange = append(seedRange, utils.Range{Start: seed.Start, End: overlapping.Start - 1})
					shiftDeltas = append(shiftDeltas, 0)
				}

				// Record delta to adjust by at end
				delta := mapRange.Destination - mapRange.SourceRange.Start
				shiftDeltas[i] = delta

				// Modify the current entry to be only the overlap (we'll adjust at end)
				seedRange[i] = overlapping
			}
		}
	}

	for i := range shiftDeltas {
		seedRange[i].Start += shiftDeltas[i]
		seedRange[i].End += shiftDeltas[i]
	}

	return seedRange
}

func findLocationRanges(seedRange []utils.Range, mappings [][]RangeMap) []utils.Range {
	for _, mapping := range mappings {
		seedRange = evaluateMap(seedRange, mapping)
	}

	return seedRange
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

func PartB(path string) int {
	seedData, maps, err := parseInput(path)
	utils.CheckError(err)

	seeds := seedsToRanges(seedData)

	seeds = findLocationRanges(seeds, maps)

	minVal := math.MaxInt
	for _, seedRange := range seeds {
		minVal = min(minVal, seedRange.Start)
	}

	return minVal
}
