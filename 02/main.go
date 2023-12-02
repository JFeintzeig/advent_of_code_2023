package main

import (
  "bufio"
  "fmt"
  "os"
  "strings"
  "strconv"

  "github.com/dlclark/regexp2"
)

const (
  MAX_RED = 12
  MAX_GREEN = 13
  MAX_BLUE = 14
)

type CubeSet struct {
  nRed int
  nGreen int
  nBlue int
}

func (cs *CubeSet) power() int {
  return cs.nRed * cs.nGreen * cs.nBlue
}

type Game struct {
  id int
  rounds []CubeSet
}

func (g *Game) isPossible() bool {
  for _, r := range g.rounds {
    if r.nRed > MAX_RED || r.nGreen > MAX_GREEN || r.nBlue > MAX_BLUE {
      return false
    }
  }
  return true
}

func (g *Game) minSet() CubeSet {
  minRed, minGreen, minBlue := 0, 0, 0
  for _, r := range g.rounds {
    if r.nRed > minRed { minRed = r.nRed }
    if r.nGreen > minGreen { minGreen = r.nGreen }
    if r.nBlue > minBlue { minBlue = r.nBlue }
  }

  return CubeSet{minRed, minGreen, minBlue}
}

func LoadFile(file string) *bufio.Scanner {
    f, err := os.Open("input.data")
    if err != nil {
      panic("cant open file")
    }

    scanner := bufio.NewScanner(f)
    return scanner
}

func extractIntFromRegexp(str string, re *regexp2.Regexp) int {
  n := 0
  nStr, _ := re.FindStringMatch(str)
  if nStr != nil {
    n, _ = strconv.Atoi(nStr.String())
  }
  return n
}

func LineToGame(line string) Game {
  game := strings.Split(line, ":")[0]
  roundString := strings.Split(strings.Split(line, ":")[1], ";")

  gID, _ := regexp2.Compile("(?<=Game )\\d+", 0)
  id := extractIntFromRegexp(game, gID)

  colors := []string{"red", "green","blue"}
  colorRE := map[string]*regexp2.Regexp{}
  for _, col := range colors {
    colorRE[col], _ = regexp2.Compile(fmt.Sprintf("\\d+(?= %s)", col), 0)
  }

  var rounds []CubeSet

  for _, r := range roundString {
    colorCount := make(map[string]int)
    for k, v := range colorRE {
      colorCount[k] = extractIntFromRegexp(r, v)
    }

    rounds = append(rounds, CubeSet{colorCount["red"], colorCount["green"], colorCount["blue"]})
  }

  return Game{id: id, rounds: rounds}
}

func main() {
  scanner := LoadFile("test1.data")

  sumPossible := 0
  sumMinPower := 0

  for scanner.Scan() {
    line := scanner.Text()
    g := LineToGame(line)
    if g.isPossible() {
      sumPossible += g.id
    }
    minSet := g.minSet()
    sumMinPower += minSet.power()
  }

  fmt.Printf("Possible: %d\n", sumPossible)
  fmt.Printf("Min Power: %d\n", sumMinPower)
}
