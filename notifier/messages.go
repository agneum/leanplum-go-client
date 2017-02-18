package notifier

import (
	"encoding/json"
	"time"

	"github.com/agneum/leanplum-go-client"
)

type Message struct {
	Data MessageContent `json:"data"`
}

type MessageContent struct {
	Time   int64             `json:"time"`
	Values map[string]string `json:"values"`
}

func SendMessage(config leanplum.Config, queryParams, values map[string]string) ([]leanplum.ServerResponse, error) {
	body, _ := createBody(values)
	return leanplum.Post(config, queryParams, body)
}

func createBody(values map[string]string) ([]byte, error) {
	return json.Marshal(Message{Data: MessageContent{Time: time.Now().Unix(), Values: values}})
}
