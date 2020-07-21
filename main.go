package main

import (
	"flag"

	"github.com/schollz/miti/src/log"
	"github.com/schollz/miti/src/play"
	"github.com/schollz/miti/src/record"
)

var flagDebug, flagTrace bool
var flagFile, flagRecord string

func init() {
	flag.BoolVar(&flagDebug, "debug", false, "debug")
	flag.BoolVar(&flagTrace, "trace", false, "trace")
	flag.StringVar(&flagRecord, "record", "", "record input to miti file")
	flag.StringVar(&flagFile, "play", "", "play sequence from miti file")
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

	var err error
	if flagRecord != "" {
		err = record.Record(flagRecord)
	} else {
		err = play.Play(flagFile)
	}
	if err != nil {
		log.Error(err)
	}

}
