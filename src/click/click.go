package click

import (
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	log "github.com/schollz/logger"
)

var sampleNum = 0.0
var pulseWidth = 2400.0  // microseconds
var sampleRate = 44100.0 // hz
var periodTime = 1.0     // seconds

func SetBPM(bpm float64) {
	periodTime = 60 / bpm
}

func Click() beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			sample := 0.0
			if sampleNum < (sampleRate * pulseWidth / 1000000) {
				sample = 1
			}
			samples[i][0] = sample
			samples[i][1] = sample
			sampleNum++
			if sampleNum > sampleRate*periodTime {
				sampleNum = 0
			}
		}
		return len(samples), true
	})
}

func Play(bpm float64) {
	sr := beep.SampleRate(int(sampleRate))
	speaker.Init(sr, sr.N(time.Second/10))
	log.Infof("starting click track at %f", bpm)
	SetBPM(bpm)
	speaker.Play(Click())
}

func Stop() {
	log.Infof("closing click track")
	speaker.Clear()
	speaker.Close()
}
