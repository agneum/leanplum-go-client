package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/google/go-querystring/query"
)

const (
	Leanplum_api_url = "https://www.leanplum.com/api"
	configFile       = "config.toml"
)

type Config struct {
	AppId      string `url:"appId"`
	ClientKey  string `url:"clientKey"`
	ApiVersion string `url:"apiVersion"`
}

func main() {
	query := ReadConfig()
	params := map[string]string{
		"action":         "setUserAttributes",
		"userAttributes": "{\"Interests\":[\"Go\",\"IT\"]}",
	}
	Get(query, params)
}

func Get(config Config, arguments map[string]string) {
	auth, _ := query.Values(config)
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

func ReadConfig() Config {
	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatal("Config file is not read: ", configFile)
	}

	var conf Config
	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		log.Fatal(err)
	}

	return conf
}

func Call(auth Config, action string, argumetns map[string]string) (bool, error) {

	return true, nil
}
