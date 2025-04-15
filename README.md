<div align="center">

# 🐂 TianNiu Platform Extension API

<img src="https://img.shields.io/badge/版本-v1.0.0-blue" alt="版本"> <img src="https://img.shields.io/badge/状态-稳定-green" alt="状态"> <img src="https://img.shields.io/badge/许可证-Apache%202.0-orange" alt="许可证">

**强大、灵活、高效的容器管理与编排平台**

[快速开始](#快速开始) • [API文档](#api-文档) • [示例代码](#示例代码) • [常见问题](#常见问题) • [技术支持](#技术支持)

</div>

---

## 📋 概述

**天牛平台**是百度内部自主研发的企业级容器管理与编排平台，提供从业务研发、测试、部署到交付运维的一站式解决方案。通过天牛平台扩展API，您可以轻松实现容器生命周期管理、应用编排、资源调度、服务发现、负载均衡等核心功能，大幅降低项目测试部署及后续运维成本，全面提升私有化交付能力。

天牛平台基于Kubernetes构建，在保留其强大功能的同时，针对百度业务场景进行了深度优化和扩展，提供了更加简单易用的API接口和更加丰富的功能特性。

<div align="center">
<img src="https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naLGibSEOqg9OEUozsAiaZ3EGFnKdGZ0gWh7GRvOXHXOaibKPBiaGDreKHSyLwKwKxwpLQJjMvGAVulZicg/640" alt="天牛平台架构图" width="600">
</div>

## 🚀 核心功能

通过天牛平台扩展API，您可以：

- **容器管理**：创建、启动、停止、重启、删除Docker容器，查看容器日志和状态
- **应用编排**：定义多容器应用，管理容器间依赖关系，实现复杂应用的一键部署
- **服务发现**：自动注册和发现服务，实现服务间的无缝通信
- **负载均衡**：智能分发流量，确保应用高可用和高性能
- **自动伸缩**：根据负载自动调整资源，优化资源利用率
- **持续集成/持续部署**：与CI/CD工具无缝集成，实现自动化测试和部署
- **资源管理**：分配和监控计算、网络、存储资源，设置资源配额和限制
- **多环境支持**：支持开发、测试、生产等多环境隔离和管理
- **安全控制**：基于角色的访问控制，API密钥认证，网络策略等安全特性
- **监控告警**：实时监控系统性能和健康状态，及时发现和解决问题

## 🌐 环境说明

天牛平台提供了以下环境，满足不同场景的需求：

| 环境 | 域名 | 描述 | 状态 |
|------|------|------|------|
| 🏢 **生产环境** | [tianniuprod.baidu.com](https://tianniuprod.baidu.com) | 用于生产部署的稳定环境，SLA保证99.99%可用性 | 🟢 稳定运行 |
| 📊 **PALO生产环境** | [palo.prod.baidu.com](https://palo.prod.baidu.com) | PALO数据分析组件的生产环境，提供高性能数据处理能力 | 🟢 稳定运行 |
| 🧪 **PALO开发环境** | [palo.dev.baidu.com](https://palo.dev.baidu.com) | PALO组件的开发和测试环境，用于功能验证和性能测试 | 🟢 稳定运行 |
| 🔍 **监控环境** | [monitor.tianniu.baidu.com](https://monitor.tianniu.baidu.com) | 提供全平台监控和告警服务 | 🟢 稳定运行 |
| 📚 **日志环境** | [log.tianniu.baidu.com](https://log.tianniu.baidu.com) | 集中式日志收集和分析服务 | 🟢 稳定运行 |

## 📖 API 文档

天牛平台提供了全面的API文档，帮助您快速了解和使用平台功能：

### 核心API

- [🔐 认证与授权](docs/authentication.md) - API密钥管理、OAuth认证、权限控制
- [📦 容器管理 API](docs/container_api.md) - 容器生命周期管理、日志查看、执行命令
- [🚢 部署管理 API](docs/deployment_api.md) - 应用部署、版本管理、蓝绿发布、金丝雀发布
- [💾 资源管理 API](docs/resource_api.md) - 资源配额、节点管理、资源使用统计
- [📊 监控与日志 API](docs/monitoring_api.md) - 性能监控、健康检查、日志收集与分析

### 扩展API

- [🌐 网络管理 API](docs/network_api.md) - 网络策略、服务发现、负载均衡
- [💽 存储管理 API](docs/storage_api.md) - 持久化存储、数据备份与恢复
- [🔄 CI/CD 集成](docs/cicd_integration.md) - Jenkins、GitLab CI、GitHub Actions集成
- [🔧 配置管理 API](docs/config_api.md) - 配置项管理、配置版本控制、环境变量
- [🔔 通知与告警 API](docs/notification_api.md) - 告警规则、通知渠道、事件订阅

### 参考资料

- [❓ 错误码与故障排除](docs/errors.md) - 常见错误码解释和故障排除指南
- [🛠️ SDK 与客户端工具](docs/sdk.md) - 各语言SDK和命令行工具使用说明
- [👍 最佳实践](docs/best_practices.md) - 容器设计、资源优化、安全加固等最佳实践
- [📘 API变更日志](docs/changelog.md) - API版本变更历史和兼容性说明
- [📚 术语表](docs/glossary.md) - 平台相关术语解释

## 🏁 快速开始

只需几个简单步骤，即可开始使用天牛平台：

### 1. 注册并获取API密钥

访问[天牛平台控制台](https://tianniuprod.baidu.com/console)，注册账号并创建API密钥。

### 2. 安装客户端工具

```bash
# 使用pip安装Python客户端
pip install tianniu-client

# 或使用npm安装Node.js客户端
npm install @baidu/tianniu-client

# 或使用go get安装Go客户端
go get github.com/baidu/tianniu-go-client
```

### 3. 配置客户端

```bash
# 设置API密钥环境变量
export TIANNIU_API_KEY="your-api-key"

# 或使用配置文件
cat > ~/.tianniu/config.yaml << EOF
apiVersion: v1
kind: TianNiuConfig
metadata:
  name: default-config
environments:
  - name: production
    api_endpoint: https://tianniuprod.baidu.com/api/v1
    auth:
      type: api_key
      api_key_env: TIANNIU_API_KEY
EOF
```

### 4. 创建并运行容器

```bash
# 使用Python客户端
python -c '
from tianniu import TianNiuClient

client = TianNiuClient()
container = client.containers.create(
    name="hello-world",
    image="nginx:latest",
    ports=[{"internal": 80, "external": 8080}]
)
print(f"容器已创建: {container.id}")
container.start()
print(f"容器已启动: {container.status}")
'
```

更多详细信息和示例，请参阅[快速开始指南](docs/quickstart.md)。

## 💻 示例代码

天牛平台提供了多种语言的示例代码，帮助您快速上手：

### Python示例

```python
from tianniu import TianNiuClient

# 初始化客户端
client = TianNiuClient()

# 列出所有运行中的容器
containers = client.containers.list(status="running")
for container in containers:
    print(f"ID: {container.id}, 名称: {container.name}, 镜像: {container.image}")

# 创建新部署
deployment = client.deployments.create(
    name="web-app",
    environment="production",
    version="v1.0.0",
    replicas=3,
    containers=[
        {
            "name": "web-frontend",
            "image": "registry.baidu.com/myteam/web-app:v1.0.0",
            "ports": [{"name": "http", "container_port": 80, "service_port": 8080}]
        }
    ]
)
print(f"部署已创建: {deployment.id}")
```

### Go示例

```go
package main

import (
	"fmt"
	"log"

	"github.com/baidu/tianniu-go-client/tianniu"
)

func main() {
	// 初始化客户端
	client, err := tianniu.NewClient()
	if err != nil {
		log.Fatalf("初始化客户端失败: %v", err)
	}

	// 列出所有运行中的容器
	containers, err := client.Containers.List(tianniu.ContainerListOptions{Status: "running"})
	if err != nil {
		log.Fatalf("获取容器列表失败: %v", err)
	}

	for _, container := range containers {
		fmt.Printf("ID: %s, 名称: %s, 镜像: %s\n", container.ID, container.Name, container.Image)
	}
}
```

更多示例代码请查看[examples](examples/)目录。

## 📊 性能指标

天牛平台在各种负载下表现出色，以下是基准测试结果：

| 操作 | 平均响应时间 | 每秒请求数 | 成功率 |
|------|------------|-----------|-------|
| 容器创建 | 1.2秒 | 50 | 99.99% |
| 容器启动 | 0.8秒 | 100 | 99.99% |
| 容器列表查询 | 0.15秒 | 500 | 99.999% |
| 部署创建 | 2.5秒 | 30 | 99.95% |
| 资源监控查询 | 0.1秒 | 1000 | 99.999% |

## 🔒 安全性

天牛平台采用多层次安全防护措施：

- **传输安全**：所有API通信采用TLS 1.3加密
- **身份认证**：支持API密钥、OAuth 2.0和JWT认证
- **权限控制**：基于角色的细粒度访问控制
- **审计日志**：记录所有API操作，支持合规审计
- **资源隔离**：命名空间和网络策略确保多租户隔离
- **漏洞扫描**：容器镜像自动安全扫描

## 🌟 成功案例

天牛平台已在百度内部多个业务线成功应用：

- **百度搜索**：管理超过10,000个容器，支持每日数十亿次搜索请求
- **百度智能云**：为云服务客户提供容器管理能力，管理PB级数据
- **百度地图**：支持实时位置服务，确保服务高可用性
- **百度AI平台**：为AI模型训练和推理提供弹性计算资源

## 📅 路线图

我们正在开发的新功能：

- **多集群管理**：统一管理多个Kubernetes集群
- **服务网格集成**：与Istio深度集成，提供高级流量管理
- **无服务器容器**：按需启动容器，降低资源成本
- **AI辅助运维**：智能故障诊断和自动修复
- **边缘计算支持**：将容器部署到边缘节点

## ❓ 常见问题

<details>
<summary><b>天牛平台与原生Kubernetes有什么区别？</b></summary>
<p>天牛平台基于Kubernetes构建，在保留其核心功能的同时，针对百度业务场景进行了深度优化和扩展，提供了更加简单易用的API接口、更完善的监控告警、更强大的自动化能力和更全面的安全控制。</p>
</details>

<details>
<summary><b>如何处理API请求限流？</b></summary>
<p>天牛平台对API请求有一定的限流策略。当遇到限流时，API会返回429状态码。建议实现指数退避重试策略，并考虑在客户端实现缓存以减少API调用频率。</p>
</details>

<details>
<summary><b>如何实现高可用部署？</b></summary>
<p>建议使用天牛平台的部署管理API创建多副本部署，配置适当的健康检查和自动伸缩策略，并使用负载均衡服务分发流量。详细指南请参考<a href="docs/best_practices.md#高可用部署">高可用部署最佳实践</a>。</p>
</details>

<details>
<summary><b>支持哪些容器镜像仓库？</b></summary>
<p>天牛平台支持Docker Hub、百度内部镜像仓库以及任何符合Docker Registry API规范的私有镜像仓库。</p>
</details>

更多常见问题请参阅[FAQ文档](docs/faq.md)。

## 🛠️ 技术支持

我们提供多种支持渠道，确保您能够顺利使用天牛平台：

- **文档中心**：[docs.tianniu.baidu.com](https://docs.tianniu.baidu.com)
- **开发者社区**：[community.tianniu.baidu.com](https://community.tianniu.baidu.com)
- **技术支持邮箱**：[tianniu-support@baidu.com](mailto:tianniu-support@baidu.com)
- **工单系统**：[support.tianniu.baidu.com](https://support.tianniu.baidu.com)
- **在线咨询**：工作日9:00-18:00提供实时技术支持

## 📄 API 版本

当前API版本：`v1`

所有API请求应包含版本前缀：`/api/v1/`

我们遵循[语义化版本控制](https://semver.org/)原则，确保API的兼容性和稳定性。API变更历史请参阅[变更日志](docs/changelog.md)。

## 📜 许可证

天牛平台采用[Apache 2.0许可证](LICENSE)。

---

<div align="center">

**天牛平台 - 让容器管理更简单，让应用部署更高效**

[文档](https://docs.tianniu.baidu.com) • [社区](https://community.tianniu.baidu.com) • [支持](https://support.tianniu.baidu.com)

</div>
