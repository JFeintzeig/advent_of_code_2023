package main

import (
	"fmt"
	"math"
	"math/bits"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type record struct {
	pattern         []string
	groups          []int
	nunknown        int
	unknownindex    []int
	nbrokentotal    int
	nbrokenexisting int
}

func CountChar(pattern string, re *regexp.Regexp) ([]int, int) {
	result := re.FindAllStringIndex(pattern, -1)
	out := make([]int, 0)
	for _, r := range result {
		out = append(out, r[0])
	}
	return out, len(result)
}

func (r *record) GetPossible() int {
	// return number of strings that satisfy groups
	count := 0
	// for each possible string
	for i := uint(0); i < uint(math.Pow(2, float64(r.nunknown))); i++ {
		// bitmask j represents whether we're filling . (0) or # (1) in each part

		if bits.OnesCount(i)+r.nbrokenexisting != r.nbrokentotal {
			continue
		}

		newpattern := make([]string, len(r.pattern))
		copy(newpattern, r.pattern)
		j := i
		for _, ind := range r.unknownindex {
			if j&0x01 == 1 {
				newpattern[ind] = "#"
			} else {
				newpattern[ind] = "."
			}
			j = j >> 1
		}
		newstring := strings.Join(newpattern, "")
		//fmt.Printf("%s %s, %v %b\n", r.pattern, newpattern, r.groups, i)
		if ValueMatchesGroups(newstring, r.groups) {
			count += 1
		}
	}
	return count
}

func ValueMatchesGroups(val string, groups []int) bool {
	count := 0
	startedgroup := false
	groupNum := 0
	for i, c := range strings.Split(val, "") {
		if c == "#" {
			startedgroup = true
			count += 1
		} else {
			if startedgroup {
				if groupNum > len(groups)-1 {
					return false
				}
				if count != groups[groupNum] {
					return false
				}
				count = 0
				startedgroup = false
				groupNum += 1
			}
		}
		if i == len(val)-1 && startedgroup {
			if groupNum > len(groups)-1 {
				return false
			}
			if count != groups[groupNum] {
				return false
			}
			groupNum += 1
		}
	}
	if len(groups) != groupNum {
		return false
	}
	return true
}

func LoadFile(fileName string, part int) []record {
	req, _ := regexp.Compile("\\?{1}")
	reb, _ := regexp.Compile("#{1}")

	file, _ := os.ReadFile(fileName)

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	records := make([]record, 0)
	for _, line := range lines {
		parts := strings.Split(line, " ")
		pattern := parts[0]

		if part == 2 {
			// gotta do this x 5
			pattern = strings.Join([]string{pattern, pattern, pattern, pattern, pattern}, "?")
		}

		splitpattern := strings.Split(pattern, "")

		unknownindex, nunknown := CountChar(pattern, req)
		_, nbrokenexisting := CountChar(pattern, reb)
		groups := make([]int, 0)
		nbrokentotal := 0
		constraints := strings.Split(parts[1], ",")

		if part == 2 {
			// gotta do this x 5
			dummy := append(constraints, constraints...)
			dummy = append(dummy, constraints...)
			dummy = append(dummy, constraints...)
			constraints = append(dummy, constraints...)
		}

		for _, val := range constraints {
			num, _ := strconv.Atoi(val)
			groups = append(groups, num)
			nbrokentotal += num
		}
		records = append(records, record{pattern: splitpattern, groups: groups, nunknown: nunknown, unknownindex: unknownindex, nbrokentotal: nbrokentotal, nbrokenexisting: nbrokenexisting})
	}
	return records
}

func main() {
	records := LoadFile("input.data", 1)

	answer1 := 0
	maxLen := 0
	for _, rec := range records {
		fmt.Printf("%v ", rec)
		npos := rec.GetPossible()
		answer1 += npos
		fmt.Printf("%d\n", npos)
		if len(rec.pattern) > maxLen {
			maxLen = rec.nunknown
		}
	}
	fmt.Printf("Part 1 Answer: %d\n", answer1)
	fmt.Printf("%d\n", maxLen)
}
