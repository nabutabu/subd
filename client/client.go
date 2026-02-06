package client

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/nabutabu/subd/types"
)

const (
	HEARTBEAT = "/v1/heatbeat"
)

type Client struct {
	URL    string
	Token  string
	Client *http.Client
}

func New(url string, token string) *Client {
	return &Client{
		URL:    url,
		Token:  token,
		Client: &http.Client{},
	}
}

func (api *Client) Heartbeat(currState types.State) (*types.State, error) {
	// Convert the struct to a JSON byte slice
	jsonData, err := json.Marshal(currState)
	if err != nil {
		log.Println(err)
	}

	// POST to Dominator with current state
	req, err := http.NewRequest("POST", api.URL+HEARTBEAT, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+api.Token)

	response, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var desiredState types.State
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(body, &desiredState)

	return &desiredState, nil
}
