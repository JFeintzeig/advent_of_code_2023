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

const MAX_UINT = ^uint(0) 
const MAX_INT = int(MAX_UINT >> 1) 

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

func LineToGame(line string) Game {
  game := strings.Split(line, ":")[0]
  roundString := strings.Split(strings.Split(line, ":")[1], ";")

  gID, _ := regexp2.Compile("(?<=Game )\\d+", 0)
  idStr, _ := gID.FindStringMatch(game)
  id, _ := strconv.Atoi(idStr.String())

  red, _ := regexp2.Compile("\\d+(?= red)", 0)
  green, _ := regexp2.Compile("\\d+(?= green)", 0)
  blue, _ := regexp2.Compile("\\d+(?= blue)", 0)

  var rounds []CubeSet

  for _, r := range roundString {
    nRed, nGreen, nBlue := 0, 0, 0

    nRedStr, _ := red.FindStringMatch(r)
    if nRedStr != nil {
      nRed, _ = strconv.Atoi(nRedStr.String())
    }

    nGreenStr, _ := green.FindStringMatch(r)
    if nGreenStr != nil {
      nGreen, _ = strconv.Atoi(nGreenStr.String())
    }

    nBlueStr, _ := blue.FindStringMatch(r)
    if nBlueStr != nil {
      nBlue, _ = strconv.Atoi(nBlueStr.String())
    }

    rounds = append(rounds, CubeSet{nRed, nGreen, nBlue})
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
