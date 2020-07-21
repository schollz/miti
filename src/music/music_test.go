package music

import (
	"fmt"
	"testing"
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
