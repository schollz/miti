package rtmidi

import (
	"testing"

	"github.com/schollz/miti/src/log"
)

func TestMIDI(t *testing.T) {
	log.SetLevel("trace")
	_, err := Init()
	if err != nil {
		t.Errorf("err: %s", err.Error())
	}
}
