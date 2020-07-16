package sequencer

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

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
	s := `section a

 instruments op-1, sh01a
 CEG
 ACE
 CEGC
 ACEA
 
 instruments nts-1
 C E G E C E G E
 A C E C A C E C
 
 section b 
 
 instruments op-1
 CEG 
 GBD`

	sections, err := Parse(s)
	assert.Nil(t, err)
	b, _ := json.MarshalIndent(sections, "", " ")
	fmt.Println(string(b))
}
