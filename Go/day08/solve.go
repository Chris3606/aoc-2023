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

func followPath(nodes map[string]NodeData, curNode string, dir byte) string {
	if dir == 'L' {
		curNode = nodes[curNode].Left
	} else {
		curNode = nodes[curNode].Right
	}

	return curNode
}

func PartA(path string) int {
	directions, nodes, err := parseInput(path)
	utils.CheckError(err)

	curNode := "AAA"
	curDirIdx := 0
	steps := 0

	for curNode != "ZZZ" {
		curDirection := directions[curDirIdx]
		curNode = followPath(nodes, curNode, curDirection)

		steps++
		curDirIdx = (curDirIdx + 1) % len(directions)
	}

	return steps
}

func PartB(path string) int {
	directions, nodes, err := parseInput(path)
	utils.CheckError(err)

	var paths []int
	for k := range nodes {
		if k[2] != 'A' {
			continue
		}

		curNode := k
		curDirIdx := 0
		steps := 0
		for curNode[2] != 'Z' {
			curDirection := directions[curDirIdx]
			curNode = followPath(nodes, curNode, curDirection)

			steps++
			curDirIdx = (curDirIdx + 1) % len(directions)
		}

		paths = append(paths, steps)
	}

	return utils.LCM(paths[0], paths[1], paths[2:]...)
}
