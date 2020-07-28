package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	log "github.com/schollz/logger"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	err := run()
	if err != nil {
		log.Error(err)
	}
}

func run() (err error) {
	err = loadChords()
	if err != nil {
		return
	}
	numChords := 4
	changes := generateNumbersThatAddTo16(numChords)
	chords := make([]string, numChords)
startover:
	chords[0] = "C"
	chords[1] = ""
	chords[2] = ""
	chords[3] = ""
	for i := 1; i < numChords; i++ {
		chordString := strings.TrimSpace(strings.Join(chords, " "))
		if _, ok := chordChanges[chordString]; !ok {
			log.Debugf("could not find '%s'", chordString)
			goto startover
		}
		chords[i] = randomWeightedChoice(chordChanges[chordString])
	}
	fmt.Println(chords)
	fmt.Println(changes)
	return
}

var chordChanges map[string]map[string]float64

func loadChords() (err error) {
	b, err := ioutil.ReadFile("chordIndexInC.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &chordChanges)
	return
}

func randomWeightedChoice(m map[string]float64) string {
	type kv struct {
		Key   string
		Value float64
	}
	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	curSum := 0.0
	target := rand.Float64() * 100
	for _, kv := range ss {
		curSum += kv.Value
		if curSum >= target {
			return kv.Key
		}
	}
	panic("could not find key")
	return ""
}

func generateNumbersThatAddTo16(numNumbers int) (nums []int) {
	var total int
	nums = make([]int, numNumbers)
tryagain:
	total = 0
	for i := range nums {
		// nums[i] = rand.Intn(12) + 1 // numbers must be in range [1,13]
		nums[i] = int(math.Round(rand.NormFloat64()*1.5+3)) + 1
		if nums[i] < 1 {
			nums[i] = 1
		}
		total += nums[i]
		if total > 16 {
			goto tryagain
		}
	}
	if total != 16 {
		goto tryagain
	}
	return
}
