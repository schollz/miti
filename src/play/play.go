package play

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/olekukonko/tablewriter"
	log "github.com/schollz/logger"
	"github.com/schollz/saps/src/midi"
	"github.com/schollz/saps/src/music"
	"github.com/schollz/saps/src/sequencer"
)

func PrintDevices() (err error) {
	devices, err := midi.Init()
	if err != nil {
		return
	}
	if len(devices) == 0 {
		err = fmt.Errorf("no devices detected, try plugging some")
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

func Play() (err error) {
	// show devices
	err = PrintDevices()
	if err != nil {
		return
	}

	// start sequencer with midi equipped
	seq := sequencer.New(func(s string, c music.Chord) {
		if c.On && len(c.Notes) > 0 {
			log.Infof("[%.5s] %s", s, c)
		}
		errMidi := midi.Midi(s, c)
		if errMidi != nil {
			log.Error(errMidi)
		}
	})

	// shutdown everything on Ctl+C
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

	// start tempo

	err = seq.Parse(`section a

	tempo 120
	instruments NTS-1 digital kit 1 SOUND
	A C E A C E A C A C E A C E A C   
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
