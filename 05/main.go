package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Transition struct {
	low    int64
	high   int64
	offset int64
}

type TransitionMap struct {
	transitions []Transition
	outputName  string
}

func NewTransitionMap(outputName string) *TransitionMap {
	t := make([]Transition, 0)
	return &TransitionMap{t, outputName}
}

type TransitionContainer map[string]*TransitionMap

func ParseSeeds(line string) []int64 {
	var output []int64
	for _, item := range strings.Split(strings.TrimSpace(strings.Replace(line, "seeds: ", "", -1)), " ") {
		myint, _ := strconv.Atoi(item)
		output = append(output, int64(myint))
	}
	return output
}

func ParseLines(lines []string) {
}

func LoadFile(fileName string) []string {
	file, _ := os.ReadFile(fileName)

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	return lines
}

func LinesToTransitionContainer(lines []string) *TransitionContainer {
	var tc TransitionContainer = make(TransitionContainer)

	var madeMap bool = false
	var currentMap *TransitionMap
	for _, line := range lines {
		if line == "\n" {
			madeMap = false
			continue
		}
		split := strings.Split(line, " ")
		if len(split) < 2 {
			continue
		}
		if split[1] == "map:" {
			inputName := strings.Split(split[0], "-")[0]
			outputName := strings.Split(split[0], "-")[2]
			currentMap = NewTransitionMap(outputName)
			tc[inputName] = currentMap
			madeMap = true
		} else if madeMap {
			destStart, _ := strconv.Atoi(split[0])
			sourceStart, _ := strconv.Atoi(split[1])
			rangeL, _ := strconv.Atoi(split[2])
			tr := Transition{low: int64(sourceStart), high: int64(sourceStart + rangeL), offset: int64(destStart - sourceStart)}
			(*currentMap).transitions = append(currentMap.transitions, tr)
		}
	}
	return &tc
}

func main() {
	lines := LoadFile("input.data")

	seeds := ParseSeeds(lines[0])

	tc := LinesToTransitionContainer(lines)
	for k, v := range *tc {
		fmt.Printf("%s %v\n\n\n", k, v)
	}

	output := make(map[int64]int64)
	for _, val := range seeds {
		output[val] = val
	}

	var inputName string = "seed"
	for {
		// we done
		if inputName == "location" {
			break
		}
		fmt.Printf("input: %s\n", inputName)

		// get transition map
		tm, ok := (*tc)[inputName]
		if !ok {
			panic("uh oh")
		}

		// for each thing
		for seed, thing := range output {
			// if no transition, fill with default
			var offset int64 = 0

			// for each transition
			for _, tr := range tm.transitions {
				// if transition matches, apply and go to next thing
				if thing >= tr.low && thing < tr.high {
					offset = tr.offset
					break
				}
			}

			// we simply overwrite value corresponding to seed
			output[seed] = thing + offset
		}

		// ready for next iteration
		inputName = tm.outputName
	}

	fmt.Printf("%v\n", output)

	const MaxUint = ^uint64(0)
	const MaxInt = int64(MaxUint >> 1)
	minLoc := MaxInt
	for _, v := range output {
		if v < minLoc {
			minLoc = v
		}
	}

	fmt.Printf("%d\n", minLoc)
}
