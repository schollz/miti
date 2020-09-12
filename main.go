package main

//go:generate git tag -af v$VERSION -m "v$VERSION"
//go:generate go run .github/updateversion.go
//go:generate git commit -am "bump $VERSION"
//go:generate git tag -af v$VERSION -m "v$VERSION"

import (
	"flag"
	"fmt"

	"github.com/schollz/miti/src/click"
	"github.com/schollz/miti/src/log"
	"github.com/schollz/miti/src/midi"
	"github.com/schollz/miti/src/play"
	"github.com/schollz/miti/src/record"
)

var flagDebug, flagTrace, flagVersion, flagList, flagWait, flagClick bool
var flagFile, flagRecord string
var flagLatency int64

// Version specifies the version
var Version string

func init() {
	flag.BoolVar(&flagDebug, "debug", false, "debug")
	flag.BoolVar(&flagTrace, "trace", false, "trace")
	flag.BoolVar(&flagList, "list", false, "list midi devices")
	flag.BoolVar(&flagVersion, "version", false, "show version")
	flag.BoolVar(&flagWait, "sync", false, "wait for midi input to start")
	flag.BoolVar(&flagClick, "click", false, "output click track with metronome")
	flag.Int64Var(&flagLatency, "latency", 2000, "latency for midi output")
	flag.Int64Var(&click.TuneLatency, "clicklag", 0, "add lag to click track to sync better")
	flag.StringVar(&flagRecord, "record", "", "record input to miti file")
	flag.StringVar(&flagFile, "play", "", "play sequence from miti file")
	if Version == "" {
		Version = "v0.5.0-6f4d05a"
	}
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
	fmt.Printf("miti %s - musical instrument textual interface\nsubmit bugs to https://github.com/schollz/miti/issues\n\n", Version)
	if flagVersion {
		return
	}
	midi.Latency = flagLatency
	play.Version = Version
	play.SyncWithMidi = flagWait
	play.ClickTrack = flagClick
	play.Latency = flagLatency
	var err error
	if flagRecord != "" {
		err = record.Record(flagRecord)
	} else if flagList {
		err = play.Play("", true)
	} else {
		err = play.Play(flagFile, false)
	}
	if err != nil {
		log.Error(err)
	}

}
