package main

import (
  "fmt"
  "os"
  "bufio"
  "strconv"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {

    intStrings := make(map[string]int)
    intStrings["zero"]=0
    intStrings["one"]=1
    intStrings["two"]=2
    intStrings["three"]=3
    intStrings["four"]=4
    intStrings["five"]=5
    intStrings["six"]=6
    intStrings["seven"]=7
    intStrings["eight"]=8
    intStrings["nine"]=9

    f, err := os.Open("input.data")
    check(err)

    scanner := bufio.NewScanner(f)

    var sum int = 0

    for scanner.Scan() {
      line := scanner.Text()

      if err != nil {
        break
      }

      var firstIndex int = 65535
      var lastIndex int = 0

      var firstNum int
      var lastNum int

      for index, val := range line {
        num, err := strconv.Atoi(string(val))

        // if single character is a digit
        if err == nil {
          if index <= firstIndex {
            firstIndex = index
            firstNum = num
          }
          if index >= lastIndex {
            lastIndex = index
            lastNum = num
          }
        } else {
          // loop over keys in map, for each key:
          // get len of key
          // get line[index:index+len]
          // check if equal to digit
          // then do first/last index/num
          for key, val := range intStrings {
            l := len(key)
            if index+l > len(line) {
              continue
            }
            if line[index:index+l] == key {
              if index <= firstIndex {
                firstIndex = index
                firstNum = val
              }
              if index >= lastIndex {
                lastIndex = index
                lastNum = val
              }
            }
          }
        }
      }

      sum += firstNum*10 + lastNum
    }

    fmt.Println(sum)
}
