package sequencer

import (
	"strings"

	log "github.com/schollz/logger"
	"github.com/schollz/midi-sequencer/src/metronome"
	"github.com/schollz/midi-sequencer/src/music"
)

const PULSES_PER_QUARTER_NOTE = 24.0

type Sequencer struct {
	metronome *metronome.Metronome
	Sections  []Section
}

type Section struct {
	Name  string
	Parts []Part
}

type Part struct {
	Instruments []string
	Measures    []Measure
}

type Measure struct {
	Chords []music.Chord
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
	log.Trace(section, measure, beat)
}

func Parse(s string) (sections []Section, err error) {
	isPart := false

	var section Section
	var part Part
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		log.Debug(line)
		if strings.HasPrefix(line, "section") {
			if len(part.Instruments) > 0 {
				section.Parts = append(section.Parts, part)
			}
			if isPart {
				sections = append(sections, section)
			}
			section = Section{Name: line}
			isPart = false
		} else if strings.HasPrefix(line, "instruments") {
			if len(part.Instruments) > 0 {
				section.Parts = append(section.Parts, part)
			}
			isPart = true
			line = strings.TrimPrefix(line, "instruments")
			line = strings.TrimPrefix(line, "instrument")
			part = Part{Instruments: strings.Split(line, ",")}
		} else if len(line) > 0 {
			measure := Measure{}
			fs := strings.Fields(line)
			for _, cluster := range fs {
				if cluster == "." {
					continue
				}
				var notes []music.Note
				if len(measure.Chords) > 0 {
					lastChord := measure.Chords[len(measure.Chords)-1]
					lastNote := lastChord.Notes[len(lastChord.Notes)-1]
					notes, err = music.ParseCluster(cluster, lastNote)
				} else {
					notes, err = music.ParseCluster(cluster)
				}
				if err != nil {
					log.Error(err)
					return
				}
				measure.Chords = append(measure.Chords, music.Chord{Notes: notes})
			}
			part.Measures = append(part.Measures, measure)
		}
	}
	return
}
