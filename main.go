package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

const (
	Leanplum_api_url = "https://www.leanplum.com/api"
	Api_version      = "1.0.6"
)

type Credentials struct {
	appId      string `url:"appId"`
	clientKey  string `url:"clientKey"`
	apiVersion string `url:"apiVersion"`
	action     string `url:"action"`
}

func main() {
	query := Credentials{"123", "123", "1.0.6", "multi"}
	params := map[string]string{
		"action":         "setUserAttributes",
		"userAttributes": "{\"Interests\":[\"Go\",\"IT\"]}",
	}
	Get(query, params)
}

func Get(credentials Credentials, arguments map[string]string) {

	auth, _ := query.Values(credentials)
	queryString := url.Values{}
	for k, v := range arguments {
		queryString.Add(k, v)
	}
	url := Leanplum_api_url + "?" + auth.Encode() + queryString.Encode()

	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Success: %s", body)

}

func Call(auth Credentials, action string, argumetns map[string]string) (bool, error) {

	return true, nil
}
