package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type registerRequest struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

func register(name, url string) (err error) {

	body, err := json.Marshal(registerRequest{
		Name: name,
		URL:  url,
	})
	if err != nil {
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Println("response", response)

	return

}
