package main

import (
	"flag"

	log "github.com/schollz/logger"
	"github.com/schollz/midi-sequencer/src/midi"
	"github.com/schollz/midi-sequencer/src/server"
)

var flagDebug, flagTrace bool

func init() {
	flag.BoolVar(&flagDebug, "debug", false, "debug")
	flag.BoolVar(&flagTrace, "trace", false, "debug")
}

func main() {
	flag.Parse()
	if flagDebug {
		log.SetLevel("debug")
	} else if flagTrace {
		log.SetLevel("trace")
	} else {
		log.SetLevel("info")
	}
	midi.Init()
	err := server.Run()
	if err != nil {
		log.Error(err)
	}

}
