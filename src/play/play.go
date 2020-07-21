package play

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/schollz/miti/src/log"
	"github.com/schollz/miti/src/music"
	"github.com/schollz/miti/src/sequencer"

	// midi "github.com/schollz/miti/src/rtmidi" // use rtmidi instead
	midi "github.com/schollz/miti/src/midi"
)

func Play(mitiFile string) (err error) {
	// show devices
	devices, err := midi.Init()
	if err != nil {
		return
	}
	if len(devices) == 0 {
		err = fmt.Errorf("no devices found")
		return
	}
	fmt.Println("Available MIDI devices:")
	for _, device := range devices {
		fmt.Printf("- %s\n", strings.ToLower(device))
	}
	fmt.Println(" ")

	if len(mitiFile) == 0 {
		return
	}

	playDone := make(chan bool)
	watcherDone := make(chan bool)
	shutdownInitiated := false

	// start sequencer with midi equipped
	// if log.GetLevel() == "info" {
	// 	tm.Clear()
	// }
	startTime := time.Now()
	seq := sequencer.New(func(s string, c music.Chord) {
		if shutdownInitiated {
			return
		}
		log.Debugf("%2.5f [%s] emitting %+v", time.Since(startTime).Seconds(), s, c)
		errMidi := midi.Midi(s, c)
		if errMidi != nil {
			log.Trace(errMidi)
		}
		// if log.GetLevel() == "info" {
		// 	tm.MoveCursor(1, 1)
		// 	tm.Printf("%s: %+v\n\n", s, c)
		// 	tm.Flush()
		// }
	})

	// shutdown everything on Ctl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			shutdownInitiated = true
			log.Debug(sig)
			log.Info("shutting down")
			go seq.Stop()
			time.Sleep(50 * time.Millisecond)
			midi.Shutdown()
			time.Sleep(50 * time.Millisecond)
			watcherDone <- true
			time.Sleep(50 * time.Millisecond)
			playDone <- true
		}
	}()

	// load miti file
	err = seq.Parse(mitiFile)
	if err != nil {
		return
	}

	// hot-reload file
	go func() {
		err = hotReloadFile(seq, mitiFile, watcherDone)
		if err != nil {
			log.Error(err)
		}
	}()

	log.Info("playing")
	seq.Start()
	<-playDone
	return
}

func hotReloadFile(seq *sequencer.Sequencer, fname string, watcherDone chan bool) (err error) {
	ticker := time.NewTicker(700 * time.Millisecond)
	lastInfo := ""
	for {
		select {
		case <-watcherDone:
			return
		case t := <-ticker.C:
			log.Tracef("checking file at %s", t)
			var statinfo os.FileInfo
			statinfo, err = os.Stat(fname)
			if err != nil {
				log.Debug(err)
				continue
			}
			currentInfo := fmt.Sprintf("%s%d", statinfo.ModTime(), statinfo.Size())
			if lastInfo != "" && lastInfo != currentInfo {
				err = seq.Parse(fname)
				if err != nil {
					log.Warnf("problem hot-reloading %s: %s", fname, err.Error())
				} else {
					midi.NotesOff()
				}
			}
			lastInfo = currentInfo
		}

	}
	return
}
