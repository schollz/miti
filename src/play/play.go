package play

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/miti/src/log"
	"github.com/schollz/miti/src/music"
	"github.com/schollz/miti/src/sequencer"
	"github.com/skratchdot/open-golang/open"

	// midi "github.com/schollz/miti/src/rtmidi" // use rtmidi instead
	midi "github.com/schollz/miti/src/midi"
)

var Version string
var SyncWithMidi = false
var ClickTrack = false
var Latency = int64(0)

func Play(mitiFile string, justShowDevices bool) (err error) {
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

	if justShowDevices {
		return
	}

	if mitiFile == "" {
		// generate a default miti file
		mitiFile, _ = filepath.Abs("miti.txt")
		_, erre := os.Stat(mitiFile)
		if os.IsNotExist(erre) {
			f, _ := os.Create(mitiFile)
			f.WriteString(`# welcome to miti (` + Version + `)!
# this is your miti file: ` + mitiFile + ` (this file).
# modify this file and save to update the sequencing on your instruments.
# <- lines beginning with "#" are comments. feel free to delete them.

# use chain to chain together patterns.
# (the following chain loops pattern 1 twice, then pattern 2 once).
chain 1 1 2

# adjust tempo.
tempo 60

# define a pattern.
# this pattern is named "1", but you can use any name.
pattern 1

# define instruments for pattern.
# multiple instruments separated commas.
# (your available instruments are already listed).
instruments ` + strings.ToLower(strings.Join(devices, ", ")) + `

# choose legato of the notes.
legato 100

# add in notes.
# notes are subidivided by number of spaces.
# each line is one measure.
CEG
ACE 
FAC 
GBD

# you can add other instruments in this pattern here, for example:
# instruments instrument2
# C E G C E G

# define another pattern.
pattern 2
instruments ` + strings.ToLower(strings.Join(devices, ", ")) + `
legato 50
E B G E B G E B G E B G 
`)
			f.Close()
		}

		open.Run(mitiFile)
	}

	if len(mitiFile) == 0 {
		return
	}

	playDone := make(chan bool)
	watcherDone := make(chan bool)
	shutdownInitiated := false

	startTime := time.Now()
	seq := sequencer.New(ClickTrack, Latency, func(s string, c music.Chord) {
		if shutdownInitiated {
			return
		}
		errMidi := midi.Midi(s, c)
		if errMidi != nil {
			log.Trace(errMidi)
		} else {
			log.Infof("%2.5f [%s] emitting %+v", time.Since(startTime).Seconds(), s, c)
		}
	})

	// shutdown everything on Ctl+C
	finished := make(chan bool)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			shutdownInitiated = true
			log.Debug(sig)
			log.Info("shutting down")
			watcherDone <- true
			time.Sleep(50 * time.Millisecond)
			go seq.Stop()
			time.Sleep(50 * time.Millisecond)
			midi.Shutdown()
			time.Sleep(50 * time.Millisecond)
			playDone <- true
			finished <- true
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

	if SyncWithMidi {
		events, errR := midi.ReadAll(finished)
		if errR != nil {
			err = errR
			return
		}
		for {
			e := <-events
			if !e.On {
				break
			}
		}
		finished <- true
	}
	log.Info("playing")
	seq.Start()
	<-playDone
	log.Info("done playing")
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
