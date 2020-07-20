package midi

import (
	"testing"
)

func TestMIDI(t *testing.T) {
	_, err := Init()
	if err != nil {
		t.Errorf("err: %s", err.Error())
	}
}

// func TestPlay(t *testing.T) {
// 	_, err := Init()
// 	assert.Nil(t, err)
// 	assert.Nil(t, Midi("NTS-1 digital kit 1 SOUND", music.Chord{On: true, Notes: []music.Note{music.Note{MIDI: 81}}}))
// 	time.Sleep(2 * time.Second)
// 	assert.Nil(t, Midi("NTS-1 digital kit 1 SOUND", music.Chord{On: false, Notes: []music.Note{music.Note{MIDI: 81}}}))
// }
