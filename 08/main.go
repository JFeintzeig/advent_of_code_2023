package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Node struct {
	name      string
	leftName  string
	rightName string
	next      []*Node
}

func ParseLines(lines []string, re *regexp.Regexp) ([]int, *Node, []*Node) {
	// parse directions
	directions := strings.Split(strings.Replace(strings.Replace(strings.Replace(lines[0], "L", "0", -1), "R", "1", -1), " ", "", -1), "")
	var out []int
	for _, val := range directions {
		num, _ := strconv.Atoi(val)
		out = append(out, num)
	}

	// parse nodes in 2 passes:
	// 1: create Node() without left and right
	// 2: loop over them again and fill in for each
	nodeMap := make(map[string]*Node)
	nodeList := make([]*Node, 0)
	for i, line := range lines {
		if i == 0 {
			continue
		}
		matches := re.FindAllString(line, -1)
		if len(matches) < 3 {
			fmt.Printf("nothing to parse in line: %s\n", line)
			continue
		}
		name := matches[0]
		leftName := matches[1]
		rightName := matches[2]

		node := &Node{name: name, leftName: leftName, rightName: rightName}
		nodeMap[name] = node

		nodeList = append(nodeList, node)

	}

	headNode := nodeMap["AAA"]

	for _, v := range nodeMap {
		leftNode, okL := nodeMap[v.leftName]
		rightNode, okR := nodeMap[v.rightName]
		if !okL || !okR {
			fmt.Printf("Node %s doesnt have next(): %s:%t %s:%t\n", v.name, v.leftName, okL, v.rightName, okR)
			v.next = []*Node{&Node{name: "INVALID"}, &Node{name: "INVALID"}}
			continue
		}
		v.next = []*Node{leftNode, rightNode}
	}

	return out, headNode, nodeList
}

func LoadFile(fileName string) []string {
	file, _ := os.ReadFile(fileName)

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	return lines
}

func PrintStuff(directions []int, headNode *Node, nodeList []*Node) {
	fmt.Printf("%v\n", directions)
	for _, v := range directions {
		if v == 0 {
			fmt.Printf("L")
		} else {
			fmt.Printf("R")
		}
	}
	fmt.Printf("\n%v\n\n", *headNode)

	//for _, n := range nodeList {
	//  fmt.Printf("%s = (%s, %s)\n", n.name, n.next[0].name, n.next[1].name)
	//}

}

func GetANodes(nodeList []*Node) map[string]*Node {
	out := make(map[string]*Node)
	for _, n := range nodeList {
		if strings.HasSuffix(n.name, "A") {
			out[n.name] = n
		}
	}
	return out
}

// i just copied this GCD // LCM routine off the internet
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func main() {
	re, _ := regexp.Compile("\\b([A-Z]){3}\\b")

	lines := LoadFile("input.data")

	directions, headNode, nodeList := ParseLines(lines, re)

	// part 1
	if false {
		i := 0
		node := headNode
		for {
			if node.name == "ZZZ" {
				break
			}

			idx := i % len(directions)

			node = node.next[directions[idx]]
			i += 1
		}
		fmt.Printf("\nPart 1 Answer: %d\n", i)
	}

	// part 2
	aNodes := GetANodes(nodeList)
	counterByNode := make(map[string]int)
	cycles := make(map[string]int)
	for k, _ := range aNodes {
		fmt.Printf("%s\n", k)
		counterByNode[k] = 0
	}
	j := 0
outer:
	for {
		for k, v := range aNodes {
			if strings.HasSuffix(aNodes[k].name, "Z") {
				fmt.Printf("%s %d to %s\n", k, counterByNode[k], aNodes[k].name)
				_, ok := cycles[k]
				if !ok {
					cycles[k] = counterByNode[k]
				}
			}

			jdx := j % len(directions)
			aNodes[k] = v.next[directions[jdx]]
			counterByNode[k] += 1
		}
		j += 1

		if len(cycles) == len(aNodes) {
			break
		}

		// too slow for all of them to line up at once
		nZ := 0
		for k, _ := range aNodes {
			if strings.HasSuffix(aNodes[k].name, "Z") {
				nZ += 1
			}
			if nZ == len(aNodes) {
				break outer
			}
		}
	}

	fmt.Printf("%v", cycles)
	cycleList := make([]int, 0)
	for _, v := range cycles {
		cycleList = append(cycleList, v)
	}
	answer := LCM(cycleList[0], cycleList[1], cycleList...)
	fmt.Printf("\nPart 2 Answer: %d\n", answer)
}
