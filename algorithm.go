package main

import (
	"errors"
)

// balanceAlgorithm represents an algorithm used by tcparity to pick which server to
// proxy requests between
type balanceAlgorithm interface {
	Balance() (*server, error)
	SetServers([]string) error
}

// roundRobinAlgorithm represents a balancing algorithm in which servers are rotated
// through in order, equally choosing each server from the set
type roundRobinAlgorithm struct {
	index   int
	servers []*server
}

// Balance uses the Round-Robin Algorithm to rotate through each server in a list,
// equally choosing each server from the set
func (a *roundRobinAlgorithm) Balance() (*server, error) {
	// Reset index if at end
	if a.index == len(a.servers) {
		a.index = 0
	}

	// Select a server, increment index
	s := a.servers[a.index]
	a.index++

	return s, nil
}

// SetServers takes an input slice and stores it for use with the Round-Robin Algorithm
func (a *roundRobinAlgorithm) SetServers(servers []string) error {
	// Verify servers not empty
	if len(servers) == 0 {
		return errors.New("roundrobin: no servers in slice")
	}

	// Initialize index
	a.index = 0

	// Check for valid servers, set them
	a.servers = make([]*server, 0)
	for _, s := range servers {
		server := &server{Host: s}

		// Add server to list
		a.servers = append(a.servers, server)
	}

	return nil
}
