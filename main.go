package main

//go:generate git tag -af v$VERSION -m "v$VERSION"
//go:generate go run .github/updateversion.go
//go:generate git commit -am "bump $VERSION"
//go:generate git tag -af v$VERSION -m "v$VERSION"

import (
	"flag"
	"fmt"
	"os"

	"github.com/schollz/miti/src/log"
	"github.com/schollz/miti/src/play"
	"github.com/schollz/miti/src/record"
	"github.com/skratchdot/open-golang/open"
)

var flagDebug, flagTrace, flagVersion bool
var flagFile, flagRecord string

// Version specifies the version
var Version string

func init() {
	flag.BoolVar(&flagDebug, "debug", false, "debug")
	flag.BoolVar(&flagTrace, "trace", false, "trace")
	flag.BoolVar(&flagVersion, "version", false, "show version")
	flag.StringVar(&flagRecord, "record", "", "record input to miti file")
	flag.StringVar(&flagFile, "play", "", "play sequence from miti file")
	if Version == "" {
		Version = "v0.3.4-66cb94c"
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
	var err error
	if flagRecord != "" {
		err = record.Record(flagRecord)
	} else {
		f, _ := os.Create("song.txt")
		f.WriteString(`chain 1 
pattern 1 
instruments 1
CEG
`)
		f.Close()
		open.Start(f.Name())
		flagFile = f.Name()
		err = play.Play(flagFile)
	}
	if err != nil {
		log.Error(err)
	}

}
