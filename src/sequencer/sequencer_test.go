package sequencer

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/kr/pretty"
	log "github.com/schollz/logger"
	"github.com/schollz/miti/src/music"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	log.SetLevel("trace")
	config := `
chain a a a a b
tempo 240 
pattern a

 instruments op-1
 CEG 

 pattern b 
 
 instruments op-1
 DF#A 
 `
	ioutil.WriteFile("temp.miti", []byte(config), 0644)
	defer func() {
		os.Remove("temp.miti")
	}()

	s := New(func(s string, c music.Chord) {
		log.Tracef("%s %s", s, pretty.Sprint(c))
	})
	err := s.Parse("temp.miti")
	assert.Nil(t, err)
	fmt.Printf(pretty.Sprint(s))
	s.Start()
	time.Sleep(20 * time.Second)
	s.Stop()
	time.Sleep(1 * time.Second)
}
