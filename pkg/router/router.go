package router

import (
	"crypto/tls"
	"errors"
	"net/url"
	log "github.com/sirupsen/logrus"
)

// RoutingTable holds the routing information for different hosts.
type RoutingTable struct {
	// For simplicity, this map could represent your routing rules or backend services.
	// In a real application, you'd likely have a more complex structure or database to manage this.
	Backends map[string]string
	// Certificates map for storing TLS certificates by server name (e.g., "example.com").
	Certificates map[string]tls.Certificate
}

// NewRoutingTable creates and initializes a new RoutingTable.
func NewRoutingTable() *RoutingTable {
	return &RoutingTable{
		Backends:     make(map[string]string),
		Certificates: make(map[string]tls.Certificate),
	}
}

// GetBackend retrieves the backend service URL for a given host and requestURI.
// In a real application, you might query a database or use more complex logic to determine the backend.
func (rt *RoutingTable) GetBackend(host, requestURI string) (*url.URL, error) {
	log.Printf("Routing request for host: %s, requestURI: %s", host, requestURI)
	backendURLString, ok := rt.Backends[host]
	if !ok {
		return nil, errors.New("backend service not found for host")
	}
	backendURL, err := url.Parse(backendURLString)
	if err != nil {
		return nil, err
	}
	return backendURL, nil
}

// AddBackend is a helper method to add a backend service URL for a host.
func (rt *RoutingTable) AddBackend(host, backendURL string) {
	rt.Backends[host] = backendURL
}