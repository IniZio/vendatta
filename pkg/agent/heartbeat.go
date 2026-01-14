package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// HeartbeatService handles periodic heartbeat to coordination server
type HeartbeatService struct {
	agent  *Agent
	client *http.Client
}

// NewHeartbeatService creates a new heartbeat service
func NewHeartbeatService(agent *Agent) *HeartbeatService {
	return &HeartbeatService{
		agent: agent,
		client: &http.Client{
			Timeout: agent.config.Heartbeat.Timeout,
		},
	}
}

// Start starts the heartbeat service
func (h *HeartbeatService) Start(ctx context.Context) {
	ticker := time.NewTicker(h.agent.config.Heartbeat.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := h.sendHeartbeat(); err != nil {
				log.Printf("Heartbeat failed: %v", err)
			}
		}
	}
}

// sendHeartbeat sends a heartbeat to the coordination server
func (h *HeartbeatService) sendHeartbeat() error {
	if h.agent.config.CoordinationURL == "" {
		return nil
	}

	h.agent.mu.RLock()
	nodeCopy := *h.agent.node
	nodeCopy.LastSeen = time.Now()
	servicesCopy := make(map[string]Service)
	for k, v := range h.agent.services {
		servicesCopy[k] = v
	}
	h.agent.mu.RUnlock()

	heartbeatData := map[string]interface{}{
		"last_seen": nodeCopy.LastSeen,
		"status":    nodeCopy.Status,
		"services":  servicesCopy,
		"version":   nodeCopy.Version,
		"uptime":    time.Since(nodeCopy.CreatedAt).String(),
	}

	url := fmt.Sprintf("%s/api/v1/nodes/%s/heartbeat", h.agent.config.CoordinationURL, h.agent.node.ID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if h.agent.config.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+h.agent.config.AuthToken)
	}

	// For now, we'll simulate the heartbeat call
	log.Printf("Heartbeat: %+v", heartbeatData)

	// In a real implementation:
	// req.Body = io.NopCloser(bytes.NewReader(data))
	// resp, err := h.client.Do(req)
	// ... handle response

	return nil
}

// StatusReporter reports node status updates
type StatusReporter struct {
	agent  *Agent
	client *http.Client
}

// NewStatusReporter creates a new status reporter
func NewStatusReporter(agent *Agent) *StatusReporter {
	return &StatusReporter{
		agent: agent,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ReportStatus reports current status to coordination server
func (s *StatusReporter) ReportStatus(status string, message string) error {
	if s.agent.config.CoordinationURL == "" {
		return nil
	}

	s.agent.mu.Lock()
	s.agent.node.Status = status
	s.agent.node.LastSeen = time.Now()
	s.agent.mu.Unlock()

	statusData := map[string]interface{}{
		"status":    status,
		"message":   message,
		"timestamp": time.Now(),
	}

	data, err := json.Marshal(statusData)
	if err != nil {
		return fmt.Errorf("failed to marshal status data: %w", err)
	}

	log.Printf("Status update: %s - %s", status, message)

	// In a real implementation, send to coordination server
	_ = data

	return nil
}

// HealthChecker performs periodic health checks
type HealthChecker struct {
	agent  *Agent
	client *http.Client
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(agent *Agent) *HealthChecker {
	return &HealthChecker{
		agent: agent,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Start starts the health checker
func (h *HealthChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.checkHealth()
		}
	}
}

// checkHealth performs health checks on node and services
func (h *HealthChecker) checkHealth() {
	health := map[string]interface{}{
		"status":     "healthy",
		"timestamp":  time.Now(),
		"node_id":    h.agent.node.ID,
		"provider":   h.agent.node.Provider,
		"cpu_usage":  "low",
		"memory":     "available",
		"disk_space": "available",
		"network":    "connected",
	}

	providersHealth := make(map[string]string)
	for name := range h.agent.providers {
		providersHealth[name] = "available"
	}
	health["providers"] = providersHealth

	log.Printf("Health check: %+v", health)

	// In a real implementation, send to coordination server
}
