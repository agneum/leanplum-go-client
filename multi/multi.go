package notifier

import (
	"encoding/json"
	"time"
)

type Message struct {
	Data MessageContent `json:"data"`
}

type MessageContent struct {
	Time   int64             `json:"time"`
	Values map[string]string `json:"values"`
}

func createBody(values map[string]string) ([]byte, error) {
	return json.Marshal(Message{Data: MessageContent{Time: time.Now().Unix(), Values: values}})
}
