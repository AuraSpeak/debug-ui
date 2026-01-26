package api

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type DatagramDirection int

const (
	ClientToServer DatagramDirection = 1
	ServerToClient DatagramDirection = 2
)

type Datagram struct {
	Direction DatagramDirection `json:"direction"`
	Message   []byte            `json:"message"`
}

func (d *Datagram) Send(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(d)
	if err != nil {
		log.WithField("caller", "web").WithError(err).Error("Can't marshal Datagram to json")
	}
	w.Write(b)
	w.Write([]byte("\n"))
}

type trace struct {
	TS time.Time `json:"ts"`
}
