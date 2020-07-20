package midi

import (
	"github.com/schollz/miti/src/log"
	"github.com/schollz/miti/src/portmidi"
)

type Event struct {
	MIDI       int
	Instrument string
}

func ReadAll(finished chan bool) (events chan Event, err error) {
	events = make(chan Event, 100)
	err = portmidi.Initialize()
	if err != nil {
		return
	}

	go func() {
		numDevices := portmidi.CountDevices()
		log.Debugf("found %d devices", numDevices)
		done := make(chan bool)
		for i := 0; i < numDevices; i++ {
			di := portmidi.Info(portmidi.DeviceID(i))
			log.Tracef("di: %+v", di)
			if di.IsInputAvailable {
				go func(id int, events chan Event, done chan bool) {
					d := portmidi.Info(portmidi.DeviceID(i))
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
							log.Tracef("event: %+v", event)
							events <- Event{int(event.Data1), d.Name}
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
		portmidi.Terminate()
	}()
	return
}
