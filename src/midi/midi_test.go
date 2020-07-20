package midi

import (
	"testing"
	"time"

	"github.com/schollz/miti/src/log"
)

func TestMIDI(t *testing.T) {
	_, err := Init()
	if err != nil {
		t.Errorf("err: %s", err.Error())
	}
}

func TestIn(t *testing.T) {
	log.SetLevel("trace")
	finished := make(chan bool)
	events, err := ReadAll(finished)
	if err != nil {
		t.Errorf("err: %s", err.Error())
	}
	go func() {
		for {
			e := <-events
			log.Debugf("e: %+v", e)
		}
	}()
	time.Sleep(2 * time.Second)
	log.Debug("sending finish")
	finished <- true
	log.Debug("sent")
	time.Sleep(100 * time.Millisecond)
}

// func TestPlay(t *testing.T) {
// 	_, err := Init()
// 	assert.Nil(t, err)
// 	assert.Nil(t, Midi("NTS-1 digital kit 1 SOUND", music.Chord{On: true, Notes: []music.Note{music.Note{MIDI: 81}}}))
// 	time.Sleep(2 * time.Second)
// 	assert.Nil(t, Midi("NTS-1 digital kit 1 SOUND", music.Chord{On: false, Notes: []music.Note{music.Note{MIDI: 81}}}))
// }
