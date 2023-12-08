package main

import (
  "fmt"
  "os"
  "sort"
  "strings"
  "strconv"
)

const part = 2
var cardTypes = make(map[string]int)

type HandList []*Hand

func (h HandList) Len() int {
  return len(h)
}

func (h HandList) Swap(i, j int) {
  h[i], h[j] = h[j], h[i]
}

func (h HandList) Less(i, j int) bool {
  return h[i].Less(h[j])
}

type HandType int

const (
    HighCard HandType = iota
    Pair
    TwoPair
    ThreeKind
    FullHouse
    FourKind
    FiveKind
)

type Hand struct {
  cards []string
  bid int
  cardCount map[string]int
  score HandType
}

func (h1 *Hand) Less(h2 *Hand) bool {
  if h1.score < h2.score {
    return true
  } else if h1.score == h2.score {
    for i, c := range h1.cards {
      if cardTypes[c] < cardTypes[h2.cards[i]] {
        return true
      } else if cardTypes[c] > cardTypes[h2.cards[i]] {
        return false
      } else {
        continue
      }
    }
    panic("what?")
    return false
  } else {
    return false
  }
}

func (h *Hand) countCards() {
  cc := make(map[string]int)
  for _, c := range h.cards {
    _, ok := cc[c]
    if ok {
      cc[c] += 1
    } else {
      cc[c] = 1
    }
  }
  h.cardCount = cc
}

func (h *Hand) calcScore() {
    if h.cardCount == nil {
      h.countCards()
    }

    var maxType1 HandType = HighCard
    var maxType2 HandType = HighCard
    for _, v := range h.cardCount {
      switch v {
        case 5:
          maxType1 = FiveKind
          break
        case 4:
          maxType1 = FourKind
          break
        case 3:
          if maxType1 == HighCard {
            maxType1 = ThreeKind
          } else {
            maxType2 = ThreeKind
          }
        case 2:
          if maxType1 == HighCard {
            maxType1 = Pair
          } else {
            maxType2 = Pair
          }
      }
    }
    if maxType1 == ThreeKind && maxType2 == Pair {
      maxType1 = FullHouse
    }
    if maxType2 == ThreeKind && maxType1 == Pair {
      maxType1 = FullHouse
    }
    if maxType1 == Pair && maxType2 == Pair {
      maxType1 = TwoPair
    }

    h.score = maxType1

    // part 2 joker routine modifies score
    if part == 2 {
      val, ok := h.cardCount["J"]
      if !ok {
        return
      }
      if val == 5 {
        return
      }
      if val == 4 {
        return
      }
      switch maxType1 {
        case FourKind:
            maxType1 = FiveKind
        case FullHouse:
            maxType1 = FiveKind
        case ThreeKind:
            // either 3 jokers + 1 other or 3 others + 1 joker
            // 3 others + 1 joker would already be a full house
            maxType1 = FourKind
        case TwoPair:
            if val == 1 {
              maxType1 = FullHouse
            } else if val == 2 {
              maxType1 = FourKind
            }
        case Pair:
            maxType1 = ThreeKind
        case HighCard:
            maxType1 = Pair
      }
      h.score = maxType1
    }
}

func NewHand(cards []string, bid int) *Hand {
  h := Hand{cards: cards, bid: bid}
  h.countCards()
  h.calcScore()
  return &h
}

func LineToHand(line string) *Hand {
  sp := strings.Split(strings.TrimSpace(line)," ")
  cards := sp[0]
  bid, _ := strconv.Atoi(sp[1])
  return NewHand(strings.Split(cards,""), bid)
}

func LoadFile(fileName string) []string {
	file, _ := os.ReadFile(fileName)

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	return lines
}

func main() {
    cardTypes["2"] = 2
    cardTypes["3"] = 3
    cardTypes["4"] = 4
    cardTypes["5"] = 5
    cardTypes["6"] = 6
    cardTypes["7"] = 7
    cardTypes["8"] = 8
    cardTypes["9"] = 9
    cardTypes["T"] = 10
    cardTypes["J"] = 11
    cardTypes["Q"] = 12
    cardTypes["K"] = 13
    cardTypes["A"] = 14

    if part == 2 {
      cardTypes["J"] = 1
    }

	lines := LoadFile("test1.data")
    
    var hands HandList
    for _, line := range lines {
      h := LineToHand(line)
      hands = append(hands, h)
    }

    sort.Sort(hands)

    out1 := 0
    for i, h := range hands {
      fmt.Printf("%v %d\n",*h, (i+1)*h.bid)
      out1 += (i+1)*h.bid
    }

    fmt.Printf("%d\n",len(hands))

    fmt.Printf("Part 1 Answer: %d\n", out1)
}
