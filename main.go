package leanplum

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type IResponse interface {
	CheckErrors() bool
}

type CommonResponse struct {
	Response []ServerResponse `json:"response"`
}

type ServerResponse struct {
	Success      bool           `json:"success"`
	Error        *ResponseError `json:"error,omitempty"`
	Warning      *ResponseError `json:"warning,omitempty"`
	Result       *string        `json:"result,omitempty"`
	MessagesSent *string        `json:"messagesSent,omitempty"`
}

type ResponseError struct {
	Message string `json:"message"`
}

func Get(config Config, arguments map[string]string) ([]ServerResponse, error) {
	uriParams, _ := query.Values(config)
	queryString := makeQueryString(arguments)

	url := Leanplum_api_url + "?" + uriParams.Encode() + "&" + queryString.Encode()

	response, err := http.Get(url)

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	commonResponser := new(CommonResponse)
	slice, err := commonResponser.processResponse(response)

	return slice, err
}

func Post(config Config, queryParams map[string]string, body []byte) ([]ServerResponse, error) {
	uriParams, _ := query.Values(config)
	queryString := makeQueryString(queryParams)
	url := Leanplum_api_url + "?" + uriParams.Encode() + "&" + queryString.Encode()

	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	commonResponser := new(CommonResponse)

	return commonResponser.processResponse(response)
}

func (responser *CommonResponse) processResponse(response *http.Response) ([]ServerResponse, error) {

	defer response.Body.Close()

	err := json.NewDecoder(response.Body).Decode(&responser)

	if err != nil {
		return nil, err
	}

	return responser.Response, nil
}

func (serverResponse *ServerResponse) CheckErrors() (bool, error) {
	if !serverResponse.Success {
		if serverResponse.Error != nil {
			log.Printf("The error message: %v\n", serverResponse.Error.Message)
			return false, fmt.Errorf("The error message: %v\n", serverResponse.Error.Message)
		}

		if serverResponse.Warning != nil {
			log.Printf("The warning message: %v\n", serverResponse.Warning.Message)
			return false, fmt.Errorf("The warning message: %v\n", serverResponse.Warning.Message)
		}
	}

	return serverResponse.Success, nil
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
