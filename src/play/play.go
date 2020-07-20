package play

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/olekukonko/tablewriter"
	log "github.com/schollz/logger"
	"github.com/schollz/miti/src/midi"
	"github.com/schollz/miti/src/music"
	"github.com/schollz/miti/src/sequencer"
)

func PrintDevices() (err error) {
	devices, err := midi.Init()
	if err != nil {
		return
	}
	if len(devices) == 0 {
		fmt.Println(`+-------------------+
| NO INSTRUMENTS :( |
+-------------------+`)
		return
	}
	data := [][]string{}
	for _, device := range devices {
		data = append(data, []string{device})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"INSTRUMENTS"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
	return
}

func Play(mitiFile string) (err error) {
	// show devices
	err = PrintDevices()
	if err != nil {
		return
	}

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
	seq := sequencer.New(func(s string, c music.Chord) {
		if shutdownInitiated {
			return
		}
		log.Tracef("[%s] forwarding emit", s)
		errMidi := midi.Midi(s, c)
		if errMidi != nil {
			log.Debug(errMidi)
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
