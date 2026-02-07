package collector

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/nabutabu/subd/types"
)

type BasicCollector struct {
	NodeID string
}

func NewBasic(nodeID string) *BasicCollector {
	return &BasicCollector{NodeID: nodeID}
}

func getServices() (map[string]types.Service, error) {
	cmd := exec.Command("systemctl", "list-units", "--type=service",
	"--state=running,active", "--no-pager")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run systemctl: %w, Output: %s", err, string(output))
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	services := make(map[string]types.Service)

	for _, line := range lines {
		// Example line format: "ServiceName.service loaded active running Description"
		fields := strings.Fields(line)
		if len(fields) > 0 {
			services[fields[0]] = types.Service{
				Name: fields[0],
				Sub: fields[3],
				Description: fields[4],
			}
		}
	}
	return services, nil
}

func (c *BasicCollector) Collect() (*types.State, error) {
	log.Println("Collecting...")
	services, err := getServices()
	if err != nil {
		return nil, err
	}

	return &types.State{
		Services: services,
	}, nil
}
