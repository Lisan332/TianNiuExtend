package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Container represents a container in the TianNiu platform
type Container struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	StartedAt time.Time `json:"started_at,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
	Ports     []struct {
		Internal int    `json:"internal"`
		External int    `json:"external"`
		Protocol string `json:"protocol"`
	} `json:"ports,omitempty"`
	Volumes []struct {
		HostPath      string `json:"host_path"`
		ContainerPath string `json:"container_path"`
		Mode          string `json:"mode"`
	} `json:"volumes,omitempty"`
	Network struct {
		Name      string `json:"name"`
		IPAddress string `json:"ip_address"`
	} `json:"network,omitempty"`
	ResourceLimits struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
	} `json:"resource_limits,omitempty"`
	ResourceUsage struct {
		CPU       string `json:"cpu"`
		Memory    string `json:"memory"`
		NetworkRX string `json:"network_rx"`
		NetworkTX string `json:"network_tx"`
	} `json:"resource_usage,omitempty"`
	EnvironmentVariables []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"environment_variables,omitempty"`
	HealthCheck struct {
		Status      string    `json:"status"`
		LastChecked time.Time `json:"last_checked"`
		Endpoint    string    `json:"endpoint"`
		Interval    string    `json:"interval"`
		Timeout     string    `json:"timeout"`
		Retries     int       `json:"retries"`
	} `json:"health_check,omitempty"`
	LogsURL string `json:"logs_url,omitempty"`
	Message string `json:"message,omitempty"`
}

// ContainerList represents a list of containers
type ContainerList struct {
	Total      int         `json:"total"`
	Limit      int         `json:"limit"`
	Offset     int         `json:"offset"`
	Containers []Container `json:"containers"`
}

// ContainerClient represents a TianNiu API client for container operations
type ContainerClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewContainerClient creates a new TianNiu API client for container operations
func NewContainerClient(baseURL, apiKey string) *ContainerClient {
	return &ContainerClient{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// ListContainers lists all containers
func (c *ContainerClient) ListContainers(status string, limit, offset int) (*ContainerList, error) {
	url := fmt.Sprintf("%s/containers?limit=%d&offset=%d", c.BaseURL, limit, offset)
	if status != "" {
		url += fmt.Sprintf("&status=%s", status)
	}
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var containerList ContainerList
	if err := json.NewDecoder(resp.Body).Decode(&containerList); err != nil {
		return nil, err
	}

	return &containerList, nil
}

// GetContainer gets a container by ID
func (c *ContainerClient) GetContainer(containerID string) (*Container, error) {
	url := fmt.Sprintf("%s/containers/%s", c.BaseURL, containerID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var container Container
	if err := json.NewDecoder(resp.Body).Decode(&container); err != nil {
		return nil, err
	}

	return &container, nil
}

// CreateContainer creates a new container
func (c *ContainerClient) CreateContainer(container *Container) (*Container, error) {
	url := fmt.Sprintf("%s/containers", c.BaseURL)
	
	containerJSON, err := json.Marshal(container)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(containerJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var createdContainer Container
	if err := json.NewDecoder(resp.Body).Decode(&createdContainer); err != nil {
		return nil, err
	}

	return &createdContainer, nil
}

// StartContainer starts a container
func (c *ContainerClient) StartContainer(containerID string) (*Container, error) {
	url := fmt.Sprintf("%s/containers/%s/start", c.BaseURL, containerID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var container Container
	if err := json.NewDecoder(resp.Body).Decode(&container); err != nil {
		return nil, err
	}

	return &container, nil
}

// StopContainer stops a container
func (c *ContainerClient) StopContainer(containerID string, timeout int) (*Container, error) {
	url := fmt.Sprintf("%s/containers/%s/stop", c.BaseURL, containerID)
	if timeout > 0 {
		url += fmt.Sprintf("?timeout=%d", timeout)
	}
	
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var container Container
	if err := json.NewDecoder(resp.Body).Decode(&container); err != nil {
		return nil, err
	}

	return &container, nil
}

// DeleteContainer deletes a container
func (c *ContainerClient) DeleteContainer(containerID string, force, removeVolumes bool) (*Container, error) {
	url := fmt.Sprintf("%s/containers/%s", c.BaseURL, containerID)
	params := []string{}
	
	if force {
		params = append(params, "force=true")
	}
	
	if removeVolumes {
		params = append(params, "remove_volumes=true")
	}
	
	if len(params) > 0 {
		url += "?" + params[0]
		for i := 1; i < len(params); i++ {
			url += "&" + params[i]
		}
	}
	
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var container Container
	if err := json.NewDecoder(resp.Body).Decode(&container); err != nil {
		return nil, err
	}

	return &container, nil
}

// Mock server for container testing
func setupContainerMockServer() *httptest.Server {
	handler := http.NewServeMux()
	
	// Mock container list endpoint
	handler.HandleFunc("/api/v1/containers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			status := r.URL.Query().Get("status")
			
			containerList := ContainerList{
				Total:  2,
				Limit:  10,
				Offset: 0,
				Containers: []Container{
					{
						ID:        "c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2",
						Name:      "web-server-1",
						Image:     "nginx:latest",
						Status:    "running",
						CreatedAt: time.Now().Add(-24 * time.Hour),
						Labels: map[string]string{
							"app":         "web",
							"environment": "production",
						},
					},
					{
						ID:        "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
						Name:      "api-service",
						Image:     "registry.baidu.com/myteam/api-service:v1.2.3",
						Status:    "running",
						CreatedAt: time.Now().Add(-48 * time.Hour),
						Labels: map[string]string{
							"app":         "api",
							"environment": "production",
							"team":        "backend",
						},
					},
				},
			}
			
			// Filter by status if provided
			if status != "" {
				filteredContainers := []Container{}
				for _, container := range containerList.Containers {
					if container.Status == status {
						filteredContainers = append(filteredContainers, container)
					}
				}
				containerList.Containers = filteredContainers
				containerList.Total = len(filteredContainers)
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(containerList)
			return
		} else if r.Method == "POST" {
			// Handle container creation
			var container Container
			if err := json.NewDecoder(r.Body).Decode(&container); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			
			// Generate a random ID for the new container
			container.ID = "new1c2o3n4t5a6i7n8e9r0id"
			container.Status = "creating"
			container.CreatedAt = time.Now()
			container.Message = "Container is being created and will start shortly"
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(container)
			return
		}
		
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
	
	// Mock container detail endpoint
	handler.HandleFunc("/api/v1/containers/c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			container := Container{
				ID:        "c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2",
				Name:      "web-server-1",
				Image:     "nginx:latest",
				Status:    "running",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				StartedAt: time.Now().Add(-24 * time.Hour).Add(2 * time.Second),
				Labels: map[string]string{
					"app":         "web",
					"environment": "production",
				},
			}
			
			container.Ports = []struct {
				Internal int    `json:"internal"`
				External int    `json:"external"`
				Protocol string `json:"protocol"`
			}{
				{
					Internal: 80,
					External: 8080,
					Protocol: "tcp",
				},
			}
			
			container.Volumes = []struct {
				HostPath      string `json:"host_path"`
				ContainerPath string `json:"container_path"`
				Mode          string `json:"mode"`
			}{
				{
					HostPath:      "/data/nginx/conf",
					ContainerPath: "/etc/nginx/conf.d",
					Mode:          "ro",
				},
				{
					HostPath:      "/data/nginx/html",
					ContainerPath: "/usr/share/nginx/html",
					Mode:          "rw",
				},
			}
			
			container.Network.Name = "frontend-network"
			container.Network.IPAddress = "172.18.0.2"
			
			container.ResourceLimits.CPU = "1.0"
			container.ResourceLimits.Memory = "512MB"
			
			container.ResourceUsage.CPU = "0.05"
			container.ResourceUsage.Memory = "128MB"
			container.ResourceUsage.NetworkRX = "1.2MB/s"
			container.ResourceUsage.NetworkTX = "0.8MB/s"
			
			container.EnvironmentVariables = []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			}{
				{
					Name:  "NGINX_HOST",
					Value: "example.com",
				},
				{
					Name:  "NGINX_PORT",
					Value: "80",
				},
			}
			
			container.HealthCheck.Status = "healthy"
			container.HealthCheck.LastChecked = time.Now().Add(-15 * time.Minute)
			container.HealthCheck.Endpoint = "http://localhost:80/health"
			container.HealthCheck.Interval = "30s"
			container.HealthCheck.Timeout = "5s"
			container.HealthCheck.Retries = 3
			
			container.LogsURL = "https://tianniuprod.baidu.com/api/v1/containers/c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2/logs"
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(container)
			return
		}
		
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
	
	// Mock container start endpoint
	handler.HandleFunc("/api/v1/containers/c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			response := struct {
				ID      string `json:"id"`
				Status  string `json:"status"`
				Message string `json:"message"`
			}{
				ID:      "c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2",
				Status:  "starting",
				Message: "Container is starting",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
	
	// Mock container stop endpoint
	handler.HandleFunc("/api/v1/containers/c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2/stop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			response := struct {
				ID      string `json:"id"`
				Status  string `json:"status"`
				Message string `json:"message"`
			}{
				ID:      "c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2",
				Status:  "stopping",
				Message: "Container is stopping",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
	
	// Mock container delete endpoint
	handler.HandleFunc("/api/v1/containers/c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			response := struct {
				ID      string `json:"id"`
				Status  string `json:"status"`
				Message string `json:"message"`
			}{
				ID:      "c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2",
				Status:  "deleted",
				Message: "Container has been deleted",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		// GET method is handled by the other handler
	})
	
	return httptest.NewServer(handler)
}

// Test ListContainers
func TestListContainers(t *testing.T) {
	server := setupContainerMockServer()
	defer server.Close()
	
	client := NewContainerClient(server.URL+"/api/v1", "test-api-key")
	
	containerList, err := client.ListContainers("", 10, 0)
	if err != nil {
		t.Fatalf("ListContainers failed: %v", err)
	}
	
	if containerList.Total != 2 {
		t.Errorf("Expected 2 containers, got %d", containerList.Total)
	}
	
	if len(containerList.Containers) != 2 {
		t.Errorf("Expected 2 containers in the list, got %d", len(containerList.Containers))
	}
	
	if containerList.Containers[0].Name != "web-server-1" {
		t.Errorf("Expected first container name to be 'web-server-1', got '%s'", containerList.Containers[0].Name)
	}
	
	if containerList.Containers[1].Name != "api-service" {
		t.Errorf("Expected second container name to be 'api-service', got '%s'", containerList.Containers[1].Name)
	}
	
	// Test filtering by status
	containerList, err = client.ListContainers("running", 10, 0)
	if err != nil {
		t.Fatalf("ListContainers with status filter failed: %v", err)
	}
	
	for _, container := range containerList.Containers {
		if container.Status != "running" {
			t.Errorf("Expected all containers to have status 'running', got '%s'", container.Status)
		}
	}
}

// Test GetContainer
func TestGetContainer(t *testing.T) {
	server := setupContainerMockServer()
	defer server.Close()
	
	client := NewContainerClient(server.URL+"/api/v1", "test-api-key")
	
	container, err := client.GetContainer("c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2")
	if err != nil {
		t.Fatalf("GetContainer failed: %v", err)
	}
	
	if container.ID != "c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2" {
		t.Errorf("Expected container ID to be 'c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2', got '%s'", container.ID)
	}
	
	if container.Name != "web-server-1" {
		t.Errorf("Expected container name to be 'web-server-1', got '%s'", container.Name)
	}
	
	if container.Image != "nginx:latest" {
		t.Errorf("Expected container image to be 'nginx:latest', got '%s'", container.Image)
	}
	
	if container.Status != "running" {
		t.Errorf("Expected container status to be 'running', got '%s'", container.Status)
	}
	
	if len(container.Ports) != 1 {
		t.Errorf("Expected 1 port mapping, got %d", len(container.Ports))
	}
	
	if container.Ports[0].Internal != 80 || container.Ports[0].External != 8080 {
		t.Errorf("Expected port mapping 80:8080, got %d:%d", container.Ports[0].Internal, container.Ports[0].External)
	}
	
	if len(container.Volumes) != 2 {
		t.Errorf("Expected 2 volume mappings, got %d", len(container.Volumes))
	}
	
	if container.Network.Name != "frontend-network" {
		t.Errorf("Expected network name to be 'frontend-network', got '%s'", container.Network.Name)
	}
	
	if container.Network.IPAddress != "172.18.0.2" {
		t.Errorf("Expected IP address to be '172.18.0.2', got '%s'", container.Network.IPAddress)
	}
	
	if container.ResourceLimits.CPU != "1.0" {
		t.Errorf("Expected CPU limit to be '1.0', got '%s'", container.ResourceLimits.CPU)
	}
	
	if container.ResourceLimits.Memory != "512MB" {
		t.Errorf("Expected memory limit to be '512MB', got '%s'", container.ResourceLimits.Memory)
	}
	
	if len(container.EnvironmentVariables) != 2 {
		t.Errorf("Expected 2 environment variables, got %d", len(container.EnvironmentVariables))
	}
	
	if container.HealthCheck.Status != "healthy" {
		t.Errorf("Expected health check status to be 'healthy', got '%s'", container.HealthCheck.Status)
	}
}

// Test CreateContainer
func TestCreateContainer(t *testing.T) {
	server := setupContainerMockServer()
	defer server.Close()
	
	client := NewContainerClient(server.URL+"/api/v1", "test-api-key")
	
	newContainer := &Container{
		Name:  "test-container",
		Image: "ubuntu:latest",
		Labels: map[string]string{
			"app":         "test",
			"environment": "development",
		},
	}
	
	createdContainer, err := client.CreateContainer(newContainer)
	if err != nil {
		t.Fatalf("CreateContainer failed: %v", err)
	}
	
	if createdContainer.ID == "" {
		t.Error("Expected container ID to be set")
	}
	
	if createdContainer.Name != "test-container" {
		t.Errorf("Expected container name to be 'test-container', got '%s'", createdContainer.Name)
	}
	
	if createdContainer.Image != "ubuntu:latest" {
		t.Errorf("Expected container image to be 'ubuntu:latest', got '%s'", createdContainer.Image)
	}
	
	if createdContainer.Status != "creating" {
		t.Errorf("Expected container status to be 'creating', got '%s'", createdContainer.Status)
	}
}

// Test StartContainer
func TestStartContainer(t *testing.T) {
	server := setupContainerMockServer()
	defer server.Close()
	
	client := NewContainerClient(server.URL+"/api/v1", "test-api-key")
	
	container, err := client.StartContainer("c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2")
	if err != nil {
		t.Fatalf("StartContainer failed: %v", err)
	}
	
	if container.ID != "c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2" {
		t.Errorf("Expected container ID to be 'c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2', got '%s'", container.ID)
	}
	
	if container.Status != "starting" {
		t.Errorf("Expected container status to be 'starting', got '%s'", container.Status)
	}
}

// Test StopContainer
func TestStopContainer(t *testing.T) {
	server := setupContainerMockServer()
	defer server.Close()
	
	client := NewContainerClient(server.URL+"/api/v1", "test-api-key")
	
	container, err := client.StopContainer("c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2", 30)
	if err != nil {
		t.Fatalf("StopContainer failed: %v", err)
	}
	
	if container.ID != "c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2" {
		t.Errorf("Expected container ID to be 'c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2', got '%s'", container.ID)
	}
	
	if container.Status != "stopping" {
		t.Errorf("Expected container status to be 'stopping', got '%s'", container.Status)
	}
}

// Test DeleteContainer
func TestDeleteContainer(t *testing.T) {
	server := setupContainerMockServer()
	defer server.Close()
	
	client := NewContainerClient(server.URL+"/api/v1", "test-api-key")
	
	container, err := client.DeleteContainer("c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2", true, true)
	if err != nil {
		t.Fatalf("DeleteContainer failed: %v", err)
	}
	
	if container.ID != "c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2" {
		t.Errorf("Expected container ID to be 'c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2', got '%s'", container.ID)
	}
	
	if container.Status != "deleted" {
		t.Errorf("Expected container status to be 'deleted', got '%s'", container.Status)
	}
}
