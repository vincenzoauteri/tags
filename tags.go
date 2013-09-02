package main

import (
  "sort"
  "fmt"
  "os"
  "bufio"
  "encoding/csv"
  "strings"
  "regexp"
  "time"
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

func getTagListFromTrainData() {
  //An artificial input source.
  f, err := os.Open("/vol/Train.csv")
  if err != nil {
    fmt.Printf("error opening file: %v\n",err)
    os.Exit(1)
  }

  r := csv.NewReader(f)
  slices, _ := r.Read()

  tagMap = make(map[string]int)

  //specialPattern := regexp.MustCompile("(\\]|\\[|^\\$|\\.|\\||\\?|\\*|\\+|\\(|\\))")

  startTime := time.Now()

  for i :=0 ; slices != nil ; slices, _ = r.Read() {
    tags := strings.Split(slices[3]," ")
    for _, tag := range tags {

      if val,ok := tagMap[tag]; ok {
        tagMap[tag]= val+1
      } else {
        tagMap[tag]= 1
      }

    }
    if i % 10000  == 0 {
        fmt.Println("tag N:", slices[0], " Elapsed time:",time.Since(startTime))
    }
    i=i+1
  }
  fmt.Println(sortedKeys(tagMap))

  //open output file
  fo, err := os.Create("trainTags.txt")
  if err != nil { panic(err) }
  // close fo on exit and check for its returned error
  defer func() {
    if err := fo.Close(); err != nil {
      panic(err)
    }
  }()

  if _, err := fo.Write([]byte(strings.Join(sortedKeys(tagMap)," "))); err != nil {
    panic(err)
  }
}

func loadTagList() []string {
  f, err := os.Open("./trainTags.txt")

  if err != nil {
    fmt.Printf("error opening file: %v\n",err)
    os.Exit(1)
  }
  /*
  scanner := bufio.NewScanner(f)

  for scanner.Scan() {
    fmt.Println(scanner.Bytes())
  }
  */
  scanner := bufio.NewReader(f)

  s,_ := scanner.ReadString('\n')

  return strings.Split(s," ")
}

func makePredictionOnTestData(tagSlice []string) {
  f, err := os.Open("./Test.csv")
  if err != nil {
    fmt.Printf("error opening file: %v\n",err)
    os.Exit(1)
  }

  //open output file
  fo, err := os.Create("Submission.txt")
  if err != nil { panic(err) }
  // close fo on exit and check for its returned error
  defer func() {
    if err := fo.Close(); err != nil {
      panic(err)
    }
  }()
  
  if _, err := fo.Write([]byte("\"Id\",\"Tags\"\n")); err != nil {
    panic(err)
  }

  r := csv.NewReader(f)

  //Skip header
  slices, _ := r.Read()

  slices, _ = r.Read()

  specialPattern := regexp.MustCompile("(\\]|\\[|^\\$|\\.|\\||\\?|\\*|\\+|\\(|\\))")

  startTime := time.Now()

  for i :=0 ; slices != nil ; slices, _ = r.Read() {

    numMatches := 0
    matches := ""
    for _, tag := range tagSlice {

      tagChunks := strings.Split(tag,"-")
      match := true
      for _, tagChunk := range tagChunks{

        escapedTagChunk := specialPattern.ReplaceAllString(tagChunk,"\\${1}")
        tagChunkPattern := regexp.MustCompile("\\W" + escapedTagChunk + "\\W")

        if !tagChunkPattern.MatchString(strings.ToLower(slices[1]))  && !tagChunkPattern.MatchString(strings.ToLower(slices[2])) {
          match = false
        }
      }
      if match {
        numMatches = numMatches + 1
        matches = matches + tag + " "
        if numMatches >=2 {
          break
        }
      }
    }
    i=i+1
    if i % 10000  == 0 {
        fmt.Println("tag N:", slices[0], " Elapsed time:",time.Since(startTime))
    }
    //fmt.Println("Matching tags:", matches, "with title", slices[1])
    if _, err := fo.Write([]byte(slices[0]+",\""+matches+"\"\n")); err != nil {
      panic(err)
    }
  }

}

func main() {
  //tagSlice:= loadTagList()
  //makePredictionOnTestData(tagSlice)
  startTime := time.Now()
  getTagListFromTrainData()
  fmt.Println("Elapsed time", time.Since(startTime))
}
