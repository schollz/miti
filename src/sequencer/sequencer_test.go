package sequencer

import (
	"testing"
	"time"

	log "github.com/schollz/logger"
)

func TestSequencer(t *testing.T) {
	log.SetLevel("trace")
	s := New()
	s.Start()
	time.Sleep(3 * time.Second)
	s.UpdateTempo(120)
	time.Sleep(3 * time.Second)
	s.Stop()
	time.Sleep(1 * time.Second)
}
