# Go Kevin

基于 Go 语言的企业级 AI 智能体 SaaS 平台，基于 tRPC-Agent-Go 框架构建。

## 技术栈

- **语言**: Go 1.22+
- **Web 框架**: Gin
- **ORM**: GORM
- **AI 框架**: tRPC-Agent-Go
- **认证**: JWT + RBAC
- **缓存**: Redis
- **消息队列**: RabbitMQ
- **日志**: Zap
- **配置**: Viper

## 快速开始

### 环境要求

- Go 1.22+
- MySQL 8.0+
- Redis 7.0+

### 安装

```bash
git clone https://github.com/kevin-ai/go-kevin.git
cd go-kevin
go mod download
```

### 配置

编辑 `configs/config.yaml`:

```yaml
database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: your_password
  dbname: kevin_app

redis:
  host: 127.0.0.1
  port: 6379
  password: your_password
```

### 运行

```bash
# 开发环境
make run

# 或直接运行
go run cmd/server/main.go
```

### 测试

```bash
# 运行所有测试
make test

# 运行单元测试
make test-unit

# 运行集成测试
make test-integration

# 运行 API 测试
make test-e2e

# 生成覆盖率报告
make test-coverage
```

### Docker 部署

```bash
# 构建镜像
make docker-build

# 启动服务
make docker-run

# 停止服务
make docker-stop
```

## 项目结构

```
go-kevin/
├── cmd/                    # 应用入口
│   └── server/             # HTTP 服务器
├── internal/               # 内部实现
│   ├── domain/             # 领域层
│   │   ├── entity/         # 实体定义
│   │   ├── repository/     # 仓储接口
│   │   └── event/          # 领域事件
│   ├── application/        # 应用层
│   │   ├── command/        # CQRS 命令
│   │   ├── query/          # CQRS 查询
│   │   └── dto/            # 数据传输对象
│   ├── infrastructure/     # 基础设施层
│   │   ├── persistence/    # 持久化实现
│   │   ├── eventbus/       # 事件总线
│   │   └── ai/             # AI 集成
│   └── interfaces/         # 接口层
│       └── http/           # HTTP 接口
│           ├── handler/    # 处理器
│           ├── middleware/ # 中间件
│           └── router/     # 路由定义
├── pkg/                    # 公共工具包
│   ├── config/             # 配置管理
│   ├── logger/             # 日志
│   ├── auth/               # JWT 认证
│   └── response/           # 统一响应
├── configs/                # 配置文件
├── docs/                   # 文档
├── test/                   # 测试
│   ├── integration/        # 集成测试
│   └── e2e/                # API 测试
├── Dockerfile
├── docker-compose.yaml
├── Makefile
└── README.md
```

## 架构设计

项目采用 DDD + Event-Driven 架构，CQRS 分离命令与查询：

```
┌─────────────────────────────────────────────────────────────┐
│                     Presentation Layer                        │
│                   (Gin HTTP Handlers)                        │
├─────────────────────────────────────────────────────────────┤
│                   Application Layer                          │
│           (Command / Query Handlers - CQRS)                 │
├─────────────────────────────────────────────────────────────┤
│                      Domain Layer                            │
│          (Entities / Repository Interfaces / Events)        │
├─────────────────────────────────────────────────────────────┤
│                   Infrastructure Layer                       │
│       (GORM / Redis / RabbitMQ / Event Bus)                │
└─────────────────────────────────────────────────────────────┘
```

## API 文档

启动服务后访问: http://localhost:9901/swagger

### 主要 API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/auth/login | 用户登录 |
| POST | /api/v1/auth/register | 用户注册 |
| POST | /api/v1/auth/refresh | 刷新 Token |
| GET | /api/v1/user | 获取用户列表 |
| GET | /api/v1/user/:id | 获取用户详情 |
| POST | /api/v1/user | 创建用户 |
| PUT | /api/v1/user/:id | 更新用户 |
| DELETE | /api/v1/user/:id | 删除用户 |
| GET | /api/v1/aiapps | AI 应用列表 |
| POST | /api/v1/aichat/sessions | 创建聊天会话 |
| POST | /api/v1/aikmss/:id/query | 知识库问答 |

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

MIT License
