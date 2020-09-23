package rtmidi

import (
	"time"

	"github.com/schollz/miti/src/log"
	driver "gitlab.com/gomidi/rtmididrv"
)

type Event struct {
	MIDI       int
	On         bool
	Instrument string
	Timestamp  time.Time
}

func ReadAll(finished chan bool) (events chan Event, err error) {
	events = make(chan Event, 100)
	if !inited {
		drv, err = driver.New()
		if err != nil {
			return
		}
		inited = true
	}

	ins, err := drv.Ins()
	if err != nil {
		return
	}

	for i := range ins {
		err = ins[i].Open()
		if err != nil {
			log.Error(err)
			continue
		}
		func(j int) {
			name := ins[j].String()
			log.Tracef("setting up %s", name)
			ins[j].SetListener(func(data []byte, deltaMicroseconds int64) {
				if len(data) == 3 {
					log.Tracef("[%s] %d %+v", name, data)
					events <- Event{int(data[1]), data[0] == 144, name, time.Now()}
				}
			})
		}(i)

	}
	go func() {
		<-finished
		for i := range ins {
			log.Debugf("[%s] closing midi input", ins[i].String())
			ins[i].StopListening()
			ins[i].Close()
		}
	}()

	return
}
