package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/schollz/logger"
	"github.com/schollz/saps/src/sequencer"
)

func Run() (err error) {
	s := sequencer.New()
	r := gin.New()
	r.POST("/msg", func(c *gin.Context) {
		var m Message
		if err := c.ShouldBindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "error", "payload": err.Error(), "success": false})
			return
		}
		log.Infof("message: %+v", m)
		go process(s, m)
		c.JSON(http.StatusOK, gin.H{"msg": "ok", "success": true})
	})
	r.Run()
	return
}

type Message struct {
	Message    string `json:"msg"`
	Payload    string `json:"data"`
	PayloadNum int    `json:"num"`
	Success    bool   `json:"success"`
}

func process(s *sequencer.Sequencer, m Message) {
	if m.Message == "tempo" {
		s.UpdateTempo(m.PayloadNum)
	} else if m.Message == "start" {
		s.Start()
	} else if m.Message == "stop" {
		s.Stop()
	}
}
