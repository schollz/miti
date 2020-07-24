package music

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCluster(t *testing.T) {
	cluster := "CEG"
	notes, err := ParseCluster(cluster)
	if err != nil {
		t.Errorf("err: %s", err.Error())
	}
	fmt.Println(notes)
}

func TestMidiToNote(t *testing.T) {
	note := MidiToNote(98)
	if note.Name != "D" && note.Octave == 7 {
		t.Errorf("got wrong note: %v", note)
	}
	note = MidiToNote(26)
	if note.Name != "D" && note.Octave == 1 {
		t.Errorf("got wrong note: %v", note)
	}
}

func TestNoteToMidi(t *testing.T) {
	note := NoteToMidi("F", 4)
	if note != 65 {
		t.Errorf("got wrong note: %v", note)
	}
}

func TestChordToNotes(t *testing.T) {
	var tts = []struct {
		chord string
		notes []Note
	}{
		{"Bbm7/F", []Note{Note{"F", 4, 65},
			Note{"Ab", 4, 68},
			Note{"Bb", 4, 70},
			Note{"Db", 5, 73}}},
	}

	for _, tt := range tts {
		notes, err := ChordToNotes(tt.chord)
		assert.Nil(t, err)
		assert.Equal(t, tt.notes, notes)
	}
}

func TestChordToLilypond(t *testing.T) {
	var tts = []struct {
		chord         string
		lilypondchord string
	}{
		{"Bbm7/F", "bes1:m7/f"},
		{"C#", "cis1"},
		{"Cdim", "c1:dim"},
	}

	for _, tt := range tts {
		lilypondchord := ChordToLilypond(tt.chord)
		assert.Equal(t, tt.lilypondchord, lilypondchord)
	}
}
