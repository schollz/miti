package main

import (
	"flag"

	"github.com/schollz/miti/src/log"
	"github.com/schollz/miti/src/play"
)

var flagDebug, flagTrace bool
var flagFile string

func init() {
	flag.BoolVar(&flagDebug, "debug", false, "debug")
	flag.BoolVar(&flagTrace, "trace", false, "trace")
	flag.StringVar(&flagFile, "file", "", "file to load")
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

	err := play.Play(flagFile)
	if err != nil {
		log.Error(err)
	}

}
