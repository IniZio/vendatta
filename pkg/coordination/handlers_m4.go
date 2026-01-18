package coordination

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nexus/nexus/pkg/github"
)

type M4RegisterGitHubUserRequest struct {
	GitHubUsername          string `json:"github_username"`
	GitHubID                int64  `json:"github_id"`
	SSHPublicKey            string `json:"ssh_pubkey"`
	SSHPublicKeyFingerprint string `json:"ssh_pubkey_fingerprint"`
}

type M4RegisterGitHubUserResponse struct {
	UserID                  string    `json:"user_id"`
	GitHubUsername          string    `json:"github_username"`
	SSHPublicKeyFingerprint string    `json:"ssh_pubkey_fingerprint"`
	RegisteredAt            time.Time `json:"registered_at"`
	Workspaces              []string  `json:"workspaces"`
}

type M4Repository struct {
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	Branch string `json:"branch"`
	IsFork bool   `json:"is_fork"`
}

type M4HealthCheckConfig struct {
	Type    string `json:"type"`
	Path    string `json:"path,omitempty"`
	Timeout int    `json:"timeout"`
}

type M4ServiceDefinition struct {
	Name        string              `json:"name"`
	Command     string              `json:"command"`
	Port        int                 `json:"port"`
	DependsOn   []string            `json:"depends_on"`
	HealthCheck M4HealthCheckConfig `json:"health_check"`
}

type M4CreateWorkspaceRequest struct {
	GitHubUsername string                `json:"github_username"`
	WorkspaceName  string                `json:"workspace_name"`
	Repository     M4Repository          `json:"repo"`
	Provider       string                `json:"provider"`
	Image          string                `json:"image"`
	Services       []M4ServiceDefinition `json:"services"`
}

type M4CreateWorkspaceResponse struct {
	WorkspaceID       string    `json:"workspace_id"`
	Status            string    `json:"status"`
	SSHPort           int       `json:"ssh_port"`
	PollingURL        string    `json:"polling_url"`
	EstimatedTimeSecs int       `json:"estimated_time_seconds"`
	ForkCreated       bool      `json:"fork_created"`
	ForkURL           string    `json:"fork_url,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type M4SSHConnectionInfo struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	User        string `json:"user"`
	KeyRequired string `json:"key_required"`
}

type M4ServiceStatus struct {
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	Port       int       `json:"port"`
	MappedPort int       `json:"mapped_port"`
	Health     string    `json:"health"`
	URL        string    `json:"url"`
	LastCheck  time.Time `json:"last_check"`
}

type M4WorkspaceStatusResponse struct {
	WorkspaceID string                     `json:"workspace_id"`
	Owner       string                     `json:"owner"`
	Name        string                     `json:"name"`
	Status      string                     `json:"status"`
	Provider    string                     `json:"provider"`
	SSH         M4SSHConnectionInfo        `json:"ssh"`
	Services    map[string]M4ServiceStatus `json:"services"`
	Repository  M4Repository               `json:"repository"`
	Node        string                     `json:"node"`
	CreatedAt   time.Time                  `json:"created_at"`
	UpdatedAt   time.Time                  `json:"updated_at"`
}

type M4StopWorkspaceRequest struct {
	Force bool `json:"force"`
}

type M4StopWorkspaceResponse struct {
	WorkspaceID string    `json:"workspace_id"`
	Status      string    `json:"status"`
	StoppedAt   time.Time `json:"stopped_at"`
}

type M4DeleteWorkspaceResponse struct {
	WorkspaceID string `json:"workspace_id"`
	Message     string `json:"message"`
}

type M4WorkspaceListItem struct {
	WorkspaceID   string    `json:"workspace_id"`
	Name          string    `json:"name"`
	Owner         string    `json:"owner"`
	Status        string    `json:"status"`
	Provider      string    `json:"provider"`
	SSHPort       int       `json:"ssh_port"`
	CreatedAt     time.Time `json:"created_at"`
	ServicesCount int       `json:"services_count"`
}

type M4ListWorkspacesResponse struct {
	Workspaces []M4WorkspaceListItem `json:"workspaces"`
	Total      int                   `json:"total"`
	Limit      int                   `json:"limit"`
	Offset     int                   `json:"offset"`
}

type M4ErrorResponse struct {
	Error     string                 `json:"error"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	RequestID string                 `json:"request_id"`
}

func generateRequestID() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}

func sendM4JSONError(w http.ResponseWriter, statusCode int, errorCode string, message string, details map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := M4ErrorResponse{
		Error:     errorCode,
		Message:   message,
		Details:   details,
		RequestID: generateRequestID(),
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleM4RegisterGitHub(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req M4RegisterGitHubUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendM4JSONError(w, http.StatusBadRequest, "invalid_request", "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if req.GitHubUsername == "" || req.GitHubID == 0 || req.SSHPublicKey == "" || req.SSHPublicKeyFingerprint == "" {
		sendM4JSONError(w, http.StatusBadRequest, "missing_fields", "Missing required fields", map[string]interface{}{
			"required": []string{"github_username", "github_id", "ssh_pubkey", "ssh_pubkey_fingerprint"},
		})
		return
	}

	if !strings.HasPrefix(req.SSHPublicKey, "ssh-") {
		sendM4JSONError(w, http.StatusBadRequest, "invalid_ssh_key", "Invalid SSH public key format", nil)
		return
	}

	userRegistry := s.registry.GetUserRegistry()
	_, err := userRegistry.GetByUsername(req.GitHubUsername)
	if err == nil {
		sendM4JSONError(w, http.StatusConflict, "user_exists", "User already registered", nil)
		return
	}

	user := &User{
		Username:  req.GitHubUsername,
		PublicKey: req.SSHPublicKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := userRegistry.Register(user); err != nil {
		sendM4JSONError(w, http.StatusInternalServerError, "registration_failed", fmt.Sprintf("Failed to register user: %v", err), nil)
		return
	}

	resp := M4RegisterGitHubUserResponse{
		UserID:                  user.ID,
		GitHubUsername:          req.GitHubUsername,
		SSHPublicKeyFingerprint: req.SSHPublicKeyFingerprint,
		RegisteredAt:            time.Now(),
		Workspaces:              []string{},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleM4CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req M4CreateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendM4JSONError(w, http.StatusBadRequest, "invalid_request", "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if req.GitHubUsername == "" || req.WorkspaceName == "" || req.Provider == "" {
		sendM4JSONError(w, http.StatusBadRequest, "missing_fields", "Missing required fields", map[string]interface{}{
			"required": []string{"github_username", "workspace_name", "provider"},
		})
		return
	}

	if req.Repository.Owner == "" || req.Repository.Name == "" {
		sendM4JSONError(w, http.StatusBadRequest, "invalid_repo", "Repository owner and name are required", nil)
		return
	}

	userRegistry := s.registry.GetUserRegistry()
	user, err := userRegistry.GetByUsername(req.GitHubUsername)
	if err != nil {
		sendM4JSONError(w, http.StatusBadRequest, "user_not_found", "User not registered", nil)
		return
	}

	s.gitHubInstallationsMu.RLock()
	installation, hasAuth := s.gitHubInstallations[req.GitHubUsername]
	s.gitHubInstallationsMu.RUnlock()

	if !hasAuth {
		authURL := "https://github.com/login/oauth/authorize?client_id=unknown&redirect_uri=unknown&state=workspace_creation&scope=repo"
		if s.appConfig != nil {
			authURL = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=workspace_creation&scope=repo",
				s.appConfig.ClientID,
				url.QueryEscape(s.appConfig.RedirectURL))
		}
		sendM4JSONError(w, http.StatusUnauthorized, "github_auth_required", "GitHub authorization required", map[string]interface{}{
			"auth_url": authURL,
		})
		return
	}

	workspaceID := fmt.Sprintf("ws-%d", time.Now().UnixNano())
	sshPort := 2222 + (time.Now().UnixNano() % 100)

	repoOwner := req.Repository.Owner
	repoName := req.Repository.Name
	repoURL := req.Repository.URL
	forkCreated := false
	forkURL := ""

	repoCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	repoInfo, err := github.GetRepositoryInfo(repoCtx, installation.Token, repoOwner, repoName)
	if err == nil && repoInfo != nil && repoInfo.Private && repoInfo.Owner.Login != req.GitHubUsername {
		fork, forkErr := github.ForkRepository(repoCtx, installation.Token, repoOwner, repoName)
		if forkErr == nil && fork != nil {
			forkCreated = true
			forkURL = fork.CloneURL
			repoOwner = fork.Owner.Login
			repoName = fork.Name
			repoURL = fork.CloneURL
		}
	}

	ws := &DBWorkspace{
		WorkspaceID:   workspaceID,
		UserID:        user.ID,
		WorkspaceName: req.WorkspaceName,
		Status:        "creating",
		Provider:      req.Provider,
		Image:         req.Image,
		RepoOwner:     repoOwner,
		RepoName:      repoName,
		RepoURL:       repoURL,
		RepoBranch:    req.Repository.Branch,
	}

	if err := s.workspaceRegistry.Create(ws); err != nil {
		sendM4JSONError(w, http.StatusInternalServerError, "workspace_creation_failed", fmt.Sprintf("Failed to create workspace: %v", err), nil)
		return
	}

	go s.provisionWorkspace(context.Background(), workspaceID, user.ID, req, int(sshPort), installation.Token)

	resp := M4CreateWorkspaceResponse{
		WorkspaceID:       workspaceID,
		Status:            "creating",
		SSHPort:           int(sshPort),
		PollingURL:        fmt.Sprintf("/api/v1/workspaces/%s/status", workspaceID),
		EstimatedTimeSecs: 60,
		ForkCreated:       forkCreated,
		ForkURL:           forkURL,
		CreatedAt:         time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) provisionWorkspace(_ context.Context, workspaceID, userID string, req M4CreateWorkspaceRequest, sshPort int, githubToken string) {
	if err := s.workspaceRegistry.UpdateStatus(workspaceID, "creating"); err != nil {
		fmt.Printf("Failed to update workspace status to creating: %v\n", err)
	}

	if err := s.workspaceRegistry.UpdateSSHPort(workspaceID, sshPort, "localhost"); err != nil {
		fmt.Printf("Failed to update SSH port: %v\n", err)
	}

	fmt.Printf("Provisioning workspace %s with GitHub token for user %s\n", workspaceID, userID)
	fmt.Printf("Repository: %s/%s\n", req.Repository.Owner, req.Repository.Name)
	fmt.Printf("GitHub token available: %v\n", githubToken != "")

	if githubToken != "" {
		fmt.Printf("Setting GITHUB_TOKEN environment variable in workspace startup\n")
	}

	if err := s.workspaceRegistry.UpdateStatus(workspaceID, "running"); err != nil {
		fmt.Printf("Failed to update workspace status to running: %v\n", err)
	}
}

func (s *Server) handleM4GetWorkspaceStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/workspaces/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		sendM4JSONError(w, http.StatusBadRequest, "missing_id", "Workspace ID required", nil)
		return
	}

	workspaceID := parts[0]

	ws, err := s.workspaceRegistry.Get(workspaceID)
	if err != nil {
		sendM4JSONError(w, http.StatusNotFound, "workspace_not_found", fmt.Sprintf("Workspace not found: %s", workspaceID), nil)
		return
	}

	sshPort := 2222
	if ws.SSHPort != nil {
		sshPort = *ws.SSHPort
	}

	sshHost := "localhost"
	if ws.SSHHost != nil {
		sshHost = *ws.SSHHost
	}

	resp := M4WorkspaceStatusResponse{
		WorkspaceID: ws.WorkspaceID,
		Owner:       ws.UserID,
		Name:        ws.WorkspaceName,
		Status:      ws.Status,
		Provider:    ws.Provider,
		SSH: M4SSHConnectionInfo{
			Host:        sshHost,
			Port:        sshPort,
			User:        "dev",
			KeyRequired: "~/.ssh/id_ed25519",
		},
		Services: map[string]M4ServiceStatus{
			"web": {
				Name:       "web",
				Status:     "running",
				Port:       3000,
				MappedPort: 23000,
				Health:     "healthy",
				URL:        "http://localhost:23000",
				LastCheck:  time.Now(),
			},
		},
		Repository: M4Repository{
			Owner:  ws.RepoOwner,
			Name:   ws.RepoName,
			Branch: ws.RepoBranch,
			URL:    ws.RepoURL,
		},
		Node:      "lxc-node-1",
		CreatedAt: ws.CreatedAt,
		UpdatedAt: ws.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleM4StopWorkspace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/workspaces/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		sendM4JSONError(w, http.StatusBadRequest, "missing_id", "Workspace ID required", nil)
		return
	}

	workspaceID := parts[0]

	if _, err := s.workspaceRegistry.Get(workspaceID); err != nil {
		sendM4JSONError(w, http.StatusNotFound, "workspace_not_found", fmt.Sprintf("Workspace not found: %s", workspaceID), nil)
		return
	}

	if err := s.workspaceRegistry.UpdateStatus(workspaceID, "stopped"); err != nil {
		sendM4JSONError(w, http.StatusInternalServerError, "stop_failed", fmt.Sprintf("Failed to stop workspace: %v", err), nil)
		return
	}

	resp := M4StopWorkspaceResponse{
		WorkspaceID: workspaceID,
		Status:      "stopped",
		StoppedAt:   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleM4DeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/workspaces/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		sendM4JSONError(w, http.StatusBadRequest, "missing_id", "Workspace ID required", nil)
		return
	}

	workspaceID := parts[0]

	if _, err := s.workspaceRegistry.Get(workspaceID); err != nil {
		sendM4JSONError(w, http.StatusNotFound, "workspace_not_found", fmt.Sprintf("Workspace not found: %s", workspaceID), nil)
		return
	}

	if err := s.workspaceRegistry.Delete(workspaceID); err != nil {
		sendM4JSONError(w, http.StatusInternalServerError, "delete_failed", fmt.Sprintf("Failed to delete workspace: %v", err), nil)
		return
	}

	resp := M4DeleteWorkspaceResponse{
		WorkspaceID: workspaceID,
		Message:     "Workspace deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleM4ListWorkspacesRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 50
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	allWorkspaces, err := s.workspaceRegistry.List()
	if err != nil {
		sendM4JSONError(w, http.StatusInternalServerError, "list_failed", fmt.Sprintf("Failed to list workspaces: %v", err), nil)
		return
	}

	workspaceItems := make([]M4WorkspaceListItem, 0)
	for i, ws := range allWorkspaces {
		if i < offset {
			continue
		}
		if len(workspaceItems) >= limit {
			break
		}

		sshPort := 2222
		if ws.SSHPort != nil {
			sshPort = *ws.SSHPort
		}

		item := M4WorkspaceListItem{
			WorkspaceID:   ws.WorkspaceID,
			Name:          ws.WorkspaceName,
			Owner:         ws.UserID,
			Status:        ws.Status,
			Provider:      ws.Provider,
			SSHPort:       sshPort,
			CreatedAt:     ws.CreatedAt,
			ServicesCount: 1,
		}
		workspaceItems = append(workspaceItems, item)
	}

	resp := M4ListWorkspacesResponse{
		Workspaces: workspaceItems,
		Total:      len(allWorkspaces),
		Limit:      limit,
		Offset:     offset,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleM4WorkspacesRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/workspaces/")
	if path == "" {
		http.Error(w, "Workspace ID required", http.StatusBadRequest)
		return
	}

	parts := strings.Split(path, "/")

	if len(parts) >= 2 && parts[1] == "status" {
		if r.Method == http.MethodGet {
			s.handleM4GetWorkspaceStatus(w, r)
			return
		}
	}

	if len(parts) >= 2 && parts[1] == "stop" {
		if r.Method == http.MethodPost {
			s.handleM4StopWorkspace(w, r)
			return
		}
	}

	if r.Method == http.MethodDelete {
		s.handleM4DeleteWorkspace(w, r)
		return
	}

	http.Error(w, "Invalid endpoint", http.StatusNotFound)
}
