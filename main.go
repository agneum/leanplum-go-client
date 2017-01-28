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

type CommonResponse struct {
	Response []Response `json:"response"`
}

type Response struct {
	Success bool           `json:"success"`
	Error   *ResponseError `json:"error,omitempty"`
	Warning *ResponseError `json:"warning,omitempty"`
	Result  *string        `json:"result,omitempty"`
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

func Get(config Config, arguments map[string]string) {
	uriParams, _ := query.Values(config)
	queryString := url.Values{}
	for k, v := range arguments {
		queryString.Add(k, v)
	}
	url := Leanplum_api_url + "?" + uriParams.Encode() + queryString.Encode()

	response, err := http.Get(url)

	defer response.Body.Close()

	resp := new(CommonResponse)
	err = json.NewDecoder(response.Body).Decode(resp)

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	if resp.Response != nil {
		log.Fatalf("The error message: %v\n", resp.Response[0].Error.Message)
	}

	log.Printf("Success is %v\n", resp.Response[0].Success)
}

func Post(config Config, body []byte, parameters map[string]string) {
	uriParams, _ := query.Values(config)
	queryString := makeQueryString(parameters)
	url := Leanplum_api_url + "?" + uriParams.Encode() + "&" + queryString.Encode()

	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		log.Fatalf("%v\n", err)
		return
	}

	resp := new(CommonResponse)
	err = json.NewDecoder(response.Body).Decode(resp)

	if err != nil {
		log.Fatalf("Error: %v\n", err)
		return
	}

	if !resp.Response[0].Success {
		if resp.Response[0].Error != nil {
			log.Fatalf("The error message: %v\n", resp.Response[0].Error.Message)
			return
		}
	}

	if resp.Response[0].Warning != nil {
		log.Fatalf("The warning message: %v\n", resp.Response[0].Warning.Message)
		return
	}

	log.Printf("Success is %v. Result is %v\n", resp.Response[0].Success, resp.Response[0].Result)
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
