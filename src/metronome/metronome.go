package metronome

import (
	"math"
	"time"

	log "github.com/schollz/logger"
)

const PULSES_PER_QUARTER_NOTE = 24.0

type Metronome struct {
	quarterNotePerMeasure         float64
	tempo                         int
	pulse, beat, measure, section float64
	sections                      []float64
	update                        chan bool
	stop                          chan bool
	on                            bool
	stepemit                      func(int, int, int)
}

func New(stepemit func(int, int, int)) (m *Metronome) {
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
	m.on = true
	go func() {
		ticker := time.NewTicker(time.Duration(1000*60/m.tempo) * time.Millisecond)

		for {
			select {
			case <-ticker.C:
				log.Trace("tick")
				go m.Step()
			case <-m.update:
				ticker.Stop()
				ticker = time.NewTicker(time.Duration(1000*60/m.tempo) * time.Millisecond)
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

func (m *Metronome) Step() {
	m.pulse++
	if m.pulse == PULSES_PER_QUARTER_NOTE {
		m.pulse = 0
		m.beat++
		if m.beat == m.quarterNotePerMeasure {
			m.beat = 0
			m.measure++
		}
		if m.measure == m.sections[int(m.section)] {
			m.section++
			m.section = math.Mod(m.section, float64(len(m.sections)))
			m.measure = 0
		}
		log.Tracef("%2.0f %2.0f %2.0f", m.section, m.measure, m.beat)
		m.stepemit(int(m.section), int(m.measure), int(m.beat))
	}
}

func (m *Metronome) Stop() {
	m.stop <- true
}

func (m *Metronome) UpdateTempo(tempo int) {
	if tempo <= 0 {
		return
	}
	m.tempo = tempo
	m.update <- true
}
