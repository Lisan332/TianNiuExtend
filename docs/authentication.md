# 认证与授权

## 概述

天牛平台使用基于令牌（Token）的认证机制，所有API请求都需要进行身份验证。认证系统支持多种身份验证方式，包括API密钥、OAuth 2.0和服务账户令牌。

## 认证方式

### API密钥认证

最简单的认证方式是使用API密钥。您可以在天牛平台控制台生成API密钥。

**请求示例：**

```bash
curl -X GET "https://tianniuprod.baidu.com/api/v1/containers" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### OAuth 2.0认证

对于需要用户级别权限的应用，天牛平台支持OAuth 2.0认证流程。

**获取授权码：**

```
GET https://tianniuprod.baidu.com/oauth/authorize
```

参数：
- `client_id`: 您的应用ID
- `redirect_uri`: 回调URL
- `response_type`: 设置为"code"
- `scope`: 请求的权限范围

**获取访问令牌：**

```
POST https://tianniuprod.baidu.com/oauth/token
```

参数：
- `grant_type`: 设置为"authorization_code"
- `code`: 授权码
- `client_id`: 您的应用ID
- `client_secret`: 您的应用密钥
- `redirect_uri`: 回调URL

### 服务账户令牌

对于自动化系统和CI/CD流程，建议使用服务账户令牌。

**创建服务账户令牌：**

```bash
curl -X POST "https://tianniuprod.baidu.com/api/v1/service-accounts" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "ci-pipeline",
    "description": "Token for CI pipeline",
    "permissions": ["container:read", "container:write", "deployment:read"]
  }'
```

## 权限模型

天牛平台使用基于角色的访问控制（RBAC）模型。每个API密钥或用户账户都与一个或多个角色关联，每个角色定义了一组权限。

### 预定义角色

| 角色 | 描述 | 权限 |
|------|------|------|
| 管理员 | 完全访问所有资源 | 所有权限 |
| 开发者 | 管理容器和部署 | container:*, deployment:*, network:read, storage:read |
| 运维 | 监控和管理资源 | container:read, deployment:read, monitoring:*, resource:* |
| 只读用户 | 只读访问 | *:read |

### 自定义角色

您可以通过API创建自定义角色：

```bash
curl -X POST "https://tianniuprod.baidu.com/api/v1/roles" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "deployment-manager",
    "description": "Can manage deployments only",
    "permissions": ["deployment:*", "container:read"]
  }'
```

## 令牌管理

### 刷新令牌

OAuth访问令牌有效期为2小时，刷新令牌有效期为30天。使用刷新令牌获取新的访问令牌：

```bash
curl -X POST "https://tianniuprod.baidu.com/oauth/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=refresh_token&refresh_token=YOUR_REFRESH_TOKEN&client_id=YOUR_CLIENT_ID&client_secret=YOUR_CLIENT_SECRET"
```

### 撤销令牌

撤销不再需要的令牌：

```bash
curl -X POST "https://tianniuprod.baidu.com/oauth/revoke" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "token=YOUR_TOKEN&token_type_hint=access_token&client_id=YOUR_CLIENT_ID&client_secret=YOUR_CLIENT_SECRET"
```

## 最佳实践

1. 定期轮换API密钥和服务账户令牌
2. 遵循最小权限原则，只授予必要的权限
3. 在生产环境中使用HTTPS进行所有API通信
4. 不要在客户端代码中硬编码API密钥
5. 使用环境变量或安全的密钥管理服务存储敏感凭据
