package day15

import (
	"aoc/utils"
	"bufio"
	"os"
	"slices"
)

type OpType int

const (
	SubOp = iota
	EqOp
)

type StartupElement struct {
	Label    string
	Op       OpType
	FocalLen int
}

type Lens struct {
	Label    string
	FocalLen int
}

func (s *StartupElement) getBox() int {
	return hash(s.Label)
}

func hash(s string) int {
	bytes := []byte(s)

	h := 0
	for _, b := range bytes {
		h += int(b)
		h *= 17
		h %= 256
	}

	return h
}

func parseInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(utils.ScanDelimiterFunc(","))

	return utils.ReadItems(scanner, func(s string) (string, error) { return s, nil }, false)
}

func parseInput2(sequence []string) []StartupElement {
	var seq []StartupElement
	for _, s := range sequence {
		ch := s[len(s)-1]
		if ch == '-' {
			seq = append(seq, StartupElement{Label: s[:len(s)-1], Op: SubOp, FocalLen: -1})
		} else {
			fLen := ch - '0'
			seq = append(seq, StartupElement{Label: s[:len(s)-2], Op: EqOp, FocalLen: int(fLen)})
		}
	}

	return seq
}

func PartA(path string) int {
	startup_seq, err := parseInput(path)
	utils.CheckError(err)

	sum := 0
	for _, s := range startup_seq {
		sum += hash(s)
	}

	return sum
}

func PartB(path string) int {
	startup_seq1, err := parseInput(path)
	utils.CheckError(err)

	startup_seq := parseInput2(startup_seq1)

	var boxes [][]Lens
	for i := 0; i < 256; i++ {
		boxes = append(boxes, nil)
	}

	for _, seq := range startup_seq {
		box := seq.getBox()

		switch seq.Op {
		case SubOp:
			x := slices.IndexFunc(boxes[box], func(l Lens) bool { return l.Label == seq.Label })
			if x != -1 {
				boxes[box] = utils.RemoveFromSlice(boxes[box], x)
			}

		case EqOp:
			lens := Lens{Label: seq.Label, FocalLen: seq.FocalLen}
			x := slices.IndexFunc(boxes[box], func(l Lens) bool { return l.Label == seq.Label })
			if x == -1 {
				boxes[box] = append(boxes[box], lens)
			} else {
				boxes[box][x] = lens
			}
		default:
			panic("Invalid operation.")
		}
	}

	sum := 0
	for box, boxData := range boxes {
		for i, lens := range boxData {
			focusingPower := (1 + box) * (1 + i) * lens.FocalLen
			sum += focusingPower
		}
	}

	return sum
}
