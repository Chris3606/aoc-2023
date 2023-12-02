package main

import (
	"aoc/day01"
	"aoc/day02"
	"aoc/day03"
	"aoc/day04"
	"aoc/day05"
	"aoc/day06"
	"aoc/day07"
	"aoc/day08"
	"aoc/day09"
	"aoc/day10"
	"aoc/day11"
	"aoc/day12"
	"aoc/day13"
	"aoc/day14"
	"aoc/day15"
	"aoc/day16"
	"aoc/day17"
	"aoc/day18"
	"aoc/day19"
	"aoc/day20"
	"aoc/day21"
	"aoc/day22"
	"aoc/day23"
	"aoc/day24"
	"aoc/day25"
	"flag"
	"fmt"
)

func printResult[T1 any, T2 any](day int, sample bool, partA T1, partB T2) {
	formatA := "Day %dA:"
	formatB := "Day %dB:"
	if sample {
		formatA = "Day %dA (sample):"
		formatB = "Day %dB (sample):"
	}
	fmt.Printf(formatA+"\n%v\n", day, partA)
	fmt.Printf("\n"+formatB+"\n%v\n", day, partB)
}

func runCode(day int, sample bool) {
	dayStr := fmt.Sprintf("%02d", day)
	file := "inputs/day" + dayStr
	if sample {
		file += "_sample"
	}
	file += ".txt"

	switch day {
	case 1:
		printResult(day, sample, day01.PartA(file), day01.PartB(file))
	case 2:
		printResult(day, sample, day02.PartA(file), day02.PartB(file))
	case 3:
		printResult(day, sample, day03.PartA(file), day03.PartB(file))
	case 4:
		printResult(day, sample, day04.PartA(file), day04.PartB(file))
	case 5:
		printResult(day, sample, day05.PartA(file), day05.PartB(file))
	case 6:
		printResult(day, sample, day06.PartA(file), day06.PartB(file))
	case 7:
		printResult(day, sample, day07.PartA(file), day07.PartB(file))
	case 8:
		printResult(day, sample, day08.PartA(file), day08.PartB(file))
	case 9:
		printResult(day, sample, day09.PartA(file), day09.PartB(file))
	case 10:
		printResult(day, sample, day10.PartA(file), day10.PartB(file))
	case 11:
		printResult(day, sample, day11.PartA(file), day11.PartB(file))
	case 12:
		printResult(day, sample, day12.PartA(file), day12.PartB(file))
	case 13:
		printResult(day, sample, day13.PartA(file), day13.PartB(file))
	case 14:
		printResult(day, sample, day14.PartA(file), day14.PartB(file))
	case 15:
		printResult(day, sample, day15.PartA(file), day15.PartB(file))
	case 16:
		printResult(day, sample, day16.PartA(file), day16.PartB(file))
	case 17:
		printResult(day, sample, day17.PartA(file), day17.PartB(file))
	case 18:
		printResult(day, sample, day18.PartA(file), day18.PartB(file))
	case 19:
		printResult(day, sample, day19.PartA(file), day19.PartB(file))
	case 20:
		printResult(day, sample, day20.PartA(file), day20.PartB(file))
	case 21:
		printResult(day, sample, day21.PartA(file), day21.PartB(file))
	case 22:
		printResult(day, sample, day22.PartA(file), day22.PartB(file))
	case 23:
		printResult(day, sample, day23.PartA(file), day23.PartB(file))
	case 24:
		printResult(day, sample, day24.PartA(file), day24.PartB(file))
	case 25:
		printResult(day, sample, day25.PartA(file), day25.PartB(file))
	default:
		panic("Unsupported day parameter.")
	}
}

func main() {
	dayPtr := flag.Int("d", 1, "The day to run.")
	samplePtr := flag.Bool("s", false, "When set, runs with the sample input instead of the real input.")

	flag.Parse()

	runCode(*dayPtr, *samplePtr)
}
