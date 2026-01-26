package ws

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type WebSocketHook struct {
	hub       *WebSocketHub
	formatter logrus.Formatter
}

func NewWebSocketHook(hub *WebSocketHub) *WebSocketHook {
	return &WebSocketHook{
		hub:       hub,
		formatter: &logrus.JSONFormatter{},
	}
}

func (wh *WebSocketHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (wh *WebSocketHook) Fire(entry *logrus.Entry) error {
	if wh.hub == nil {
		return nil
	}
	b, err := wh.formatter.Format(entry)
	if err != nil {
		fmt.Println(err)
		return err
	}
	wh.hub.Broadcast([]byte(b))
	return nil
}
