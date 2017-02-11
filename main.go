package leanplum

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/google/go-querystring/query"
)

const (
	Leanplum_api_url = "https://www.leanplum.com/api"
)

type Config struct {
	AppId      string `url:"appId"`
	ClientKey  string `url:"clientKey"`
	ApiVersion string `url:"apiVersion"`
}

type ICommonResponse interface {
	ProcessResponse() error
}

type ServerResponse struct {
	Success bool           `json:"response[0].success"`
	Error   *ResponseError `json:"response[0].error,omitempty"`
	Warning *ResponseError `json:"response[0].warning,omitempty"`
	Result  *string        `json:"response[0].result,omitempty"`
	// MessagesSent *string        `json:"messagesSent,omitempty"`
}

type ResponseError struct {
	Message string `json:"message"`
}

type Message struct {
	Data MessageContent `json:"data"`
}

type MessageContent struct {
	Time   int64             `json:"time"`
	Values map[string]string `json:"values"`
}

func Get(config Config, arguments map[string]string) *ServerResponse {
	uriParams, _ := query.Values(config)
	queryString := makeQueryString(arguments)

	url := Leanplum_api_url + "?" + uriParams.Encode() + queryString.Encode()

	response, err := http.Get(url)

	defer response.Body.Close()

	resp := new(ServerResponse)
	err = json.NewDecoder(response.Body).Decode(resp)

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	if resp.Error != nil {
		log.Fatalf("The error message: %v\n", resp.Error.Message)
	}

	log.Printf("Success is %v\n", resp.Success)

	return resp
}

func Post(config Config, body []byte, parameters map[string]string) {
	uriParams, _ := query.Values(config)
	queryString := makeQueryString(parameters)
	url := Leanplum_api_url + "?" + uriParams.Encode() + "&" + queryString.Encode()

	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	resp := new(ServerResponse)
	_ = resp.processResponse(response)

	if !resp.Success {
		if resp.Error != nil {
			log.Fatalf("The error message: %v\n", resp.Error.Message)
			return
		}
	}

	if resp.Warning != nil {
		log.Fatalf("The warning message: %v\n", resp.Warning.Message)
		return
	}

	log.Printf("Success is %v. Result is %v\n", resp.Success, resp.Result)
}

func (responser *ServerResponse) processResponse(response *http.Response) error {
	if response != nil {
		defer response.Body.Close()
	}
	err := json.NewDecoder(response.Body).Decode(responser)

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// TODO: processing errors and warnings
	if responser.Error != nil {
		log.Fatalf("The error message: %v\n", responser.Error.Message)
	}

	return err
}

// Read config from toml-file
func ReadConfig(configFile string) Config {
	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatal("Config file is not read: ", configFile)
	}

	var config Config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}

	return config
}

func makeQueryString(parameters map[string]string) url.Values {
	queryString := url.Values{}
	for k, v := range parameters {
		queryString.Add(k, v)
	}

	return queryString
}
