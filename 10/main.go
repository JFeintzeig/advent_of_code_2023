package main

import (
	"fmt"
	"os"
	"strings"
)

type Coord struct {
	i, j int
}

type NodeTraverser interface {
	before() (i, j int)
	after() (i, j int)
	getinsideba() []Coord
	getinsideab() []Coord
}

func GetNext(n *Node, i int, j int) (ni, nj int) {
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

func GetInside(n *Node, i int, j int) []Coord {
    iBefore, jBefore := n.before()
    iAfter, jAfter := n.after()
    if iBefore == i && jBefore == j {
        return n.getinsideba()
    } else if iAfter == i && jAfter == j {
        return n.getinsideab()
    } else {
        panic("problem")
        return make([]Coord, 0)
    }
}

type Node struct {
    NodeTraverser
	name string
}

type Pipe struct {
    Coord
}

func (n Pipe) before() (i, j int) {
	return n.i - 1, n.j
}

func (n Pipe) after() (i, j int) {
	return n.i + 1, n.j
}

func (n Pipe) getinsideba() []Coord {
    return []Coord{Coord{n.i, n.j-1}}
}

func (n Pipe) getinsideab() []Coord {
    return []Coord{Coord{n.i, n.j+1}}
}

type Dash struct {
	Coord
}

func (n Dash) before() (i, j int) {
	return n.i, n.j - 1
}

func (n Dash) after() (i, j int) {
	return n.i, n.j + 1
}

func (n Dash) getinsideba() []Coord {
    return []Coord{Coord{n.i+1, n.j}}
}

func (n Dash) getinsideab() []Coord {
    return []Coord{Coord{n.i-1, n.j}}
}

type L struct {
	Coord
}

func (n L) before() (i, j int) {
	return n.i - 1, n.j
}

func (n L) after() (i, j int) {
	return n.i, n.j + 1
}

func (n L) getinsideba() []Coord {
    return []Coord{Coord{n.i, n.j-1},Coord{n.i+1, n.j}}
}

func (n L) getinsideab() []Coord {
    return []Coord{Coord{n.i, n.j+1}}
}

type J struct {
    Coord
}

func (n J) before() (i, j int) {
	return n.i - 1, n.j
}

func (n J) after() (i, j int) {
	return n.i, n.j - 1
}

func (n J) getinsideba() []Coord {
    return []Coord{Coord{n.i, n.j-1}}
}

func (n J) getinsideab() []Coord {
    return []Coord{Coord{n.i+1, n.j}, Coord{n.i, n.j+1}}
}

type Seven struct {
	Coord
}

func (n Seven) before() (i, j int) {
	return n.i, n.j - 1
}

func (n Seven) after() (i, j int) {
	return n.i + 1, n.j
}

func (n Seven) getinsideba() []Coord {
    return []Coord{Coord{n.i, n.j-1}}
}

func (n Seven) getinsideab() []Coord {
    return []Coord{Coord{n.i, n.j+1}, Coord{n.i-1, n.j}}
}

type F struct {
	Coord
}

func (n F) before() (i, j int) {
	return n.i + 1, n.j
}

func (n F) after() (i, j int) {
	return n.i, n.j + 1
}

func (n F) getinsideba() []Coord {
    return []Coord{Coord{n.i, n.j+1}}
}

func (n F) getinsideab() []Coord {
    return []Coord{Coord{n.i, n.j-1}, Coord{n.i-1, n.j}}
}

type Period struct {
	Coord
}

func (n Period) before() (i, j int) {
	panic("no before() for Period")
	return
}

func (n Period) after() (i, j int) {
	panic("no after() for Period")
	return
}

func (n Period) getinsideba() []Coord {
    panic("no inside for Period")
    return make([]Coord, 0)
}

func (n Period) getinsideab() []Coord {
    panic("no inside for Period")
    return make([]Coord, 0)
}

func NewNode(i int, j int, s string) *Node {
    c := Coord{i, j}
	switch s {
	case "|":
		return &Node{Pipe{c}, s}
	case "-":
		return &Node{Dash{c}, s}
	case "L":
		return &Node{L{c}, s}
	case "J":
		return &Node{J{c}, s}
	case "7":
		return &Node{Seven{c}, s}
	case "F":
		return &Node{F{c}, s}
	case ".":
		return &Node{Period{c}, s}
	// NB: we replace this later, hardcoded + dumb
	case "S":
		return &Node{Dash{c}, s}
	default:
		panic("uh oh unknown char")
		return &Node{Pipe{c}, s}
	}
}

func LoadFile(fileName string) ([][]*Node, int, int) {
	file, _ := os.ReadFile(fileName)

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")

	var iStart, jStart int
	nodeMatrix := make([][]*Node, 0)
	for i, line := range lines {
		nodeRow := make([]*Node, 0)
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
	if fileName == "test3c.data" {
		nodeMatrix[iStart][jStart] = &Node{Seven{Coord{iStart, jStart}}, "S: 7"}
	} else if strings.Contains(fileName, "test") {
		nodeMatrix[iStart][jStart] = &Node{F{Coord{iStart, jStart}}, "S: F"}
	} else if fileName == "input.data" {
		nodeMatrix[iStart][jStart] = &Node{Dash{Coord{iStart, jStart}}, "S: -"}
	}

	return nodeMatrix, iStart, jStart
}

func TraverseLoop(nodeMatrix [][]*Node, iStart int, jStart int, iPrev int, jPrev int) (int, map[Coord]bool, map[Coord]bool) {
	i := iStart
	j := jStart
	nodePartOfLoop := make(map[Coord]bool)
    insides := make(map[Coord]bool)

	nSteps := 0
	for !(i == iStart && j == jStart && nSteps != 0) {
		node := nodeMatrix[i][j]
		iNew, jNew := GetNext(node, iPrev, jPrev)
        for _, val := range GetInside(node, iPrev, jPrev) {
            insides[val] = true
        }

		iPrev, jPrev = i, j
		i, j = iNew, jNew
		nSteps += 1

		nodePartOfLoop[Coord{i, j}] = true
	}
	return nSteps, nodePartOfLoop, insides
}

func RecursiveGetNeighbors(coords map[Coord]bool, loop map[Coord]bool) map[Coord]bool {
    neighbors := make(map[Coord]bool)
    for c, _ := range coords {
        neighbors[c] = true
        north := Coord{c.i-1, c.j}
        south := Coord{c.i+1, c.j}
        east := Coord{c.i, c.j+1}
        west := Coord{c.i, c.j-1}
        for _, n:= range []Coord{north, south, east, west} {
            _, neighborOnLoop := loop[n]
            if !neighborOnLoop {
                neighbors[n] = true
            } else {
            }
        }
    }
    if len(neighbors) == len(coords) {
        return neighbors
    } else {
        return RecursiveGetNeighbors(neighbors, loop)
    }
}

func main() {
	//nodeMatrix, iStart, jStart := LoadFile("test3c.data")
	nodeMatrix, iStart, jStart := LoadFile("input.data")
	fmt.Printf("Start: (%d, %d)\n", iStart, jStart)

	// hardcoding the "entry direction" for test data
	//nSteps1, nodePartOfLoop, insides := TraverseLoop(nodeMatrix, iStart, jStart, iStart, jStart-1)
	nSteps1, nodePartOfLoop, insides := TraverseLoop(nodeMatrix, iStart, jStart, iStart, jStart-1)
	fmt.Printf("Part 1 Direction 1: %d\n", nSteps1)
	fmt.Printf("Part 1 Answer: %d\n", nSteps1/2)

    fmt.Printf("Insides:\n")
    for c, _ := range insides {
        _, ok := nodePartOfLoop[c]
        if !ok {

            m := make(map[Coord]bool)
            m[c] = true
            neighbors := RecursiveGetNeighbors(m, nodePartOfLoop)
            for k, _ := range neighbors {
                insides[k] = true
            }
        } else {
            delete(insides, c)
        }
    }

	for i, row := range nodeMatrix {
		for j, _ := range row {
            c := Coord{i, j}
			_, ok := nodePartOfLoop[c]
			if ok {
				fmt.Printf("X")
			} else {
                _, in := insides[c]
                if in {
                    fmt.Printf("I")
                } else {
				fmt.Printf("O")
                }
			}
		}
		fmt.Printf("\n")
	}

    fmt.Printf("Part 2 Answer: %d\n", len(insides))
}
