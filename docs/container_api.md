# 容器管理 API

## 概述

容器管理API允许您创建、管理、监控和删除Docker容器。这些API支持单个容器操作以及批量操作，并提供了丰富的配置选项。

基础URL: `https://tianniuprod.baidu.com/api/v1/containers`

## 容器列表

### 获取所有容器

```
GET /api/v1/containers
```

查询参数:
- `status` (可选): 按状态筛选 (running, stopped, paused)
- `label` (可选): 按标签筛选 (格式: key=value)
- `limit` (可选): 返回结果数量限制 (默认: 20, 最大: 100)
- `offset` (可选): 分页偏移量 (默认: 0)

**请求示例:**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/containers?status=running&limit=10" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "total": 42,
  "limit": 10,
  "offset": 0,
  "containers": [
    {
      "id": "c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2",
      "name": "web-server-1",
      "image": "nginx:latest",
      "status": "running",
      "created_at": "2023-06-15T08:30:45Z",
      "labels": {
        "app": "web",
        "environment": "production"
      },
      "ports": [
        {
          "internal": 80,
          "external": 8080,
          "protocol": "tcp"
        }
      ],
      "resource_usage": {
        "cpu": "0.05",
        "memory": "128MB",
        "network_rx": "1.2MB/s",
        "network_tx": "0.8MB/s"
      }
    },
    // 更多容器...
  ]
}
```

## 容器详情

### 获取单个容器详情

```
GET /api/v1/containers/{container_id}
```

**请求示例:**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/containers/c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "id": "c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2",
  "name": "web-server-1",
  "image": "nginx:latest",
  "status": "running",
  "created_at": "2023-06-15T08:30:45Z",
  "started_at": "2023-06-15T08:30:47Z",
  "labels": {
    "app": "web",
    "environment": "production"
  },
  "ports": [
    {
      "internal": 80,
      "external": 8080,
      "protocol": "tcp"
    }
  ],
  "volumes": [
    {
      "host_path": "/data/nginx/conf",
      "container_path": "/etc/nginx/conf.d",
      "mode": "ro"
    },
    {
      "host_path": "/data/nginx/html",
      "container_path": "/usr/share/nginx/html",
      "mode": "rw"
    }
  ],
  "network": {
    "name": "frontend-network",
    "ip_address": "172.18.0.2"
  },
  "resource_limits": {
    "cpu": "1.0",
    "memory": "512MB"
  },
  "resource_usage": {
    "cpu": "0.05",
    "memory": "128MB",
    "network_rx": "1.2MB/s",
    "network_tx": "0.8MB/s"
  },
  "environment_variables": [
    {
      "name": "NGINX_HOST",
      "value": "example.com"
    },
    {
      "name": "NGINX_PORT",
      "value": "80"
    }
  ],
  "health_check": {
    "status": "healthy",
    "last_checked": "2023-06-15T10:45:12Z",
    "endpoint": "http://localhost:80/health",
    "interval": "30s",
    "timeout": "5s",
    "retries": 3
  },
  "logs_url": "https://tianniuprod.baidu.com/api/v1/containers/c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2/logs"
}
```

## 创建容器

### 创建新容器

```
POST /api/v1/containers
```

**请求体:**

```json
{
  "name": "api-service",
  "image": "registry.baidu.com/myteam/api-service:v1.2.3",
  "labels": {
    "app": "api",
    "environment": "production",
    "team": "backend"
  },
  "ports": [
    {
      "internal": 8080,
      "external": 9000,
      "protocol": "tcp"
    }
  ],
  "volumes": [
    {
      "host_path": "/data/api/config",
      "container_path": "/app/config",
      "mode": "ro"
    },
    {
      "host_path": "/data/api/logs",
      "container_path": "/app/logs",
      "mode": "rw"
    }
  ],
  "network": "backend-network",
  "resource_limits": {
    "cpu": "2.0",
    "memory": "1GB"
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
      "name": "API_KEY",
      "value_from": {
        "secret_name": "api-service-secrets",
        "key": "api-key"
      }
    }
  ],
  "health_check": {
    "endpoint": "http://localhost:8080/health",
    "interval": "15s",
    "timeout": "3s",
    "retries": 5
  },
  "restart_policy": {
    "type": "on-failure",
    "max_retries": 3
  }
}
```

**响应示例:**

```json
{
  "id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "name": "api-service",
  "status": "creating",
  "created_at": "2023-06-16T14:22:36Z",
  "image": "registry.baidu.com/myteam/api-service:v1.2.3",
  "message": "Container is being created and will start shortly"
}
```

## 容器操作

### 启动容器

```
POST /api/v1/containers/{container_id}/start
```

**请求示例:**

```bash
curl -X POST "https://tianniuprod.baidu.com/api/v1/containers/a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6/start" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "status": "starting",
  "message": "Container is starting"
}
```

### 停止容器

```
POST /api/v1/containers/{container_id}/stop
```

查询参数:
- `timeout` (可选): 等待容器优雅停止的时间(秒) (默认: 10)

**请求示例:**

```bash
curl -X POST "https://tianniuprod.baidu.com/api/v1/containers/a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6/stop?timeout=30" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "status": "stopping",
  "message": "Container is stopping"
}
```

### 重启容器

```
POST /api/v1/containers/{container_id}/restart
```

查询参数:
- `timeout` (可选): 等待容器优雅停止的时间(秒) (默认: 10)

**请求示例:**

```bash
curl -X POST "https://tianniuprod.baidu.com/api/v1/containers/a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6/restart" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "status": "restarting",
  "message": "Container is restarting"
}
```

### 暂停容器

```
POST /api/v1/containers/{container_id}/pause
```

**请求示例:**

```bash
curl -X POST "https://tianniuprod.baidu.com/api/v1/containers/a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6/pause" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "status": "pausing",
  "message": "Container is being paused"
}
```

### 恢复容器

```
POST /api/v1/containers/{container_id}/unpause
```

**请求示例:**

```bash
curl -X POST "https://tianniuprod.baidu.com/api/v1/containers/a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6/unpause" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "status": "unpausing",
  "message": "Container is being unpaused"
}
```

## 容器日志

### 获取容器日志

```
GET /api/v1/containers/{container_id}/logs
```

查询参数:
- `tail` (可选): 返回最后N行日志 (默认: 100)
- `since` (可选): 返回指定时间戳之后的日志 (格式: ISO 8601)
- `until` (可选): 返回指定时间戳之前的日志 (格式: ISO 8601)
- `follow` (可选): 是否持续获取新日志 (默认: false)

**请求示例:**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/containers/a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6/logs?tail=50" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "container_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "logs": [
    {
      "timestamp": "2023-06-16T14:30:12Z",
      "stream": "stdout",
      "message": "Server started on port 8080"
    },
    {
      "timestamp": "2023-06-16T14:30:15Z",
      "stream": "stdout",
      "message": "Connected to database at db.palo.prod.baidu.com:5432"
    },
    {
      "timestamp": "2023-06-16T14:30:18Z",
      "stream": "stderr",
      "message": "Warning: High memory usage detected (75%)"
    }
    // 更多日志...
  ]
}
```

## 删除容器

### 删除容器

```
DELETE /api/v1/containers/{container_id}
```

查询参数:
- `force` (可选): 是否强制删除运行中的容器 (默认: false)
- `remove_volumes` (可选): 是否同时删除关联的匿名卷 (默认: false)

**请求示例:**

```bash
curl -X DELETE "https://tianniuprod.baidu.com/api/v1/containers/a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6?force=true" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "status": "deleted",
  "message": "Container has been deleted"
}
```

## 批量操作

### 批量启动容器

```
POST /api/v1/containers/batch/start
```

**请求体:**

```json
{
  "container_ids": [
    "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
    "b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7"
  ]
}
```

**响应示例:**

```json
{
  "results": [
    {
      "id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
      "status": "starting",
      "message": "Container is starting"
    },
    {
      "id": "b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7",
      "status": "starting",
      "message": "Container is starting"
    }
  ]
}
```

### 批量停止容器

```
POST /api/v1/containers/batch/stop
```

**请求体:**

```json
{
  "container_ids": [
    "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
    "b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7"
  ],
  "timeout": 30
}
```

**响应示例:**

```json
{
  "results": [
    {
      "id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
      "status": "stopping",
      "message": "Container is stopping"
    },
    {
      "id": "b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7",
      "status": "stopping",
      "message": "Container is stopping"
    }
  ]
}
```

## 容器统计信息

### 获取容器资源使用统计

```
GET /api/v1/containers/{container_id}/stats
```

查询参数:
- `interval` (可选): 统计数据的时间间隔 (默认: 1s)

**请求示例:**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/containers/a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6/stats" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "container_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2023-06-16T15:30:45Z",
  "cpu": {
    "usage_percent": 12.5,
    "system_cpu_usage": 45.2,
    "online_cpus": 4
  },
  "memory": {
    "usage": 256000000,
    "limit": 1073741824,
    "usage_percent": 23.8
  },
  "network": {
    "rx_bytes": 1024000,
    "rx_packets": 1500,
    "rx_errors": 0,
    "rx_dropped": 0,
    "tx_bytes": 512000,
    "tx_packets": 1200,
    "tx_errors": 0,
    "tx_dropped": 0
  },
  "io": {
    "read_ops": 120,
    "write_ops": 80,
    "read_bytes": 2048000,
    "write_bytes": 1024000
  }
}
```

## 容器执行命令

### 在容器中执行命令

```
POST /api/v1/containers/{container_id}/exec
```

**请求体:**

```json
{
  "command": ["ls", "-la", "/app"],
  "env": [
    {
      "name": "DEBUG",
      "value": "true"
    }
  ],
  "working_dir": "/app",
  "user": "app-user",
  "tty": false,
  "attach_stdin": false,
  "attach_stdout": true,
  "attach_stderr": true
}
```

**响应示例:**

```json
{
  "exit_code": 0,
  "stdout": "total 24\ndrwxr-xr-x 5 app-user app-user 4096 Jun 16 15:45 .\ndrwxr-xr-x 1 root root 4096 Jun 16 14:22 ..\n-rw-r--r-- 1 app-user app-user 2048 Jun 16 15:40 app.js\ndrwxr-xr-x 2 app-user app-user 4096 Jun 16 15:30 config\ndrwxr-xr-x 2 app-user app-user 4096 Jun 16 15:35 logs\ndrwxr-xr-x 3 app-user app-user 4096 Jun 16 15:25 node_modules\n",
  "stderr": ""
}
```

## 错误响应

所有API错误响应都遵循以下格式:

```json
{
  "error": {
    "code": "CONTAINER_NOT_FOUND",
    "message": "Container with ID 'a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6' not found",
    "details": {
      "container_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
    }
  }
}
```

常见错误代码:

| 错误代码 | HTTP状态码 | 描述 |
|----------|------------|------|
| CONTAINER_NOT_FOUND | 404 | 指定的容器不存在 |
| CONTAINER_ALREADY_EXISTS | 409 | 容器名称已存在 |
| INVALID_CONTAINER_STATE | 400 | 容器状态不允许执行请求的操作 |
| IMAGE_NOT_FOUND | 404 | 指定的镜像不存在 |
| NETWORK_NOT_FOUND | 404 | 指定的网络不存在 |
| VOLUME_NOT_FOUND | 404 | 指定的卷不存在 |
| RESOURCE_LIMIT_EXCEEDED | 403 | 超出资源限制 |
| PERMISSION_DENIED | 403 | 没有执行操作的权限 |
| INTERNAL_ERROR | 500 | 内部服务器错误 |
