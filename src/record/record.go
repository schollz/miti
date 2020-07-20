package record

import (
	"os"
	"os/signal"

	"github.com/schollz/miti/src/log"
)

func Record() (err error) {
	log.Debug("init recording")

	// shutdown everything on Ctl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Debug(sig)
			log.Info("shutting down")
		}
	}()

	return
}
