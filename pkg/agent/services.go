package agent

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ServiceManager handles service management on the node
type ServiceManager struct {
	agent    *Agent
	services map[string]Service
	mu       sync.RWMutex
}

// NewServiceManager creates a new service manager
func NewServiceManager(agent *Agent) *ServiceManager {
	return &ServiceManager{
		agent:    agent,
		services: make(map[string]Service),
	}
}

// RegisterService registers a new service
func (sm *ServiceManager) RegisterService(service Service) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	service.Status = "starting"
	service.Health = "unknown"
	sm.services[service.Name] = service

	log.Printf("Registered service: %s", service.Name)

	return sm.updateServiceInAgent(service.Name, service)
}

// UnregisterService unregisters a service
func (sm *ServiceManager) UnregisterService(name string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	service, exists := sm.services[name]
	if !exists {
		return fmt.Errorf("service %s not found", name)
	}

	service.Status = "stopped"
	service.Health = "stopped"

	log.Printf("Unregistered service: %s", name)

	delete(sm.services, name)
	return sm.updateServiceInAgent(name, service)
}

// UpdateService updates a service
func (sm *ServiceManager) UpdateService(name string, service Service) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	_, exists := sm.services[name]
	if !exists {
		return fmt.Errorf("service %s not found", name)
	}

	sm.services[name] = service
	log.Printf("Updated service: %s", name)

	return sm.updateServiceInAgent(name, service)
}

// GetService gets a service by name
func (sm *ServiceManager) GetService(name string) (Service, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	service, exists := sm.services[name]
	if !exists {
		return Service{}, fmt.Errorf("service %s not found", name)
	}

	return service, nil
}

// ListServices lists all services
func (sm *ServiceManager) ListServices() map[string]Service {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	result := make(map[string]Service)
	for name, service := range sm.services {
		result[name] = service
	}

	return result
}

// StartService starts a service
func (sm *ServiceManager) StartService(name string) error {
	sm.mu.Lock()
	service, exists := sm.services[name]
	if !exists {
		sm.mu.Unlock()
		return fmt.Errorf("service %s not found", name)
	}

	service.Status = "starting"
	sm.services[name] = service
	sm.mu.Unlock()

	log.Printf("Starting service: %s", name)

	if err := sm.updateServiceInAgent(name, service); err != nil {
		log.Printf("Failed to update service %s in agent: %v", name, err)
	}

	sm.mu.Lock()
	service.Status = "running"
	service.Health = "healthy"
	sm.services[name] = service
	sm.mu.Unlock()

	return sm.updateServiceInAgent(name, service)
}

// StopService stops a service
func (sm *ServiceManager) StopService(name string) error {
	sm.mu.Lock()
	service, exists := sm.services[name]
	if !exists {
		sm.mu.Unlock()
		return fmt.Errorf("service %s not found", name)
	}

	service.Status = "stopping"
	sm.services[name] = service
	sm.mu.Unlock()

	log.Printf("Stopping service: %s", name)

	sm.mu.Lock()
	service.Status = "stopped"
	service.Health = "stopped"
	sm.services[name] = service
	sm.mu.Unlock()

	return sm.updateServiceInAgent(name, service)
}

// CheckServiceHealth performs health check on a service
func (sm *ServiceManager) CheckServiceHealth(name string) error {
	sm.mu.RLock()
	service, exists := sm.services[name]
	if !exists {
		sm.mu.RUnlock()
		return fmt.Errorf("service %s not found", name)
	}
	sm.mu.RUnlock()

	health := sm.performHealthCheck(service)

	sm.mu.Lock()
	service.Health = health
	sm.services[name] = service
	sm.mu.Unlock()

	return sm.updateServiceInAgent(name, service)
}

// performHealthCheck performs actual health check for a service
func (sm *ServiceManager) performHealthCheck(service Service) string {
	if service.Port <= 0 {
		return "no_port"
	}

	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://localhost:%d/health", service.Port)

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Health check failed for %s: %v", service.Name, err)
		return "unhealthy"
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return "healthy"
	}

	return "unhealthy"
}

// updateServiceInAgent updates service in the agent's service map
func (sm *ServiceManager) updateServiceInAgent(name string, service Service) error {
	sm.agent.mu.Lock()
	defer sm.agent.mu.Unlock()

	sm.agent.services[name] = service
	return nil
}

// ServiceDiscovery handles service discovery and registration
type ServiceDiscovery struct {
	manager *ServiceManager
	agent   *Agent
}

// NewServiceDiscovery creates a new service discovery
func NewServiceDiscovery(manager *ServiceManager, agent *Agent) *ServiceDiscovery {
	return &ServiceDiscovery{
		manager: manager,
		agent:   agent,
	}
}

// DiscoverServices discovers services on the node
func (sd *ServiceDiscovery) DiscoverServices() error {
	sessions := sd.agent.sessions

	for sessionID, session := range sessions {
		service := Service{
			Name:   sessionID,
			Type:   session.Provider,
			Status: "running",
			Port:   session.BridgePort,
			Health: "unknown",
			Labels: session.Labels,
			Metadata: map[string]interface{}{
				"session_id": sessionID,
				"provider":   session.Provider,
				"created_at": time.Now(),
			},
		}

		if err := sd.manager.RegisterService(service); err != nil {
			log.Printf("Failed to register discovered service %s: %v", sessionID, err)
		}
	}

	log.Printf("Service discovery completed")
	return nil
}

// MonitorServices monitors services for changes
func (sd *ServiceDiscovery) MonitorServices(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sd.manager.mu.RLock()
			services := make(map[string]Service)
			for name, service := range sd.manager.services {
				services[name] = service
			}
			sd.manager.mu.RUnlock()

			for name := range services {
				if err := sd.manager.CheckServiceHealth(name); err != nil {
					log.Printf("Health check failed for service %s: %v", name, err)
				}
			}
		}
	}
}
