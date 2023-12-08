package day08

import (
	"aoc/utils"
	"bufio"
	"os"
)

type NodeData struct {
	Left  string
	Right string
}

func parseInput(path string) (string, map[string]NodeData, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	directions, err := utils.ReadStringFromScanner(scanner)
	if err != nil {
		return "", nil, err
	}

	// Skip blank line
	_, err = utils.ReadStringFromScanner(scanner)
	if err != nil {
		return "", nil, err
	}

	nodes := map[string]NodeData{}
	for scanner.Scan() {
		text := scanner.Text()
		nodes[text[:3]] = NodeData{text[7:10], text[12:15]}
	}

	return directions, nodes, nil
}

func followDirection(nodes map[string]NodeData, curNode string, dir byte) string {
	if dir == 'L' {
		curNode = nodes[curNode].Left
	} else {
		curNode = nodes[curNode].Right
	}

	return curNode
}

func getStepsForPath(nodes map[string]NodeData, startNode string, isEndNode func(string) bool, directions string) int {
	curNode := startNode
	curDirIdx := 0
	steps := 0

	for !isEndNode(curNode) {
		curDirection := directions[curDirIdx]
		curNode = followDirection(nodes, curNode, curDirection)

		steps++
		curDirIdx = (curDirIdx + 1) % len(directions)
	}

	return steps
}

func PartA(path string) int {
	directions, nodes, err := parseInput(path)
	utils.CheckError(err)

	return getStepsForPath(nodes, "AAA", func(s string) bool { return s == "ZZZ" }, directions)
}

func PartB(path string) int {
	directions, nodes, err := parseInput(path)
	utils.CheckError(err)

	var paths []int
	for k := range nodes {
		if k[2] != 'A' {
			continue
		}

		paths = append(paths, getStepsForPath(nodes, k, func(s string) bool { return s[2] == 'Z' }, directions))
	}

	return utils.LCM(paths[0], paths[1], paths[2:]...)
}
