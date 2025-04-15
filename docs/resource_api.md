# 资源管理 API

## 概述

资源管理API允许您管理和监控天牛平台上的计算、网络和存储资源。通过这些API，您可以分配资源、设置配额、监控使用情况并优化资源利用率。

基础URL: `https://tianniuprod.baidu.com/api/v1/resources`

## 资源配额

### 获取资源配额

```
GET /api/v1/resources/quotas
```

查询参数:
- `namespace` (可选): 按命名空间筛选
- `resource_type` (可选): 按资源类型筛选 (cpu, memory, storage, network)

**请求示例:**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/resources/quotas?namespace=production" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "namespace": "production",
  "quotas": [
    {
      "resource_type": "cpu",
      "limit": 100,
      "used": 45,
      "available": 55,
      "unit": "cores"
    },
    {
      "resource_type": "memory",
      "limit": 256,
      "used": 128,
      "available": 128,
      "unit": "GB"
    },
    {
      "resource_type": "storage",
      "limit": 1000,
      "used": 350,
      "available": 650,
      "unit": "GB"
    },
    {
      "resource_type": "network",
      "limit": 100,
      "used": 25,
      "available": 75,
      "unit": "Mbps"
    }
  ]
}
```

### 更新资源配额

```
PUT /api/v1/resources/quotas/{namespace}
```

**请求体:**

```json
{
  "quotas": [
    {
      "resource_type": "cpu",
      "limit": 150
    },
    {
      "resource_type": "memory",
      "limit": 512
    }
  ]
}
```

**响应示例:**

```json
{
  "namespace": "production",
  "updated_quotas": [
    {
      "resource_type": "cpu",
      "limit": 150,
      "used": 45,
      "available": 105,
      "unit": "cores"
    },
    {
      "resource_type": "memory",
      "limit": 512,
      "used": 128,
      "available": 384,
      "unit": "GB"
    }
  ],
  "message": "Resource quotas updated successfully"
}
```

## 节点管理

### 获取节点列表

```
GET /api/v1/resources/nodes
```

查询参数:
- `status` (可选): 按状态筛选 (ready, not_ready, cordoned)
- `role` (可选): 按角色筛选 (master, worker)
- `limit` (可选): 返回结果数量限制 (默认: 20, 最大: 100)
- `offset` (可选): 分页偏移量 (默认: 0)

**请求示例:**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/resources/nodes?status=ready&role=worker" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "total": 15,
  "limit": 20,
  "offset": 0,
  "nodes": [
    {
      "id": "node-1",
      "name": "worker-01.palo.prod.baidu.com",
      "status": "ready",
      "role": "worker",
      "created_at": "2023-01-15T08:30:00Z",
      "updated_at": "2023-05-10T12:45:30Z",
      "ip_address": "10.0.1.101",
      "resources": {
        "cpu": {
          "capacity": 32,
          "allocatable": 30,
          "allocated": 25,
          "available": 5
        },
        "memory": {
          "capacity": 128,
          "allocatable": 120,
          "allocated": 100,
          "available": 20,
          "unit": "GB"
        },
        "pods": {
          "capacity": 110,
          "allocatable": 110,
          "allocated": 85,
          "available": 25
        }
      },
      "conditions": [
        {
          "type": "Ready",
          "status": "True",
          "last_transition_time": "2023-05-10T12:45:30Z",
          "reason": "KubeletReady",
          "message": "kubelet is posting ready status"
        }
      ],
      "labels": {
        "region": "beijing",
        "zone": "zone-a",
        "instance-type": "high-memory"
      }
    },
    // 更多节点...
  ]
}
```

### 获取节点详情

```
GET /api/v1/resources/nodes/{node_id}
```

**请求示例:**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/resources/nodes/node-1" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "id": "node-1",
  "name": "worker-01.palo.prod.baidu.com",
  "status": "ready",
  "role": "worker",
  "created_at": "2023-01-15T08:30:00Z",
  "updated_at": "2023-05-10T12:45:30Z",
  "ip_address": "10.0.1.101",
  "resources": {
    "cpu": {
      "capacity": 32,
      "allocatable": 30,
      "allocated": 25,
      "available": 5
    },
    "memory": {
      "capacity": 128,
      "allocatable": 120,
      "allocated": 100,
      "available": 20,
      "unit": "GB"
    },
    "pods": {
      "capacity": 110,
      "allocatable": 110,
      "allocated": 85,
      "available": 25
    }
  },
  "conditions": [
    {
      "type": "Ready",
      "status": "True",
      "last_transition_time": "2023-05-10T12:45:30Z",
      "reason": "KubeletReady",
      "message": "kubelet is posting ready status"
    },
    {
      "type": "DiskPressure",
      "status": "False",
      "last_transition_time": "2023-05-10T12:45:30Z",
      "reason": "KubeletHasSufficientDisk",
      "message": "kubelet has sufficient disk space available"
    },
    {
      "type": "MemoryPressure",
      "status": "False",
      "last_transition_time": "2023-05-10T12:45:30Z",
      "reason": "KubeletHasSufficientMemory",
      "message": "kubelet has sufficient memory available"
    },
    {
      "type": "PIDPressure",
      "status": "False",
      "last_transition_time": "2023-05-10T12:45:30Z",
      "reason": "KubeletHasSufficientPID",
      "message": "kubelet has sufficient PID available"
    },
    {
      "type": "NetworkUnavailable",
      "status": "False",
      "last_transition_time": "2023-01-15T08:35:00Z",
      "reason": "RouteCreated",
      "message": "RouteController created a route"
    }
  ],
  "system_info": {
    "kernel_version": "5.4.0-1045-aws",
    "os_image": "Ubuntu 20.04.4 LTS",
    "container_runtime_version": "containerd://1.5.11",
    "kubelet_version": "v1.23.6",
    "architecture": "amd64"
  },
  "labels": {
    "region": "beijing",
    "zone": "zone-a",
    "instance-type": "high-memory"
  },
  "taints": [],
  "pods": [
    {
      "name": "web-frontend-5d8fb97d55-2xvqz",
      "namespace": "production",
      "status": "running",
      "cpu_request": 0.5,
      "cpu_limit": 1.0,
      "memory_request": "512Mi",
      "memory_limit": "1Gi"
    },
    // 更多Pod...
  ]
}
```

### 节点操作

#### 隔离节点

```
POST /api/v1/resources/nodes/{node_id}/cordon
```

**响应示例:**

```json
{
  "id": "node-1",
  "name": "worker-01.palo.prod.baidu.com",
  "status": "cordoned",
  "message": "Node has been cordoned successfully"
}
```

#### 恢复节点

```
POST /api/v1/resources/nodes/{node_id}/uncordon
```

**响应示例:**

```json
{
  "id": "node-1",
  "name": "worker-01.palo.prod.baidu.com",
  "status": "ready",
  "message": "Node has been uncordoned successfully"
}
```

#### 排空节点

```
POST /api/v1/resources/nodes/{node_id}/drain
```

查询参数:
- `grace_period` (可选): 优雅终止期限(秒) (默认: 30)
- `force` (可选): 是否强制排空 (默认: false)

**请求示例:**

```bash
curl -X POST "https://tianniuprod.baidu.com/api/v1/resources/nodes/node-1/drain?grace_period=60" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "id": "node-1",
  "name": "worker-01.palo.prod.baidu.com",
  "status": "draining",
  "message": "Node drain operation has started"
}
```

## 资源使用情况

### 获取集群资源使用情况

```
GET /api/v1/resources/usage
```

查询参数:
- `period` (可选): 统计周期 (hour, day, week, month) (默认: hour)
- `start_time` (可选): 开始时间 (ISO 8601格式)
- `end_time` (可选): 结束时间 (ISO 8601格式)

**请求示例:**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/resources/usage?period=day" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "period": "day",
  "start_time": "2023-05-15T00:00:00Z",
  "end_time": "2023-05-16T00:00:00Z",
  "interval": "1h",
  "resources": {
    "cpu": {
      "capacity": 320,
      "data_points": [
        {
          "timestamp": "2023-05-15T00:00:00Z",
          "usage": 180,
          "usage_percent": 56.25
        },
        {
          "timestamp": "2023-05-15T01:00:00Z",
          "usage": 185,
          "usage_percent": 57.81
        },
        // 更多数据点...
      ]
    },
    "memory": {
      "capacity": 1024,
      "unit": "GB",
      "data_points": [
        {
          "timestamp": "2023-05-15T00:00:00Z",
          "usage": 720,
          "usage_percent": 70.31
        },
        {
          "timestamp": "2023-05-15T01:00:00Z",
          "usage": 735,
          "usage_percent": 71.78
        },
        // 更多数据点...
      ]
    },
    "storage": {
      "capacity": 10000,
      "unit": "GB",
      "data_points": [
        {
          "timestamp": "2023-05-15T00:00:00Z",
          "usage": 6500,
          "usage_percent": 65.00
        },
        {
          "timestamp": "2023-05-15T01:00:00Z",
          "usage": 6520,
          "usage_percent": 65.20
        },
        // 更多数据点...
      ]
    },
    "network": {
      "capacity": 1000,
      "unit": "Mbps",
      "data_points": [
        {
          "timestamp": "2023-05-15T00:00:00Z",
          "ingress": 350,
          "egress": 280,
          "total": 630,
          "usage_percent": 63.00
        },
        {
          "timestamp": "2023-05-15T01:00:00Z",
          "ingress": 360,
          "egress": 290,
          "total": 650,
          "usage_percent": 65.00
        },
        // 更多数据点...
      ]
    }
  }
}
```

### 获取命名空间资源使用情况

```
GET /api/v1/resources/usage/{namespace}
```

查询参数:
- `period` (可选): 统计周期 (hour, day, week, month) (默认: hour)
- `start_time` (可选): 开始时间 (ISO 8601格式)
- `end_time` (可选): 结束时间 (ISO 8601格式)

**请求示例:**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/resources/usage/production?period=day" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "namespace": "production",
  "period": "day",
  "start_time": "2023-05-15T00:00:00Z",
  "end_time": "2023-05-16T00:00:00Z",
  "interval": "1h",
  "resources": {
    "cpu": {
      "quota": 100,
      "data_points": [
        {
          "timestamp": "2023-05-15T00:00:00Z",
          "usage": 45,
          "usage_percent": 45.00
        },
        {
          "timestamp": "2023-05-15T01:00:00Z",
          "usage": 47,
          "usage_percent": 47.00
        },
        // 更多数据点...
      ]
    },
    "memory": {
      "quota": 256,
      "unit": "GB",
      "data_points": [
        {
          "timestamp": "2023-05-15T00:00:00Z",
          "usage": 128,
          "usage_percent": 50.00
        },
        {
          "timestamp": "2023-05-15T01:00:00Z",
          "usage": 130,
          "usage_percent": 50.78
        },
        // 更多数据点...
      ]
    },
    "storage": {
      "quota": 1000,
      "unit": "GB",
      "data_points": [
        {
          "timestamp": "2023-05-15T00:00:00Z",
          "usage": 350,
          "usage_percent": 35.00
        },
        {
          "timestamp": "2023-05-15T01:00:00Z",
          "usage": 352,
          "usage_percent": 35.20
        },
        // 更多数据点...
      ]
    },
    "network": {
      "quota": 100,
      "unit": "Mbps",
      "data_points": [
        {
          "timestamp": "2023-05-15T00:00:00Z",
          "ingress": 15,
          "egress": 10,
          "total": 25,
          "usage_percent": 25.00
        },
        {
          "timestamp": "2023-05-15T01:00:00Z",
          "ingress": 16,
          "egress": 11,
          "total": 27,
          "usage_percent": 27.00
        },
        // 更多数据点...
      ]
    }
  }
}
```

## 资源推荐

### 获取资源优化建议

```
GET /api/v1/resources/recommendations
```

查询参数:
- `namespace` (可选): 按命名空间筛选
- `resource_type` (可选): 按资源类型筛选 (cpu, memory, storage, network)
- `recommendation_type` (可选): 按建议类型筛选 (resize, consolidate, distribute)

**请求示例:**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/resources/recommendations?namespace=production" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**响应示例:**

```json
{
  "namespace": "production",
  "recommendations": [
    {
      "id": "rec-1",
      "resource_type": "cpu",
      "recommendation_type": "resize",
      "target": {
        "type": "deployment",
        "name": "api-backend",
        "namespace": "production"
      },
      "current_value": {
        "request": 2.0,
        "limit": 4.0
      },
      "recommended_value": {
        "request": 1.0,
        "limit": 2.0
      },
      "potential_savings": {
        "cpu": 2.0,
        "cost": "¥200/month"
      },
      "confidence": "high",
      "reason": "CPU utilization has been consistently below 30% for the past 30 days",
      "observation_period": {
        "start_time": "2023-04-15T00:00:00Z",
        "end_time": "2023-05-15T00:00:00Z"
      }
    },
    {
      "id": "rec-2",
      "resource_type": "memory",
      "recommendation_type": "resize",
      "target": {
        "type": "deployment",
        "name": "web-frontend",
        "namespace": "production"
      },
      "current_value": {
        "request": "1Gi",
        "limit": "2Gi"
      },
      "recommended_value": {
        "request": "512Mi",
        "limit": "1Gi"
      },
      "potential_savings": {
        "memory": "1Gi",
        "cost": "¥100/month"
      },
      "confidence": "medium",
      "reason": "Memory utilization has been consistently below 40% for the past 30 days",
      "observation_period": {
        "start_time": "2023-04-15T00:00:00Z",
        "end_time": "2023-05-15T00:00:00Z"
      }
    },
    {
      "id": "rec-3",
      "resource_type": "cpu",
      "recommendation_type": "consolidate",
      "target": {
        "type": "node_group",
        "name": "worker-nodes",
        "namespace": "system"
      },
      "current_value": {
        "nodes": 10,
        "total_cpu": 320,
        "used_cpu": 160
      },
      "recommended_value": {
        "nodes": 8,
        "total_cpu": 256,
        "used_cpu": 160
      },
      "potential_savings": {
        "nodes": 2,
        "cpu": 64,
        "cost": "¥1000/month"
      },
      "confidence": "high",
      "reason": "Cluster CPU utilization has been consistently below 50% for the past 30 days",
      "observation_period": {
        "start_time": "2023-04-15T00:00:00Z",
        "end_time": "2023-05-15T00:00:00Z"
      }
    }
  ]
}
```

### 应用资源优化建议

```
POST /api/v1/resources/recommendations/{recommendation_id}/apply
```

**响应示例:**

```json
{
  "id": "rec-1",
  "status": "applying",
  "message": "Resource optimization is being applied",
  "estimated_completion_time": "2023-05-16T15:30:00Z"
}
```

## 错误响应

所有API错误响应都遵循以下格式:

```json
{
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "Resource with ID 'node-1' not found",
    "details": {
      "resource_id": "node-1",
      "resource_type": "node"
    }
  }
}
```

常见错误代码:

| 错误代码 | HTTP状态码 | 描述 |
|----------|------------|------|
| RESOURCE_NOT_FOUND | 404 | 指定的资源不存在 |
| QUOTA_EXCEEDED | 403 | 超出资源配额 |
| INVALID_RESOURCE_STATE | 400 | 资源状态不允许执行请求的操作 |
| NAMESPACE_NOT_FOUND | 404 | 指定的命名空间不存在 |
| PERMISSION_DENIED | 403 | 没有执行操作的权限 |
| INTERNAL_ERROR | 500 | 内部服务器错误 |
