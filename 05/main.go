package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type ThingRange struct {
	low  int64
	high int64
}

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

func ParseThingRanges(line string) []*ThingRange {
	var output []*ThingRange
	listOfSeeds := strings.Split(strings.TrimSpace(strings.Replace(line, "seeds: ", "", -1)), " ")
	for i := 0; i < len(listOfSeeds); i = i + 2 {
		myint, _ := strconv.Atoi(listOfSeeds[i])
		myrange, _ := strconv.Atoi(listOfSeeds[i+1])
		output = append(output, &ThingRange{int64(myint), int64(myint + myrange - 1)})
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

	var seeds []int64
	seeds = ParseSeeds(lines[0])

	tc := LinesToTransitionContainer(lines)

	// part 1
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

	const MaxUint = ^uint64(0)
	const MaxInt = int64(MaxUint >> 1)
	minLoc := MaxInt
	for _, v := range output {
		if v < minLoc {
			minLoc = v
		}
	}

	fmt.Printf("Part 1 Answer: %d\n\n\n\n", minLoc)

	// part 2
	for k, v := range *tc {
		fmt.Printf("%s %v\n\n\n", k, v)
	}

	seedRanges := ParseThingRanges(lines[0])
	output2 := make(map[string][]*ThingRange, 0)
	output2["seed"] = seedRanges

	var inputName2 string = "seed"
	for {
		// we done
		if inputName2 == "location" {
			break
		}
		fmt.Printf("input: %s\n", inputName2)

		// get transition map
		tm, ok := (*tc)[inputName2]
		if !ok {
			panic("uh oh")
		}

		// for each thing
		out := make([]*ThingRange, 0)
		for _, tRange := range output2[inputName2] {

			// for each transition, look for overlap with our ThingRange
			// as this proceeds, we (a) append out with the new ThingRange
			// if a transition affects the original, and (b) mutate the
			// original so only the remaining non-transformed portion
			// is evaluated for the other transitions
			// we do this by setting low/high to -1 and filter it out
			// at the end. super janky.
			for _, tr := range tm.transitions {
				if tRange.low >= tr.low && tRange.high <= tr.high {
					// complete overlap: full change
					fmt.Printf("full overlap: %v %v\n", tRange, tr)
					out = append(out, &ThingRange{low: tRange.low + tr.offset, high: tRange.high + tr.offset})
					tRange.low = -1
					tRange.high = -1
					break
				} else if tRange.high < tr.low {
					// if thing range is entirely below low end of
					// transformation then no change
					fmt.Printf("below: %v %v\n", tRange, tr)
				} else if tRange.low > tr.high {
					// if thing range is entirely above high end of
					// transformation then no change
					fmt.Printf("above: %v %v\n", tRange, tr)
				} else if tRange.low < tr.low && tRange.high >= tr.low {
					// partial overlap low
					fmt.Printf("partial overlap low: %v %v\n", tRange, tr)
					out = append(out, &ThingRange{low: tr.low + tr.offset, high: tRange.high + tr.offset})
					tRange.high = tr.low - 1
					if tRange.low == tRange.high {
						tRange.low = -1
						tRange.high = -1
						break
					}
				} else if tRange.low >= tr.low && tRange.high > tr.high {
					// partial overlap high
					fmt.Printf("partial overlap high: %v %v\n", tRange, tr)
					out = append(out, &ThingRange{low: tRange.low + tr.offset, high: tr.high + tr.offset - 1})
					tRange.low = tr.high
					if tRange.low == tRange.high {
						tRange.low = -1
						tRange.high = -1
						break
					}
				}
			}
			if tRange.low != -1 && tRange.high != -1 {
				out = append(out, tRange)
			}
		}
		output2[tm.outputName] = out
		fmt.Printf("%s: ", tm.outputName)
		for _, val := range out {
			fmt.Printf("%v ", *val)
		}
		fmt.Printf("\n")

		// ready for next iteration
		inputName2 = tm.outputName
	}

	for k, v := range output2 {
		fmt.Printf("%s: ", k)
		for _, val := range v {
			fmt.Printf("%v ", *val)
		}
		fmt.Printf("\n")
	}

	answer := make([]int64, 0)
	for _, tRange := range output2["location"] {
		answer = append(answer, tRange.low)
	}
	fmt.Printf("Part 2 Answer: %d\n", slices.Min(answer))
}
