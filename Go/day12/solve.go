package day12

import (
	"aoc/utils"
	"bufio"
	"os"
	"slices"
	"strconv"
)

type Row struct {
	Springs []byte
	Groups  []int
}

func (row *Row) GroupsMatch() bool {
	groups := findGroupsForSprings(row.Springs)
	return utils.CompareSlicesElementwise(groups, row.Groups)
}

func findGroupsForSprings(springs []byte) []int {
	var groups []int

	curGroup := 0
	for _, s := range springs {
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

	return groups
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
// func doTest(row *Row, questions []int, curQuestionIdx int) int {
// 	matches := 0
// 	row.Springs[questions[curQuestionIdx]] = '.'
// 	curGroups := findGroupsForSprings(row.Springs[0 : questions[curQuestionIdx]+1])
// 	if len(row.Groups) >= len(curGroups) && utils.CompareSlicesElementwise(curGroups, row.Groups[:len(curGroups)]) {
// 		if curQuestionIdx == len(questions)-1 {
// 			if row.GroupsMatch() {
// 				matches++
// 			}
// 		} else {
// 			matches += doTest(row, questions, curQuestionIdx+1)
// 		}
// 	}

// 	row.Springs[questions[curQuestionIdx]] = '#'
// 	curGroups = findGroupsForSprings(row.Springs[0 : questions[curQuestionIdx]+1])
// 	curGroups = curGroups[:len(curGroups)-1] // Cut off last group since we might still add to it later
// 	if len(row.Groups) >= len(curGroups) && utils.CompareSlicesElementwise(curGroups, row.Groups[:len(curGroups)]) {
// 		if curQuestionIdx == len(questions)-1 {
// 			if row.GroupsMatch() {
// 				matches++
// 			}
// 		} else {
// 			matches += doTest(row, questions, curQuestionIdx+1)
// 		}
// 	}

// 	row.Springs[questions[curQuestionIdx]] = '?'

// 	return matches
// }

func countSpringPermutations(unprocessed_springs []byte, cur_groups []int) int {
	// Base case: no more springs; so if we've satisfied all the groups, we're good; if there's groups
	// left, we haven't matched
	if len(unprocessed_springs) == 0 {
		if len(cur_groups) == 0 {
			return 1
		} else {
			return 0
		}
	}

	// First spring
	ch := unprocessed_springs[0]

	switch ch {
	// Working springs don't contribute to groups, so they are irrelevant; we'll skip past them
	case '.':
		idx := slices.IndexFunc(unprocessed_springs, func(spring byte) bool { return spring != '.' })
		if idx == -1 {
			unprocessed_springs = nil
		} else {
			unprocessed_springs = unprocessed_springs[idx:]
		}

		return countSpringPermutations(unprocessed_springs, cur_groups)

	// Unknown spring; count permutations with it both ways
	case '?':
		permutations := 0
		unprocessed_springs[0] = '.'
		permutations += countSpringPermutations(unprocessed_springs, cur_groups)

		unprocessed_springs[0] = '#'
		permutations += countSpringPermutations(unprocessed_springs, cur_groups)

		unprocessed_springs[0] = '?'
		return permutations

	// Damaged spring; so see if the current group is possible and recurse if needed
	case '#':
		// We found a damaged spring, but there is no group for it to belong to
		if len(cur_groups) == 0 {
			return 0
		}

		// The current group requires x springs, so including this one there must be at least
		// x springs left
		cur_group := cur_groups[0]
		springs_left := len(unprocessed_springs)
		if cur_group > springs_left {
			return 0
		}

		// Specifically, of the springs remaining, there must be at least x consecutive that are
		// damaged or potentially damaged
		group_found := 0
		for i := 0; i < cur_group; i++ {
			if unprocessed_springs[i] != '.' {
				group_found++
			}
		}
		if group_found < cur_group { // Too few springs in group
			return 0
		}

		// The spring after the current group (if any), must NOT be broken; otherwise the group is
		// _bigger_ than the one we're satisfying
		unprocessed_springs = unprocessed_springs[cur_group:]
		if len(unprocessed_springs) > 0 {
			if unprocessed_springs[0] == '#' {
				return 0
			}

			// The next one _has_ to be considered unbroken; so we can just skip past it
			unprocessed_springs = unprocessed_springs[1:]
		}

		// The current group is satisfied; so recurse on the ones after this group
		return countSpringPermutations(unprocessed_springs, cur_groups[1:])

	default:
		panic("Invalid spring state.")
	}
}

func expandRows(rows []Row, factor int) {
	for idx, row := range rows {
		length := len(row.Springs)
		for i := 2; i <= factor; i++ {
			row.Springs = append(row.Springs, '?')
			row.Springs = append(row.Springs, row.Springs[0:length]...)
		}

		length = len(row.Groups)
		for i := 2; i <= factor; i++ {
			row.Groups = append(row.Groups, row.Groups[0:length]...)
		}

		rows[idx] = row
	}
}

func PartA(path string) int {
	rows, err := parseInput(path)
	utils.CheckError(err)

	sum := 0
	for _, row := range rows {
		// var questions []int
		// for i, c := range row.Springs {
		// 	if c == '?' {
		// 		questions = append(questions, i)
		// 	}
		// }

		// matches := countSpringPermutations(row.Springs, row.Groups)
		// fmt.Println(matches)
		sum += countSpringPermutations(row.Springs, row.Groups)
	}

	return sum
}

func PartB(path string) int {
	rows, err := parseInput(path)
	utils.CheckError(err)

	expandRows(rows, 5)

	sum := 0
	for _, row := range rows {
		sum += countSpringPermutations(row.Springs, row.Groups)
	}

	return sum
}
