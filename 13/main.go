package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type pattern []string

func FindReflection(p pattern) (int, int) {
	score1 := 0
	n1 := 0

	score2 := 0
	n2 := 0
	for i, _ := range p {
		pure, offBy := IsReflection(p, i)
		if pure {
			score1 += 100 * i
			n1 += 1
		}
		if offBy == 1 {
			score2 += 100 * i
			n2 += 1
		}
	}
	transposedLines := Transpose(p)
	for i, _ := range transposedLines {
		pure, offBy := IsReflection(transposedLines, i)
		if pure {
			score1 += i
			n1 += 1
		}
		if offBy == 1 {
			score2 += i
			n2 += 1
		}
	}
	if score1 == 0 {
		panic("no p1 reflectino found\n")
	}
	if score2 == 0 {
		panic("no p2 reflectino found\n")
	}
	if n1 > 1 {
		fmt.Printf("%d solutions to part 1 found\n", n1)
	}
	if n2 > 1 {
		fmt.Printf("%d solutions to part 2 found\n", n2)
	}
	return score1, score2
}

func Transpose(p pattern) []string {
	nRows := len(p)
	nCols := len(p[0])
	temp := make([][]string, nCols)
	out := make([]string, nCols)
	for i := range temp {
		temp[i] = make([]string, nRows)
	}

	for i := 0; i < nRows; i++ {
		row := strings.Split(p[i], "")
		for j := 0; j < nCols; j++ {
			temp[j][i] = row[j]
		}
	}

	for i, _ := range temp {
		out[i] = strings.Join(temp[i], "")
	}

	return out
}
func IsReflection(lines []string, idx int) (bool, int) {
	if idx == 0 || idx == len(lines) {
		return false, -1
	}

	isPureReflection := true
	offBy := 0

	maxI := min(idx+1, len(lines)-idx)
	for i := 0; i < maxI; i++ {
		c1 := idx - i - 1
		c2 := idx + i
		if (c1 < 0) || c2 > len(lines)-1 {
			break
		}
		l1 := lines[c1]
		l2 := lines[c2]
		for i, _ := range l1 {
			if l1[i] != l2[i] {
				offBy += 1
			}
		}
		if !reflect.DeepEqual(l1, l2) {
			isPureReflection = false
		}
	}
	return isPureReflection, offBy
}

func LoadFile(fileName string) []pattern {
	file, _ := os.ReadFile(fileName)

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")

	patterns := make([]pattern, 0)
	plines := make([]string, 0)
	for _, l := range lines {
		if l == "" {
			patterns = append(patterns, plines)
			plines = make([]string, 0)
		} else {
			plines = append(plines, l)
		}
	}
	patterns = append(patterns, plines)

	return patterns
}

func MakeSmudge(p pattern, i int, j int) pattern {
	return make(pattern, 0)
}

func main() {
	patterns := LoadFile("input.data")

	answer1 := 0
	answer2 := 0
	for _, p := range patterns {
		s1, s2 := FindReflection(p)
		answer1 += s1
		answer2 += s2
	}
	fmt.Printf("Part 1 Answer: %d\n", answer1)
	fmt.Printf("Part 2 Answer: %d\n", answer2)
}
