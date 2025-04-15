# 天牛平台常见问题解答 (FAQ)

## 目录

- [一般问题](#一般问题)
- [账号与认证](#账号与认证)
- [容器管理](#容器管理)
- [部署管理](#部署管理)
- [资源管理](#资源管理)
- [网络与存储](#网络与存储)
- [监控与日志](#监控与日志)
- [安全性](#安全性)
- [性能与扩展性](#性能与扩展性)
- [故障排除](#故障排除)

## 一般问题

### 什么是天牛平台？

天牛平台是百度内部自主研发的企业级容器管理与编排平台，提供从业务研发、测试、部署到交付运维的一站式解决方案。它基于Kubernetes构建，针对百度业务场景进行了深度优化和扩展，提供了更加简单易用的API接口和更加丰富的功能特性。

### 天牛平台与原生Kubernetes有什么区别？

天牛平台基于Kubernetes构建，在保留其核心功能的同时，针对百度业务场景进行了深度优化和扩展，提供了以下优势：

1. **简化的API接口**：更加直观和易用的API设计，降低学习成本
2. **增强的监控能力**：集成了更全面的监控和告警功能
3. **自动化运维**：提供智能故障诊断和自动修复能力
4. **多环境管理**：统一管理开发、测试、生产等多个环境
5. **安全加固**：增强的安全控制和合规审计功能
6. **百度生态集成**：与百度内部其他服务的无缝集成

### 天牛平台支持哪些编程语言的SDK？

天牛平台目前提供以下编程语言的官方SDK：

- Python
- Go
- Node.js
- Java
- .NET

其他语言可以通过RESTful API直接调用。

### 如何获取天牛平台的访问权限？

如果您是百度内部员工，可以通过内部OA系统申请访问权限。如果您是外部合作伙伴，请联系您的百度对接人申请访问权限。

## 账号与认证

### 如何创建API密钥？

1. 登录[天牛平台控制台](https://tianniuprod.baidu.com/console)
2. 导航至"个人设置" > "API密钥"
3. 点击"创建新密钥"按钮
4. 输入密钥名称和描述
5. 选择适当的权限范围
6. 点击"创建"按钮
7. 保存生成的API密钥（注意：密钥只会显示一次）

### API密钥有效期是多久？

默认情况下，API密钥的有效期为1年。您可以在创建密钥时自定义有效期，最长可设置为5年。

### 如何轮换API密钥？

为了安全起见，建议定期轮换API密钥。轮换步骤如下：

1. 创建新的API密钥
2. 更新应用程序配置，使用新的API密钥
3. 验证应用程序正常工作
4. 删除旧的API密钥

### 如何实现基于角色的访问控制？

天牛平台支持基于角色的访问控制（RBAC）。您可以：

1. 创建自定义角色，定义特定的权限集
2. 将用户分配到一个或多个角色
3. 为API密钥分配特定的权限范围

详细配置请参考[认证与授权](authentication.md)文档。

## 容器管理

### 支持哪些容器镜像格式？

天牛平台支持标准的Docker容器镜像格式，以及兼容OCI（Open Container Initiative）规范的镜像。

### 如何将私有镜像仓库与天牛平台集成？

您可以通过以下步骤将私有镜像仓库与天牛平台集成：

1. 在天牛平台控制台中导航至"设置" > "镜像仓库"
2. 点击"添加镜像仓库"
3. 输入仓库URL、认证信息和其他配置
4. 点击"保存"

之后，您就可以在创建容器或部署时使用私有仓库中的镜像。

### 容器日志保留多长时间？

默认情况下，容器日志保留7天。您可以根据需要调整日志保留策略，最长可保留90天。

### 如何设置容器的资源限制？

您可以在创建容器时设置CPU和内存的请求（requests）和限制（limits）：

```json
{
  "resources": {
    "limits": {
      "cpu": "2.0",
      "memory": "2Gi"
    },
    "requests": {
      "cpu": "1.0",
      "memory": "1Gi"
    }
  }
}
```

这确保容器至少获得1个CPU核心和1GB内存，最多可使用2个CPU核心和2GB内存。

### 如何在容器中执行命令？

您可以使用API或SDK在运行中的容器中执行命令：

```python
from tianniu import TianNiuClient

client = TianNiuClient(api_key="YOUR_API_KEY")
container = client.containers.get("CONTAINER_ID")

# 执行命令
result = container.exec(["ls", "-la", "/app"])
print(result.stdout)
```

## 部署管理

### 什么是蓝绿部署？

蓝绿部署是一种应用发布策略，通过同时维护两个生产环境来减少停机时间和风险。"蓝"环境是当前的生产环境，而"绿"环境是新版本的环境。在验证"绿"环境正常工作后，流量会从"蓝"环境切换到"绿"环境。

### 什么是金丝雀发布？

金丝雀发布是一种渐进式发布策略，先将新版本部署到一小部分服务器，并将少量用户流量引导到这些服务器。如果新版本运行正常，则逐步增加流量比例，直到所有流量都切换到新版本。

### 如何回滚失败的部署？

您可以使用API或SDK回滚到之前的部署版本：

```python
from tianniu import TianNiuClient

client = TianNiuClient(api_key="YOUR_API_KEY")
deployment = client.deployments.get("DEPLOYMENT_ID")

# 回滚到指定版本
deployment.rollback(version="v1.2.3")
```

### 部署的最大副本数是多少？

默认情况下，单个部署的最大副本数为100。如果您需要更多副本，请联系管理员调整配额。

### 如何设置自动伸缩？

您可以为部署配置水平自动伸缩（HPA）和垂直自动伸缩（VPA）：

```python
from tianniu import TianNiuClient

client = TianNiuClient(api_key="YOUR_API_KEY")

# 创建部署时设置自动伸缩
deployment = client.deployments.create(
    name="web-app",
    # ... 其他配置 ...
    auto_scaling={
        "horizontal": {
            "min_replicas": 2,
            "max_replicas": 10,
            "metrics": [
                {
                    "type": "Resource",
                    "resource": {
                        "name": "cpu",
                        "target_average_utilization": 80
                    }
                }
            ]
        },
        "vertical": {
            "update_mode": "Auto",
            "min_allowed": {
                "cpu": "0.1",
                "memory": "128Mi"
            },
            "max_allowed": {
                "cpu": "4.0",
                "memory": "8Gi"
            }
        }
    }
)
```

## 资源管理

### 如何查看资源使用情况？

您可以通过API或SDK查看资源使用情况：

```python
from tianniu import TianNiuClient

client = TianNiuClient(api_key="YOUR_API_KEY")

# 获取集群资源使用情况
usage = client.resources.get_usage(period="day")
print(f"CPU使用率: {usage['cpu']['data_points'][-1]['usage_percent']}%")
print(f"内存使用率: {usage['memory']['data_points'][-1]['usage_percent']}%")

# 获取命名空间资源使用情况
namespace_usage = client.resources.get_usage("production", period="day")
print(f"生产环境CPU使用率: {namespace_usage['cpu']['data_points'][-1]['usage_percent']}%")
```

### 如何设置资源配额？

您可以为命名空间设置资源配额：

```python
from tianniu import TianNiuClient

client = TianNiuClient(api_key="YOUR_API_KEY")

# 更新资源配额
client.resources.update_quotas("production", {
    "cpu": 150,
    "memory": 512
})
```

### 如何获取资源优化建议？

天牛平台提供资源优化建议，帮助您优化资源利用率：

```python
from tianniu import TianNiuClient

client = TianNiuClient(api_key="YOUR_API_KEY")

# 获取资源优化建议
recommendations = client.resources.get_recommendations(namespace="production")
for rec in recommendations:
    print(f"建议: {rec['recommendation_type']} {rec['resource_type']}")
    print(f"目标: {rec['target']['type']} {rec['target']['name']}")
    print(f"当前值: {rec['current_value']}")
    print(f"建议值: {rec['recommended_value']}")
    print(f"潜在节省: {rec['potential_savings']}")
```

## 网络与存储

### 支持哪些网络模式？

天牛平台支持以下网络模式：

- **ClusterIP**：仅集群内部可访问的服务
- **NodePort**：通过节点端口暴露服务
- **LoadBalancer**：使用云提供商的负载均衡器暴露服务
- **ExternalName**：将服务映射到外部DNS名称

### 如何配置服务发现？

天牛平台自动为每个服务创建DNS记录，格式为`<service-name>.<namespace>.svc.cluster.local`。您可以在容器中使用这些DNS名称访问其他服务。

### 支持哪些存储类型？

天牛平台支持以下存储类型：

- **临时存储**：容器重启后数据会丢失
- **持久卷**：容器重启后数据仍然保留
- **共享存储**：多个容器可以同时访问的存储

### 如何备份和恢复数据？

天牛平台提供数据备份和恢复功能：

```python
from tianniu import TianNiuClient

client = TianNiuClient(api_key="YOUR_API_KEY")

# 创建数据备份
backup = client.storage.create_backup(
    name="database-backup",
    volume_id="vol-123456",
    description="Daily database backup"
)

# 从备份恢复数据
restore = client.storage.restore_from_backup(
    backup_id=backup.id,
    target_volume_id="vol-654321"
)
```

## 监控与日志

### 如何查看容器日志？

您可以通过API或SDK查看容器日志：

```python
from tianniu import TianNiuClient

client = TianNiuClient(api_key="YOUR_API_KEY")
container = client.containers.get("CONTAINER_ID")

# 获取最近100行日志
logs = container.logs(tail=100)
for log in logs:
    print(f"[{log.timestamp}] {log.message}")

# 持续获取新日志
for log in container.logs(follow=True):
    print(f"[{log.timestamp}] {log.message}")
```

### 如何设置监控告警？

您可以通过API或SDK设置监控告警：

```python
from tianniu import TianNiuClient

client = TianNiuClient(api_key="YOUR_API_KEY")

# 创建告警规则
alert = client.monitoring.create_alert(
    name="high-cpu-usage",
    description="Alert when CPU usage is high",
    target={
        "type": "deployment",
        "name": "web-app",
        "namespace": "production"
    },
    condition={
        "metric": "cpu_usage_percent",
        "operator": ">",
        "threshold": 90,
        "duration": "5m"
    },
    severity="warning",
    notifications=[
        {
            "type": "email",
            "recipients": ["admin@example.com"]
        },
        {
            "type": "webhook",
            "url": "https://example.com/webhook"
        }
    ]
)
```

### 支持哪些监控指标？

天牛平台支持以下监控指标：

- **CPU使用率**：容器CPU使用百分比
- **内存使用率**：容器内存使用百分比
- **网络流量**：入站和出站网络流量
- **磁盘I/O**：读写操作和吞吐量
- **请求延迟**：API请求的响应时间
- **错误率**：API请求的错误率
- **自定义指标**：通过Prometheus导出器暴露的自定义指标

## 安全性

### 如何保护敏感数据？

天牛平台提供以下机制保护敏感数据：

1. **密钥管理**：安全存储密码、API密钥等敏感信息
2. **环境变量**：通过环境变量注入敏感配置
3. **加密存储**：自动加密存储中的敏感数据
4. **网络隔离**：通过网络策略限制容器间通信

### 支持哪些认证方式？

天牛平台支持以下认证方式：

- **API密钥**：用于服务到服务的认证
- **OAuth 2.0**：用于用户级别的认证
- **JWT**：用于单点登录集成
- **服务账户**：用于容器内部访问API

### 如何实现网络隔离？

您可以使用网络策略实现容器间的网络隔离：

```python
from tianniu import TianNiuClient

client = TianNiuClient(api_key="YOUR_API_KEY")

# 创建网络策略
policy = client.network.create_policy(
    name="frontend-backend-only",
    namespace="production",
    pod_selector={"app": "frontend"},
    ingress_rules=[
        {
            "from": [
                {"pod_selector": {"app": "backend"}}
            ],
            "ports": [
                {"port": 80, "protocol": "TCP"}
            ]
        }
    ]
)
```

### 如何进行容器镜像安全扫描？

天牛平台集成了容器镜像安全扫描功能，可以自动扫描镜像中的安全漏洞：

```python
from tianniu import TianNiuClient

client = TianNiuClient(api_key="YOUR_API_KEY")

# 扫描镜像
scan_result = client.security.scan_image("nginx:latest")
for vulnerability in scan_result.vulnerabilities:
    print(f"漏洞: {vulnerability.id}, 严重性: {vulnerability.severity}")
```

## 性能与扩展性

### 天牛平台的API请求限制是多少？

天牛平台对API请求有以下限制：

- 认证API：每分钟100次请求
- 读取API：每分钟1000次请求
- 写入API：每分钟300次请求

如果您需要更高的限制，请联系管理员。

### 如何处理API请求限流？

当遇到API请求限流时，您应该实现指数退避重试策略：

```python
import time
import random

def api_request_with_retry(func, max_retries=5, base_delay=1):
    retries = 0
    while retries < max_retries:
        try:
            return func()
        except RateLimitExceededError:
            retries += 1
            if retries == max_retries:
                raise
            delay = base_delay * (2 ** retries) + random.uniform(0, 0.5)
            time.sleep(delay)
```

### 单个集群支持的最大节点数是多少？

天牛平台单个集群支持的最大节点数为1000个。如果您需要更大规模的集群，建议使用多集群管理功能。

### 如何优化API性能？

优化API性能的建议：

1. 使用批量操作API而不是多次调用单个操作API
2. 实现客户端缓存，减少重复请求
3. 只请求必要的字段，减少响应大小
4. 使用分页API获取大量数据，避免一次请求过多数据
5. 使用适当的索引和筛选条件优化查询

## 故障排除

### 如何排查容器启动失败？

排查容器启动失败的步骤：

1. 检查容器日志：`client.containers.get(id).logs()`
2. 检查容器事件：`client.containers.get(id).events()`
3. 检查容器状态详情：`client.containers.get(id).status_details`
4. 检查镜像是否存在且可访问
5. 检查资源限制是否合理
6. 检查健康检查配置是否正确

### 如何排查网络连接问题？

排查网络连接问题的步骤：

1. 检查网络策略是否正确配置
2. 检查服务和端口映射是否正确
3. 使用网络诊断工具：`client.network.diagnose(source_pod, destination_pod, port)`
4. 检查DNS解析是否正常
5. 检查负载均衡器配置

### 如何获取API请求的详细错误信息？

当API请求失败时，响应中会包含详细的错误信息：

```json
{
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "Container with ID 'a1b2c3d4e5f6' not found",
    "details": {
      "resource_id": "a1b2c3d4e5f6",
      "resource_type": "container"
    },
    "request_id": "req-1234567890"
  }
}
```

在提交支持工单时，请务必包含`request_id`，这有助于技术支持团队快速定位问题。

### 如何查看API请求历史？

您可以在天牛平台控制台中查看API请求历史：

1. 登录[天牛平台控制台](https://tianniuprod.baidu.com/console)
2. 导航至"个人设置" > "API活动"
3. 查看API请求历史，包括请求时间、方法、路径、状态码等信息

### 如何联系技术支持？

如果您遇到无法解决的问题，可以通过以下方式联系技术支持：

- 邮箱：[tianniu-support@baidu.com](mailto:tianniu-support@baidu.com)
- 工单系统：[support.tianniu.baidu.com](https://support.tianniu.baidu.com)
- 在线咨询：工作日9:00-18:00提供实时技术支持
