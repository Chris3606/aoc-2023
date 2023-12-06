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

	timeData, err := utils.ReadStringFromScanner(scanner)
	if err != nil {
		return nil, err
	}

	distData, err := utils.ReadStringFromScanner(scanner)
	if err != nil {
		return nil, err
	}

	times, err := utils.ReadItems(utils.NewStringDelimiterScanner(timeData[len("Time:"):], " "), strconv.Atoi, true)
	if err != nil {
		return nil, err
	}

	recordDistances, err := utils.ReadItems(utils.NewStringDelimiterScanner(distData[len("Distance:"):], " "), strconv.Atoi, true)
	if err != nil {
		return nil, err
	}

	if len(times) != len(recordDistances) {
		return nil, errors.New("times and distances did not match in length")
	}

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

	timeData = strings.ReplaceAll(timeData, " ", "")
	distData = strings.ReplaceAll(distData, " ", "")

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

func PartA(path string) int {
	races, err := parseInput(path)
	utils.CheckError(err)

	prod := 1
	for _, race := range races {
		waysToWin := 0
		for timeIdx := 0; timeIdx < race.Time; timeIdx++ {
			// Figure out how far we would get if we started movement at this time unit
			dist := timeIdx * (race.Time - timeIdx)

			if dist > race.RecordDistance {
				waysToWin++
			}
		}
		prod *= waysToWin
	}

	return prod
}

func PartB(path string) int {
	race, err := parseInput2(path)
	utils.CheckError(err)

	waysToWin := 0
	for timeIdx := 0; timeIdx < race.Time; timeIdx++ {
		// Figure out how far we would get if we started movement at this time unit
		dist := timeIdx * (race.Time - timeIdx)

		if dist > race.RecordDistance {
			waysToWin++
		}
	}

	return waysToWin
}
