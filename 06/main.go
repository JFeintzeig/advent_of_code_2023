package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func ParseRaces(lines []string, respectSpaces bool) ([]int, []int) {
	var time []int
	var distance []int
	for _, line := range lines {
		sp := strings.Split(line, ":")
		quant := sp[0]
		var input []string
		if respectSpaces {
			input = strings.Split(strings.Replace(strings.TrimSpace(sp[1]), "  ", " ", -1), " ")
		} else {
			input = strings.Split(strings.Replace(strings.TrimSpace(sp[1]), " ", "", -1), " ")
		}
		for _, v := range input {
			vInt, _ := strconv.Atoi(v)
			if vInt == 0 {
				continue
			}
			if quant == "Time" {
				time = append(time, vInt)
			} else if quant == "Distance" {
				distance = append(distance, vInt)
			}
		}
	}
	fmt.Printf("%v\n", time)
	fmt.Printf("%v\n", distance)
	return time, distance
}

func LoadFile(fileName string) []string {
	file, _ := os.ReadFile(fileName)

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	return lines
}

func main() {
	lines := LoadFile("input.data")

	part := 2

	var respectSpaces bool
	if part == 1 {
		respectSpaces = true
	} else if part == 2 {
		respectSpaces = false
	} else {
		panic("part must be 1 or 2")
	}

	time, distance := ParseRaces(lines, respectSpaces)

	// distance = (time - x_hold)*x_hold
	// distance = t*x - x^2
	// 0 = -x^2 + t*x - d
	// x = [ -b +/- sqrt(b2-4ac) ] / 2a
	// x = [ -t +/- sqrt(t^2 - 4d) / -2]
	// # of winning solutions = # of integers between x's?

	output := 1
	for i, ti := range time {
		d := float64(distance[i])
		t := float64(ti)
		x1 := (-t + math.Sqrt(math.Pow(t, 2)-4*d)) / -2
		x2 := (-t - math.Sqrt(math.Pow(t, 2)-4*d)) / -2
		fmt.Printf("%3.3f %3.3f\n", x1, x2)
		fmt.Printf("%d\n", int(x2)-int(x1))
		output *= (int(x2) - int(x1))
	}

	fmt.Printf("Part %d Answer: %d\n", part, output)

}
