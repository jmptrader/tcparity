package main

import (
	"errors"

	"github.com/mdlayher/goset"
)

// balanceAlgorithm represents an algorithm used by tcparity to pick which server to
// proxy requests between
type balanceAlgorithm interface {
	Balance() (*server, error)
	SetServers(*set.Set) error
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

// SetServers takes an input set and stores it for use with the Round-Robin Algorithm
func (a *roundRobinAlgorithm) SetServers(servers *set.Set) error {
	// Verify servers not empty
	if servers.Size() == 0 {
		return errors.New("roundrobin: no servers in set")
	}

	// Initialize index
	a.index = 0

	// Check for valid servers, set them
	a.servers = make([]*server, 0)
	for _, s := range servers.Enumerate() {
		host, ok := s.(string)
		if !ok {
			return errors.New("roundrobin: invalid server")
		}

		server := &server{Host: host}

		// Add server to list
		a.servers = append(a.servers, server)
	}

	return nil
}
