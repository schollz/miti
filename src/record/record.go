package record

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/schollz/miti/src/log"
	"github.com/schollz/miti/src/midi"
	"github.com/schollz/miti/src/music"
)

func Record(fname string) (err error) {
	log.Debug("init recording")

	f, err := os.Create(fname)
	if err != nil {
		return
	}
	defer f.Close()

	finished := make(chan bool)

	events, err := midi.ReadAll(finished)
	if err != nil {
		return
	}

	patterns := 0
	currentPattern := ""
	currentState := ""
	previousNote := music.NewNote("C", 4)
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
				if currentPattern == "" {
					instrument := notes[0].Instrument
					fs := strings.Fields(instrument)
					if len(fs) > 2 {
						fs = fs[:2]
					}
					currentPattern += fmt.Sprintf("pattern %d\n", patterns)
					currentPattern += "instruments " + strings.ToLower(strings.Join(fs, " ")) + "\n"
					fmt.Println(currentPattern)
					f.WriteString(currentPattern)
				}

				for _, e := range notes {
					log.Debugf("e: %+v", e)
					note := music.MidiToNote(e.MIDI)
					closestNote := music.ClosestNote(note.Name, previousNote)
					if closestNote.Octave == note.Octave {
						currentState += fmt.Sprintf("%s", note.Name)
					} else {
						currentState += fmt.Sprintf("%s%d", note.Name, note.Octave)
					}
					previousNote = note
				}
				currentState += " "
				notes = []midi.Event{}

				fmt.Print("\r" + currentState)
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
	fmt.Println("Press backspace to delete last")
	fmt.Println("Press Ctl+C to quit")
	fmt.Println("---------------------------")
	for {
		char, key, err := keyboard.GetKey()
		log.Tracef("char: %v, key: %v", char, key)
		if err != nil {
			panic(err)
		}
		if key == 127 {
			fs := strings.Fields(currentState)
			fmt.Print("\r                  ")
			if len(fs) <= 1 {
				currentState = ""
			} else {
				fs = fs[:len(fs)-1]
				currentState = strings.Join(fs, " ") + " "
			}
			fmt.Printf("\r%s", currentState)
		} else if key == 3 {
			f.WriteString(currentState)
			fmt.Println("---------------------------")
			fmt.Printf("\n\nwrote to '%s'\n", fname)
			break
		}
		if char == rune('m') {
			f.WriteString(currentState)
			f.WriteString("\n")
			fmt.Print("\n")
			currentState = ""
		}
		if char == rune('p') {
			f.WriteString(currentState)
			f.WriteString("\n\n\n")
			fmt.Print("\n\n\n")
			currentState = ""
			currentPattern = ""
			patterns++
		}
	}
	finished <- true
	time.Sleep(500 * time.Millisecond)
	return
}
