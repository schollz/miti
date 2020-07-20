package midi

import (
	"strings"
	"time"

	"github.com/schollz/miti/src/log"
	"github.com/schollz/miti/src/portmidi"
)

type Event struct {
	MIDI       int
	On         bool
	Instrument string
	Timestamp  time.Time
}

func ReadAll(finished chan bool) (events chan Event, err error) {
	events = make(chan Event, 100)
	err = portmidi.Initialize()
	if err != nil {
		return
	}

	go func() {
		numDevices := 0
		log.Debugf("found %d devices", numDevices)
		done := make(chan bool)
		for i := 0; i < portmidi.CountDevices(); i++ {
			di := portmidi.Info(portmidi.DeviceID(i))
			log.Tracef("di: %+v", di)
			if strings.Contains(di.Name, "Wavetable Synth") {
				continue
			}
			numDevices++
			if di.IsInputAvailable {
				log.Tracef("%s input available", di.Name)
				go func(id int, events chan Event, done chan bool) {
					d := portmidi.Info(portmidi.DeviceID(id))
					log.Debugf("reading from %s", d.Name)
					in, err := portmidi.NewInputStream(portmidi.DeviceID(id), 1024)
					if err != nil {
						panic(err)
					}
					defer in.Close()

					ch := in.Listen()
					for {
						select {
						case event := <-ch:
							if event.Data1 == 0 {
								continue
							}
							log.Tracef("event: %+v", event)
							events <- Event{int(event.Data1), event.Status == 144, d.Name, time.Now()}
						case <-done:
							log.Debug("closing midi read")
							return
						}
					}
				}(i, events, done)
			}
		}

		<-finished
		for i := 0; i < numDevices; i++ {
			done <- true
		}
		time.Sleep(100 * time.Millisecond)
		log.Debug("terminating")
		portmidi.Terminate()
	}()
	return
}
