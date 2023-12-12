package day12

import (
	"aoc/utils"
	"bufio"
	"os"
	"strconv"
)

type Row struct {
	Springs []byte
	Groups  []int
}

func (row *Row) GroupsMatch() bool {
	var groups []int

	curGroup := 0
	for _, s := range row.Springs {
		if s == '#' {
			curGroup += 1
		} else {
			if curGroup != 0 {
				groups = append(groups, curGroup)
				curGroup = 0
			}
		}
	}

	if curGroup != 0 {
		groups = append(groups, curGroup)
	}

	if len(groups) != len(row.Groups) {
		return false
	}

	for i := range groups {
		if groups[i] != row.Groups[i] {
			return false
		}
	}

	return true
}

func parseInput(path string) ([]Row, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var rows []Row
	for scanner.Scan() {
		parts := utils.NewStringDelimiterScanner(scanner.Text(), " ")
		springs, err := utils.ReadStringFromScanner(parts)
		if err != nil {
			return rows, err
		}

		groupData, err := utils.ReadStringFromScanner(parts)
		if err != nil {
			return rows, err
		}
		groups, err := utils.ReadItems(utils.NewStringDelimiterScanner(groupData, ","), strconv.Atoi, false)
		if err != nil {
			return rows, err
		}

		rows = append(rows, Row{Springs: []byte(springs), Groups: groups})
	}

	return rows, nil
}

// Fill question with . or #, then recurse on next
func doTest(row *Row, questions []int, curQuestionIdx int) int {
	matches := 0
	row.Springs[questions[curQuestionIdx]] = '.'
	if curQuestionIdx == len(questions)-1 {
		if row.GroupsMatch() {
			matches++
		}
	} else {
		matches += doTest(row, questions, curQuestionIdx+1)
	}

	row.Springs[questions[curQuestionIdx]] = '#'
	if curQuestionIdx == len(questions)-1 {
		if row.GroupsMatch() {
			matches++
		}
	} else {
		matches += doTest(row, questions, curQuestionIdx+1)
	}

	row.Springs[questions[curQuestionIdx]] = '?'

	return matches
}

func PartA(path string) int {
	rows, err := parseInput(path)
	utils.CheckError(err)

	sum := 0
	for _, row := range rows {
		var questions []int
		for i, c := range row.Springs {
			if c == '?' {
				questions = append(questions, i)
			}
		}

		sum += doTest(&row, questions, 0)
	}

	return sum
}

func PartB(path string) string {
	_, err := parseInput(path)
	utils.CheckError(err)

	return "Not implemented"
}
