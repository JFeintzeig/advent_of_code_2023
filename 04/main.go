package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type set map[int]bool

type Card struct {
	num     int
	winning set
	mine    set
}

func (c *Card) FindOverlap() []int {
	var out []int
	for k, _ := range c.mine {
		_, ok := c.winning[k]
		if ok {
			out = append(out, k)
		}
	}
	return out
}

func EvalCard(cards map[int]*Card, nCards int, toEval []int) int {
	var newToEval []int
	for _, num := range toEval {
		card := cards[num]
		overlap := card.FindOverlap()

		for i := 1; i <= len(overlap); i++ {
			newToEval = append(newToEval, card.num+i)
		}
	}

	nCards += len(toEval)
	if len(newToEval) > 0 {
		return EvalCard(cards, nCards, newToEval)
	} else {
		return nCards
	}
}

func StringToIntMap(str string) set {
	out := make(set)
	for _, s := range strings.Split(strings.Replace(strings.TrimSpace(str), "  ", " ", -1), " ") {
		trimmed := strings.TrimSpace(s)
		myint, _ := strconv.Atoi(trimmed)
		out[myint] = true
	}
	return out
}

func CardTo2Sets(line string) (int, []set) {
	cardNum, _ := strconv.Atoi(strings.TrimSpace(strings.Replace(strings.Split(line, ":")[0], "Card", "", -1)))
	pair := strings.Split(strings.Split(line, ":")[1], "|")
	winning := StringToIntMap(pair[0])
	mine := StringToIntMap(pair[1])
	return cardNum, []set{winning, mine}
}

func LoadFile(fileName string) []string {
	file, _ := os.ReadFile(fileName)

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	return lines
}

func main() {
	lines := LoadFile("input.data")

	// part 1
	var answer1 int
	cards := make(map[int]*Card)
	for _, line := range lines {
		cardNum, numSets := CardTo2Sets(line)
		cards[cardNum] = &Card{cardNum, numSets[0], numSets[1]}
		overlap := cards[cardNum].FindOverlap()
		//fmt.Printf("%v\n", overlap)
		if len(overlap) > 0 {
			answer1 += 0x1 << (len(overlap) - 1)
		}
	}

	fmt.Printf("Part 1: %d\n", answer1)

	// part 2
	start := time.Now()
	var total int
	for cardNum, _ := range cards {
		total += EvalCard(cards, 0, []int{cardNum})
	}

	fmt.Printf("Part 2: %d\n", total)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Time: %3.2f\n", float64(elapsed)/1e9)
}
