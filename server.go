package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

type registerRequest struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

func register(name, url string) (err error) {

	log.Info().Str("name", name).Str("url", url).Msg("making request to existing cluster")

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

	log.Debug().Int("status code", resp.StatusCode).Msg("server responded")
	fmt.Println("response", response)

	parsed := map[string]string{}
	err = json.Unmarshal(response, &parsed)
	if err != nil {
		return err
	}

	for k, v := range parsed {
		nodes[k] = v
	}

	return

}
