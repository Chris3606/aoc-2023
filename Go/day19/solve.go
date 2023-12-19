package day19

import (
	"aoc/utils"
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Cat int

const (
	XCat Cat = iota
	MCat
	ACat
	SCat
)

type Op int

const (
	OpGreater Op = iota
	OpLess
)

type Condition struct {
	Category  Cat
	Operation Op
	Value     int
}

func (cnd *Condition) appliesTo(part Part) bool {
	var val int
	switch cnd.Category {
	case XCat:
		val = part.X
	case MCat:
		val = part.M
	case ACat:
		val = part.A
	case SCat:
		val = part.S
	default:
		panic("Unsupported op")
	}

	switch cnd.Operation {
	case OpGreater:
		return val > cnd.Value
	case OpLess:
		return val < cnd.Value
	default:
		panic("Unsupported op")
	}
}

type Rule struct {
	Cnd Condition
	Dst string
}

type Workflow struct {
	ID          string
	Rules       []Rule
	DefaultRule string
}

type Part struct {
	X int
	M int
	A int
	S int
}

func parseRule(rule string) (Rule, error) {
	scanner := utils.NewStringDelimiterScanner(rule, ":")

	condition, err := utils.ReadStringFromScanner(scanner)
	if err != nil {
		return Rule{}, err
	}

	var cat Cat
	switch condition[0] {
	case 'x':
		cat = XCat
	case 'm':
		cat = MCat
	case 'a':
		cat = ACat
	case 's':
		cat = SCat
	default:
		return Rule{}, errors.New("invalid category")
	}

	var op Op
	switch condition[1] {
	case '<':
		op = OpLess
	case '>':
		op = OpGreater
	default:
		return Rule{}, errors.New("invalid op")
	}

	val, err := strconv.Atoi(condition[2:])
	if err != nil {
		return Rule{}, err
	}

	dst, err := utils.ReadStringFromScanner(scanner)
	if err != nil {
		return Rule{}, err
	}

	return Rule{Cnd: Condition{Category: cat, Operation: op, Value: val}, Dst: dst}, nil
}

func parseWorkflow(workflow string) (Workflow, error) {
	idx := strings.IndexRune(workflow, '{')
	if idx == -1 {
		return Workflow{}, errors.New("workflow ID not found")
	}

	id := workflow[:idx]

	data := workflow[idx+1 : len(workflow)-1]
	idx = strings.LastIndex(data, ",")

	// Parse out everything
	var rules []Rule
	var defaultRule string
	if idx == -1 {
		defaultRule = data
	} else {
		defaultRule = data[idx+1:]
		r, err := utils.ReadItems(utils.NewStringDelimiterScanner(data[:idx], ","), parseRule, false)
		if err != nil {
			return Workflow{}, err
		}

		rules = r
	}

	return Workflow{ID: id, Rules: rules, DefaultRule: defaultRule}, nil
}

func parsePartValue(partValueData string) (int, error) {
	return strconv.Atoi(partValueData[2:])
}

func parsePart(part string) (Part, error) {
	scanner := utils.NewStringDelimiterScanner(part[1:len(part)-1], ",")

	partVals, err := utils.ReadItems(scanner, parsePartValue, false)
	if err != nil {
		return Part{}, nil
	}

	if len(partVals) != 4 {
		return Part{}, errors.New("incorrect part value specification")
	}

	return Part{X: partVals[0], M: partVals[1], A: partVals[2], S: partVals[3]}, nil
}

func parseInput(path string) ([]Workflow, []Part, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(utils.ScanDoubleNewLines)

	wfData, err := utils.ReadStringFromScanner(scanner)
	if err != nil {
		return nil, nil, err
	}
	workflows, err := utils.ReadItemsFromString(wfData, bufio.ScanLines, parseWorkflow, false)
	if err != nil {
		return workflows, nil, err
	}

	partData, err := utils.ReadStringFromScanner(scanner)
	if err != nil {
		return workflows, nil, err
	}
	parts, err := utils.ReadItemsFromString(partData, bufio.ScanLines, parsePart, false)
	if err != nil {
		return workflows, parts, err
	}

	return workflows, parts, nil
}

func genWorkflowMap(workflows []Workflow) map[string]Workflow {
	wfMap := map[string]Workflow{}
	for _, w := range workflows {
		wfMap[w.ID] = w
	}

	return wfMap
}

type ValRange struct {
	X utils.Range
	M utils.Range
	A utils.Range
	S utils.Range
}

func countAccepted(ranges ValRange, workflows map[string]Workflow, curFlow string) int {
	if curFlow == "A" {
		prod := 1
		prod *= ranges.X.End - ranges.X.Start + 1
		prod *= ranges.M.End - ranges.M.Start + 1
		prod *= ranges.A.End - ranges.A.Start + 1
		prod *= ranges.S.End - ranges.S.Start + 1

		return prod
	} else if curFlow == "R" {
		return 0
	}

	var count int

	flow := workflows[curFlow]
	allValuesCovered := false
	for _, r := range flow.Rules {
		var cur utils.Range
		switch r.Cnd.Category {
		case XCat:
			cur = ranges.X
		case MCat:
			cur = ranges.M
		case ACat:
			cur = ranges.A
		case SCat:
			cur = ranges.S
		default:
			panic("Unsupported category")
		}

		var trueRange, falseRange utils.Range
		switch r.Cnd.Operation {
		case OpLess:
			trueRange = utils.Range{Start: cur.Start, End: r.Cnd.Value - 1}
			falseRange = utils.Range{Start: r.Cnd.Value, End: cur.End}
		case OpGreater:
			trueRange = utils.Range{Start: r.Cnd.Value + 1, End: cur.End}
			falseRange = utils.Range{Start: cur.Start, End: r.Cnd.Value}
		default:
			panic("Unsupported operation")
		}

		if trueRange.Start <= trueRange.End {
			rng := ranges
			switch r.Cnd.Category {
			case XCat:
				rng.X = trueRange
			case MCat:
				rng.M = trueRange
			case ACat:
				rng.A = trueRange
			case SCat:
				rng.S = trueRange
			}
			count += countAccepted(rng, workflows, r.Dst)
		}

		if falseRange.Start <= falseRange.End {
			rng := ranges
			switch r.Cnd.Category {
			case XCat:
				rng.X = falseRange
			case MCat:
				rng.M = falseRange
			case ACat:
				rng.A = falseRange
			case SCat:
				rng.S = falseRange
			}

			ranges = rng
		} else {
			allValuesCovered = true
			break
		}

	}

	if !allValuesCovered {
		count += countAccepted(ranges, workflows, flow.DefaultRule)
	}

	return count
}

func PartA(path string) int {
	workflowList, parts, err := parseInput(path)
	utils.CheckError(err)

	workflows := genWorkflowMap(workflowList)

	sum := 0
	for _, part := range parts {
		curFlow := "in"

		accepted := false
		for {
			ruleFound := false
			for _, rule := range workflows[curFlow].Rules {
				if rule.Cnd.appliesTo(part) {
					ruleFound = true
					curFlow = rule.Dst
					break
				}
			}

			if !ruleFound {
				curFlow = workflows[curFlow].DefaultRule
			}

			if curFlow == "R" {
				break
			} else if curFlow == "A" {
				accepted = true
				break
			}
		}

		if accepted {
			sum += part.X + part.M + part.A + part.S
		}
	}

	return sum
}

func PartB(path string) int {
	workflowList, _, err := parseInput(path)
	utils.CheckError(err)

	workflows := genWorkflowMap(workflowList)

	return countAccepted(ValRange{X: utils.Range{Start: 1, End: 4000}, M: utils.Range{Start: 1, End: 4000}, A: utils.Range{Start: 1, End: 4000}, S: utils.Range{Start: 1, End: 4000}}, workflows, "in")
}
