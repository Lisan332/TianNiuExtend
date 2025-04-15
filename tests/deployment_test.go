package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Deployment represents a deployment in the TianNiu platform
type Deployment struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	Environment string    `json:"environment"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Version     string    `json:"version"`
	Replicas    int       `json:"replicas"`
	Strategy    struct {
		Type           string `json:"type"`
		MaxSurge       int    `json:"max_surge,omitempty"`
		MaxUnavailable int    `json:"max_unavailable,omitempty"`
	} `json:"strategy,omitempty"`
	Containers []struct {
		Name      string `json:"name"`
		Image     string `json:"image"`
		Ports     []struct {
			Name          string `json:"name"`
			ContainerPort int    `json:"container_port"`
			ServicePort   int    `json:"service_port"`
		} `json:"ports,omitempty"`
		Resources struct {
			Limits struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"limits,omitempty"`
			Requests struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"requests,omitempty"`
		} `json:"resources,omitempty"`
		EnvironmentVariables []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"environment_variables,omitempty"`
		HealthCheck struct {
			HTTPPath            string `json:"http_path"`
			Port                int    `json:"port"`
			InitialDelaySeconds int    `json:"initial_delay_seconds"`
			PeriodSeconds       int    `json:"period_seconds"`
			TimeoutSeconds      int    `json:"timeout_seconds"`
			SuccessThreshold    int    `json:"success_threshold"`
			FailureThreshold    int    `json:"failure_threshold"`
		} `json:"health_check,omitempty"`
	} `json:"containers,omitempty"`
	Services []struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Ports    []struct {
			Name       string `json:"name"`
			Port       int    `json:"port"`
			TargetPort int    `json:"target_port"`
		} `json:"ports,omitempty"`
		ExternalEndpoints []string `json:"external_endpoints,omitempty"`
	} `json:"services,omitempty"`
	Message string `json:"message,omitempty"`
}

// DeploymentList represents a list of deployments
type DeploymentList struct {
	Total       int          `json:"total"`
	Limit       int          `json:"limit"`
	Offset      int          `json:"offset"`
	Deployments []Deployment `json:"deployments"`
}

// Client represents a TianNiu API client
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new TianNiu API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// ListDeployments lists all deployments
func (c *Client) ListDeployments(environment string, limit, offset int) (*DeploymentList, error) {
	url := fmt.Sprintf("%s/deployments?environment=%s&limit=%d&offset=%d", c.BaseURL, environment, limit, offset)
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

	var deploymentList DeploymentList
	if err := json.NewDecoder(resp.Body).Decode(&deploymentList); err != nil {
		return nil, err
	}

	return &deploymentList, nil
}

// GetDeployment gets a deployment by ID
func (c *Client) GetDeployment(deploymentID string) (*Deployment, error) {
	url := fmt.Sprintf("%s/deployments/%s", c.BaseURL, deploymentID)
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

	var deployment Deployment
	if err := json.NewDecoder(resp.Body).Decode(&deployment); err != nil {
		return nil, err
	}

	return &deployment, nil
}

// Mock server for testing
func setupMockServer() *httptest.Server {
	handler := http.NewServeMux()
	
	// Mock deployment list endpoint
	handler.HandleFunc("/api/v1/deployments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		environment := r.URL.Query().Get("environment")
		if environment == "" {
			http.Error(w, "Missing environment parameter", http.StatusBadRequest)
			return
		}
		
		deploymentList := DeploymentList{
			Total:  2,
			Limit:  10,
			Offset: 0,
			Deployments: []Deployment{
				{
					ID:          "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
					Name:        "web-frontend",
					Status:      "active",
					Environment: environment,
					CreatedAt:   time.Now().Add(-24 * time.Hour),
					UpdatedAt:   time.Now().Add(-12 * time.Hour),
					Version:     "v2.3.1",
					Replicas:    3,
				},
				{
					ID:          "f1e2d3c4b5a6f7e8d9c0b1a2d3e4f5c6",
					Name:        "api-backend",
					Status:      "active",
					Environment: environment,
					CreatedAt:   time.Now().Add(-48 * time.Hour),
					UpdatedAt:   time.Now().Add(-24 * time.Hour),
					Version:     "v1.5.1",
					Replicas:    2,
				},
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deploymentList)
	})
	
	// Mock deployment detail endpoint
	handler.HandleFunc("/api/v1/deployments/d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		deployment := Deployment{
			ID:          "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
			Name:        "web-frontend",
			Description: "Web前端应用",
			Status:      "active",
			Environment: "production",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now().Add(-12 * time.Hour),
			Version:     "v2.3.1",
			Replicas:    3,
		}
		
		deployment.Strategy.Type = "rolling-update"
		deployment.Strategy.MaxSurge = 1
		deployment.Strategy.MaxUnavailable = 0
		
		container := struct {
			Name      string `json:"name"`
			Image     string `json:"image"`
			Ports     []struct {
				Name          string `json:"name"`
				ContainerPort int    `json:"container_port"`
				ServicePort   int    `json:"service_port"`
			} `json:"ports,omitempty"`
			Resources struct {
				Limits struct {
					CPU    string `json:"cpu"`
					Memory string `json:"memory"`
				} `json:"limits,omitempty"`
				Requests struct {
					CPU    string `json:"cpu"`
					Memory string `json:"memory"`
				} `json:"requests,omitempty"`
			} `json:"resources,omitempty"`
			EnvironmentVariables []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"environment_variables,omitempty"`
			HealthCheck struct {
				HTTPPath            string `json:"http_path"`
				Port                int    `json:"port"`
				InitialDelaySeconds int    `json:"initial_delay_seconds"`
				PeriodSeconds       int    `json:"period_seconds"`
				TimeoutSeconds      int    `json:"timeout_seconds"`
				SuccessThreshold    int    `json:"success_threshold"`
				FailureThreshold    int    `json:"failure_threshold"`
			} `json:"health_check,omitempty"`
		}{
			Name:  "web-frontend",
			Image: "registry.baidu.com/frontend/web-app:v2.3.1",
		}
		
		container.Ports = []struct {
			Name          string `json:"name"`
			ContainerPort int    `json:"container_port"`
			ServicePort   int    `json:"service_port"`
		}{
			{
				Name:          "http",
				ContainerPort: 80,
				ServicePort:   8080,
			},
		}
		
		container.Resources.Limits.CPU = "1.0"
		container.Resources.Limits.Memory = "1Gi"
		container.Resources.Requests.CPU = "0.5"
		container.Resources.Requests.Memory = "512Mi"
		
		container.EnvironmentVariables = []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{
			{
				Name:  "API_ENDPOINT",
				Value: "https://api.palo.prod.baidu.com/v1",
			},
			{
				Name:  "LOG_LEVEL",
				Value: "info",
			},
		}
		
		container.HealthCheck.HTTPPath = "/health"
		container.HealthCheck.Port = 80
		container.HealthCheck.InitialDelaySeconds = 10
		container.HealthCheck.PeriodSeconds = 30
		container.HealthCheck.TimeoutSeconds = 5
		container.HealthCheck.SuccessThreshold = 1
		container.HealthCheck.FailureThreshold = 3
		
		deployment.Containers = append(deployment.Containers, container)
		
		service := struct {
			Name     string `json:"name"`
			Type     string `json:"type"`
			Ports    []struct {
				Name       string `json:"name"`
				Port       int    `json:"port"`
				TargetPort int    `json:"target_port"`
			} `json:"ports,omitempty"`
			ExternalEndpoints []string `json:"external_endpoints,omitempty"`
		}{
			Name: "web-frontend-svc",
			Type: "LoadBalancer",
		}
		
		service.Ports = []struct {
			Name       string `json:"name"`
			Port       int    `json:"port"`
			TargetPort int    `json:"target_port"`
		}{
			{
				Name:       "http",
				Port:       80,
				TargetPort: 8080,
			},
		}
		
		service.ExternalEndpoints = []string{"web-frontend.palo.prod.baidu.com"}
		
		deployment.Services = append(deployment.Services, service)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deployment)
	})
	
	return httptest.NewServer(handler)
}

// Test ListDeployments
func TestListDeployments(t *testing.T) {
	server := setupMockServer()
	defer server.Close()
	
	client := NewClient(server.URL+"/api/v1", "test-api-key")
	
	deploymentList, err := client.ListDeployments("production", 10, 0)
	if err != nil {
		t.Fatalf("ListDeployments failed: %v", err)
	}
	
	if deploymentList.Total != 2 {
		t.Errorf("Expected 2 deployments, got %d", deploymentList.Total)
	}
	
	if len(deploymentList.Deployments) != 2 {
		t.Errorf("Expected 2 deployments in the list, got %d", len(deploymentList.Deployments))
	}
	
	if deploymentList.Deployments[0].Name != "web-frontend" {
		t.Errorf("Expected first deployment name to be 'web-frontend', got '%s'", deploymentList.Deployments[0].Name)
	}
	
	if deploymentList.Deployments[1].Name != "api-backend" {
		t.Errorf("Expected second deployment name to be 'api-backend', got '%s'", deploymentList.Deployments[1].Name)
	}
}

// Test GetDeployment
func TestGetDeployment(t *testing.T) {
	server := setupMockServer()
	defer server.Close()
	
	client := NewClient(server.URL+"/api/v1", "test-api-key")
	
	deployment, err := client.GetDeployment("d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6")
	if err != nil {
		t.Fatalf("GetDeployment failed: %v", err)
	}
	
	if deployment.ID != "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6" {
		t.Errorf("Expected deployment ID to be 'd1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6', got '%s'", deployment.ID)
	}
	
	if deployment.Name != "web-frontend" {
		t.Errorf("Expected deployment name to be 'web-frontend', got '%s'", deployment.Name)
	}
	
	if deployment.Environment != "production" {
		t.Errorf("Expected deployment environment to be 'production', got '%s'", deployment.Environment)
	}
	
	if deployment.Version != "v2.3.1" {
		t.Errorf("Expected deployment version to be 'v2.3.1', got '%s'", deployment.Version)
	}
	
	if deployment.Replicas != 3 {
		t.Errorf("Expected deployment replicas to be 3, got %d", deployment.Replicas)
	}
	
	if len(deployment.Containers) != 1 {
		t.Errorf("Expected 1 container, got %d", len(deployment.Containers))
	}
	
	if deployment.Containers[0].Name != "web-frontend" {
		t.Errorf("Expected container name to be 'web-frontend', got '%s'", deployment.Containers[0].Name)
	}
	
	if deployment.Containers[0].Image != "registry.baidu.com/frontend/web-app:v2.3.1" {
		t.Errorf("Expected container image to be 'registry.baidu.com/frontend/web-app:v2.3.1', got '%s'", deployment.Containers[0].Image)
	}
	
	if len(deployment.Services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(deployment.Services))
	}
	
	if deployment.Services[0].Name != "web-frontend-svc" {
		t.Errorf("Expected service name to be 'web-frontend-svc', got '%s'", deployment.Services[0].Name)
	}
	
	if deployment.Services[0].Type != "LoadBalancer" {
		t.Errorf("Expected service type to be 'LoadBalancer', got '%s'", deployment.Services[0].Type)
	}
	
	if len(deployment.Services[0].ExternalEndpoints) != 1 {
		t.Errorf("Expected 1 external endpoint, got %d", len(deployment.Services[0].ExternalEndpoints))
	}
	
	if deployment.Services[0].ExternalEndpoints[0] != "web-frontend.palo.prod.baidu.com" {
		t.Errorf("Expected external endpoint to be 'web-frontend.palo.prod.baidu.com', got '%s'", deployment.Services[0].ExternalEndpoints[0])
	}
}
