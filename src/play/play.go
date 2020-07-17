package play

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/tidi/src/midi"
	"github.com/schollz/tidi/src/music"
	"github.com/schollz/tidi/src/sequencer"
	log "github.com/schollz/logger"
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

func Play(tidiFile string) (err error) {
	// show devices
	err = PrintDevices()
	if err != nil {
		return
	}

	if len(tidiFile) == 0 {
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
			log.Error(errMidi)
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

	// load tidi file
	err = seq.Parse(tidiFile)
	if err != nil {
		return
	}

	// hot-reload file
	go func() {
		err = hotReloadFile(seq, tidiFile, watcherDone)
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
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

	lastEvent := time.Now()
	done := make(chan bool)
	go func() {
		for {
			select {
			case _ = <-watcherDone:
				done <- true
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Debugf("event: %+v", event)
				if event.Op&fsnotify.Write == fsnotify.Write && time.Since(lastEvent).Seconds() > 1 {
					lastEvent = time.Now()
					log.Infof("reloading: %+v", event.Name)
					time.Sleep(100 * time.Millisecond)
					err = seq.Parse(fname)
					if err != nil {
						log.Warnf("problem hot-reloading %s: %s", fname, err.Error())
					} else {
						midi.NotesOff()
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error(err)
			}
		}
	}()

	err = watcher.Add(fname)
	if err != nil {
		return
	}
	<-done
	log.Debug("watcher done")
	return
}
