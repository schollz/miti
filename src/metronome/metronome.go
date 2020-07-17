package metronome

import (
	"time"

	log "github.com/schollz/logger"
)

const PULSES_PER_QUARTER_NOTE = 24.0
const QUARTER_NOTES_PER_MEASURE = 4.0

type Metronome struct {
	quarterNotePerMeasure float64
	tempo                 int
	pulse                 float64
	sections              []float64
	update                chan bool
	stop                  chan bool
	on                    bool
	stepemit              func(int)
}

func New(stepemit func(int)) (m *Metronome) {
	m = new(Metronome)
	m.tempo = 60
	m.quarterNotePerMeasure = 4
	m.sections = []float64{4}
	m.update = make(chan bool)
	m.stop = make(chan bool)
	m.stepemit = stepemit
	return
}

func (m *Metronome) Start() {
	if m.on {
		log.Debug("metronome already running")
		return
	}
	log.Debug("starting metronome")
	m.on = true
	m.pulse = -1
	go func() {
		ticker := time.NewTicker(time.Duration(1000000*60/m.tempo/PULSES_PER_QUARTER_NOTE) * time.Microsecond)
		log.Tracef("ticker time: %+v", time.Duration(1000000*60/m.tempo/PULSES_PER_QUARTER_NOTE)*time.Microsecond)

		for {
			select {
			case <-ticker.C:
				m.pulse++
				if m.pulse == PULSES_PER_QUARTER_NOTE*QUARTER_NOTES_PER_MEASURE {
					m.pulse = 0
				}
				m.stepemit(int(m.pulse))
			case <-m.update:
				ticker.Stop()
				ticker = time.NewTicker(time.Duration(1000000*60/m.tempo/PULSES_PER_QUARTER_NOTE) * time.Microsecond)
			case <-m.stop:
				ticker.Stop()
				log.Debug("..ticker stopped!")
				m.on = false
				return
			}
		}
	}()
	return
}

func (m *Metronome) Stop() {
	m.stop <- true
}

func (m *Metronome) UpdateTempo(tempo int) {
	if tempo <= 0 || tempo == m.tempo {
		return
	}
	m.tempo = tempo
	if m.on {
		m.update <- true
	}
}
