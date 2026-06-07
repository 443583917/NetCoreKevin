# NetCoreKevin Go 重构设计方案

**文档版本**: v1.0  
**创建日期**: 2026/06/07  
**状态**: 已批准

---

## 目录

- [1. 项目概述](#1-项目概述)
- [2. 架构选型](#2-架构选型)
- [3. 项目目录结构](#3-项目目录结构)
- [4. 技术栈与依赖库](#4-技术栈与依赖库)
- [5. 领域模型设计](#5-领域模型设计)
- [6. 应用层设计](#6-应用层设计)
- [7. tRPC-Agent-Go 集成](#7-trpc-agent-go-集成)
- [8. HTTP 接口层](#8-http-接口层)
- [9. 数据库迁移](#9-数据库迁移)
- [10. 配置与部署](#10-配置与部署)
- [11. 测试策略](#11-测试策略)

---

## 1. 项目概述

### 1.1 重构目标

将 NetCoreKevin 项目后端从 .NET 9 迁移到 Go 语言，基于 tRPC-Agent-Go 框架构建企业级 AI 智能体 SaaS 平台。

### 1.2 重构策略

- **完整重构**：一次性将所有模块迁移到 Go
- **数据库迁移**：从 EF Core 迁移到 GORM，保留现有数据
- **API 适度调整**：保持核心接口兼容，允许小幅优化

### 1.3 核心模块

| 模块 | 优先级 | 说明 |
|------|--------|------|
| AI 智能体 | P0 | tRPC-Agent-Go 核心集成 |
| AI 知识库/RAG | P0 | 向量数据库集成 |
| 用户权限 | P1 | 基础设施模块 |
| 认证授权 | P1 | JWT + RBAC |
| 任务调度 | P2 | Cron 定时任务 |
| 消息队列 | P2 | RabbitMQ 集成 |
| 文件存储 | P3 | 多云存储支持 |

---

## 2. 架构选型

### 2.1 架构模式

**DDD + Event-Driven (领域驱动 + 事件驱动)**

```
┌─────────────────────────────────────────────────────────────┐
│                     Presentation Layer                        │
│                   (tRPC / REST API)                          │
├─────────────────────────────────────────────────────────────┤
│                   Application Layer                          │
│           (Command / Query Handlers - CQRS)                 │
├─────────────────────────────────────────────────────────────┤
│                      Domain Layer                            │
│          (Aggregate Root / Domain Events)                   │
├─────────────────────────────────────────────────────────────┤
│                   Infrastructure Layer                       │
│       (Event Store / Read Model / External Services)        │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 设计原则

1. **依赖反转**：领域层不依赖任何框架
2. **CQRS 分离**：命令和查询职责分离
3. **事件驱动**：通过领域事件解耦模块
4. **模块化单体**：清晰边界，未来可拆分微服务

---

## 3. 项目目录结构

```
go-kevin/
├── cmd/                            # 应用入口
│   └── server/
│       └── main.go
│
├── internal/                       # 内部实现
│   ├── domain/                     # 领域层
│   │   ├── entity/                 # 实体定义
│   │   ├── repository/             # 仓储接口
│   │   ├── event/                  # 领域事件
│   │   └── service/                # 领域服务
│   │
│   ├── application/                # 应用层
│   │   ├── command/                # 命令处理器
│   │   ├── query/                  # 查询处理器
│   │   └── dto/                    # 数据传输对象
│   │
│   ├── infrastructure/             # 基础设施层
│   │   ├── persistence/            # 持久化实现
│   │   │   ├── gorm/
│   │   │   │   ├── model/          # GORM 模型
│   │   │   │   └── repository/     # 仓储实现
│   │   │   └── migration/          # 数据库迁移
│   │   ├── cache/                  # 缓存实现
│   │   ├── message/                # 消息队列
│   │   ├── ai/                     # AI 集成
│   │   │   ├── agent/              # tRPC-Agent-Go
│   │   │   └── llm/                # LLM 云服务
│   │   └── eventbus/               # 事件总线
│   │
│   └── interfaces/                 # 接口层
│       ├── http/                   # HTTP 接口
│       │   ├── handler/            # 处理器
│       │   ├── middleware/         # 中间件
│       │   └── router/             # 路由定义
│       └── trpc/                   # tRPC 接口
│
├── pkg/                            # 公共工具包
│   ├── config/                     # 配置管理
│   ├── logger/                     # 日志
│   ├── auth/                       # JWT 认证
│   ├── response/                   # 统一响应
│   └── utils/                      # 工具函数
│
├── api/                            # API 定义
├── configs/                        # 配置文件
├── scripts/                        # 脚本
├── docs/                           # 文档
├── test/                           # 测试
├── go.mod
├── Makefile
└── README.md
```

---

## 4. 技术栈与依赖库

### 4.1 核心框架

| 组件 | 技术选型 | 版本 | 用途 |
|------|----------|------|------|
| Web 框架 | Gin | v1.9+ | HTTP API 服务 |
| ORM | GORM | v1.25+ | 数据库访问 |
| AI 框架 | tRPC-Agent-Go | latest | AI Agent 核心 |
| 认证 | golang-jwt/jwt | v5 | JWT 生成验证 |
| 权限 | Casbin | v2 | RBAC 权限控制 |
| 缓存 | go-redis | v9 | Redis 客户端 |
| 消息队列 | amqp091-go | latest | RabbitMQ 客户端 |
| 任务调度 | robfig/cron | v3 | 定时任务 |
| 日志 | zap | v1.26+ | 高性能日志 |
| 配置 | viper | v1.18+ | 配置管理 |

### 4.2 go.mod 依赖

```go
module github.com/kevin-ai/go-kevin

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    gorm.io/gorm v1.25.5
    gorm.io/driver/mysql v1.5.2
    gorm.io/plugin/soft_delete v1.2.1
    trpc.group/trpc-go/trpc-agent-go v0.1.0
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/casbin/casbin/v2 v2.82.0
    github.com/redis/go-redis/v9 v9.4.0
    github.com/rabbitmq/amqp091-go v1.9.0
    github.com/robfig/cron/v3 v3.0.1
    go.uber.org/zap v1.26.0
    github.com/spf13/viper v1.18.2
    github.com/google/uuid v1.5.0
    golang.org/x/crypto v0.17.0
)
```

---

## 5. 领域模型设计

### 5.1 用户权限实体

```go
// User 用户聚合根
type User struct {
    ID         int64     `json:"id"`
    UserName   string    `json:"userName"`
    Password   string    `json:"-"`
    RealName   string    `json:"realName"`
    Email      string    `json:"email"`
    Phone      string    `json:"phone"`
    TenantID   int64     `json:"tenantId"`
    IsDelete   bool      `json:"isDelete"`
    CreateTime time.Time `json:"createTime"`
    UpdateTime time.Time `json:"updateTime"`
    Roles      []*Role   `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}

// Role 角色实体
type Role struct {
    ID          int64         `json:"id"`
    RoleName    string        `json:"roleName"`
    RoleCode    string        `json:"roleCode"`
    TenantID    int64         `json:"tenantId"`
    Permissions []*Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}

// Permission 权限实体
type Permission struct {
    ID       int64  `json:"id"`
    ParentID int64  `json:"parentId"`
    Name     string `json:"name"`
    Code     string `json:"code"`
    Type     int    `json:"type"` // 1:菜单 2:按钮 3:API
    URL      string `json:"url"`
}
```

### 5.2 AI 领域实体

```go
// AIApp AI应用聚合根
type AIApp struct {
    ID           int64     `json:"id"`
    AppName      string    `json:"appName"`
    AppDesc      string    `json:"appDesc"`
    ModelID      int64     `json:"modelId"`
    SystemPrompt string    `json:"systemPrompt"`
    TenantID     int64     `json:"tenantId"`
    Model        *AIModel  `json:"model,omitempty"`
    Skills       []*Skill  `json:"skills,omitempty" gorm:"many2many:app_skills;"`
}

// AIModel AI模型配置
type AIModel struct {
    ID        int64  `json:"id"`
    ModelName string `json:"modelName"`
    Provider  string `json:"provider"` // openai/claude/wenxin
    APIKey    string `json:"-"`
    BaseURL   string `json:"baseUrl"`
}

// ChatSession 聊天会话
type ChatSession struct {
    ID         int64     `json:"id"`
    UserID     int64     `json:"userId"`
    AppID      int64     `json:"appId"`
    Title      string    `json:"title"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
    ID        int64  `json:"id"`
    SessionID int64  `json:"sessionId"`
    Role      string `json:"role"` // user/assistant/system
    Content   string `json:"content"`
}

// KnowledgeBase 知识库
type KnowledgeBase struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    VectorModel string `json:"vectorModel"`
    TenantID    int64  `json:"tenantId"`
}
```

### 5.3 领域事件

```go
// DomainEvent 领域事件接口
type DomainEvent interface {
    EventName() string
    OccurredOn() time.Time
}

// UserCreatedEvent 用户创建事件
type UserCreatedEvent struct {
    UserID    int64
    UserName  string
    Timestamp time.Time
}

// ChatMessageSentEvent 聊天消息发送事件
type ChatMessageSentEvent struct {
    SessionID int64
    MessageID int64
    Content   string
    Timestamp time.Time
}
```

---

## 6. 应用层设计

### 6.1 CQRS 命令

```go
// CreateUserCommand 创建用户命令
type CreateUserCommand struct {
    UserName string `json:"userName" binding:"required"`
    Password string `json:"password" binding:"required"`
    RealName string `json:"realName"`
    Email    string `json:"email"`
    Phone    string `json:"phone"`
}

// CreateUserHandler 命令处理器
type CreateUserHandler struct {
    userRepo repository.UserRepository
    eventBus event.EventBus
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd *CreateUserCommand) (*CreateUserResult, error) {
    // 密码加密
    // 创建用户实体
    // 持久化
    // 发布领域事件
}
```

### 6.2 CQRS 查询

```go
// GetUserQuery 获取用户查询
type GetUserQuery struct {
    ID int64 `json:"id"`
}

// GetUserHandler 查询处理器
type GetUserHandler struct {
    userRepo repository.UserRepository
}

func (h *GetUserHandler) Handle(ctx context.Context, query *GetUserQuery) (*GetUserResult, error) {
    // 查询用户
    // 转换为 DTO
}
```

---

## 7. tRPC-Agent-Go 集成

### 7.1 AI Agent 核心

```go
// KevinAgent 自定义 AI Agent
type KevinAgent struct {
    agent.BaseAgent
    model  model.Model
    tools  []tool.Tool
    config *entity.AIApp
}

func (a *KevinAgent) Run(ctx context.Context, message string) (string, error) {
    // 构建消息
    // 调用模型
    // 处理工具调用
    // 返回结果
}
```

### 7.2 LLM 云服务适配

```go
// Provider LLM 提供商接口
type Provider interface {
    Chat(ctx context.Context, req *model.ChatRequest) (*model.ChatResponse, error)
    StreamChat(ctx context.Context, req *model.ChatRequest) (<-chan *model.ChatChunk, error)
}

// 支持的提供商
// - OpenAIProvider
// - ClaudeProvider
// - WenxinProvider (百度文心一言)
```

### 7.3 RAG 知识库集成

```go
// RAGService RAG 服务
type RAGService struct {
    vectorStore  VectorStore
    embedding    EmbeddingService
    chatProvider Provider
}

func (s *RAGService) Query(ctx context.Context, kbID int64, question string) (string, error) {
    // 搜索相关文档片段
    // 构建上下文
    // 调用 LLM 生成答案
}
```

---

## 8. HTTP 接口层

### 8.1 API 路由

```
/api/v1/
├── /auth/
│   ├── POST /login          # 登录
│   ├── POST /register       # 注册
│   └── POST /refresh        # 刷新 Token
├── /user/
│   ├── GET /                # 用户列表
│   ├── GET /:id             # 用户详情
│   ├── POST /               # 创建用户
│   ├── PUT /:id             # 更新用户
│   └── DELETE /:id          # 删除用户
├── /role/                   # 角色管理
├── /permission/             # 权限管理
├── /aiapps/                 # AI 应用管理
├── /aichat/                 # AI 聊天
│   ├── POST /sessions       # 创建会话
│   ├── POST /sessions/:id/messages  # 发送消息
│   └── POST /sessions/:id/stream    # 流式响应
├── /aikmss/                 # 知识库管理
│   ├── POST /:id/documents  # 上传文档
│   └── POST /:id/query      # 知识库问答
└── /file/                   # 文件管理
```

### 8.2 中间件

- **JWTAuth** - JWT 认证中间件
- **TenantContext** - 租户上下文中间件
- **RBAC** - 权限检查中间件
- **RateLimit** - 限流中间件

### 8.3 统一响应格式

```json
{
    "code": 200,
    "message": "操作成功",
    "data": {},
    "timestamp": 1234567890
}
```

---

## 9. 数据库迁移

### 9.1 GORM 模型映射

| 原表名 | GORM 模型 | 说明 |
|--------|-----------|------|
| T_User | UserGORM | 用户表 |
| T_Role | RoleGORM | 角色表 |
| T_Permission | PermissionGORM | 权限表 |
| T_AI_Apps | AIAppGORM | AI应用表 |
| T_AI_Chat_Sessions | ChatSessionGORM | 聊天会话表 |
| T_AI_Chat_Messages | ChatMessageGORM | 聊天消息表 |
| T_AI_Knowledge_Bases | KnowledgeBaseGORM | 知识库表 |
| T_AI_Documents | DocumentGORM | 文档表 |

### 9.2 迁移策略

1. 使用 GORM AutoMigrate 自动迁移表结构
2. 保留现有数据，不丢失
3. 软删除使用 `gorm.io/plugin/soft_delete`

---

## 10. 配置与部署

### 10.1 配置文件

```yaml
# configs/config.yaml
server:
  port: 9901
  mode: debug

database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: admin123
  dbname: kevin_app

redis:
  host: 127.0.0.1
  port: 6379
  password: "123456"

jwt:
  secret: your-secret-key
  expireHour: 24

ai:
  defaultProvider: openai
  providers:
    - name: openai
      apiKey: ${OPENAI_API_KEY}
      baseUrl: https://api.openai.com/v1
      model: gpt-4

rabbitmq:
  host: 127.0.0.1
  port: 5672
  user: guest
  password: guest
```

### 10.2 Docker 部署

```yaml
# docker-compose.yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "9901:9901"
    depends_on:
      - mysql
      - redis
  
  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"
  
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
  
  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"
```

---

## 11. 测试策略

### 11.1 测试层次

| 层次 | 覆盖率要求 | 工具 |
|------|-----------|------|
| 单元测试 | ≥ 80% | testify, mock |
| 集成测试 | ≥ 70% | httptest, testcontainers |
| API 测试 | ≥ 90% | httptest |

### 11.2 CI/CD

- GitHub Actions 自动运行测试
- 代码覆盖率报告
- 自动构建 Docker 镜像

---

## 附录

### A. 术语对照

| .NET 术语 | Go 术语 | 说明 |
|-----------|---------|------|
| Controller | Handler | HTTP 处理器 |
| Service | Service | 业务服务 |
| Repository | Repository | 仓储接口 |
| DbContext | *gorm.DB | 数据库上下文 |
| ILogger | *zap.Logger | 日志实例 |
| IConfiguration | *config.Config | 配置实例 |

### B. 参考资源

- [tRPC-Agent-Go 文档](https://github.com/trpc-group/trpc-agent-go)
- [GORM 文档](https://gorm.io/docs/)
- [Gin 文档](https://gin-gonic.com/docs/)

---

**文档维护者**: NetCoreKevin 开发团队  
**最后更新**: 2026/06/07
