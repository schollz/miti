package play

import (
	"os"
	"os/signal"
	"time"

	log "github.com/schollz/logger"
	"github.com/schollz/saps/src/midi"
	"github.com/schollz/saps/src/music"
	"github.com/schollz/saps/src/sequencer"
)

func Play() (err error) {
	err = midi.Init()
	if err != nil {
		return
	}

	seq := sequencer.New(func(s string, c music.Chord) {
		errMidi := midi.Midi(s, c)
		if errMidi != nil {
			log.Error(errMidi)
		}
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Debug(sig)
			seq.Stop()
			err = midi.Shutdown()
			os.Exit(1)
		}
	}()

	seq.UpdateTempo(60)

	err = seq.Parse(`section a

	instruments NTS-1 digital kit 1 SOUND
	C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4
	C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4
	
	instruments Boutique SH-01A
	A3CE  
	C4EG 
	A3CE  
	C4EG 
	
	section b
		
	
	instruments Boutique SH-01A
	DF#A
	DF#A
	DF#A
	DF#A
	
 `)
	if err != nil {
		return
	}

	seq.Start()
	time.Sleep(5 * time.Hour)

	return
}
