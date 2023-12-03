package main

import (
  "fmt"
  "os"
  "strconv"
  "strings"
)

type Matrix struct {
  values [][]string
  digits map[string]bool
  symbols map[string]bool
  neighborOffsets [][]int
}

func (m *Matrix) nRows() int {
  return len(m.values)
}

func (m *Matrix) nCols() int {
  return len(m.values[0])
}

func (m *Matrix) isDigit(x int, y int) bool {
  _, ok := m.digits[m.values[y][x]]
  return ok
}

func (m *Matrix) isSymbol(x int, y int) bool {
  _, ok := m.symbols[m.values[y][x]]
  return ok
}

func (m *Matrix) hasSymbolNeighbor(x int, y int) bool {
  return len(m.getSymbolNeighbors(x, y)) > 0
}

func (m *Matrix) getSymbolNeighbors(x int, y int) [][]int {
  out := make([][]int,0)
  for _, coord := range m.neighborOffsets { 
    dx := coord[0]
    dy := coord[1]
    if x+dx >= 0 && x+dx < m.nCols() && y+dy >=0 && y+dy < m.nRows() {
      if m.isSymbol(x+dx, y+dy) {
        out = append(out, []int{x+dx, y+dy})
      }
    }
  }
  return out
}

func NewMatrix(values [][]string) *Matrix {
  digits := map[string]bool{"0": true, "1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true, "8": true, "9": true}
  symbols := map[string]bool{"!":true, "@": true, "#": true, "$": true, "%": true, "^": true, "&": true, "*": true, "(": true, ")": true, "_": true, "+": true, "=": true, "-": true, "?": true, "/": true}
  neighborOffsets := [][]int{{-1,-1},{0,-1},{1,-1},{1,0},{1,1},{0,1},{-1,1},{-1,0}}
  return &Matrix{values, digits, symbols, neighborOffsets}
}

// just for part 2
type Point struct {
  x int
  y int
}

func LoadFileAsMatrix(fileName string) *Matrix {
  file, _ := os.ReadFile(fileName)

  lines := strings.Split(strings.TrimSpace(string(file)),"\n")
  matrix := make([][]string, len(lines))

  for i, line := range lines {
    matrix[i] = strings.Split(line, "")
  }
  return NewMatrix(matrix)
}

func main() {
  matrix := LoadFileAsMatrix("input.data")

  // part 1
  // for each entry:
  // (1) if its a digit
  //    (a) append to running digit string
  //    (b) check neighbors for symbols, if so mark hasSymbol true
  // (2) if its not a digit:
  //     (a) if hasSymbol is true, add digitString to sum
  //     (b) no matter what, reset digit string and hasSymbol

  // part 2
  //  (1)(c) append (x,y) coords of symbols to running list
  //  (2)(c) for each symbol that is a *, add the num to a running
  //         list for that symbol (in a map)
  //  (3) loop over symbol map, find symbols that have 2 nums, multiply and add

  var answer1 int
  var sb strings.Builder
  hasSymbol := false

  symbols := make([][]int, 0)
  symbolNumMap := make(map[Point][]int)

  for y, row := range matrix.values {
    for x, item := range row {
      if matrix.isDigit(x, y) {
        // append to running digit string
        sb.WriteString(item)
        // eval symbols
        hasSymbol = hasSymbol || matrix.hasSymbolNeighbor(x, y)
        // append to symbols list
        symbols = append(symbols, matrix.getSymbolNeighbors(x, y)...)
      } else {
        // end of digit string, now we eval
        if hasSymbol {
          // part 1
          num, _ := strconv.Atoi(sb.String())
          answer1 += num

          // part 2
          for _, symXY := range symbols {
            if matrix.values[symXY[1]][symXY[0]] != "*" {
              continue
            }
            sym := Point{symXY[0], symXY[1]}

            _, ok := symbolNumMap[sym]
            if ok {
              // if its already in there, dont add it again
              if symbolNumMap[sym][len(symbolNumMap[sym])-1] != num {
                symbolNumMap[sym] = append(symbolNumMap[sym], num)
              }
            } else {
              symbolNumMap[sym] = []int{num}
            }
          }
        }
        // reset for new digits
        sb.Reset()
        hasSymbol = false
        symbols = make([][]int, 0)
      }
    }
  }

  fmt.Printf("Part1: %d\n", answer1)

  var answer2 int 
  for _, val := range symbolNumMap {
    if len(val)==2 {
      answer2 += val[0]*val[1]
    }
  }

  fmt.Printf("Part2: %d\n", answer2)
}
