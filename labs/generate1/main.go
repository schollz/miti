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
	"gopkg.in/music-theory.v0/chord"
	"gopkg.in/music-theory.v0/scale"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	err := run()
	if err != nil {
		log.Error(err)
	}
}

func run() (err error) {
	fmt.Println("chord")
	s1 := chord.Of("A(add9)")
	for _, note := range s1.Notes() {
		fmt.Println(note.Class.String(s1.AdjSymbol))
	}
	fmt.Println("notes")
	s := scale.Of("A(add9)")
	for _, note := range s.Notes() {
		fmt.Println(note.Class.String(s.AdjSymbol))
	}
	err = loadChords()
	if err != nil {
		return
	}

	// determine chords
	numChords := 4
	// chord changes add up to 16 beats
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
	mitiChords := ""
	totalBeats := 0
	for j, beats := range changes {
		for i := 0; i < beats; i++ {
			if totalBeats == 4 {
				mitiChords += "\n"
				totalBeats = 0
			}
			mitiChords += chords[j]
			if i < beats-1 {
				mitiChords += "-"
			}
			mitiChords += " "
			totalBeats++
		}
	}
	fmt.Println(mitiChords)

	// determine melody based on chords
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
