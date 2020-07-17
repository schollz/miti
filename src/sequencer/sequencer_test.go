package sequencer

import (
	"fmt"
	"testing"
	"time"

	"github.com/kr/pretty"
	log "github.com/schollz/logger"
	"github.com/schollz/saps/src/music"
	"github.com/stretchr/testify/assert"
)

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

	s := New(func(s string, c music.Chord) {
		log.Tracef("%s %s", s, pretty.Sprint(c))
	})
	err := s.Parse(config)
	assert.Nil(t, err)
	fmt.Printf(pretty.Sprint(s))
	s.Start()
	time.Sleep(20 * time.Second)
	s.Stop()
	time.Sleep(1 * time.Second)
}
