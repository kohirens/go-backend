package backend

import "fmt"

func NewServiceManager() ServiceManager {
	return &Servicer{
		services: make(map[string]interface{}),
	}
}

// ServiceManager handles storing and retrieval of dependencies when an endpoint
// handler function is called. Granting access to objects such as Session,
// Database, and what ever else may be useful.
type ServiceManager interface {
	Add(name string, service interface{})
	Get(name string) (interface{}, error)
}

// Servicer A default service manager provided by the backend package.
type Servicer struct {
	services map[string]interface{}
}

// Add adds a service to a map for access from an endpoint function.
func (m *Servicer) Add(name string, service interface{}) {
	m.services[name] = service
}

// Get Retrieves a service from.
func (m *Servicer) Get(name string) (interface{}, error) {
	x, ok := m.services[name]
	if !ok {
		return nil, fmt.Errorf(stderr.ServiceNotFound, name)
	}

	return x, nil
}
