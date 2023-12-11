package main

import (
	"fmt"
	"os"
	"strings"
)

// coord of a galaxy
type Coord struct {
	i, j int
}

func ExpandRows(matrix [][]string) ([][]string, []int) {
	out := make([][]string, 0)
	indices := make([]int, 0)
	// expand rows
	for i, row := range matrix {
		rowVals := make(map[string]int)
		for _, val := range row {
			rowVals[val] += 1
		}
		out = append(out, row)
		nDots, _ := rowVals["."]
		if len(row) == nDots {
			out = append(out, row)
			indices = append(indices, i)
		}
	}
	return out, indices
}

func Transpose(matrix [][]string) [][]string {
	nRows := len(matrix)
	nCols := len(matrix[0])
	out := make([][]string, nCols)
	for i := range out {
		out[i] = make([]string, nRows)
	}

	for i := 0; i < nRows; i++ {
		for j := 0; j < nCols; j++ {
			out[j][i] = matrix[i][j]
		}
	}
	return out
}

func Expand(matrix [][]string) ([][]string, []int, []int) {
	m1, rowIndices := ExpandRows(matrix)
	m2 := Transpose(m1)
	m3, colIndices := ExpandRows(m2)
	m4 := Transpose(m3)
	return m4, rowIndices, colIndices
}

func LoadFile(fileName string, part int) []Coord {
	file, _ := os.ReadFile(fileName)

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")

	matrix := make([][]string, 0)
	for _, line := range lines {
		row := make([]string, 0)
		for _, char := range strings.Split(line, "") {
			row = append(row, char)
		}
		matrix = append(matrix, row)
	}

	expandedMatrix, rowIndices, colIndices := Expand(matrix)
	out := make([]Coord, 0)

	if part == 1 {
		for i, row := range expandedMatrix {
			for j, char := range row {
				if char == "#" {
					out = append(out, Coord{i, j})
				}
			}
		}
	} else {
		for i, row := range matrix {
			for j, char := range row {
				coordI := 0
				coordJ := 0
				for _, rowVal := range rowIndices {
					if rowVal < i {
						coordI += 1000000 - 1
					}
				}
				for _, colVal := range colIndices {
					if colVal < j {
						coordJ += 1000000 - 1
					}
				}
				if char == "#" {
					out = append(out, Coord{i + coordI, j + coordJ})
				}
			}
		}
	}
	return out
}

func Abs(x int) int {
	if x > 0 {
		return x
	} else {
		return -1 * x
	}
}

func CalcDistances(coords []Coord) int {
	answer := 0
	for i, c1 := range coords {
		for j, c2 := range coords {
			if i != j {
				answer += Abs(c2.i-c1.i) + Abs(c2.j-c1.j)
			}
		}
	}
	return answer / 2
}

func main() {
	filename := "input.data"

	coords1 := LoadFile(filename, 1)
	fmt.Printf("%v\n", coords1)
	fmt.Printf("Part 1 Answer: %d\n", CalcDistances(coords1))

	coords2 := LoadFile(filename, 2)
	fmt.Printf("%v\n", coords2)
	fmt.Printf("Part 2 Answer: %d\n", CalcDistances(coords2))
}
