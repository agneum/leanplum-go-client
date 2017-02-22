package leanplum_users

import (
	"encoding/json"
	"reflect"

	"github.com/agneum/leanplum-go-client"
)

const setAttributeAction = "setUserAttributes"

type AttributeContent struct {
	Action                    string `url:"action"`
	UserId                    string `url:"userId"`
	UserAttributes            string `url:"userAttributes,omitempty"`
	UserAttributesToAdd       string `url:"userAttributeValuesToAdd,omitempty"`
	UserAttributesToRemove    string `url:"userAttributeValuesToRemove,omitempty"`
	UserAttributesToIncrement string `url:"userAttributeValuesToIncrement,omitempty"`
}

func NewAttributeContent(userId string) AttributeContent {
	var userAttribute AttributeContent
	userAttribute.Action = setAttributeAction
	userAttribute.UserId = userId

	return userAttribute
}

func Start(config leanplum.Config, arguments map[string]string) {
}

func Stop() {}

func SendAttributes(config leanplum.Config, queryParams AttributeContent) ([]leanplum.ServerResponse, error) {
	return leanplum.Get(config, leanplum.MakeEncodedQueryStringFromStruct(queryParams))
}

func (attribute *AttributeContent) SetAttribute(field string, values map[string]string) error {
	v := reflect.ValueOf(attribute).Elem().FieldByName(field)

	if v.IsValid() {
		val, err := json.Marshal(values)
		if err != nil {
			return err
		}
		v.SetString(string(val))
	}

	return nil
}

func (attribute *AttributeContent) SetSliceOfAttributes(field string, values map[string][]string) error {
	v := reflect.ValueOf(attribute).Elem().FieldByName(field)

	if v.IsValid() {
		val, err := json.Marshal(values)
		if err != nil {
			return err
		}
		v.SetString(string(val))
	}

	return nil
}

func ExportsUsers() {}
