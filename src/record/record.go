package record

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/schollz/miti/src/log"
	"github.com/schollz/miti/src/midi"
	"github.com/schollz/miti/src/music"
)

func Record() (err error) {
	log.Debug("init recording")
	finished := make(chan bool)

	// shutdown everything on Ctl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Debug(sig)
			log.Info("shutting down")
			finished <- true
			time.Sleep(200 * time.Millisecond)
		}
	}()

	events, err := midi.ReadAll(finished)
	if err != nil {
		return
	}

	currentState := ""
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		notes := []midi.Event{}
		for {
			select {
			case e := <-events:
				if !e.On {
					continue
				}
				notes = append(notes, e)
			case <-ticker.C:
				if len(notes) == 0 {
					continue
				}
				for _, e := range notes {
					log.Debugf("e: %+v", e)
					note := music.MidiToNote(e.MIDI)
					currentState += fmt.Sprintf("%s%d", note.Name, note.Octave)
				}
				currentState += " "
				fmt.Println("\r" + currentState)
				notes = []midi.Event{}
			}
		}
	}()

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press p to make new pattern")
	fmt.Println("Press m to make new measure")
	fmt.Println("Press ESC to quit")
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)
		if key == keyboard.KeyEsc {
			break
		}
		if char == rune('m') {
			fmt.Println("presed m!")
		}
	}
	finished <- true
	time.Sleep(500 * time.Millisecond)
	return
}
