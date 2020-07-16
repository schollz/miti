package sequencer

import (
	log "github.com/schollz/logger"
	"github.com/schollz/midi-sequencer/src/metronome"
)

const PULSES_PER_QUARTER_NOTE = 24.0

type Sequencer struct {
	metronome *metronome.Metronome
}

func New() (s *Sequencer) {
	s = new(Sequencer)
	s.metronome = metronome.New(s.Emit)
	return
}

func (s *Sequencer) Start() {
	s.metronome.Start()
}

func (s *Sequencer) Stop() {
	s.metronome.Stop()
}

func (s *Sequencer) UpdateTempo(tempo int) {
	s.metronome.UpdateTempo(tempo)
}

func (s *Sequencer) Emit(section int, measure int, beat int) {
	log.Debug(section, measure, beat)
}
