package rtmidi

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/schollz/miti/src/log"
	"github.com/schollz/miti/src/music"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/rtmididrv"
)

// outputChannelsMap keeps track of channels
var outputChannelsMap map[string]int

var outputChannelsMapMatch map[string]int

// outputChannels allows global access to channels
var outputChannels []chan music.Chord

// channelLock is nessecary for Linux systems to
// only transmit one thing at a time
var channelLock sync.Mutex
var inited bool
var drv *driver.Driver

func Init() (devices []string, err error) {
	defer func() {
		if err == nil {
			inited = true
		}
	}()

	drv, err = driver.New()
	if err != nil {
		return
	}

	outs, err := drv.Outs()
	if err != nil {
		return
	}
	log.Debugf("found %d devices", len(outs))

	outputChannelsMap = make(map[string]int)
	outputChannelsMapMatch = make(map[string]int)
	for _, out := range outs {
		log.Debugf("device: %s", out.String())
		// create a buffered channel for each instrument
		outputChannelsMap[out.String()] = len(outputChannels)
		outputChannels = append(outputChannels, make(chan music.Chord, 1000))
		// create a go-routine for each instrument
		go func(instrument string, channelNum int, midio midi.Out) {
			log.Debugf("[%s] opening stream", instrument)
			err := midio.Open()
			if err != nil {
				panic(err)
			}
			log.Debugf("[%s] making writer", instrument)
			wr := writer.New(midio)
			log.Debugf("[%s] opened stream", instrument)
			notesOn := make(map[uint8]bool)
			for {
				log.Tracef("[%s] waiting for chord", instrument)
				chord := <-outputChannels[channelNum]
				log.Tracef("[%s] got chord: %+v", instrument, chord)

				// special things
				// midi note -1 turns off all on notes
				// midi note -2 turns off all on notes and shuts down
				if chord.Notes[0].MIDI < 0 {
					// turn off all notes
					channelLock.Lock()
					for note := range notesOn {
						if notesOn[note] {
							writer.NoteOff(wr, uint8(note))
						}
					}
					if chord.Notes[0].MIDI == -2 {
						// shutdown
						midio.Close()
						channelLock.Unlock()
						go func() {
							time.Sleep(50 * time.Millisecond)
							drv.Close()
						}()
						return
					}
					channelLock.Unlock()
				}
				channelLock.Lock()
				for _, n := range chord.Notes {
					midinote := uint8(n.MIDI)
					if onState, ok := notesOn[midinote]; ok {
						if onState && chord.On {
							// this note already has this state
							log.Tracef("already played")
							continue
						}
					}
					if chord.On {
						writer.NoteOn(wr, midinote, 100)
					} else {
						writer.NoteOff(wr, midinote)
					}
				}
				channelLock.Unlock()
			}
		}(out.String(), len(outputChannels)-1, out)
	}
	time.Sleep(3 * time.Second)

	return
}

func Shutdown() (err error) {
	inited = false
	for out := range outputChannels {
		outputChannels[out] <- music.Chord{Notes: []music.Note{music.Note{MIDI: -2}}, On: false}
	}
	return drv.Close()
}

func NotesOff() (err error) {
	for out := range outputChannels {
		outputChannels[out] <- music.Chord{Notes: []music.Note{music.Note{MIDI: -1}}, On: false}
	}
	return
}

func Midi(msg string, chord music.Chord) (err error) {
	if !inited {
		err = fmt.Errorf("not initialized")
		return
	}
	if len(chord.Notes) == 0 {
		return
	}
	channelID, ok := outputChannelsMap[msg]
	if !ok {
		channelID, ok = outputChannelsMapMatch[msg]
		if !ok {
			found := false
			for m := range outputChannelsMap {
				if strings.Contains(strings.ToLower(m), strings.ToLower(msg)) {
					outputChannelsMapMatch[msg] = outputChannelsMap[m]
					found = true
					log.Infof("mapping '%s' -> '%s'", msg, m)
				}
			}
			if !found {
				err = fmt.Errorf("no such device: %s", msg)
				return
			}
		}
	}
	outputChannels[channelID] <- chord
	return
}
