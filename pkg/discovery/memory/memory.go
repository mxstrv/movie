package memory

import (
	"context"
	"sync"
	"time"
)

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

type Registry struct {
	sync.RWMutex
	serviceAddrs map[string]map[string]*serviceInstance
}

func NewRegistry() *Registry {
	return &Registry{serviceAddrs: map[string]map[string]*serviceInstance{}}
}

// Register creates a service record in the registry.
func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		r.serviceAddrs[serviceName] = map[string]*serviceInstance{}
	}
	r.serviceAddrs[serviceName][instanceID] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

// Deregister removes a service record in the registry
func (r *Registry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	return nil
	// TODO: implement
}

func (r *Registry) ServiceAddresses(ctx context.Context, serviceID string) ([]string, error) {
	return nil, nil
	// TODO: implement
}

func (r *Registry) ReportHealthyState(instanceID string, serviceName string) error {
	return nil
	//TODO: implement
}
