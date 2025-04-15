# 天牛平台快速开始指南

本指南将帮助您快速上手天牛平台，了解如何使用API进行容器和应用的管理。

## 目录

- [前提条件](#前提条件)
- [获取API密钥](#获取api密钥)
- [安装客户端工具](#安装客户端工具)
- [基本操作](#基本操作)
  - [容器管理](#容器管理)
  - [部署管理](#部署管理)
  - [资源管理](#资源管理)
- [高级功能](#高级功能)
- [故障排除](#故障排除)

## 前提条件

在开始使用天牛平台之前，您需要：

1. 拥有天牛平台账号（如果没有，请联系管理员创建）
2. 安装以下工具：
   - cURL（用于API请求）
   - 您选择的编程语言环境（Python、Go、Node.js等）

## 获取API密钥

1. 登录[天牛平台控制台](https://tianniuprod.baidu.com/console)
2. 导航至"个人设置" > "API密钥"
3. 点击"创建新密钥"按钮
4. 输入密钥名称和描述
5. 选择适当的权限范围
6. 点击"创建"按钮
7. 保存生成的API密钥（注意：密钥只会显示一次）

## 安装客户端工具

天牛平台提供了多种语言的客户端SDK，您可以选择适合您的编程语言：

### Python

```bash
pip install tianniu-client
```

### Go

```bash
go get github.com/baidu/tianniu-go-client
```

### Node.js

```bash
npm install @baidu/tianniu-client
```

### 命令行工具

```bash
# 使用npm安装
npm install -g tianniu-cli

# 或使用二进制安装
curl -sSL https://tianniuprod.baidu.com/download/cli/install.sh | bash
```

## 基本操作

### 容器管理

#### 创建容器

**使用cURL：**

```bash
curl -X POST "https://tianniuprod.baidu.com/api/v1/containers" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "web-server",
    "image": "nginx:latest",
    "ports": [
      {
        "internal": 80,
        "external": 8080,
        "protocol": "tcp"
      }
    ],
    "environment_variables": [
      {
        "name": "NGINX_HOST",
        "value": "example.com"
      }
    ]
  }'
```

**使用Python SDK：**

```python
from tianniu import TianNiuClient

# 初始化客户端
client = TianNiuClient(api_key="YOUR_API_KEY")

# 创建容器
container = client.containers.create(
    name="web-server",
    image="nginx:latest",
    ports=[{"internal": 80, "external": 8080, "protocol": "tcp"}],
    environment_variables=[{"name": "NGINX_HOST", "value": "example.com"}]
)

print(f"容器已创建: {container.id}")
```

**使用Go SDK：**

```go
package main

import (
    "fmt"
    "log"

    "github.com/baidu/tianniu-go-client/tianniu"
)

func main() {
    // 初始化客户端
    client, err := tianniu.NewClient(tianniu.WithAPIKey("YOUR_API_KEY"))
    if err != nil {
        log.Fatalf("初始化客户端失败: %v", err)
    }

    // 创建容器
    container, err := client.Containers.Create(tianniu.ContainerCreateOptions{
        Name:  "web-server",
        Image: "nginx:latest",
        Ports: []tianniu.Port{
            {
                Internal: 80,
                External: 8080,
                Protocol: "tcp",
            },
        },
        EnvironmentVariables: []tianniu.EnvVar{
            {
                Name:  "NGINX_HOST",
                Value: "example.com",
            },
        },
    })
    if err != nil {
        log.Fatalf("创建容器失败: %v", err)
    }

    fmt.Printf("容器已创建: %s\n", container.ID)
}
```

#### 列出容器

**使用cURL：**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/containers?status=running" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**使用Python SDK：**

```python
from tianniu import TianNiuClient

# 初始化客户端
client = TianNiuClient(api_key="YOUR_API_KEY")

# 列出运行中的容器
containers = client.containers.list(status="running")

for container in containers:
    print(f"ID: {container.id}, 名称: {container.name}, 状态: {container.status}")
```

#### 启动容器

**使用cURL：**

```bash
curl -X POST "https://tianniuprod.baidu.com/api/v1/containers/CONTAINER_ID/start" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**使用Python SDK：**

```python
from tianniu import TianNiuClient

# 初始化客户端
client = TianNiuClient(api_key="YOUR_API_KEY")

# 获取容器
container = client.containers.get("CONTAINER_ID")

# 启动容器
container.start()
print(f"容器已启动: {container.status}")
```

### 部署管理

#### 创建部署

**使用cURL：**

```bash
curl -X POST "https://tianniuprod.baidu.com/api/v1/deployments" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "web-app",
    "environment": "production",
    "version": "v1.0.0",
    "replicas": 3,
    "containers": [
      {
        "name": "web-frontend",
        "image": "registry.baidu.com/myteam/web-app:v1.0.0",
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
        }
      }
    ],
    "services": [
      {
        "name": "web-app-svc",
        "type": "LoadBalancer",
        "ports": [
          {
            "name": "http",
            "port": 80,
            "target_port": 8080
          }
        ]
      }
    ]
  }'
```

**使用Python SDK：**

```python
from tianniu import TianNiuClient

# 初始化客户端
client = TianNiuClient(api_key="YOUR_API_KEY")

# 创建部署
deployment = client.deployments.create(
    name="web-app",
    environment="production",
    version="v1.0.0",
    replicas=3,
    containers=[
        {
            "name": "web-frontend",
            "image": "registry.baidu.com/myteam/web-app:v1.0.0",
            "ports": [{"name": "http", "container_port": 80, "service_port": 8080}],
            "resources": {
                "limits": {"cpu": "1.0", "memory": "1Gi"},
                "requests": {"cpu": "0.5", "memory": "512Mi"}
            }
        }
    ],
    services=[
        {
            "name": "web-app-svc",
            "type": "LoadBalancer",
            "ports": [{"name": "http", "port": 80, "target_port": 8080}]
        }
    ]
)

print(f"部署已创建: {deployment.id}")
```

#### 扩缩容部署

**使用cURL：**

```bash
curl -X POST "https://tianniuprod.baidu.com/api/v1/deployments/DEPLOYMENT_ID/scale" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "replicas": 5
  }'
```

**使用Python SDK：**

```python
from tianniu import TianNiuClient

# 初始化客户端
client = TianNiuClient(api_key="YOUR_API_KEY")

# 获取部署
deployment = client.deployments.get("DEPLOYMENT_ID")

# 扩缩容部署
deployment.scale(replicas=5)
print(f"部署已扩展至 {deployment.replicas} 个副本")
```

### 资源管理

#### 获取资源配额

**使用cURL：**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/resources/quotas?namespace=production" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**使用Python SDK：**

```python
from tianniu import TianNiuClient

# 初始化客户端
client = TianNiuClient(api_key="YOUR_API_KEY")

# 获取资源配额
quotas = client.resources.get_quotas(namespace="production")

for resource_type, quota in quotas.items():
    print(f"{resource_type}: 使用 {quota['used']}/{quota['limit']} {quota['unit']}")
```

## 高级功能

### 蓝绿部署

```python
from tianniu import TianNiuClient

# 初始化客户端
client = TianNiuClient(api_key="YOUR_API_KEY")

# 创建蓝绿部署
blue_green_deployment = client.deployments.create_blue_green(
    name="payment-service",
    environment="production",
    blue_version={
        "version": "v2.0.0",
        "image": "registry.baidu.com/payment/service:v2.0.0",
        "replicas": 3
    },
    green_version={
        "version": "v2.1.0",
        "image": "registry.baidu.com/payment/service:v2.1.0",
        "replicas": 3
    },
    test_traffic_percentage=20,
    auto_promote_after="2h",
    rollback_threshold={
        "error_rate": 1.0,
        "latency_p99_ms": 500
    }
)

print(f"蓝绿部署已创建: {blue_green_deployment.id}")
```

### 金丝雀发布

```python
from tianniu import TianNiuClient

# 初始化客户端
client = TianNiuClient(api_key="YOUR_API_KEY")

# 创建金丝雀发布
canary_deployment = client.deployments.create_canary(
    name="search-service",
    environment="production",
    base_version={
        "version": "v3.2.0",
        "image": "registry.baidu.com/search/service:v3.2.0",
        "replicas": 10
    },
    canary_version={
        "version": "v3.3.0",
        "image": "registry.baidu.com/search/service:v3.3.0"
    },
    stages=[
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
    auto_promote=True
)

print(f"金丝雀发布已创建: {canary_deployment.id}")
```

## 故障排除

### 常见错误

| 错误代码 | 描述 | 解决方案 |
|----------|------|----------|
| AUTHENTICATION_FAILED | API密钥无效或已过期 | 检查API密钥是否正确，必要时重新生成 |
| RESOURCE_NOT_FOUND | 请求的资源不存在 | 确认资源ID是否正确 |
| QUOTA_EXCEEDED | 超出资源配额限制 | 请求增加配额或释放未使用的资源 |
| RATE_LIMIT_EXCEEDED | 超出API请求频率限制 | 减少请求频率，实现指数退避重试 |
| INVALID_PARAMETER | 请求参数无效 | 检查请求参数是否符合API规范 |

### 获取帮助

如果您遇到问题，可以通过以下方式获取帮助：

- 查阅[API文档](https://docs.tianniu.baidu.com)
- 访问[开发者社区](https://community.tianniu.baidu.com)
- 联系[技术支持](mailto:tianniu-support@baidu.com)
- 提交[工单](https://support.tianniu.baidu.com)

## 下一步

- 了解[最佳实践](best_practices.md)
- 探索[高级功能](advanced_features.md)
- 查看[示例项目](https://github.com/baidu/tianniu-examples)
