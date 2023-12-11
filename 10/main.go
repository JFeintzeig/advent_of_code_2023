package main

import (
	"fmt"
	"os"
	"strings"
)

type Coord struct {
	i, j int
}

type Node interface {
	before() (i, j int)
	after() (i, j int)
	//getinsideba() (i, j int)
	//getinsideab() (i, j int)
}

func GetNext(n Node, i int, j int) (ni, nj int) {
	iBefore, jBefore := n.before()
	iAfter, jAfter := n.after()
	if iBefore == i && jBefore == j {
		return n.after()
	} else if iAfter == i && jAfter == j {
		return n.before()
	} else {
		msg := fmt.Sprintf("problem with GetNext(): %v %d %d", n, i, j)
		panic(msg)
		return
	}
}

//func GetInside(n Node, i int, j int) (ini, inj int) {
//    iBefore, jBefore := n.before()
//    iAfter, jAfter := n.after()
//    if iBefore == i && jBefore == j {
//        return n.getinsideba()
//    } else if iAfter == i && jAfter == j {
//        return n.getinsideab()
//    } else {
//        panic("problem")
//        return
//    }
//}

type Pipe struct {
	i, j int
	name string
}

func (n *Pipe) before() (i, j int) {
	return n.i - 1, n.j
}

func (n *Pipe) after() (i, j int) {
	return n.i + 1, n.j
}

type Dash struct {
	i, j int
	name string
}

func (n *Dash) before() (i, j int) {
	return n.i, n.j - 1
}

func (n *Dash) after() (i, j int) {
	return n.i, n.j + 1
}

type L struct {
	i, j int
	name string
}

func (n *L) before() (i, j int) {
	return n.i - 1, n.j
}

func (n *L) after() (i, j int) {
	return n.i, n.j + 1
}

type J struct {
	i, j int
	name string
}

func (n *J) before() (i, j int) {
	return n.i - 1, n.j
}

func (n *J) after() (i, j int) {
	return n.i, n.j - 1
}

type Seven struct {
	i, j int
	name string
}

func (n *Seven) before() (i, j int) {
	return n.i, n.j - 1
}

func (n *Seven) after() (i, j int) {
	return n.i + 1, n.j
}

type F struct {
	i, j int
	name string
}

func (n *F) before() (i, j int) {
	return n.i + 1, n.j
}

func (n *F) after() (i, j int) {
	return n.i, n.j + 1
}

type Period struct {
	i, j int
	name string
}

func (n *Period) before() (i, j int) {
	panic("no before() for Period")
	return
}

func (n *Period) after() (i, j int) {
	panic("no after() for Period")
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

func TraverseLoop(nodeMatrix [][]Node, iStart int, jStart int, iPrev int, jPrev int) (int, map[Coord]bool) {
	i := iStart
	j := jStart
	nodePartOfLoop := make(map[Coord]bool)

	nSteps := 0
	for !(i == iStart && j == jStart && nSteps != 0) {
		node := nodeMatrix[i][j]
		iNew, jNew := GetNext(node, iPrev, jPrev)
		iPrev, jPrev = i, j
		i, j = iNew, jNew
		nSteps += 1

		nodePartOfLoop[Coord{i, j}] = true
	}
	return nSteps, nodePartOfLoop
}

func main() {
	nodeMatrix, iStart, jStart := LoadFile("input.data")
	fmt.Printf("(%d, %d)\n", iStart, jStart)

	// hardcoding the "entry direction" for test data
	//nSteps1, nodePartOfLoop := TraverseLoop(nodeMatrix, iStart, jStart, iStart+1, jStart)
	nSteps1, nodePartOfLoop := TraverseLoop(nodeMatrix, iStart, jStart, iStart, jStart-1)
	fmt.Printf("Part 1 Direction 1: %d\n", nSteps1)
	//nSteps2, nodePartOfLoop := TraverseLoop(nodeMatrix, iStart, jStart, iStart, jStart+1)
	nSteps2, nodePartOfLoop := TraverseLoop(nodeMatrix, iStart, jStart, iStart, jStart+1)
	fmt.Printf("Part 1 Direction 2: %d\n", nSteps2)
	fmt.Printf("Part 1 Answer: %d\n", nSteps1/2)

	// part 2: first visualize nodes that are not in the loop
	// TODO: if i can implement getinsideab() and getinsideba()
	// for ** one ** direction around the loop (clockwise?)
	// AND figure out which direction represents clockwise
	// for the real input, can i then count pixels inside?
	// basic algo: if the x,y from getinsideab() (or ba(), whichever appropriate)
	// is a pixel that's not part of the loop, then increment count
	// how to deal with contiguous spaces of inside pixels? maybe
	// iteratively apply the transformation to get new x,y's
	// farther inside and test if they're not part of loop?
	for i, row := range nodeMatrix {
		for j, _ := range row {
			_, ok := nodePartOfLoop[Coord{i, j}]
			if ok {
				fmt.Printf("X")
			} else {
				fmt.Printf("O")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}
