package day01

import (
	"aoc/utils"
	"bufio"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func parseInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func getDigits(s string) string {
	var digits strings.Builder
	for _, c := range s {
		if unicode.IsDigit(c) {
			digits.WriteRune(c)
		}
	}

	return digits.String()
}

func PartA(path string) int {
	lines, err := parseInput(path)
	utils.CheckError(err)

	sum := 0
	for _, line := range lines {
		var str strings.Builder
		digits := getDigits(line)
		str.WriteByte(digits[0])
		str.WriteByte(digits[len(digits)-1])

		number, err := strconv.Atoi(str.String())
		utils.CheckError(err)

		sum += number
	}

	return sum
}

func PartB(path string) string {
	_, err := parseInput(path)
	utils.CheckError(err)

	return "Not implemented"
}
