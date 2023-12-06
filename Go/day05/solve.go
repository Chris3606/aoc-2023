package day05

import (
	"aoc/utils"
	"bufio"
	"errors"
	"math"
	"os"
	"strconv"
	"strings"
)

// A structure representing a source range of values, and a set of destination values to map them to.
type RangeMap struct {
	SourceRange utils.Range
	// The start of the destination; precisely SourceRange.Length() elements are included
	// in the destination range, starting from this location
	Destination int
}

// Gets the "delta" value between the source and destination.  If the source range is offset
// by this many values, it becomes the destination range.
func (r RangeMap) GetDelta() int {
	return r.Destination - r.SourceRange.Start
}

// Parses in a RangeMap from the following structure:
// 0 15 37
// The first value is the destination to map the range to, the second value is the source range
// starting location, and the third value is the length
func parseRangeMap(mapData string) (RangeMap, error) {
	numbers, err := utils.ReadItems(utils.NewStringDelimiterScanner(mapData, " "), strconv.Atoi, true)
	if err != nil {
		return RangeMap{}, err
	}

	if len(numbers) != 3 {
		return RangeMap{}, errors.New("incorrect number of integers")
	}

	return RangeMap{SourceRange: utils.NewRange(numbers[1], numbers[2]), Destination: numbers[0]}, nil
}

// Parses the starting seed values and series of map values
func parseInput(path string) ([]int, [][]RangeMap, error) {
	// Open file
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	// Scan through groups separated by new-line new-line
	scanner := bufio.NewScanner(f)
	scanner.Split(utils.ScanDoubleNewLines)

	// Read in list of seed numbers
	seeds, err := utils.ReadItemFromScanner(scanner, func(s string) ([]int, error) {
		return utils.ReadItems(utils.NewStringDelimiterScanner(s[len("seeds: "):], " "), strconv.Atoi, false)
	})
	if err != nil {
		return seeds, nil, err
	}

	// The rest of the double-newline separated groups are maps.  We record them in a slice,
	// since we know the maps apply sequentially to a given value/range
	var maps [][]RangeMap
	for scanner.Scan() {
		// Iterate over each line within a group
		groupIt := bufio.NewScanner(strings.NewReader(scanner.Text()))
		groupIt.Split(bufio.ScanLines)

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

// Maps a series of seed numbers as read in by parseInput to a series of ranges defined by
// pairs of numbers (the part 2 seed input).  For example, [79, 14, 55, 13] becomes
// ranges with start/length pairs: (79, 14), (55, 13)
func seedsToRanges(seedValues []int) []utils.Range {
	var ranges []utils.Range
	for i := 0; i < len(seedValues); i += 2 {
		ranges = append(ranges, utils.NewRange(seedValues[i], seedValues[i+1]))
	}

	return ranges
}

// Given a seed value, apply the sequence of maps to it and produce the result after all maps
// have been applied
func findLocation(seed int, maps [][]RangeMap) int {
	curVal := seed
	for _, mapValues := range maps {
		for _, r := range mapValues {
			if r.SourceRange.ContainsNum(curVal) {
				curVal += r.GetDelta()
				break
			}
		}
	}

	return curVal
}

// Given a single map, (aka a series of ranges that apply source values to destination values),
// apply it to the range of seeds given.
//
// The key here is that we apply maps to RANGES of seed values, NOT individual seed values;
// otherwise the performance is very slow.
//
// Effectively, the algorithm will apply the maps to the seed ranges given, producing a new set
// of ranges which represents the seeds after the map has been applied
//
// In order to produce the modified seed ranges, for each range in the map, we do the following:
//  1. Determine if the map range overlaps with any seed ranges
//  2. For any seeds to which it does overlap, calculate what the overlap is, and add the following
//     to the resulting seed list:
//     - One range which is the overlap of the current map range, but shifted to be in the _destination_
//     range.
//     - Any non-overlapping pieces of the current seed range (unmodified, since the map range doesn't)
//     apply
//
// This implementation performs the algorithm in-place in order to be more efficient, so the
// implementation has to worry about a few caveats:
//  1. The pieces of a seed range that ARE NOT affected by a given RangeMap, could be affected by other
//     RangeMaps within the same map.
//  2. The RangeMap structures within a given map all apply to the _original_ source data, not
//     data modified by previous range maps.
func evaluateMap(seedRange []utils.Range, mapping []RangeMap) []utils.Range {
	// We must not actually shift any ranges until the end, so that all RangeMaps within the same
	// map apply to _original_ values, rather than shifted values.  So, we record the offsets to
	// apply to each range in a separate array, where the indices correspond to indices in seedRange.
	//
	// All offsets start at 0
	var shiftDeltas []int
	for i := 0; i < len(seedRange); i++ {
		shiftDeltas = append(shiftDeltas, 0)
	}

	for _, mapRange := range mapping {
		// We will be modifying the seed range as we go through. The modifications will be one or more
		// of the following:
		//    - Modifying the current element
		//    - Adding elements to the end
		// In either case, we want to ensure that we do _not_ process elements modified or added
		// within the current iteration, because once we add elements we know the current RangeMap
		// can't apply to what we added.  So, we will only iterate over the current elements.
		l := len(seedRange)
		for i := 0; i < l; i++ {
			seed := seedRange[i]
			// If the current RangeMap overlaps, we need to do 3 things:
			// 1. Modify the original range to represent JUST the overlap
			// 2. Modify shiftDeltas to record the delta we must apply to the overlapping range in order
			//    to translate it to the destination range
			// 3. Add any portions of the original range which did NOT overlap with the RangeMap
			//    to the list with a shiftDelta of 0, since the RangeMap does not apply to them.
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
				shiftDeltas[i] = mapRange.GetDelta()

				// Modify the current entry to be only the overlap (we'll adjust at end)
				seedRange[i] = overlapping
			}
		}
	}

	// Now that we've applied all the ranges within the map, apply the shift deltas to the ranges
	// and return the result
	for i := range shiftDeltas {
		seedRange[i].Start += shiftDeltas[i]
		seedRange[i].End += shiftDeltas[i]
	}

	return seedRange
}

// Given a starting set of seed ranges, apply the given set of maps to them and return the results.
func findLocationRanges(seedRange []utils.Range, mappings [][]RangeMap) []utils.Range {
	for _, mapping := range mappings {
		seedRange = evaluateMap(seedRange, mapping)
	}

	return seedRange
}

func PartA(path string) int {
	seeds, maps, err := parseInput(path)
	utils.CheckError(err)

	// Find the minimum location value
	minVal := math.MaxInt
	for _, seed := range seeds {
		minVal = min(minVal, findLocation(seed, maps))
	}

	return minVal
}

func PartB(path string) int {
	seedData, maps, err := parseInput(path)
	utils.CheckError(err)

	// Translate the seed numbers into ranges per the part 2 definition
	seeds := seedsToRanges(seedData)

	// Apply the maps to the seed ranges
	locations := findLocationRanges(seeds, maps)

	// Find the minimum location value
	minVal := math.MaxInt
	for _, r := range locations {
		minVal = min(minVal, r.Start)
	}

	return minVal
}
