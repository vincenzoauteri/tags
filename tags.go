package main

import (
  "sort"
  "fmt"
  "os"
  "encoding/csv"
  "strings"
  "regexp"
)

type sortedMap struct {
  m map[string]int
  s []string
}

func (sm *sortedMap) Len() int {
  return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
  return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
  sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}


func sortedKeys(m map[string]int) []string {
  sm := new(sortedMap)
  sm.m = m
  sm.s = make([]string, len(m))
  i := 0
  for key, _ := range m {
    sm.s[i] = key
    i++
  }
  sort.Sort(sm)
  return sm.s
}

var tagMap map[string]int

func main() {
  //An artificial input source.
  f, err := os.Open("./Train.csv")
  if err != nil {
    fmt.Printf("error opening file: %v\n",err)
    os.Exit(1)
  }

  r := csv.NewReader(f)
  slices, _ := r.Read()

  tagMap = make(map[string]int)

  for i :=0 ; slices != nil ; slices, _ = r.Read() {
    tags := strings.Split(slices[3]," ")
    for _, tag := range tags {

      if val,ok := tagMap[tag]; ok {
        tagMap[tag]= val+1
      } else {
        tagMap[tag]= 1
      }

      specialPattern := regexp.MustCompile("(\\]|\\[|^\\$|\\.|\\||\\?|\\*|\\+|\\(|\\))")
      escapedTag := specialPattern.ReplaceAllString(tag,"\\${1}")

      tagPattern := regexp.MustCompile(escapedTag)

      if tagPattern.MatchString(strings.ToLower(slices[1])) {
        //fmt.Println("Matches tag on title:" , tag)
      }

      if tagPattern.MatchString(strings.ToLower(slices[2])) {
        //fmt.Println("Matches tag on body:" , tag)
      }
    }
    i=i+1
    if i % 100000  == 0{
      fmt.Println(slices[0])
    }
  }
  fmt.Println(sortedKeys(tagMap))

  // Set the split function for the scanning operation.
  //scanner.Split(bufio.ScanWords)
  // Count the words.
  //if err := scanner.Err(); err != nil {
  //  fmt.Fprintln(os.Stderr, "reading input:", err)
  //}
  //fmt.Printf("%d\n", count)
}
