package day09

import (
	"aoc/utils"
	"bufio"
	"os"
	"strconv"
)

func parseInput(path string) ([][]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var result [][]int
	for scanner.Scan() {
		nums, err := utils.ReadItems(utils.NewStringDelimiterScanner(scanner.Text(), " "), strconv.Atoi, true)
		if err != nil {
			return result, err
		}

		result = append(result, nums)
	}

	return result, nil
}

func getDiffSequence(seq []int) []int {
	var diffSeq []int
	for i := 1; i < len(seq); i += 1 {
		diffSeq = append(diffSeq, seq[i]-seq[i-1])
	}

	return diffSeq
}

func isZeroSeq(seq []int) bool {
	for _, v := range seq {
		if v != 0 {
			return false
		}
	}

	return true
}

func PartA(path string) int {
	nums, err := parseInput(path)
	utils.CheckError(err)

	sum := 0
	for _, seq := range nums {
		diffs := [][]int{seq}
		for {
			diffs = append(diffs, getDiffSequence(diffs[len(diffs)-1]))
			if isZeroSeq(diffs[len(diffs)-1]) {
				break
			}
		}

		diffs[len(diffs)-1] = append(diffs[len(diffs)-1], 0)
		for i := len(diffs) - 2; i >= 0; i-- {
			diffs[i] = append(diffs[i], diffs[i][len(diffs[i])-1]+diffs[i+1][len(diffs[i+1])-1])
		}

		sum += diffs[0][len(diffs[0])-1]
	}

	return sum
}

func PartB(path string) int {
	nums, err := parseInput(path)
	utils.CheckError(err)

	sum := 0
	for _, seq := range nums {
		diffs := [][]int{seq}
		for {
			diffs = append(diffs, getDiffSequence(diffs[len(diffs)-1]))
			if isZeroSeq(diffs[len(diffs)-1]) {
				break
			}
		}

		diffs[len(diffs)-1] = utils.InsertToSlice(diffs[len(diffs)-1], 0, 0)
		for i := len(diffs) - 2; i >= 0; i-- {
			diffs[i] = utils.InsertToSlice(diffs[i], 0, diffs[i][0]-diffs[i+1][0])
		}

		sum += diffs[0][0]
	}

	return sum
}
