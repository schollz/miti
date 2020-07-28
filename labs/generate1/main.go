package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
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
	fmt.Println(randomWeightedChoice(chordChanges["C"]))
	fmt.Println(generateNumbersThatAddTo16(4))
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
	curSum := 0.0
	target := rand.Float64() * 100
	lastKey := ""
	for key := range m {
		curSum += m[key]
		if curSum >= target {
			return key
		}
		lastKey = key
	}
	return lastKey
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
