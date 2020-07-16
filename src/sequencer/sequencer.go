package sequencer

import (
	"math"
	"time"

	log "github.com/schollz/logger"
)

const PULSES_PER_QUARTER_NOTE = 24.0

type Sequencer struct {
	metronome *Metronome
}

type Metronome struct {
	quarterNotePerMeasure         float64
	tempo                         int
	pulse, beat, measure, section float64
	sections                      []float64
	update                        chan bool
	stop                          chan bool
	on                            bool
}

func New() (s *Sequencer) {
	s = new(Sequencer)
	s.metronome = new(Metronome)
	s.metronome.tempo = 60
	s.metronome.quarterNotePerMeasure = 4
	s.metronome.sections = []float64{4}
	s.metronome.update = make(chan bool)
	s.metronome.stop = make(chan bool)
	return
}

func (s *Sequencer) Start() {
	s.metronome.Start()
}

func (s *Sequencer) Stop() {
	s.metronome.Stop()
}

func (s *Metronome) Stop() {
	s.stop <- true
}

func (s *Metronome) Start() {
	if s.metronome_on {
		log.Debug("metronome already running")
		return
	}
	s.metronome_on = true
	go func() {
		ticker := time.NewTicker(time.Duration(1000*60/s.tempo) * time.Millisecond)

		for {
			select {
			case <-ticker.C:
				log.Trace("tick")
				go s.metronome.step()
			case <-s.update:
				ticker.Stop()
				ticker = time.NewTicker(time.Duration(1000*60/s.tempo) * time.Millisecond)
			case <-s.stop:
				ticker.Stop()
				log.Debug("..ticker stopped!")
				s.metronome_on = false
				return
			}
		}
	}()
	return
}

func (s *Metronome) step() {
	s.pulse++
	if s.pulse == PULSES_PER_QUARTER_NOTE {
		s.pulse = 0
		s.beat++
		if s.beat == s.quarterNotePerMeasure {
			s.beat = 0
			s.measure++
		}
		if s.measure == s.sections[s.section] {
			s.section++
			s.section = math.Mod(s.section, float64(len(s.sections)))
			s.measure = 0
		}
		log.Tracef("%2.0f %2.0f %2.0f", s.section, s.measure, s.beat)
	}
}

func (s *Metronome) UpdateTempo(tempo int) {
	if tempo <= 0 {
		return
	}
	s.tempo = tempo
	s.update <- true
}
