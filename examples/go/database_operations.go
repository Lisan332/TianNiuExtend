package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

// MySQLConfig represents the MySQL configuration
type MySQLConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"metadata"`
	Environments []struct {
		Name              string `yaml:"name"`
		Host              string `yaml:"host"`
		Port              int    `yaml:"port"`
		Database          string `yaml:"database"`
		UsernameEnv       string `yaml:"username_env"`
		PasswordEnv       string `yaml:"password_env"`
		MaxConnections    int    `yaml:"max_connections"`
		ConnectionTimeout string `yaml:"connection_timeout"`
		ReadTimeout       string `yaml:"read_timeout"`
		WriteTimeout      string `yaml:"write_timeout"`
		MaxIdleConns      int    `yaml:"max_idle_connections"`
		MaxOpenConns      int    `yaml:"max_open_connections"`
		ConnMaxLifetime   string `yaml:"connection_max_lifetime"`
		SSLMode           string `yaml:"ssl_mode"`
		SSLCA             string `yaml:"ssl_ca,omitempty"`
		SSLCert           string `yaml:"ssl_cert,omitempty"`
		SSLKey            string `yaml:"ssl_key,omitempty"`
	} `yaml:"environments"`
	Defaults struct {
		Charset   string `yaml:"charset"`
		Collation string `yaml:"collation"`
		Timezone  string `yaml:"timezone"`
	} `yaml:"defaults"`
	Tables []struct {
		Name    string `yaml:"name"`
		Columns []struct {
			Name       string `yaml:"name"`
			Type       string `yaml:"type"`
			PrimaryKey bool   `yaml:"primary_key,omitempty"`
			Nullable   bool   `yaml:"nullable"`
			Default    string `yaml:"default,omitempty"`
			Index      bool   `yaml:"index,omitempty"`
			Unique     bool   `yaml:"unique,omitempty"`
			OnUpdate   string `yaml:"on_update,omitempty"`
			ForeignKey struct {
				Table    string `yaml:"table"`
				Column   string `yaml:"column"`
				OnDelete string `yaml:"on_delete"`
			} `yaml:"foreign_key,omitempty"`
		} `yaml:"columns"`
	} `yaml:"tables"`
}

// Container represents a container in the TianNiu platform
type Container struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Labels    string    `json:"labels"`
}

// Deployment represents a deployment in the TianNiu platform
type Deployment struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Environment string    `json:"environment"`
	CreatedAt   time.Time `json:"created_at"`
	Version     string    `json:"version"`
	Replicas    int       `json:"replicas"`
}

// DBClient represents a database client
type DBClient struct {
	DB     *sql.DB
	Config *MySQLConfig
	Env    string
}

// NewDBClient creates a new database client
func NewDBClient(configPath, env string) (*DBClient, error) {
	// Load configuration
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	// Find environment configuration
	var envConfig *struct {
		Name              string `yaml:"name"`
		Host              string `yaml:"host"`
		Port              int    `yaml:"port"`
		Database          string `yaml:"database"`
		UsernameEnv       string `yaml:"username_env"`
		PasswordEnv       string `yaml:"password_env"`
		MaxConnections    int    `yaml:"max_connections"`
		ConnectionTimeout string `yaml:"connection_timeout"`
		ReadTimeout       string `yaml:"read_timeout"`
		WriteTimeout      string `yaml:"write_timeout"`
		MaxIdleConns      int    `yaml:"max_idle_connections"`
		MaxOpenConns      int    `yaml:"max_open_connections"`
		ConnMaxLifetime   string `yaml:"connection_max_lifetime"`
		SSLMode           string `yaml:"ssl_mode"`
		SSLCA             string `yaml:"ssl_ca,omitempty"`
		SSLCert           string `yaml:"ssl_cert,omitempty"`
		SSLKey            string `yaml:"ssl_key,omitempty"`
	}

	for _, e := range config.Environments {
		if e.Name == env {
			envConfig = &e
			break
		}
	}

	if envConfig == nil {
		return nil, fmt.Errorf("environment %s not found in configuration", env)
	}

	// Get credentials from environment variables
	username := os.Getenv(envConfig.UsernameEnv)
	if username == "" {
		return nil, fmt.Errorf("environment variable %s not set", envConfig.UsernameEnv)
	}

	password := os.Getenv(envConfig.PasswordEnv)
	if password == "" {
		return nil, fmt.Errorf("environment variable %s not set", envConfig.PasswordEnv)
	}

	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&parseTime=true&loc=%s",
		username, password, envConfig.Host, envConfig.Port, envConfig.Database,
		config.Defaults.Charset, config.Defaults.Collation, config.Defaults.Timezone)

	// Add SSL options if required
	if envConfig.SSLMode == "require" || envConfig.SSLMode == "verify-ca" || envConfig.SSLMode == "verify-full" {
		dsn += "&tls=true"
		if envConfig.SSLCA != "" {
			dsn += fmt.Sprintf("&sslca=%s", envConfig.SSLCA)
		}
		if envConfig.SSLCert != "" {
			dsn += fmt.Sprintf("&sslcert=%s", envConfig.SSLCert)
		}
		if envConfig.SSLKey != "" {
			dsn += fmt.Sprintf("&sslkey=%s", envConfig.SSLKey)
		}
	}

	// Connect to database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Configure connection pool
	db.SetMaxIdleConns(envConfig.MaxIdleConns)
	db.SetMaxOpenConns(envConfig.MaxOpenConns)

	connMaxLifetime, err := time.ParseDuration(envConfig.ConnMaxLifetime)
	if err != nil {
		return nil, fmt.Errorf("invalid connection max lifetime: %v", err)
	}
	db.SetConnMaxLifetime(connMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &DBClient{
		DB:     db,
		Config: config,
		Env:    env,
	}, nil
}

// Close closes the database connection
func (c *DBClient) Close() error {
	return c.DB.Close()
}

// GetContainers gets all containers
func (c *DBClient) GetContainers(limit int) ([]Container, error) {
	query := "SELECT id, name, image, status, created_at, labels FROM containers ORDER BY created_at DESC LIMIT ?"
	rows, err := c.DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query containers: %v", err)
	}
	defer rows.Close()

	var containers []Container
	for rows.Next() {
		var container Container
		if err := rows.Scan(&container.ID, &container.Name, &container.Image, &container.Status, &container.CreatedAt, &container.Labels); err != nil {
			return nil, fmt.Errorf("failed to scan container row: %v", err)
		}
		containers = append(containers, container)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating container rows: %v", err)
	}

	return containers, nil
}

// GetContainerByID gets a container by ID
func (c *DBClient) GetContainerByID(id string) (*Container, error) {
	query := "SELECT id, name, image, status, created_at, labels FROM containers WHERE id = ?"
	row := c.DB.QueryRow(query, id)

	var container Container
	if err := row.Scan(&container.ID, &container.Name, &container.Image, &container.Status, &container.CreatedAt, &container.Labels); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("container with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to scan container row: %v", err)
	}

	return &container, nil
}

// CreateContainer creates a new container
func (c *DBClient) CreateContainer(container *Container) error {
	query := "INSERT INTO containers (id, name, image, status, created_at, labels) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := c.DB.Exec(query, container.ID, container.Name, container.Image, container.Status, container.CreatedAt, container.Labels)
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}
	return nil
}

// UpdateContainerStatus updates a container's status
func (c *DBClient) UpdateContainerStatus(id, status string) error {
	query := "UPDATE containers SET status = ? WHERE id = ?"
	result, err := c.DB.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update container status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("container with ID %s not found", id)
	}

	return nil
}

// DeleteContainer deletes a container
func (c *DBClient) DeleteContainer(id string) error {
	query := "DELETE FROM containers WHERE id = ?"
	result, err := c.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete container: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("container with ID %s not found", id)
	}

	return nil
}

// GetDeployments gets all deployments
func (c *DBClient) GetDeployments(limit int) ([]Deployment, error) {
	query := "SELECT id, name, description, status, environment, created_at, version, replicas FROM deployments ORDER BY created_at DESC LIMIT ?"
	rows, err := c.DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query deployments: %v", err)
	}
	defer rows.Close()

	var deployments []Deployment
	for rows.Next() {
		var deployment Deployment
		if err := rows.Scan(&deployment.ID, &deployment.Name, &deployment.Description, &deployment.Status, &deployment.Environment, &deployment.CreatedAt, &deployment.Version, &deployment.Replicas); err != nil {
			return nil, fmt.Errorf("failed to scan deployment row: %v", err)
		}
		deployments = append(deployments, deployment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating deployment rows: %v", err)
	}

	return deployments, nil
}

// GetDeploymentByID gets a deployment by ID
func (c *DBClient) GetDeploymentByID(id string) (*Deployment, error) {
	query := "SELECT id, name, description, status, environment, created_at, version, replicas FROM deployments WHERE id = ?"
	row := c.DB.QueryRow(query, id)

	var deployment Deployment
	if err := row.Scan(&deployment.ID, &deployment.Name, &deployment.Description, &deployment.Status, &deployment.Environment, &deployment.CreatedAt, &deployment.Version, &deployment.Replicas); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("deployment with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to scan deployment row: %v", err)
	}

	return &deployment, nil
}

// CreateDeployment creates a new deployment
func (c *DBClient) CreateDeployment(deployment *Deployment) error {
	query := "INSERT INTO deployments (id, name, description, status, environment, created_at, version, replicas) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := c.DB.Exec(query, deployment.ID, deployment.Name, deployment.Description, deployment.Status, deployment.Environment, deployment.CreatedAt, deployment.Version, deployment.Replicas)
	if err != nil {
		return fmt.Errorf("failed to create deployment: %v", err)
	}
	return nil
}

// UpdateDeploymentStatus updates a deployment's status
func (c *DBClient) UpdateDeploymentStatus(id, status string) error {
	query := "UPDATE deployments SET status = ? WHERE id = ?"
	result, err := c.DB.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update deployment status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("deployment with ID %s not found", id)
	}

	return nil
}

// ScaleDeployment scales a deployment
func (c *DBClient) ScaleDeployment(id string, replicas int) error {
	query := "UPDATE deployments SET replicas = ? WHERE id = ?"
	result, err := c.DB.Exec(query, replicas, id)
	if err != nil {
		return fmt.Errorf("failed to scale deployment: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("deployment with ID %s not found", id)
	}

	return nil
}

// DeleteDeployment deletes a deployment
func (c *DBClient) DeleteDeployment(id string) error {
	query := "DELETE FROM deployments WHERE id = ?"
	result, err := c.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete deployment: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("deployment with ID %s not found", id)
	}

	return nil
}

// Helper function to load MySQL configuration
func loadConfig(configPath string) (*MySQLConfig, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config MySQLConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	// Parse command line arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: database_operations <command> <args>")
		fmt.Println("Commands:")
		fmt.Println("  list-containers [limit]")
		fmt.Println("  get-container <container_id>")
		fmt.Println("  list-deployments [limit]")
		fmt.Println("  get-deployment <deployment_id>")
		os.Exit(1)
	}

	command := os.Args[1]
	
	// Set environment variables for testing
	os.Setenv("MYSQL_PROD_USERNAME", "tianniu_user")
	os.Setenv("MYSQL_PROD_PASSWORD", "tianniu_password")

	// Create database client
	client, err := NewDBClient("../../config/mysql-config.yaml", "production")
	if err != nil {
		log.Fatalf("Failed to create database client: %v", err)
	}
	defer client.Close()

	// Execute command
	switch command {
	case "list-containers":
		limit := 10
		if len(os.Args) > 2 {
			fmt.Sscanf(os.Args[2], "%d", &limit)
		}
		
		containers, err := client.GetContainers(limit)
		if err != nil {
			log.Fatalf("Failed to get containers: %v", err)
		}
		
		fmt.Printf("Found %d containers:\n", len(containers))
		for _, container := range containers {
			fmt.Printf("ID: %s, Name: %s, Image: %s, Status: %s, Created: %s\n",
				container.ID, container.Name, container.Image, container.Status, container.CreatedAt)
		}
		
	case "get-container":
		if len(os.Args) < 3 {
			log.Fatal("Container ID required")
		}
		
		containerID := os.Args[2]
		container, err := client.GetContainerByID(containerID)
		if err != nil {
			log.Fatalf("Failed to get container: %v", err)
		}
		
		fmt.Printf("Container details:\n")
		fmt.Printf("ID: %s\n", container.ID)
		fmt.Printf("Name: %s\n", container.Name)
		fmt.Printf("Image: %s\n", container.Image)
		fmt.Printf("Status: %s\n", container.Status)
		fmt.Printf("Created: %s\n", container.CreatedAt)
		fmt.Printf("Labels: %s\n", container.Labels)
		
	case "list-deployments":
		limit := 10
		if len(os.Args) > 2 {
			fmt.Sscanf(os.Args[2], "%d", &limit)
		}
		
		deployments, err := client.GetDeployments(limit)
		if err != nil {
			log.Fatalf("Failed to get deployments: %v", err)
		}
		
		fmt.Printf("Found %d deployments:\n", len(deployments))
		for _, deployment := range deployments {
			fmt.Printf("ID: %s, Name: %s, Environment: %s, Status: %s, Version: %s, Replicas: %d\n",
				deployment.ID, deployment.Name, deployment.Environment, deployment.Status, deployment.Version, deployment.Replicas)
		}
		
	case "get-deployment":
		if len(os.Args) < 3 {
			log.Fatal("Deployment ID required")
		}
		
		deploymentID := os.Args[2]
		deployment, err := client.GetDeploymentByID(deploymentID)
		if err != nil {
			log.Fatalf("Failed to get deployment: %v", err)
		}
		
		fmt.Printf("Deployment details:\n")
		fmt.Printf("ID: %s\n", deployment.ID)
		fmt.Printf("Name: %s\n", deployment.Name)
		fmt.Printf("Description: %s\n", deployment.Description)
		fmt.Printf("Status: %s\n", deployment.Status)
		fmt.Printf("Environment: %s\n", deployment.Environment)
		fmt.Printf("Created: %s\n", deployment.CreatedAt)
		fmt.Printf("Version: %s\n", deployment.Version)
		fmt.Printf("Replicas: %d\n", deployment.Replicas)
		
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
