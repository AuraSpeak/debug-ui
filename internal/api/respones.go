package api

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ApiResponse interface {
	Send(w http.ResponseWriter)
}

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"`
}

func (e *ApiError) Send(w http.ResponseWriter) {
	w.WriteHeader(e.Code)
	b, err := json.Marshal(e)
	if err != nil {
		log.WithField("caller", "web").WithError(err).Error("Can't marshal ApiError to json")
	}
	w.Write(b)
	w.Write([]byte("\n"))
}

type ApiSuccess struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (s *ApiSuccess) Send(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(s)
	if err != nil {
		log.WithField("caller", "web").WithError(err).Error("Can't marshal ApiSuccess to json")
	}
	w.Write(b)
	w.Write([]byte("\n"))
}

type ServerStateResponse struct {
	ShouldStop bool `json:"shouldStop"`
	IsAlive    bool `json:"isAlive"`
}

func (s *ServerStateResponse) Send(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(s)
	if err != nil {
		log.WithField("caller", "web").WithError(err).Error("Can't marshal ServerStateResponse to json")
	}
	w.Write(b)
	w.Write([]byte("\n"))
}

type UDPClientResponse struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func (s *UDPClientResponse) Send(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(s)
	if err != nil {
		log.WithField("caller", "web").WithError(err).Error("Can't marshal UDPClientResponse to json")
	}
	w.Write(b)
	w.Write([]byte("\n"))
}

type UDPClientStateResponse struct {
	Id        int        `json:"id"`
	Running   bool       `json:"running"`
	Datagrams []Datagram `json:"datagrams"`
}

func (s *UDPClientStateResponse) Send(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(s)
	if err != nil {
		log.WithField("caller", "web").WithError(err).Error("Can't marshal UDPClientStateResponse to json")
	}
	w.Write(b)
	w.Write([]byte("\n"))
}

type AllUDPClientResponse struct {
	UDPClients []UDPClientResponse `json:"udpClients"`
}

func (a *AllUDPClientResponse) Send(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(a)
	if err != nil {
		log.WithField("caller", "web").WithError(err).Error("Can't marshal AllUDPClientResponse to json")
	}
	w.Write(b)
	w.Write([]byte("\n"))
}

type UDPClientListItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UDPClientPaginatedRespone struct {
	Items    []UDPClientListItem `json:"items"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
	Total    int                 `json:"total"`
}

func (p *UDPClientPaginatedRespone) Send(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(p)
	if err != nil {
		log.WithField("caller", "web").WithError(err).Error("Can't marshal PaginatedResponse to json")
	}
	w.Write(b)
	w.Write([]byte("\n"))
}

type SendDatagramRequest struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
	Format  string `json:"format"` // "hex" or "text"
}

type SendDatagramResponse struct {
	Message string `json:"message"`
}

func (s *SendDatagramResponse) Send(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(s)
	if err != nil {
		log.WithField("caller", "web").WithError(err).Error("Can't marshal SendDatagramResponse to json")
	}
	w.Write(b)
	w.Write([]byte("\n"))
}

type MermaidResponse struct {
	Heading string `json:"heading"`
	Diagram string `json:"diagram"`
}

func (m *MermaidResponse) Send(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(m)
	if err != nil {
		log.WithField("caller", "web").WithError(err).Error("Can't marshal MermaidDiagram to json")
	}
	w.Write(b)
	w.Write([]byte("\n"))
}
