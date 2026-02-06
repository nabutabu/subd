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

func getRunningServices() ([]string, error) {
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--state=running,active", "--no-pager")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run systemctl: %w, Output: %s", err, string(output))
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var services []string
	for _, line := range lines {
		// Example line format: "ServiceName.service loaded active running Description"
		fields := strings.Fields(line)
		if len(fields) > 0 {
			services = append(services, fields[0]) // Append the service name (first field)
		}
	}
	return services, nil
}

func (c *BasicCollector) Collect() (*types.State, error) {
	log.Println("Collecting...")
	services_raw, err := getRunningServices()
	if err != nil {
		return nil, err
	}

	var services []types.Service
	for _, service := range services_raw {
		new_service := types.Service{
			Name: service,
		}
		services = append(services, new_service)
	}

	return &types.State{
		Services: services,
	}, nil
}
