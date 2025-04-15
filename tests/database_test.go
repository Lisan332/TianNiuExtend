package tests

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
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

// NewDBClient creates a new database client for testing
func NewTestDBClient() (*DBClient, error) {
	// Create a temporary database for testing
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	// Create test database
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS tianniu_test")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create test database: %v", err)
	}

	// Close connection and connect to test database
	db.Close()
	db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/tianniu_test")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %v", err)
	}

	// Create test schema
	schemaSQL, err := ioutil.ReadFile("../config/schema.sql")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to read schema file: %v", err)
	}

	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create test schema: %v", err)
	}

	// Create mock config
	config := &MySQLConfig{}
	config.Defaults.Charset = "utf8mb4"
	config.Defaults.Collation = "utf8mb4_unicode_ci"
	config.Defaults.Timezone = "Asia/Shanghai"

	return &DBClient{
		DB:     db,
		Config: config,
		Env:    "test",
	}, nil
}

// Close closes the database connection
func (c *DBClient) Close() error {
	if c.DB != nil {
		// Drop test database
		_, err := c.DB.Exec("DROP DATABASE IF EXISTS tianniu_test")
		if err != nil {
			return fmt.Errorf("failed to drop test database: %v", err)
		}
		return c.DB.Close()
	}
	return nil
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

// TestContainerOperations tests container database operations
func TestContainerOperations(t *testing.T) {
	// Skip if not running in CI environment
	if os.Getenv("CI") != "true" {
		t.Skip("Skipping test in non-CI environment")
	}

	// Create test client
	client, err := NewTestDBClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	defer client.Close()

	// Test creating a container
	container := &Container{
		ID:        "test-container-id",
		Name:      "test-container",
		Image:     "nginx:latest",
		Status:    "created",
		CreatedAt: time.Now(),
		Labels:    `{"app": "test", "environment": "testing"}`,
	}

	err = client.CreateContainer(container)
	assert.NoError(t, err, "Failed to create container")

	// Test getting a container by ID
	retrievedContainer, err := client.GetContainerByID(container.ID)
	assert.NoError(t, err, "Failed to get container by ID")
	assert.Equal(t, container.ID, retrievedContainer.ID, "Container ID mismatch")
	assert.Equal(t, container.Name, retrievedContainer.Name, "Container name mismatch")
	assert.Equal(t, container.Image, retrievedContainer.Image, "Container image mismatch")
	assert.Equal(t, container.Status, retrievedContainer.Status, "Container status mismatch")
	assert.Equal(t, container.Labels, retrievedContainer.Labels, "Container labels mismatch")

	// Test updating container status
	err = client.UpdateContainerStatus(container.ID, "running")
	assert.NoError(t, err, "Failed to update container status")

	// Verify status update
	retrievedContainer, err = client.GetContainerByID(container.ID)
	assert.NoError(t, err, "Failed to get container by ID after status update")
	assert.Equal(t, "running", retrievedContainer.Status, "Container status not updated")

	// Test getting all containers
	containers, err := client.GetContainers(10)
	assert.NoError(t, err, "Failed to get containers")
	assert.GreaterOrEqual(t, len(containers), 1, "Expected at least one container")
	
	// Find our test container in the list
	found := false
	for _, c := range containers {
		if c.ID == container.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Test container not found in container list")

	// Test deleting a container
	err = client.DeleteContainer(container.ID)
	assert.NoError(t, err, "Failed to delete container")

	// Verify container deletion
	_, err = client.GetContainerByID(container.ID)
	assert.Error(t, err, "Container should not exist after deletion")
}

// TestDeploymentOperations tests deployment database operations
func TestDeploymentOperations(t *testing.T) {
	// Skip if not running in CI environment
	if os.Getenv("CI") != "true" {
		t.Skip("Skipping test in non-CI environment")
	}

	// Create test client
	client, err := NewTestDBClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	defer client.Close()

	// Test creating a deployment
	deployment := &Deployment{
		ID:          "test-deployment-id",
		Name:        "test-deployment",
		Description: "Test deployment for unit tests",
		Status:      "created",
		Environment: "testing",
		CreatedAt:   time.Now(),
		Version:     "v1.0.0",
		Replicas:    2,
	}

	err = client.CreateDeployment(deployment)
	assert.NoError(t, err, "Failed to create deployment")

	// Test getting a deployment by ID
	retrievedDeployment, err := client.GetDeploymentByID(deployment.ID)
	assert.NoError(t, err, "Failed to get deployment by ID")
	assert.Equal(t, deployment.ID, retrievedDeployment.ID, "Deployment ID mismatch")
	assert.Equal(t, deployment.Name, retrievedDeployment.Name, "Deployment name mismatch")
	assert.Equal(t, deployment.Description, retrievedDeployment.Description, "Deployment description mismatch")
	assert.Equal(t, deployment.Status, retrievedDeployment.Status, "Deployment status mismatch")
	assert.Equal(t, deployment.Environment, retrievedDeployment.Environment, "Deployment environment mismatch")
	assert.Equal(t, deployment.Version, retrievedDeployment.Version, "Deployment version mismatch")
	assert.Equal(t, deployment.Replicas, retrievedDeployment.Replicas, "Deployment replicas mismatch")

	// Test updating deployment status
	err = client.UpdateDeploymentStatus(deployment.ID, "active")
	assert.NoError(t, err, "Failed to update deployment status")

	// Verify status update
	retrievedDeployment, err = client.GetDeploymentByID(deployment.ID)
	assert.NoError(t, err, "Failed to get deployment by ID after status update")
	assert.Equal(t, "active", retrievedDeployment.Status, "Deployment status not updated")

	// Test scaling deployment
	err = client.ScaleDeployment(deployment.ID, 5)
	assert.NoError(t, err, "Failed to scale deployment")

	// Verify scaling
	retrievedDeployment, err = client.GetDeploymentByID(deployment.ID)
	assert.NoError(t, err, "Failed to get deployment by ID after scaling")
	assert.Equal(t, 5, retrievedDeployment.Replicas, "Deployment replicas not updated")

	// Test getting all deployments
	deployments, err := client.GetDeployments(10)
	assert.NoError(t, err, "Failed to get deployments")
	assert.GreaterOrEqual(t, len(deployments), 1, "Expected at least one deployment")
	
	// Find our test deployment in the list
	found := false
	for _, d := range deployments {
		if d.ID == deployment.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Test deployment not found in deployment list")

	// Test deleting a deployment
	err = client.DeleteDeployment(deployment.ID)
	assert.NoError(t, err, "Failed to delete deployment")

	// Verify deployment deletion
	_, err = client.GetDeploymentByID(deployment.ID)
	assert.Error(t, err, "Deployment should not exist after deletion")
}
