package midi

import (
	"testing"
	"time"

	"github.com/schollz/saps/src/music"
	"github.com/stretchr/testify/assert"
)

func TestMIDI(t *testing.T) {
	assert.Nil(t, Init())
}

func TestPlay(t *testing.T) {
	assert.Nil(t, Init())
	assert.Nil(t, Midi("NTS-1 digital kit 1 SOUND", music.Chord{On: true, Notes: []music.Note{music.Note{MIDI: 81}}}))
	time.Sleep(2 * time.Second)
	assert.Nil(t, Midi("NTS-1 digital kit 1 SOUND", music.Chord{On: false, Notes: []music.Note{music.Note{MIDI: 81}}}))
}
