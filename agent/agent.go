package agent

import (
	"log"
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
	_, err = a.dominator.Heartbeat(*currState)
	if err != nil {
		log.Println(err)
	}

	// Receive desired state in api response

	// Compute diff

	// Reconcile if needed
}
