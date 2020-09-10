package midi

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/miti/src/log"
	"github.com/schollz/miti/src/music"
	"github.com/schollz/miti/src/portmidi"
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

var Latency = int64(2000) // milliseconds

func Init() (devices []string, err error) {
	defer func() {
		if err == nil {
			inited = true
		}
	}()
	err = portmidi.Initialize()
	if err != nil {
		return
	}
	log.Debugf("found %d devices", portmidi.CountDevices())

	outputChannelsMap = make(map[string]int)
	outputChannelsMapMatch = make(map[string]int)
	for i := 0; i < portmidi.CountDevices(); i++ {
		di := portmidi.Info(portmidi.DeviceID(i))
		log.Debugf("device %d: '%s', i/o: %v/%v", i, di.Name, di.IsInputAvailable, di.IsOutputAvailable)
		if di.IsOutputAvailable && !strings.Contains(di.Name, "Wavetable Synth") {
			devices = append(devices, di.Name)

			// create a buffered channel for each instrument
			outputChannelsMap[di.Name] = len(outputChannels)
			outputChannels = append(outputChannels, make(chan music.Chord, 1000))
			// create a go-routine for each instrument
			go func(instrument string, deviceID int, channelNum int) {
				defer func() {
					if r := recover(); r != nil {
						log.Debug("recovered panic")
					}
				}()
				log.Debugf("[%s] opening stream with latency %d", instrument, Latency)
				outputStream, err := portmidi.NewOutputStream(portmidi.DeviceID(deviceID), 4096, Latency)
				if err != nil {
					panic(err)
				}
				log.Debugf("[%s] opened stream", instrument)
				midis := make([]int64, 100)
				velocities := make([]int64, 100)
				notesOn := make(map[int64]bool)
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
								outputStream.WriteShort(0x80, note, 0)
							}
						}
						if chord.Notes[0].MIDI == -2 {
							// shutdown
							outputStream.Close()
							channelLock.Unlock()
							return
						}
						channelLock.Unlock()
					}
					j := 0
					for _, n := range chord.Notes {
						midis[j] = int64(n.MIDI)
						if onState, ok := notesOn[midis[j]]; ok {
							if onState && chord.On {
								// this note already has this state
								log.Tracef("already played")
								continue
							}
						}
						notesOn[midis[j]] = chord.On
						velocities[j] = 100
						j++
					}
					if j == 0 {
						continue
					}
					channelLock.Lock()
					if chord.On {
						err = outputStream.WriteShorts(0x90, midis[:j], velocities[:j])
					} else {
						err = outputStream.WriteShorts(0x80, midis[:j], velocities[:j])
					}
					channelLock.Unlock()
					if err != nil {
						log.Errorf("[%s]: %s, could not send: %+v", instrument, err.Error(), midis[:j])
					} else {
						log.Tracef("[%s]: wrote %+v", instrument, midis[:j])
					}
				}
			}(di.Name, i, len(outputChannels)-1)
			if err != nil {
				err = fmt.Errorf("could not get output from: '%s'", di.Name)
				return
			}
		}
	}

	return
}

func Shutdown() (err error) {
	inited = false
	for out := range outputChannels {
		outputChannels[out] <- music.Chord{Notes: []music.Note{music.Note{MIDI: -2}}, On: false}
	}
	return portmidi.Terminate()
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
	log.Trace("got emit")
	outputChannels[channelID] <- chord
	log.Trace("emitted")
	// log.Trace("building midi")
	// midis := make([]int64, len(chord.Notes))
	// velocities := make([]int64, len(chord.Notes))
	// for i, n := range chord.Notes {
	// 	midis[i] = int64(n.MIDI)
	// 	encounteredNotes[midis[i]] = struct{}{}
	// 	velocities[i] = 100
	// }
	// log.Trace("sending midi")
	// if chord.On {
	// 	log.Tracef("[%s] %+v", msg, midis)
	// 	err = outputStreams[msg].WriteShorts(0x90, midis, velocities)
	// } else {
	// 	err = outputStreams[msg].WriteShorts(0x80, midis, velocities)
	// }
	// log.Trace("sent")
	// if err != nil {
	// 	return
	// }
	return
}
