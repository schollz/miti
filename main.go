package main

import (
	"flag"

	log "github.com/schollz/logger"
	"github.com/schollz/saps/src/play"
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

	err := play.Play()
	if err != nil {
		log.Error(err)
	}

}
