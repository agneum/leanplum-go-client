package notifier

import (
	"encoding/json"

	"github.com/agneum/leanplum-go-client"
)

const sendMessageAction = "sendMessage"

type MessageContent struct {
	Action            string `url:"action"`
	UserId            string `url:"userId"`
	MessageId         string `url:"messageId"`
	CreateDisposition string `url:"createDisposition,omitempty"`
	Force             bool   `url:"force,omitempty"`
	Values            string `url:"values,omitempty"`
}

func NewMessage(userId, messageId string) MessageContent {
	var message MessageContent
	message.Action = sendMessageAction
	message.UserId = userId
	message.MessageId = messageId

	return message
}

func (message *MessageContent) SetMessageValues(values map[string]string) error {
	val, err := json.Marshal(values)

	message.Values = string(val)

	return err
}

func SendMessage(config leanplum.Config, queryParams MessageContent) ([]leanplum.ServerResponse, error) {
	return leanplum.Get(config, leanplum.MakeEncodedQueryStringFromStruct(queryParams))
}
