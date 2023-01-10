// Package config holds all config related entities
package config

import (
	"fmt"
	"os"
)

// IsLocal returns whether or not the running application is operating locally
func IsLocal() bool {
	return os.Getenv("IS_LOCAL") == "true"
}

// GetToken returns a GitHub access token for the user
func GetToken() (*string, error) {
	token := os.Getenv("GIT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("no token specified")
	}
	return &token, nil
}

// GetMachineToken returns a GitHub machine access token for machine actions
func GetMachineToken() (*string, error) {
	token := os.Getenv("GIT_MACHINE_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("no machine token specified")
	}
	return &token, nil
}

// GetTrackingRepo returns the GitHub repository to use as a backing store
func GetTrackingRepo() (*string, error) {
	repo := os.Getenv("TRACKING_REPOSITORY")
	if repo == "" {
		return nil, fmt.Errorf("no tracking repository specified")
	}
	return &repo, nil
}
