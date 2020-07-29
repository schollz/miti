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
			mitiChords += ":" + chords[j] + ":3"
			if i < beats-1 {
				mitiChords += "-"
			}
			mitiChords += " "
			totalBeats++
		}
	}
	fmt.Println(mitiChords)
	//	determine melody based on chords
	melody := make([]string, 64)
	mi := 0
	numNotes := 2
	for i, chord := range chords {
		// get random notes
		beats := changes[i] * 4
		if beats > 4 {
			numNotes = 4
		}
		if beats > 8 {
			numNotes = 8
		}
		if numNotes > beats {
			numNotes = beats / 2
		}
		if numNotes <= 0 {
			numNotes = 1
		}

		notes := getRandomNotes(chords[0], chord, numNotes)
		// get random subdivisions
		log.Debugf("generate %d numbers that add up to %d", numNotes, beats)
		subdivisions := generateNumbersThatAddTo(numNotes, beats)
		for j, subdivision := range subdivisions {
			for k := 0; k < subdivision; k++ {
				melody[mi] = notes[j]
				mi++
			}
		}
	}
	// print meldoy
	mitiMelody := ""
	for i, note := range melody {
		if i%16 == 0 {
			mitiMelody += "\n"
		}
		mitiMelody += note + "5"
		// if i == 0 || note != melody[i-1] {
		// }
		if i < len(melody)-1 && note == melody[i+1] {
			mitiMelody += "-"
		}
		mitiMelody += " "
	}
	fmt.Println(mitiMelody)

	f, err := os.Create("miti.txt")
	if err != nil {
		return
	}
	defer f.Close()

	f.WriteString(`
# generate1
# seed ` + fmt.Sprint(seed) + `

tempo 90

pattern 1

instruments sh-01a
legato 90
` + mitiChords + `

instruments op-1
legato 90
` + mitiMelody + `
`)
	return
}

func getRandomNotes(key string, chord string, numNotes int) (notes []string) {
	// get possible notes, weighted by importance
	// important notes = notes in key, notes in chord, then notes in scale of chord
	noteChanges := make(map[string]float64)
	for _, note := range notesInScale(key) {
		noteChanges[note] = 1
	}
	for _, note := range notesInScale(chord) {
		if _, ok := noteChanges[note]; !ok {
			noteChanges[note] = 0
		}
		noteChanges[note]++
	}
	for _, note := range notesInChords(chord) {
		if _, ok := noteChanges[note]; !ok {
			noteChanges[note] = 1
		}
		noteChanges[note]++
	}
	//  normalize note changes
	total := 0.0
	for note := range noteChanges {
		total += noteChanges[note]
	}
	for note := range noteChanges {
		noteChanges[note] = noteChanges[note] / total * 100
	}

	notes = make([]string, numNotes)
	for i := 0; i < numNotes; i++ {
		notes[i] = randomWeightedChoice(noteChanges)
	}
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
	if numNumbers == 1 {
		return []int{addUpTo}
	}
	var total int
	nums = make([]int, numNumbers)
	tries := 0
	for {
		total = 0
		for i := range nums {
			if total > addUpTo {
				continue
			}
			if tries > 5 {
				nums[i] = rand.Intn(int(addUpTo/numNumbers+1)) + 1 // numbers must be in range [1,13]
			} else {
				nums[i] = int(math.Floor(rrand.NormFloat64()*float64(addUpTo)/10+float64(addUpTo)/4-1)) + 1
			}
			if nums[i] < 1 {
				nums[i] = 1
			}
			total += nums[i]
		}
		if total == addUpTo {
			break
		}
		tries++
	}
	return
}
