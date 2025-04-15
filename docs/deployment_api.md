# 部署管理 API

## 概述

部署管理API允许您创建、管理和监控应用部署。这些API支持多环境部署、蓝绿部署、金丝雀发布等高级部署策略。

基础URL: `https://tianniuprod.baidu.com/api/v1/deployments`

## 部署列表

### 获取所有部署

```
GET /api/v1/deployments
```

查询参数:
- `status` (可选): 按状态筛选 (active, failed, pending)
- `environment` (可选): 按环境筛选 (production, staging, development)
- `limit` (可选): 返回结果数量限制 (默认: 20, 最大: 100)
- `offset` (可选): 分页偏移量 (默认: 0)

**请求示例:**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/deployments?environment=production&limit=10" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "total": 25,
  "limit": 10,
  "offset": 0,
  "deployments": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "name": "web-frontend",
      "status": "active",
      "environment": "production",
      "created_at": "2023-05-10T14:30:25Z",
      "updated_at": "2023-05-10T14:35:12Z",
      "version": "v2.3.1",
      "replicas": 3,
      "health_status": "healthy"
    },
    // 更多部署...
  ]
}
```

## 部署详情

### 获取单个部署详情

```
GET /api/v1/deployments/{deployment_id}
```

**请求示例:**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/deployments/d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
  "name": "web-frontend",
  "description": "Web前端应用",
  "status": "active",
  "environment": "production",
  "created_at": "2023-05-10T14:30:25Z",
  "updated_at": "2023-05-10T14:35:12Z",
  "version": "v2.3.1",
  "replicas": 3,
  "available_replicas": 3,
  "strategy": {
    "type": "rolling-update",
    "max_surge": 1,
    "max_unavailable": 0
  },
  "containers": [
    {
      "name": "web-frontend",
      "image": "registry.baidu.com/frontend/web-app:v2.3.1",
      "ports": [
        {
          "name": "http",
          "container_port": 80,
          "service_port": 8080
        }
      ],
      "resources": {
        "limits": {
          "cpu": "1.0",
          "memory": "1Gi"
        },
        "requests": {
          "cpu": "0.5",
          "memory": "512Mi"
        }
      },
      "environment_variables": [
        {
          "name": "API_ENDPOINT",
          "value": "https://api.palo.prod.baidu.com/v1"
        },
        {
          "name": "LOG_LEVEL",
          "value": "info"
        }
      ],
      "health_check": {
        "http_path": "/health",
        "port": 80,
        "initial_delay_seconds": 10,
        "period_seconds": 30,
        "timeout_seconds": 5,
        "success_threshold": 1,
        "failure_threshold": 3
      }
    }
  ],
  "services": [
    {
      "name": "web-frontend-svc",
      "type": "LoadBalancer",
      "ports": [
        {
          "name": "http",
          "port": 80,
          "target_port": 8080
        }
      ],
      "external_endpoints": [
        "web-frontend.palo.prod.baidu.com"
      ]
    }
  ],
  "config_maps": [
    {
      "name": "web-frontend-config",
      "mounted_path": "/app/config",
      "data": {
        "app.conf": "server {\n  listen 80;\n  server_name web-frontend.palo.prod.baidu.com;\n  location / {\n    root /usr/share/nginx/html;\n    index index.html;\n  }\n}"
      }
    }
  ],
  "secrets": [
    {
      "name": "web-frontend-secrets",
      "mounted_path": "/app/secrets"
    }
  ],
  "history": [
    {
      "version": "v2.3.1",
      "deployed_at": "2023-05-10T14:30:25Z",
      "status": "active",
      "deployed_by": "user@baidu.com"
    },
    {
      "version": "v2.3.0",
      "deployed_at": "2023-05-01T10:15:30Z",
      "status": "superseded",
      "deployed_by": "user@baidu.com"
    }
  ],
  "health_status": {
    "status": "healthy",
    "last_checked": "2023-05-10T15:30:25Z",
    "details": {
      "readiness": "3/3",
      "liveness": "3/3"
    }
  }
}
```

## 创建部署

### 创建新部署

```
POST /api/v1/deployments
```

**请求体:**

```json
{
  "name": "api-backend",
  "description": "API后端服务",
  "environment": "production",
  "version": "v1.5.0",
  "replicas": 2,
  "strategy": {
    "type": "rolling-update",
    "max_surge": 1,
    "max_unavailable": 0
  },
  "containers": [
    {
      "name": "api-backend",
      "image": "registry.baidu.com/backend/api-service:v1.5.0",
      "ports": [
        {
          "name": "http",
          "container_port": 8080,
          "service_port": 80
        }
      ],
      "resources": {
        "limits": {
          "cpu": "2.0",
          "memory": "2Gi"
        },
        "requests": {
          "cpu": "1.0",
          "memory": "1Gi"
        }
      },
      "environment_variables": [
        {
          "name": "DB_HOST",
          "value": "db.palo.prod.baidu.com"
        },
        {
          "name": "DB_PORT",
          "value": "5432"
        },
        {
          "name": "LOG_LEVEL",
          "value": "info"
        }
      ],
      "health_check": {
        "http_path": "/health",
        "port": 8080,
        "initial_delay_seconds": 15,
        "period_seconds": 30,
        "timeout_seconds": 5,
        "success_threshold": 1,
        "failure_threshold": 3
      }
    }
  ],
  "services": [
    {
      "name": "api-backend-svc",
      "type": "ClusterIP",
      "ports": [
        {
          "name": "http",
          "port": 80,
          "target_port": 8080
        }
      ]
    }
  ],
  "config_maps": [
    {
      "name": "api-backend-config",
      "mounted_path": "/app/config",
      "data": {
        "app.yaml": "server:\n  port: 8080\n  timeout: 30s\ndatabase:\n  host: db.palo.prod.baidu.com\n  port: 5432\n  name: api_db\n  pool_size: 10"
      }
    }
  ],
  "secrets": [
    {
      "name": "api-backend-secrets",
      "mounted_path": "/app/secrets"
    }
  ]
}
```

**响应示例:**

```json
{
  "id": "f1e2d3c4b5a6f7e8d9c0b1a2d3e4f5c6",
  "name": "api-backend",
  "status": "pending",
  "environment": "production",
  "created_at": "2023-05-15T09:45:30Z",
  "version": "v1.5.0",
  "message": "Deployment is being created"
}
```

## 部署操作

### 更新部署

```
PUT /api/v1/deployments/{deployment_id}
```

**请求体:**

```json
{
  "version": "v1.5.1",
  "replicas": 3,
  "containers": [
    {
      "name": "api-backend",
      "image": "registry.baidu.com/backend/api-service:v1.5.1",
      "environment_variables": [
        {
          "name": "LOG_LEVEL",
          "value": "debug"
        }
      ]
    }
  ]
}
```

**响应示例:**

```json
{
  "id": "f1e2d3c4b5a6f7e8d9c0b1a2d3e4f5c6",
  "name": "api-backend",
  "status": "updating",
  "environment": "production",
  "updated_at": "2023-05-16T11:20:15Z",
  "version": "v1.5.1",
  "message": "Deployment is being updated"
}
```

### 扩缩容部署

```
POST /api/v1/deployments/{deployment_id}/scale
```

**请求体:**

```json
{
  "replicas": 5
}
```

**响应示例:**

```json
{
  "id": "f1e2d3c4b5a6f7e8d9c0b1a2d3e4f5c6",
  "name": "api-backend",
  "status": "scaling",
  "environment": "production",
  "updated_at": "2023-05-16T14:10:25Z",
  "replicas": 5,
  "message": "Deployment is scaling from 3 to 5 replicas"
}
```

### 回滚部署

```
POST /api/v1/deployments/{deployment_id}/rollback
```

**请求体:**

```json
{
  "version": "v1.5.0"
}
```

**响应示例:**

```json
{
  "id": "f1e2d3c4b5a6f7e8d9c0b1a2d3e4f5c6",
  "name": "api-backend",
  "status": "rolling-back",
  "environment": "production",
  "updated_at": "2023-05-16T16:45:10Z",
  "version": "v1.5.0",
  "message": "Deployment is rolling back to version v1.5.0"
}
```

### 暂停部署

```
POST /api/v1/deployments/{deployment_id}/pause
```

**响应示例:**

```json
{
  "id": "f1e2d3c4b5a6f7e8d9c0b1a2d3e4f5c6",
  "name": "api-backend",
  "status": "paused",
  "environment": "production",
  "updated_at": "2023-05-16T17:30:45Z",
  "message": "Deployment has been paused"
}
```

### 恢复部署

```
POST /api/v1/deployments/{deployment_id}/resume
```

**响应示例:**

```json
{
  "id": "f1e2d3c4b5a6f7e8d9c0b1a2d3e4f5c6",
  "name": "api-backend",
  "status": "active",
  "environment": "production",
  "updated_at": "2023-05-16T18:15:20Z",
  "message": "Deployment has been resumed"
}
```

## 部署历史

### 获取部署历史

```
GET /api/v1/deployments/{deployment_id}/history
```

**请求示例:**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/deployments/f1e2d3c4b5a6f7e8d9c0b1a2d3e4f5c6/history" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "deployment_id": "f1e2d3c4b5a6f7e8d9c0b1a2d3e4f5c6",
  "name": "api-backend",
  "history": [
    {
      "version": "v1.5.1",
      "deployed_at": "2023-05-16T11:20:15Z",
      "status": "active",
      "deployed_by": "user@baidu.com",
      "changes": [
        "Updated image to v1.5.1",
        "Changed LOG_LEVEL to debug",
        "Scaled from 2 to 3 replicas"
      ]
    },
    {
      "version": "v1.5.0",
      "deployed_at": "2023-05-15T09:45:30Z",
      "status": "superseded",
      "deployed_by": "user@baidu.com",
      "changes": [
        "Initial deployment"
      ]
    }
  ]
}
```

## 删除部署

### 删除部署

```
DELETE /api/v1/deployments/{deployment_id}
```

查询参数:
- `force` (可选): 是否强制删除 (默认: false)

**请求示例:**

```bash
curl -X DELETE "https://tianniuprod.baidu.com/api/v1/deployments/f1e2d3c4b5a6f7e8d9c0b1a2d3e4f5c6" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "id": "f1e2d3c4b5a6f7e8d9c0b1a2d3e4f5c6",
  "status": "deleting",
  "message": "Deployment is being deleted"
}
```

## 高级部署策略

### 创建蓝绿部署

```
POST /api/v1/deployments/blue-green
```

**请求体:**

```json
{
  "name": "payment-service",
  "environment": "production",
  "blue_version": {
    "version": "v2.0.0",
    "image": "registry.baidu.com/payment/service:v2.0.0",
    "replicas": 3
  },
  "green_version": {
    "version": "v2.1.0",
    "image": "registry.baidu.com/payment/service:v2.1.0",
    "replicas": 3
  },
  "test_traffic_percentage": 20,
  "auto_promote_after": "2h",
  "rollback_threshold": {
    "error_rate": 1.0,
    "latency_p99_ms": 500
  }
}
```

**响应示例:**

```json
{
  "id": "b1g2d3e4p5l6o7y8m9e0n1t2",
  "name": "payment-service",
  "status": "creating",
  "environment": "production",
  "created_at": "2023-05-20T10:30:15Z",
  "strategy": "blue-green",
  "message": "Blue-green deployment is being created"
}
```

### 创建金丝雀发布

```
POST /api/v1/deployments/canary
```

**请求体:**

```json
{
  "name": "search-service",
  "environment": "production",
  "base_version": {
    "version": "v3.2.0",
    "image": "registry.baidu.com/search/service:v3.2.0",
    "replicas": 10
  },
  "canary_version": {
    "version": "v3.3.0",
    "image": "registry.baidu.com/search/service:v3.3.0"
  },
  "stages": [
    {
      "traffic_percentage": 5,
      "replicas": 1,
      "duration": "30m",
      "analysis": {
        "metrics": ["error_rate", "latency_p95"],
        "max_error_rate": 1.0,
        "max_latency_p95_ms": 300
      }
    },
    {
      "traffic_percentage": 20,
      "replicas": 3,
      "duration": "1h",
      "analysis": {
        "metrics": ["error_rate", "latency_p95"],
        "max_error_rate": 0.5,
        "max_latency_p95_ms": 250
      }
    },
    {
      "traffic_percentage": 50,
      "replicas": 5,
      "duration": "2h",
      "analysis": {
        "metrics": ["error_rate", "latency_p95"],
        "max_error_rate": 0.2,
        "max_latency_p95_ms": 200
      }
    }
  ],
  "auto_promote": true
}
```

**响应示例:**

```json
{
  "id": "c1a2n3a4r5y6d7e8p9l0o1y2",
  "name": "search-service",
  "status": "creating",
  "environment": "production",
  "created_at": "2023-05-22T09:15:30Z",
  "strategy": "canary",
  "message": "Canary deployment is being created"
}
```

## 错误响应

所有API错误响应都遵循以下格式:

```json
{
  "error": {
    "code": "DEPLOYMENT_NOT_FOUND",
    "message": "Deployment with ID 'f1e2d3c4b5a6f7e8d9c0b1a2d3e4f5c6' not found",
    "details": {
      "deployment_id": "f1e2d3c4b5a6f7e8d9c0b1a2d3e4f5c6"
    }
  }
}
```

常见错误代码:

| 错误代码 | HTTP状态码 | 描述 |
|----------|------------|------|
| DEPLOYMENT_NOT_FOUND | 404 | 指定的部署不存在 |
| DEPLOYMENT_ALREADY_EXISTS | 409 | 部署名称已存在 |
| INVALID_DEPLOYMENT_STATE | 400 | 部署状态不允许执行请求的操作 |
| IMAGE_NOT_FOUND | 404 | 指定的镜像不存在 |
| RESOURCE_LIMIT_EXCEEDED | 403 | 超出资源限制 |
| PERMISSION_DENIED | 403 | 没有执行操作的权限 |
| INTERNAL_ERROR | 500 | 内部服务器错误 |
