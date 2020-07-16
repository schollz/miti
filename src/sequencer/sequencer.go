package sequencer

import (
	"time"

	log "github.com/schollz/logger"
)

type Sequencer struct {
	Tempo           int
	metronome_tempo chan bool
	metronome_stop  chan bool
	metronome_on    bool
}

func New() (s *Sequencer) {
	s = new(Sequencer)
	s.Tempo = 60
	s.metronome_stop = make(chan bool)
	s.metronome_tempo = make(chan bool)
	return
}

func (s *Sequencer) Start() {
	if s.metronome_on {
		log.Debug("metronome already running")
		return
	}
	s.metronome_on = true
	go func() {
		ticker := time.NewTicker(time.Duration(1000*60/s.Tempo) * time.Millisecond)

		for {
			select {
			case <-ticker.C:
				log.Trace("tick")
			case <-s.metronome_tempo:
				ticker.Stop()
				ticker = time.NewTicker(time.Duration(1000*60/s.Tempo) * time.Millisecond)
			case <-s.metronome_stop:
				ticker.Stop()
				log.Debug("..ticker stopped!")
				s.metronome_on = false
				return
			}
		}
	}()
	return
}

func (s *Sequencer) Stop() {
	s.metronome_stop <- true
}

func (s *Sequencer) UpdateTempo(tempo int) {
	if tempo <= 0 {
		return
	}
	s.Tempo = tempo
	s.metronome_tempo <- true
}
