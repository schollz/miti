package click

import (
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/schollz/miti/src/log"
)

var TuneLatency = int64(100)
var sampleNum = 0.0
var pulseWidth = 4800.0  // microseconds
var sampleRate = 44100.0 // hz
var periodTime = 1.0     // seconds
var activated = false
var activate chan bool

func SetBPM(bpm float64) {
	periodTime = 60 / bpm
}

func Click(latency ...int64) {
	go func() {
		if TuneLatency > 0 {
			time.Sleep(time.Duration(TuneLatency) * time.Millisecond)
		}
		if len(latency) > 0 && latency[0] > 0 {
			log.Tracef("activating click with latency %d", latency)
			time.Sleep(time.Duration(latency[0]) * time.Millisecond)
		}
		activate <- true
		log.Trace("activated click")
	}()
}

func click() beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			select {
			case <-activate:
				log.Tracef("clicking in %d samples", len(samples))
				activated = true
				sampleNum = 0
			default:
			}
			sample := 0.0
			window := (sampleRate * pulseWidth / 1000000)
			if sampleNum < window && activated {
				sample = 1
			}
			samples[i][0] = sample
			samples[i][1] = sample
			sampleNum++
			if sampleNum > sampleRate*periodTime {
				sampleNum = 0
				activated = false
			}
		}
		return len(samples), true
	})
}

func Play(bpm float64) {
	activate = make(chan bool, 10)
	sr := beep.SampleRate(int(sampleRate))
	speaker.Init(sr, sr.N(time.Second/400))
	log.Infof("starting click track at %f", bpm)
	SetBPM(bpm)
	speaker.Play(click())
}

func Stop() {
	log.Infof("closing click track")
	speaker.Clear()
	speaker.Close()
}

func Reset() {
	sampleNum = (sampleRate * pulseWidth / 1000000)
}
