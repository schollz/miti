package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/schollz/logger"
	"gopkg.in/music-theory.v0/chord"
	"gopkg.in/music-theory.v0/scale"
)

func init() {
	// load chord transitions
	err := loadChords()
	if err != nil {
		panic(err)
	}
}

var startChord string
var seed int64
var srand rand.Source
var rrand *rand.Rand

func main() {
	// get starting chord
	startChord = "C"
	if len(os.Args) > 1 {
		startChord = os.Args[1]
	}

	// get starting seed
	seed = int64(0)
	if len(os.Args) > 2 {
		seed, _ = strconv.ParseInt(os.Args[2], 10, 64)
	}
	if seed <= 0 {
		seed = time.Now().UnixNano()
	}
	srand = rand.NewSource(seed)
	rrand = rand.New(srand)
	fmt.Printf("seed: %d\n\n", seed)

	err := run()
	if err != nil {
		log.Error(err)
	}
}

func run() (err error) {
	// determine chords
	numChords := 4
	// chord changes add up to 16 beats
	changes := generateNumbersThatAddTo(numChords, 16)
	chords := make([]string, numChords)
startover:
	chords[0] = startChord
	chords[1] = ""
	chords[2] = ""
	chords[3] = ""
	for i := 1; i < numChords; i++ {
		chordString := strings.TrimSpace(strings.Join(chords, " "))
		if _, ok := chordChanges[chordString]; !ok {
			// log.Debugf("could not find '%s'", chordString)
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

func notesInChords(chordName string) []string {
	s1 := chord.Of(chordName)
	notes := make([]string, len(s1.Notes()))
	for i, note := range s1.Notes() {
		notes[i] = note.Class.String(s1.AdjSymbol)
	}
	return notes
}

func notesInScale(scaleName string) []string {
	s1 := scale.Of(scaleName)
	notes := make([]string, len(s1.Notes()))
	for i, note := range s1.Notes() {
		notes[i] = note.Class.String(s1.AdjSymbol)
	}
	return notes

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
		return ss[i].Key > ss[j].Key
	})

	curSum := 0.0
	target := rrand.Float64() * 100
	for _, kv := range ss {
		curSum += kv.Value
		if curSum >= target {
			return kv.Key
		}
	}
	panic("could not find key")
	return ""
}

func generateNumbersThatAddTo(numNumbers int, addUpTo int) (nums []int) {
	var total int
	nums = make([]int, numNumbers)
	for {
		total = 0
		for i := range nums {
			if total > addUpTo {
				continue
			}
			// nums[i] = rand.Intn(12) + 1 // numbers must be in range [1,13]
			nums[i] = int(math.Floor(rrand.NormFloat64()*float64(addUpTo)/10+float64(addUpTo)/4-1)) + 1
			if nums[i] < 1 {
				nums[i] = 1
			}
			total += nums[i]
		}
		if total == addUpTo {
			break
		}
	}
	return
}
