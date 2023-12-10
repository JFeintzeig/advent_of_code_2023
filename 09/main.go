package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type sequence []int

func TakeDerivative(seq sequence) sequence {
	out := make(sequence, 0)
	for i := 1; i < len(seq); i++ {
		out = append(out, seq[i]-seq[i-1])
	}
	return out
}

func IsAllZeros(seq sequence) bool {
	for _, val := range seq {
		if val != 0 {
			return false
		}
	}
	return true
}

func SumSeq(seq sequence) int {
	out := 0
	for _, val := range seq {
		out += val
	}
	return out
}

func LoadFile(fileName string) []sequence {
	file, _ := os.ReadFile(fileName)

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")

	container := make([]sequence, 0)
	for _, line := range lines {
		sequence := make(sequence, 0)
		numStrings := strings.Split(line, " ")
		for _, numS := range numStrings {
			num, _ := strconv.Atoi(numS)
			sequence = append(sequence, num)
		}
		container = append(container, sequence)
	}
	return container
}

func main() {
	seqs := LoadFile("input.data")

	answer1 := 0
	answer2 := 0
	for _, s := range seqs {
		fmt.Printf("*** Seq: %v ***\n", s)
		currentSeq := s
		lastVals := make(sequence, 0)
		firstVals := make(sequence, 0)
		for !IsAllZeros(currentSeq) {
			fmt.Printf("%v\n", currentSeq)
			lastVals = append(lastVals, currentSeq[len(currentSeq)-1])
			firstVals = append(firstVals, currentSeq[0])
			currentSeq = TakeDerivative(currentSeq)
		}
		// part 1
		answer1 += SumSeq(lastVals)

		// part 2
		prev := 0
		fmt.Printf("%v\n", firstVals)
		for i := len(firstVals) - 1; i >= 0; i-- {
			prev = firstVals[i] - prev
			fmt.Printf("%d\n", prev)
		}
		answer2 += prev
	}
	fmt.Printf("Part 1 Answer: %d\n", answer1)
	fmt.Printf("Part 2 Answer: %d\n", answer2)
}
