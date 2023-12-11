package main

import (
	"fmt"
	"os"
	"strings"
)

type Node interface {
	next(i, j int) (ni, nj int)
}

type Pipe struct {
	i, j int
	name string
}

func (n *Pipe) next(i, j int) (ni, nj int) {
	if i == n.i-1 && j == n.j {
		return n.i + 1, n.j
	} else if i == n.i+1 && j == n.j {
		return n.i - 1, n.j
	} else {
		msg := fmt.Sprintf("(%d, %d) is not valid for %s at (%d, %d)", i, j, n.name, n.i, n.j)
		panic(msg)
		return
	}
}

type Dash struct {
	i, j int
	name string
}

func (n *Dash) next(i, j int) (ni, nj int) {
	if i == n.i && j == n.j-1 {
		return n.i, n.j + 1
	} else if i == n.i && j == n.j+1 {
		return n.i, n.j - 1
	} else {
		msg := fmt.Sprintf("(%d, %d) is not valid for %s at (%d, %d)", i, j, n.name, n.i, n.j)
		panic(msg)
		return
	}
}

type L struct {
	i, j int
	name string
}

func (n *L) next(i, j int) (ni, nj int) {
	if i == n.i-1 && j == n.j {
		return n.i, n.j + 1
	} else if i == n.i && j == n.j+1 {
		return n.i - 1, n.j
	} else {
		msg := fmt.Sprintf("(%d, %d) is not valid for %s at (%d, %d)", i, j, n.name, n.i, n.j)
		panic(msg)
		return
	}
}

type J struct {
	i, j int
	name string
}

func (n *J) next(i, j int) (ni, nj int) {
	if i == n.i-1 && j == n.j {
		return n.i, n.j - 1
	} else if i == n.i && j == n.j-1 {
		return n.i - 1, n.j
	} else {
		msg := fmt.Sprintf("(%d, %d) is not valid for %s at (%d, %d)", i, j, n.name, n.i, n.j)
		panic(msg)
		return
	}
}

type Seven struct {
	i, j int
	name string
}

func (n *Seven) next(i, j int) (ni, nj int) {
	if i == n.i && j == n.j-1 {
		return n.i + 1, n.j
	} else if i == n.i+1 && j == n.j {
		return n.i, n.j - 1
	} else {
		msg := fmt.Sprintf("(%d, %d) is not valid for %s at (%d, %d)", i, j, n.name, n.i, n.j)
		panic(msg)
		return
	}
}

type F struct {
	i, j int
	name string
}

func (n *F) next(i, j int) (ni, nj int) {
	if i == n.i+1 && j == n.j {
		return n.i, n.j + 1
	} else if i == n.i && j == n.j+1 {
		return n.i + 1, n.j
	} else {
		msg := fmt.Sprintf("(%d, %d) is not valid for %s at (%d, %d)", i, j, n.name, n.i, n.j)
		panic(msg)
		return
	}
}

type Period struct {
	i, j int
	name string
}

func (n *Period) next(i, j int) (ni, nj int) {
	msg := fmt.Sprintf("(%d, %d) is not valid for %s at (%d, %d)", i, j, n.name, n.i, n.j)
	panic(msg)
	return
}

func NewNode(i int, j int, s string) Node {
	switch s {
	case "|":
		return &Pipe{i, j, s}
	case "-":
		return &Dash{i, j, s}
	case "L":
		return &L{i, j, s}
	case "J":
		return &J{i, j, s}
	case "7":
		return &Seven{i, j, s}
	case "F":
		return &F{i, j, s}
	case ".":
		return &Period{i, j, s}
	// NB: we replace this later, hardcoded + dumb
	case "S":
		return &Dash{i, j, s}
	default:
		panic("uh oh unknown char")
		return &Pipe{i, j, s}
	}
}

func LoadFile(fileName string) ([][]Node, int, int) {
	file, _ := os.ReadFile(fileName)

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")

	var iStart, jStart int
	nodeMatrix := make([][]Node, 0)
	for i, line := range lines {
		nodeRow := make([]Node, 0)
		chars := strings.Split(line, "")
		for j, char := range chars {
			if char == "S" {
				iStart = i
				jStart = j
			}
			nodeRow = append(nodeRow, NewNode(i, j, char))
		}
		nodeMatrix = append(nodeMatrix, nodeRow)
	}

	// NB: replace start node, hardcoding values in my files
	if fileName == "test1a.data" || fileName == "test1b.data" {
		nodeMatrix[iStart][jStart] = &F{iStart, jStart, "S: F"}
	} else if fileName == "test2a.data" || fileName == "test2b.data" {
		nodeMatrix[iStart][jStart] = &F{iStart, jStart, "S: F"}
	} else if fileName == "input.data" {
		nodeMatrix[iStart][jStart] = &Dash{iStart, jStart, "S: -"}
	}

	return nodeMatrix, iStart, jStart
}

func TraverseLoop(nodeMatrix [][]Node, iStart int, jStart int, iPrev int, jPrev int) int {
	i := iStart
	j := jStart
	// hardcoding the "entry" position for test data
	nSteps := 0
	for !(i == iStart && j == jStart && nSteps != 0) {
		node := nodeMatrix[i][j]
		iNew, jNew := node.next(iPrev, jPrev)
		iPrev, jPrev = i, j
		i, j = iNew, jNew
		nSteps += 1
	}
	return nSteps
}

func main() {
	fmt.Printf("hello world\n")

	nodeMatrix, iStart, jStart := LoadFile("input.data")
	fmt.Printf("(%d, %d)\n", iStart, jStart)

	// hardcoding the "entry direction" for test data
	//nSteps1 := TraverseLoop(nodeMatrix, iStart, jStart, iStart+1, jStart)
	nSteps1 := TraverseLoop(nodeMatrix, iStart, jStart, iStart, jStart-1)
	fmt.Printf("Part 1 Direction 1: %d\n", nSteps1)
	//nSteps2 := TraverseLoop(nodeMatrix, iStart, jStart, iStart, jStart+1)
	nSteps2 := TraverseLoop(nodeMatrix, iStart, jStart, iStart, jStart+1)
	fmt.Printf("Part 1 Direction 2: %d\n", nSteps2)
	fmt.Printf("Part 1 Answer: %d\n", nSteps1/2)
}
