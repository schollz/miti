package sequencer

import (
	"fmt"
	"testing"
	"time"

	"github.com/kr/pretty"
	log "github.com/schollz/logger"
	"github.com/stretchr/testify/assert"
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

func TestParse(t *testing.T) {
	config := `section a

 instruments op-1, sh01a
 CEG
 ACE
 
 instruments nts-1
 C E
 
 section b 
 
 instruments op-1
 DF#A `

	s := New()
	err := s.Parse(config)
	assert.Nil(t, err)
	fmt.Printf(pretty.Sprint(s))
	s.Start()
	time.Sleep(20 * time.Second)
	s.Stop()
	time.Sleep(1 * time.Second)
}
