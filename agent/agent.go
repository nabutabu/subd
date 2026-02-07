package agent

import (
	"errors"
	"log"
	"os/exec"
	"time"

	"github.com/nabutabu/subd/client"
	"github.com/nabutabu/subd/collector"
	"github.com/nabutabu/subd/types"
)

type Agent struct {
	state        types.State
	dominator    client.Client
	Collector    collector.BasicCollector
	tickInterval time.Duration
}

func New(dominator client.Client) *Agent {
	return &Agent{
		dominator:    dominator,
		tickInterval: 2 * time.Second,
	}
}

func (a *Agent) Run() {
	ticker := time.NewTicker(a.tickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.Check()
		}
	}
}

func (a *Agent) Check() {
	log.Println(a.state)
	// Observe local state from collector
	currState, err := a.Collector.Collect()
	if err != nil {
		log.Println(err)
	}

	// Send heartbeat through api
	desiredState, err := a.dominator.Heartbeat(*currState)
	if err != nil {
		log.Println(err)
	}

	// Receive desired state in api response and Compute diff
	a.diff(desiredState, currState)

	// Reconcile if needed
}

/*
*
@params
desiredServices: services that should be active on the Host
currServices: services currently active on the Host
@return
[]string: array of systemctl actions that either start or stop services
*/
func diffServices(desiredServices, currServices map[string]types.Service) []string {
	var result []string
	for _, service := range desiredServices {
		if _, found := currServices[service.Name]; !found {
			// service not found, add to Host
			result = append(result, "systemctl start "+service.Name)
		}
	}

	// check for extra services
	for _, service := range currServices {
		if _, found := desiredServices[service.Name]; !found {
			// this service should not exist, kill it
			result = append(result, "systemctl stop "+service.Name)
		}
	}

	return result
}

func performActions(actions []string) error {
	var flag bool = false
	for _, action := range actions {
		cmd := exec.Command(action)
		_, err := cmd.CombinedOutput()
		if err != nil {
			flag = true
			log.Printf("failed to reconcile service with the following command: %s", action)
		}
	}

	if flag {
		return errors.New("Failed to do all actions")
	}

	return nil
}

func (a *Agent) diff(desiredState, currState *types.State) {
	// check if services are different
	actions := diffServices(desiredState.Services, currState.Services)

	performActions(actions)
}
