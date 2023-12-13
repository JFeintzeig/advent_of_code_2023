package main

import (
    "fmt"
    "math"
    "os"
    "reflect"
    "regexp"
    "strings"
    "strconv"
)

type record struct {
    pattern []string
    groups []int
    nunknown int
    unknownindex []int
}

func GetUnknown(pattern string, re *regexp.Regexp) ([]int, int) {
    result := re.FindAllStringIndex(pattern,-1)
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
    for i := 0; i < int(math.Pow(2, float64(r.nunknown))); i++ {
        // i will represent a bitmask of whether we're filling . (0) or # (1) in each part
        newpattern := make([]string, len(r.pattern))
        copy(newpattern, r.pattern)
        j := i
        for _, ind := range r.unknownindex {
            if j & 0x01 == 1 {
                newpattern[ind] = "#"
            } else {
                newpattern[ind] = "."
            }
            j = j >> 1
        }
        newstring := strings.Join(newpattern,"")
        //fmt.Printf("%s %s, %v %b\n", r.pattern, newpattern, r.groups, i)
        if ValueMatchesGroups(newstring, r.groups) {
            count += 1
        }
    }
    return count
}

func ValueMatchesGroups(val string, groups []int) bool {
    foundgroups := make([]int,0)

    count := 0
    startedgroup := false
    for i, c := range strings.Split(val, "") {
        if c == "#" {
            startedgroup = true
            count += 1
        } else {
            if startedgroup {
                foundgroups = append(foundgroups, count)
                count = 0
                startedgroup = false
            }
        }
        if i == len(val)-1 && startedgroup {
            foundgroups = append(foundgroups, count)
        }
    }
    return reflect.DeepEqual(foundgroups, groups)
}

func LoadFile(fileName string) []record {
    re, _ := regexp.Compile("\\?{1}")

	file, _ := os.ReadFile(fileName)

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
    records := make([]record, 0)
    for _, line := range lines {
        parts := strings.Split(line, " ")
        pattern := strings.Split(parts[0], "")
        unknownindex, nunknown := GetUnknown(parts[0], re)
        groups := make([]int, 0)
        for _, val := range strings.Split(parts[1],",") {
            num, _ := strconv.Atoi(val)
            groups = append(groups, num)
        }
        records = append(records, record{pattern:pattern, groups: groups, nunknown: nunknown, unknownindex: unknownindex})
    }
    return records
}

func main() {
    records := LoadFile("input.data")

    answer1 := 0
    for _, rec := range records {
        //fmt.Printf("%v ------ %d \n", rec, rec.GetPossible())
        answer1 += rec.GetPossible()
    }
    fmt.Printf("Part 1 Answer: %d\n", answer1)
}
