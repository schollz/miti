package music

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Chord struct {
	Notes []Note
	On    bool
}

func (c Chord) String() string {
	s := ""
	for i, n := range c.Notes {
		if i == 0 {
			s += fmt.Sprintf("%s%d", n.Name, n.Octave)
		} else {
			s += n.Name
		}
	}
	return s
}

type Note struct {
	Name   string
	MIDI   int
	Octave int
}

func NewNote(name string, octave int) Note {
	return Note{
		Name:   name,
		Octave: octave,
		MIDI:   NoteToMidi(name, octave),
	}
}

func NoteToMidi(note string, octave int) int {
	if _, ok := allowedNotes[note]; !ok {
		return -1
	}
	return allowedNotes[note] + 12*octave
}

func ParseCluster(cluster string, lastNote0 ...Note) (ns []Note, err error) {
	lastNote := NewNote("C", 4)
	if len(lastNote0) > 0 {
		lastNote = lastNote0[0]
	}

	cc, nn, err := popFirstNote(cluster, lastNote)
	if err != nil {
		return
	}
	ns = append(ns, nn)
	lastNote = nn

	it := 0
	for {
		it++
		if it > 10 {
			err = fmt.Errorf("too many notes in cluster: %s", cluster)
			return
		}
		if cc == "" {
			break
		}
		cc, nn, err = popFirstNote(cc, lastNote)
		if err != nil {
			return
		}
		ns = append(ns, nn)
		lastNote = nn
	}

	return
}

func popFirstNote(s string, lastNote Note) (s2 string, n Note, err error) {
	defer func() {
		n.MIDI = NoteToMidi(n.Name, n.Octave)
	}()
	if isValidNote(s) {
		n = closestNote(s, lastNote)
		return
	}
	for i := 1; i <= len(s); i++ {
		nn := string(s[:i])
		nn1 := nn[len(nn)-1:]
		if isNumber(nn1) {
			s2 = strings.TrimPrefix(s, nn)
			octave, _ := strconv.Atoi(nn1)
			n = Note{Name: strings.TrimSuffix(nn, nn1), Octave: octave}
			return
		}
		if !isValidNote(nn) {
			nn = string(s[:i-1])
			s2 = strings.TrimPrefix(s, nn)
			n = closestNote(nn, lastNote)
			return
		}
	}
	err = fmt.Errorf("could not parse: '%s'", s)
	return
}

func isNumber(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

func isValidNote(s string) bool {
	if _, ok := allowedNotes[s]; ok {
		return true
	}
	return false
}

func closestNote(name string, n Note) (n2 Note) {
	possibleNotes := []Note{NewNote(name, n.Octave), NewNote(name, n.Octave-1), NewNote(name, n.Octave+1)}
	midiDiff := 10000.0
	for i := 0; i < 3; i++ {
		d := math.Abs(float64(possibleNotes[i].MIDI - n.MIDI))
		if d < midiDiff {
			midiDiff = d
			n2 = possibleNotes[i]
		}
	}
	return
}

var allowedNotes = map[string]int{
	"C":  24,
	"C#": 25,
	"Cb": 23,
	"D":  26,
	"D#": 27,
	"Db": 25,
	"E":  28,
	"E#": 29,
	"Eb": 27,
	"F":  29,
	"F#": 30,
	"Fb": 28,
	"G":  31,
	"G#": 32,
	"Gb": 30,
	"A":  33,
	"A#": 34,
	"Ab": 32,
	"B":  35,
	"B#": 26,
	"Bb": 34,
}

var c0notes = map[string]int{
	"C":  24,
	"Db": 25,
	"D":  26,
	"Eb": 27,
	"E":  28,
	"F":  29,
	"Gb": 30,
	"G":  31,
	"Ab": 32,
	"A":  33,
	"Bb": 34,
	"B":  35,
}
var midiToNote map[int]Note

func init() {
	midiToNote = make(map[int]Note)
	for octave := 1; octave < 8; octave++ {
		for note := range c0notes {
			midiToNote[c0notes[note]+12*(octave-1)] = Note{Name: note, Octave: octave, MIDI: c0notes[note] + 12*(octave-1)}
		}
	}
}

func MidiToNote(midi int) Note {
	if _, ok := midiToNote[midi]; !ok {
		return Note{}
	}
	return midiToNote[midi]
}
