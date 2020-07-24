package music

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kbinani/midi"
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
	Octave int
	MIDI   int
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
	return allowedNotes[note] + 12*(octave-1)
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
		n = ClosestNote(s, lastNote)
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
			n = ClosestNote(nn, lastNote)
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

func ClosestNote(name string, n Note) (n2 Note) {
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
var chordToNotesCache map[string][]Note

func init() {
	chordToNotesCache = make(map[string][]Note)
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

// ChordToNotes converts chords to notes using lilypond
func ChordToNotes(c string) (notes []Note, err error) {
	if _, ok := chordToNotesCache[c]; ok {
		notes = chordToNotesCache[c]
		return
	}
	defer func() {
		if err == nil {
			chordToNotesCache[c] = notes
		}
	}()

	tmpfile, err := ioutil.TempFile("", "lilypond")
	if err != nil {
		return
	}
	defer os.Remove(tmpfile.Name()) // clean up

	_, err = tmpfile.WriteString(`
\score {
\chordmode { ` + ChordToLilypond(c) + ` }
  \midi { }
}`)
	if err != nil {
		return
	}
	err = tmpfile.Close()
	if err != nil {
		return
	}

	cmd := exec.Command("lilypond", "-o", tmpfile.Name(), tmpfile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("lilypond error: %s\n\ndata: '%s'", err.Error(), output)
		return
	}
	defer os.Remove(tmpfile.Name() + ".midi")

	f, err := os.Open(tmpfile.Name() + ".midi")
	if err != nil {
		return
	}
	defer f.Close()

	file, err := midi.Read(f)
	if err != nil {
		return
	}
	for i, track := range file.Tracks {
		if i != 1 {
			continue
		}
		for _, n := range track.Events {
			if len(n.Messages) < 3 {
				continue
			}
			if n.Tick == 0 && n.Messages[2] == 90 {
				notes = append(notes, MidiToNote(int(n.Messages[1])))
			}
		}
	}
	return
}

// ChordToLilypond takes a chord like Bbm7/F and converts it
// into suitable lilypond (bes:m7/f)
func ChordToLilypond(c string, beats ...float64) (d string) {
	beats0 := 4.0
	if len(beats) > 0 {
		beats0 = beats[0]
	}
	c = strings.Replace(strings.ToLower(c), " ", "", -1)
	if len(c) > 1 && string(c[1]) == "b" {
		d = fmt.Sprintf("%ses%d", string(c[0]), int(4.0/beats0))
		if len(c) > 2 {

		}
		if len(c) > 2 {
			d += ":" + c[2:]
		}
	} else if len(c) > 1 && string(c[1]) == "#" {
		d = fmt.Sprintf("%sis%d", string(c[0]), int(4.0/beats0))
		if len(c) > 2 {
			d += ":" + c[2:]
		}
	} else {
		d = fmt.Sprintf("%s%d", string(c[0]), int(4.0/beats0))
		if len(c) > 1 {
			d += ":" + c[1:]
		}
	}
	for _, a := range []string{"a", "b", "c", "d", "e", "f", "g"} {
		d = strings.Replace(d, a+"b", a+"es", -1)
		d = strings.Replace(d, a+"#", a+"is", -1)
	}
	d = strings.Replace(d, "(add9)", "5.9", -1)
	d = strings.Replace(d, "(b5)", "7.5-", -1)
	d = strings.Replace(d, "o", "dim", -1)
	return
}
