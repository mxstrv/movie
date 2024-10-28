package memory

import (
	"context"
	"errors"
	"movieapp/pkg/discovery"
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
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return nil
	}
	delete(r.serviceAddrs[serviceName], instanceID)
	return nil
}

// ServiceAddresses returns list of addresses of active instances of a given service.
func (r *Registry) ServiceAddresses(ctx context.Context, serviceID string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.serviceAddrs[serviceID]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, ins := range r.serviceAddrs[serviceID] {
		// Hardcoded 5 seconds
		if ins.lastActive.Before(time.Now().Add(5 * time.Second)) {
			continue
		}
		res = append(res, ins.hostPort)
	}
	return res, nil
}

func (r *Registry) ReportHealthyState(instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return errors.New("service is not registered")
	}
	if _, ok := r.serviceAddrs[serviceName][instanceID]; !ok {
		return errors.New("instance is not registered")
	}
	r.serviceAddrs[serviceName][instanceID].lastActive = time.Now()
	return nil
}
