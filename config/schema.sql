-- TianNiu Platform Database Schema

-- Set character set and collation
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
SET collation_connection = utf8mb4_unicode_ci;

-- Create containers table
CREATE TABLE IF NOT EXISTS containers (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    image VARCHAR(255) NOT NULL,
    status VARCHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP NULL,
    labels JSON NULL,
    ports JSON NULL,
    volumes JSON NULL,
    network JSON NULL,
    resource_limits JSON NULL,
    resource_usage JSON NULL,
    environment_variables JSON NULL,
    health_check JSON NULL,
    logs_url VARCHAR(255) NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_container_name (name),
    INDEX idx_container_status (status),
    INDEX idx_container_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create deployments table
CREATE TABLE IF NOT EXISTS deployments (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    status VARCHAR(32) NOT NULL,
    environment VARCHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    version VARCHAR(32) NOT NULL,
    replicas INT NOT NULL DEFAULT 1,
    strategy JSON NULL,
    containers JSON NOT NULL,
    services JSON NULL,
    config_maps JSON NULL,
    secrets JSON NULL,
    history JSON NULL,
    health_status JSON NULL,
    INDEX idx_deployment_name (name),
    INDEX idx_deployment_status (status),
    INDEX idx_deployment_environment (environment),
    INDEX idx_deployment_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create resources table
CREATE TABLE IF NOT EXISTS resources (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(32) NOT NULL,
    namespace VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    quota JSON NULL,
    usage JSON NULL,
    status VARCHAR(32) NOT NULL,
    details JSON NULL,
    INDEX idx_resource_name (name),
    INDEX idx_resource_type (type),
    INDEX idx_resource_namespace (namespace)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(64) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_login TIMESTAMP NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    roles JSON NOT NULL,
    UNIQUE KEY idx_user_username (username),
    UNIQUE KEY idx_user_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create API keys table
CREATE TABLE IF NOT EXISTS api_keys (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    key_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NULL,
    last_used_at TIMESTAMP NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    permissions JSON NOT NULL,
    INDEX idx_api_key_user_id (user_id),
    CONSTRAINT fk_api_key_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create audit logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NULL,
    api_key_id VARCHAR(64) NULL,
    action VARCHAR(255) NOT NULL,
    resource_type VARCHAR(64) NOT NULL,
    resource_id VARCHAR(64) NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45) NULL,
    user_agent VARCHAR(255) NULL,
    request_details JSON NULL,
    response_status INT NULL,
    response_details JSON NULL,
    INDEX idx_audit_log_user_id (user_id),
    INDEX idx_audit_log_api_key_id (api_key_id),
    INDEX idx_audit_log_action (action),
    INDEX idx_audit_log_resource_type (resource_type),
    INDEX idx_audit_log_resource_id (resource_id),
    INDEX idx_audit_log_timestamp (timestamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create sample data for testing

-- Insert sample users
INSERT INTO users (id, username, email, password_hash, full_name, status, roles)
VALUES
    ('usr_001', 'admin', 'admin_tianniu@baidu.com', '$2a$10$hACwQ5/HQI6FhbIISOUVeusy3sKyUDhSq36fF5d/54aULe9Vg5WNi', 'System Administrator', 'active', '["admin"]'),
    ('usr_002', 'developer', 'developer@baidu.com', '$2a$10$hACwQ5/HQI6FhbIISOUVeusy3sKyUDhSq36fF5d/54aULe9Vg5WNi', 'Developer User', 'active', '["developer"]'),
    ('usr_003', 'operator', 'operator@baidu.com', '$2a$10$hACwQ5/HQI6FhbIISOUVeusy3sKyUDhSq36fF5d/54aULe9Vg5WNi', 'Operations User', 'active', '["operator"]');

-- Insert sample API keys
INSERT INTO api_keys (id, user_id, name, key_hash, expires_at, status, permissions)
VALUES
    ('key_001', 'usr_001', 'Admin API Key', 'a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6', DATE_ADD(NOW(), INTERVAL 1 YEAR), 'active', '["*:*"]'),
    ('key_002', 'usr_002', 'Developer API Key', 'b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7', DATE_ADD(NOW(), INTERVAL 6 MONTH), 'active', '["container:read", "container:write", "deployment:read", "deployment:write"]'),
    ('key_003', 'usr_003', 'Operator API Key', 'c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8', DATE_ADD(NOW(), INTERVAL 3 MONTH), 'active', '["container:read", "deployment:read", "resource:read"]');

-- Insert sample containers
INSERT INTO containers (id, name, image, status, created_at, started_at, labels, ports, volumes, network, resource_limits, resource_usage)
VALUES
    ('c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2', 'web-server-1', 'nginx:latest', 'running', DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), '{"app": "web", "environment": "production"}', '[{"internal": 80, "external": 8080, "protocol": "tcp"}]', '[{"host_path": "/data/nginx/conf", "container_path": "/etc/nginx/conf.d", "mode": "ro"}, {"host_path": "/data/nginx/html", "container_path": "/usr/share/nginx/html", "mode": "rw"}]', '{"name": "frontend-network", "ip_address": "172.18.0.2"}', '{"cpu": "1.0", "memory": "512MB"}', '{"cpu": "0.05", "memory": "128MB", "network_rx": "1.2MB/s", "network_tx": "0.8MB/s"}'),
    ('a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6', 'api-service', 'registry.baidu.com/myteam/api-service:v1.2.3', 'running', DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), '{"app": "api", "environment": "production", "team": "backend"}', '[{"internal": 8080, "external": 9000, "protocol": "tcp"}]', '[{"host_path": "/data/api/config", "container_path": "/app/config", "mode": "ro"}, {"host_path": "/data/api/logs", "container_path": "/app/logs", "mode": "rw"}]', '{"name": "backend-network", "ip_address": "172.18.1.3"}', '{"cpu": "2.0", "memory": "1GB"}', '{"cpu": "0.75", "memory": "512MB", "network_rx": "0.8MB/s", "network_tx": "1.5MB/s"}'),
    ('b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7', 'db-server', 'mysql:8.0', 'running', DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_SUB(NOW(), INTERVAL 7 DAY), '{"app": "db", "environment": "production", "team": "data"}', '[{"internal": 3306, "external": 3306, "protocol": "tcp"}]', '[{"host_path": "/data/mysql", "container_path": "/var/lib/mysql", "mode": "rw"}]', '{"name": "backend-network", "ip_address": "172.18.1.4"}', '{"cpu": "4.0", "memory": "8GB"}', '{"cpu": "1.5", "memory": "4GB", "network_rx": "0.5MB/s", "network_tx": "0.3MB/s"}');

-- Insert sample deployments
INSERT INTO deployments (id, name, description, status, environment, created_at, version, replicas, strategy, containers, services)
VALUES
    ('d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6', 'web-frontend', 'Web前端应用', 'active', 'production', DATE_SUB(NOW(), INTERVAL 10 DAY), 'v2.3.1', 3, '{"type": "rolling-update", "max_surge": 1, "max_unavailable": 0}', '[{"name": "web-frontend", "image": "registry.baidu.com/frontend/web-app:v2.3.1", "ports": [{"name": "http", "container_port": 80, "service_port": 8080}], "resources": {"limits": {"cpu": "1.0", "memory": "1Gi"}, "requests": {"cpu": "0.5", "memory": "512Mi"}}, "environment_variables": [{"name": "API_ENDPOINT", "value": "https://api.palo.prod.baidu.com/v1"}, {"name": "LOG_LEVEL", "value": "info"}], "health_check": {"http_path": "/health", "port": 80, "initial_delay_seconds": 10, "period_seconds": 30, "timeout_seconds": 5, "success_threshold": 1, "failure_threshold": 3}}]', '[{"name": "web-frontend-svc", "type": "LoadBalancer", "ports": [{"name": "http", "port": 80, "target_port": 8080}], "external_endpoints": ["web-frontend.palo.prod.baidu.com"]}]'),
    ('f1e2d3c4b5a6f7e8d9c0b1a2d3e4f5c6', 'api-backend', 'API后端服务', 'active', 'production', DATE_SUB(NOW(), INTERVAL 15 DAY), 'v1.5.1', 2, '{"type": "rolling-update", "max_surge": 1, "max_unavailable": 0}', '[{"name": "api-backend", "image": "registry.baidu.com/backend/api-service:v1.5.1", "ports": [{"name": "http", "container_port": 8080, "service_port": 80}], "resources": {"limits": {"cpu": "2.0", "memory": "2Gi"}, "requests": {"cpu": "1.0", "memory": "1Gi"}}, "environment_variables": [{"name": "DB_HOST", "value": "db.palo.prod.baidu.com"}, {"name": "DB_PORT", "value": "5432"}, {"name": "LOG_LEVEL", "value": "debug"}], "health_check": {"http_path": "/health", "port": 8080, "initial_delay_seconds": 15, "period_seconds": 30, "timeout_seconds": 5, "success_threshold": 1, "failure_threshold": 3}}]', '[{"name": "api-backend-svc", "type": "ClusterIP", "ports": [{"name": "http", "port": 80, "target_port": 8080}]}]'),
    ('g1h2i3j4k5l6m7n8o9p0q1r2s3t4u5v6', 'database', '数据库服务', 'active', 'production', DATE_SUB(NOW(), INTERVAL 20 DAY), 'v1.0.0', 1, '{"type": "recreate"}', '[{"name": "database", "image": "registry.baidu.com/database/postgres:13", "ports": [{"name": "postgres", "container_port": 5432, "service_port": 5432}], "resources": {"limits": {"cpu": "4.0", "memory": "8Gi"}, "requests": {"cpu": "2.0", "memory": "4Gi"}}, "environment_variables": [{"name": "POSTGRES_DB", "value": "appdb"}, {"name": "POSTGRES_USER", "value": "appuser"}], "health_check": {"tcp_port": 5432, "initial_delay_seconds": 30, "period_seconds": 60, "timeout_seconds": 10, "success_threshold": 1, "failure_threshold": 3}}]', '[{"name": "database-svc", "type": "ClusterIP", "ports": [{"name": "postgres", "port": 5432, "target_port": 5432}]}]');

-- Insert sample resources
INSERT INTO resources (id, name, type, namespace, created_at, quota, usage, status, details)
VALUES
    ('res_001', 'production-quota', 'quota', 'production', DATE_SUB(NOW(), INTERVAL 30 DAY), '{"cpu": 100, "memory": 256, "storage": 1000, "network": 100}', '{"cpu": 45, "memory": 128, "storage": 350, "network": 25}', 'active', '{"description": "Production environment resource quota"}'),
    ('res_002', 'worker-01', 'node', 'system', DATE_SUB(NOW(), INTERVAL 60 DAY), '{"cpu": 32, "memory": 128, "pods": 110}', '{"cpu": 25, "memory": 100, "pods": 85}', 'ready', '{"region": "beijing", "zone": "zone-a", "instance-type": "high-memory"}'),
    ('res_003', 'frontend-network', 'network', 'production', DATE_SUB(NOW(), INTERVAL 25 DAY), '{"bandwidth": 1000}', '{"bandwidth": 350}', 'active', '{"subnet": "172.18.0.0/24", "gateway": "172.18.0.1"}');

-- Insert sample audit logs
INSERT INTO audit_logs (id, user_id, api_key_id, action, resource_type, resource_id, timestamp, ip_address, user_agent, request_details, response_status)
VALUES
    ('log_001', 'usr_001', NULL, 'create', 'container', 'c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2', DATE_SUB(NOW(), INTERVAL 3 DAY), '192.168.1.100', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36', '{"method": "POST", "path": "/api/v1/containers", "body": {"name": "web-server-1", "image": "nginx:latest"}}', 200),
    ('log_002', NULL, 'key_002', 'create', 'deployment', 'd1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6', DATE_SUB(NOW(), INTERVAL 10 DAY), '192.168.1.101', 'tianniu-client/1.0.0', '{"method": "POST", "path": "/api/v1/deployments", "body": {"name": "web-frontend", "environment": "production"}}', 200),
    ('log_003', 'usr_003', NULL, 'get', 'resource', 'res_001', DATE_SUB(NOW(), INTERVAL 1 DAY), '192.168.1.102', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)', '{"method": "GET", "path": "/api/v1/resources/quotas?namespace=production"}', 200);
