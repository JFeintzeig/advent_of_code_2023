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

var req = regexp.MustCompile("\\?{1}")
var reb = regexp.MustCompile("#{1}")

type record struct {
	pattern         []string
	groups          []int
}

func CountChar(pattern string, re *regexp.Regexp) ([]int, int) {
	result := re.FindAllStringIndex(pattern, -1)
	out := make([]int, 0)
	for _, r := range result {
		out = append(out, r[0])
	}
	return out, len(result)
}

func GetPossible(pattern []string, groups []int) int {
    unknownindex, nunknown := CountChar(strings.Join(pattern, ""), req)
    _, nbrokenexisting := CountChar(strings.Join(pattern, ""), reb)
    nbrokentotal := 0
    for _, g := range groups {
        nbrokentotal += g
    }

	// return number of strings that satisfy groups
	count := 0
	// for each possible string
	for i := uint(0); i < uint(math.Pow(2, float64(nunknown))); i++ {
		// bitmask j represents whether we're filling . (0) or # (1) in each part

		if bits.OnesCount(i) + nbrokenexisting != nbrokentotal {
			continue
		}

		newpattern := make([]string, len(pattern))
		copy(newpattern, pattern)
		j := i
		for _, ind := range unknownindex {
			if j&0x01 == 1 {
				newpattern[ind] = "#"
			} else {
				newpattern[ind] = "."
			}
			j = j >> 1
		}
		newstring := strings.Join(newpattern, "")
		if ValueMatchesGroups(newstring, groups) {
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

		groups := make([]int, 0)
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
		}
		records = append(records, record{pattern: splitpattern, groups: groups})
	}
	return records
}

func Find(r record, start int, groupNum int) int {
    if !(groupNum < len(r.groups) && start < len(r.pattern)) {
        return 1
    }
    group := r.groups[groupNum]
    fmt.Printf("group %d=%d, start %d\n", groupNum, group, start)
    nBadRow := 0
    newStart := 0
    newGroupNum := 0
    nPossible := 0
    // grow a subsegment of the pattern, keeping track
    // of the number of bad unknown springs in a row
    for j := start; j < len(r.pattern); j++ {
        if r.pattern[j] != "." {
            nBadRow += 1
        } else {
            nBadRow = 0
        }
        // if have enough in a row to match the group
        if nBadRow >= group {
            // make sure next spring is a possible end of group!
            if (j+1 < len(r.pattern) && r.pattern[j+1] != "#") {
                newStart = j + 2
                newGroupNum = groupNum + 1
                nPossible = GetPossible(r.pattern[start:newStart], []int{group})
                fmt.Printf("%v %d %d\n", r.pattern[start:j+2], nBadRow, nPossible)
                break
            }
            // or we're at the end of the pattern
            if j+1 == len(r.pattern) {
                newStart = j+2
                newGroupNum = groupNum + 1
                nPossible = GetPossible(r.pattern[start:], []int{group})
                fmt.Printf("END: %v %d %d\n", r.pattern[start:], nBadRow, nPossible)
                break
            }
        }
    }
    if !(newGroupNum < len(r.groups) && newStart < len(r.pattern)) {
        return nPossible * Find(r, newStart, newGroupNum)
    }
    return nPossible * Find(r, newStart, newGroupNum) + Find(r, newStart, groupNum)
}

func main() {
	records := LoadFile("input.data", 1)

	answer1 := 0
	for _, rec := range records {
		npos := GetPossible(rec.pattern, rec.groups)
		answer1 += npos
	}
	fmt.Printf("Part 1 Answer: %d\n", answer1)

    // part 2
	records = LoadFile("input.data", 1)

    r := records[4]

    fmt.Printf("%v\n", r)
    s := 0
    g := 0
    //np := 1
    npos := 1
    npos = Find(r, s, g)
    //for s < len(r.pattern) && g < len(r.groups) {
    //    np, s, g = Find(r, s, g)

    //    // failed
    //    if s < 0 || g < 0 {
    //        fmt.Printf("failed\n")
    //        break
    //    }
    //    npos *= np
    //    fmt.Printf("%d\n", np)
    //}
    fmt.Printf("%d\n", npos)
    fmt.Printf("true: %d\n", GetPossible(r.pattern, r.groups))
}
