package sequencer

import (
	"strings"

	log "github.com/schollz/logger"
	"github.com/schollz/midi-sequencer/src/metronome"
	"github.com/schollz/midi-sequencer/src/music"
)

const PULSES_PER_QUARTER_NOTE = 24.0
const QUARTERNOTES_PER_MEASURE = 4.0

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
	Emit   map[int][]music.Chord
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

func (s *Sequencer) Emit(section int, measure int, beat int, pulse int) {
	log.Trace(section, measure, beat, beat*int(PULSES_PER_QUARTER_NOTE)+pulse)
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
			measure := Measure{Emit: make(map[int][]music.Chord)}
			fs := strings.Fields(line)
			for i, cluster := range fs {
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
				startPulse := float64(i) / float64(len(fs)) * (QUARTERNOTES_PER_MEASURE*PULSES_PER_QUARTER_NOTE - 1)
				endPulse := startPulse + 1/float64(len(fs))*(QUARTERNOTES_PER_MEASURE*PULSES_PER_QUARTER_NOTE-1)
				// TODO: add in legato
				if _, ok := measure.Emit[int(startPulse)]; !ok {
					measure.Emit[int(startPulse)] = []music.Chord{}
				}
				measure.Emit[int(startPulse)] = append(measure.Emit[int(startPulse)], music.Chord{Notes: notes, On: true})
				if _, ok := measure.Emit[int(endPulse)]; !ok {
					measure.Emit[int(endPulse)] = []music.Chord{}
				}
				measure.Emit[int(endPulse)] = append(measure.Emit[int(endPulse)], music.Chord{Notes: notes, On: false})

			}
			part.Measures = append(part.Measures, measure)
		}
	}
	return
}
