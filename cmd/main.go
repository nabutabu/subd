package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nabutabu/subd/agent"
	"github.com/nabutabu/subd/client"
)

const (
	APP_SETTINGS_FILE = "appsettings.json"
	IDENTITY_FILE     = "identity.json"
)

type Configuration struct {
	DominatorURL string
	Token        string
}

func LoadConfig() Configuration {
	file, err := os.Open(APP_SETTINGS_FILE)
	if err != nil {
		log.Fatal("Error opening config file:", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("Error decoding config file:", err)
	}

	return configuration
}

func main() {
	config := LoadConfig()

	dom := *client.New(config.DominatorURL, config.Token)

	agent := agent.New(dom)

	agent.Run()
}
