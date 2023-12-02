package day14

import (
	"aoc/utils"
	"bufio"
	"os"
)

func parseInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		// Parse lines
	}

	return []string{"not", "implemented"}, nil
}

func PartA(path string) string {
	_, err := parseInput(path)
	utils.CheckError(err)

	return "Not implemented"
}

func PartB(path string) string {
	_, err := parseInput(path)
	utils.CheckError(err)

	return "Not implemented"
}
