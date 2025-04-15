package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// TianNiuConfig represents the configuration for the TianNiu platform
type TianNiuConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"metadata"`
	Environments []struct {
		Name        string `yaml:"name"`
		APIEndpoint string `yaml:"api_endpoint"`
		Auth        struct {
			Type       string `yaml:"type"`
			APIKeyEnv  string `yaml:"api_key_env"`
			ClientID   string `yaml:"client_id_env"`
			ClientSecret string `yaml:"client_secret_env"`
		} `yaml:"auth"`
		Kubeconfig string `yaml:"kubeconfig"`
		Default    bool   `yaml:"default"`
	} `yaml:"environments"`
}

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

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Details interface{} `json:"details,omitempty"`
	} `json:"error"`
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

// LoadConfig loads the TianNiu configuration from a YAML file
func LoadConfig(configPath string) (*TianNiuConfig, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config TianNiuConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetDefaultEnvironment returns the default environment from the configuration
func (c *TianNiuConfig) GetDefaultEnvironment() (string, string, error) {
	for _, env := range c.Environments {
		if env.Default {
			return env.Name, env.APIEndpoint, nil
		}
	}
	return "", "", fmt.Errorf("no default environment found in configuration")
}

// GetEnvironment returns the environment with the given name
func (c *TianNiuConfig) GetEnvironment(name string) (string, error) {
	for _, env := range c.Environments {
		if env.Name == name {
			return env.APIEndpoint, nil
		}
	}
	return "", fmt.Errorf("environment %s not found in configuration", name)
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
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("API error: %s", resp.Status)
		}
		return nil, fmt.Errorf("API error: %s - %s", errResp.Error.Code, errResp.Error.Message)
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
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("API error: %s", resp.Status)
		}
		return nil, fmt.Errorf("API error: %s - %s", errResp.Error.Code, errResp.Error.Message)
	}

	var deployment Deployment
	if err := json.NewDecoder(resp.Body).Decode(&deployment); err != nil {
		return nil, err
	}

	return &deployment, nil
}

// CreateDeployment creates a new deployment
func (c *Client) CreateDeployment(deployment *Deployment) (*Deployment, error) {
	url := fmt.Sprintf("%s/deployments", c.BaseURL)
	
	deploymentJSON, err := json.Marshal(deployment)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(deploymentJSON))
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
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("API error: %s", resp.Status)
		}
		return nil, fmt.Errorf("API error: %s - %s", errResp.Error.Code, errResp.Error.Message)
	}

	var createdDeployment Deployment
	if err := json.NewDecoder(resp.Body).Decode(&createdDeployment); err != nil {
		return nil, err
	}

	return &createdDeployment, nil
}

// UpdateDeployment updates an existing deployment
func (c *Client) UpdateDeployment(deploymentID string, deployment *Deployment) (*Deployment, error) {
	url := fmt.Sprintf("%s/deployments/%s", c.BaseURL, deploymentID)
	
	deploymentJSON, err := json.Marshal(deployment)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(deploymentJSON))
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
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("API error: %s", resp.Status)
		}
		return nil, fmt.Errorf("API error: %s - %s", errResp.Error.Code, errResp.Error.Message)
	}

	var updatedDeployment Deployment
	if err := json.NewDecoder(resp.Body).Decode(&updatedDeployment); err != nil {
		return nil, err
	}

	return &updatedDeployment, nil
}

// ScaleDeployment scales a deployment
func (c *Client) ScaleDeployment(deploymentID string, replicas int) (*Deployment, error) {
	url := fmt.Sprintf("%s/deployments/%s/scale", c.BaseURL, deploymentID)
	
	scaleJSON, err := json.Marshal(map[string]int{"replicas": replicas})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(scaleJSON))
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
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("API error: %s", resp.Status)
		}
		return nil, fmt.Errorf("API error: %s - %s", errResp.Error.Code, errResp.Error.Message)
	}

	var scaledDeployment Deployment
	if err := json.NewDecoder(resp.Body).Decode(&scaledDeployment); err != nil {
		return nil, err
	}

	return &scaledDeployment, nil
}

// DeleteDeployment deletes a deployment
func (c *Client) DeleteDeployment(deploymentID string, force bool) error {
	url := fmt.Sprintf("%s/deployments/%s", c.BaseURL, deploymentID)
	if force {
		url += "?force=true"
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("API error: %s", resp.Status)
		}
		return fmt.Errorf("API error: %s - %s", errResp.Error.Code, errResp.Error.Message)
	}

	return nil
}

func main() {
	// Define command line flags
	configPath := flag.String("config", "../../config/tianniu-config.yaml", "Path to TianNiu configuration file")
	environment := flag.String("env", "", "Environment to use (defaults to the default environment in config)")
	
	// Subcommands
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listLimit := listCmd.Int("limit", 20, "Limit number of results")
	listOffset := listCmd.Int("offset", 0, "Offset for pagination")
	
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createFile := createCmd.String("file", "", "Path to deployment JSON file")
	
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateFile := updateCmd.String("file", "", "Path to deployment JSON file")
	
	scaleCmd := flag.NewFlagSet("scale", flag.ExitOnError)
	scaleReplicas := scaleCmd.Int("replicas", 1, "Number of replicas")
	
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteForce := deleteCmd.Bool("force", false, "Force deletion")
	
	flag.Parse()
	
	if len(os.Args) < 2 {
		fmt.Println("Expected 'list', 'get', 'create', 'update', 'scale', or 'delete' subcommand")
		os.Exit(1)
	}
	
	// Load configuration
	config, err := LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}
	
	// Determine environment and API endpoint
	var envName, apiEndpoint string
	if *environment == "" {
		envName, apiEndpoint, err = config.GetDefaultEnvironment()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		envName = *environment
		apiEndpoint, err = config.GetEnvironment(*environment)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
	
	// Get API key from environment variable
	apiKey := os.Getenv("TIANNIU_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: TIANNIU_API_KEY environment variable not set")
		os.Exit(1)
	}
	
	// Create client
	client := NewClient(apiEndpoint, apiKey)
	
	// Handle subcommands
	switch os.Args[1] {
	case "list":
		listCmd.Parse(os.Args[2:])
		deployments, err := client.ListDeployments(envName, *listLimit, *listOffset)
		if err != nil {
			fmt.Printf("Error listing deployments: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Total deployments: %d\n\n", deployments.Total)
		for _, d := range deployments.Deployments {
			fmt.Printf("ID: %s\nName: %s\nStatus: %s\nEnvironment: %s\nVersion: %s\nReplicas: %d\n\n",
				d.ID, d.Name, d.Status, d.Environment, d.Version, d.Replicas)
		}
		
	case "get":
		getCmd.Parse(os.Args[2:])
		if getCmd.NArg() < 1 {
			fmt.Println("Error: deployment ID required")
			os.Exit(1)
		}
		
		deploymentID := getCmd.Arg(0)
		deployment, err := client.GetDeployment(deploymentID)
		if err != nil {
			fmt.Printf("Error getting deployment: %v\n", err)
			os.Exit(1)
		}
		
		deploymentJSON, _ := json.MarshalIndent(deployment, "", "  ")
		fmt.Println(string(deploymentJSON))
		
	case "create":
		createCmd.Parse(os.Args[2:])
		if *createFile == "" {
			fmt.Println("Error: deployment file required")
			os.Exit(1)
		}
		
		deploymentData, err := ioutil.ReadFile(*createFile)
		if err != nil {
			fmt.Printf("Error reading deployment file: %v\n", err)
			os.Exit(1)
		}
		
		var deployment Deployment
		if err := json.Unmarshal(deploymentData, &deployment); err != nil {
			fmt.Printf("Error parsing deployment JSON: %v\n", err)
			os.Exit(1)
		}
		
		createdDeployment, err := client.CreateDeployment(&deployment)
		if err != nil {
			fmt.Printf("Error creating deployment: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Deployment created successfully:")
		deploymentJSON, _ := json.MarshalIndent(createdDeployment, "", "  ")
		fmt.Println(string(deploymentJSON))
		
	case "update":
		updateCmd.Parse(os.Args[2:])
		if updateCmd.NArg() < 1 {
			fmt.Println("Error: deployment ID required")
			os.Exit(1)
		}
		if *updateFile == "" {
			fmt.Println("Error: deployment file required")
			os.Exit(1)
		}
		
		deploymentID := updateCmd.Arg(0)
		deploymentData, err := ioutil.ReadFile(*updateFile)
		if err != nil {
			fmt.Printf("Error reading deployment file: %v\n", err)
			os.Exit(1)
		}
		
		var deployment Deployment
		if err := json.Unmarshal(deploymentData, &deployment); err != nil {
			fmt.Printf("Error parsing deployment JSON: %v\n", err)
			os.Exit(1)
		}
		
		updatedDeployment, err := client.UpdateDeployment(deploymentID, &deployment)
		if err != nil {
			fmt.Printf("Error updating deployment: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Deployment updated successfully:")
		deploymentJSON, _ := json.MarshalIndent(updatedDeployment, "", "  ")
		fmt.Println(string(deploymentJSON))
		
	case "scale":
		scaleCmd.Parse(os.Args[2:])
		if scaleCmd.NArg() < 1 {
			fmt.Println("Error: deployment ID required")
			os.Exit(1)
		}
		
		deploymentID := scaleCmd.Arg(0)
		scaledDeployment, err := client.ScaleDeployment(deploymentID, *scaleReplicas)
		if err != nil {
			fmt.Printf("Error scaling deployment: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Deployment scaled successfully:")
		deploymentJSON, _ := json.MarshalIndent(scaledDeployment, "", "  ")
		fmt.Println(string(deploymentJSON))
		
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if deleteCmd.NArg() < 1 {
			fmt.Println("Error: deployment ID required")
			os.Exit(1)
		}
		
		deploymentID := deleteCmd.Arg(0)
		if err := client.DeleteDeployment(deploymentID, *deleteForce); err != nil {
			fmt.Printf("Error deleting deployment: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Deployment deleted successfully")
		
	default:
		fmt.Println("Expected 'list', 'get', 'create', 'update', 'scale', or 'delete' subcommand")
		os.Exit(1)
	}
}
