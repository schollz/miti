package midi

import (
	"github.com/rakyll/portmidi"
	log "github.com/schollz/logger"
)

func Init() {
	portmidi.Initialize()
	log.Debugf("found %d devices", portmidi.CountDevices())
	for i := 0; i < portmidi.CountDevices(); i++ {
		di := portmidi.Info(portmidi.DeviceID(i))
		log.Debugf("device %d: '%s', i/o: %v/%v", i, di.Name, di.IsInputAvailable, di.IsOutputAvailable)
	}
}
