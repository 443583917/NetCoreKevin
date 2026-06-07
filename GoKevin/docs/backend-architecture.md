# GoKevin 后端架构功能说明

**文档版本**: v1.0
**创建日期**: 2026/06/07
**最后更新**: 2026/06/07

---

## 目录

- [1. 架构概述](#1-架构概述)
- [2. 分层架构设计](#2-分层架构设计)
- [3. 领域层 (Domain Layer)](#3-领域层-domain-layer)
- [4. 应用层 (Application Layer)](#4-应用层-application-layer)
- [5. 基础设施层 (Infrastructure Layer)](#5-基础设施层-infrastructure-layer)
- [6. 接口层 (Interface Layer)](#6-接口层-interface-layer)
- [7. 公共工具包 (pkg)](#7-公共工具包-pkg)
- [8. 核心功能模块](#8-核心功能模块)
- [9. 数据库设计](#9-数据库设计)
- [10. API 接口设计](#10-api-接口设计)

---

## 1. 架构概述

GoKevin 后端采用 **DDD (领域驱动设计) + CQRS (命令查询职责分离)** 架构模式，结合 **事件驱动设计**，构建了一个高内聚、低耦合的企业级应用。

### 1.1 架构原则

| 原则 | 说明 |
|------|------|
| **依赖反转** | 领域层不依赖任何框架和外部服务 |
| **CQRS 分离** | 命令（写操作）和查询（读操作）职责分离 |
| **事件驱动** | 通过领域事件实现模块间解耦 |
| **模块化单体** | 清晰的模块边界，未来可拆分为微服务 |
| **接口隔离** | 通过接口定义依赖关系，便于测试和替换 |

### 1.2 技术栈

| 类别 | 技术 | 版本 | 说明 |
|------|------|------|------|
| 语言 | Go | 1.22+ | 主要开发语言 |
| Web 框架 | Gin | 1.12+ | HTTP 路由和中间件 |
| ORM | GORM | 1.31+ | 数据库访问层 |
| 日志 | Zap | 1.26+ | 高性能结构化日志 |
| 配置 | Viper | 1.18+ | 配置文件管理 |
| 认证 | JWT | - | Token 认证 |
| 缓存 | Redis | 7.0+ | 分布式缓存 |
| 消息队列 | RabbitMQ | 3.x | 异步消息处理 |

---

## 2. 分层架构设计

```
┌─────────────────────────────────────────────────────────────────┐
│                      接口层 (Interface Layer)                      │
│                    HTTP Handlers / Middleware / Router              │
├─────────────────────────────────────────────────────────────────┤
│                      应用层 (Application Layer)                    │
│                Command Handlers / Query Handlers / DTOs            │
├─────────────────────────────────────────────────────────────────┤
│                        领域层 (Domain Layer)                       │
│              Entities / Repository Interfaces / Domain Events      │
├─────────────────────────────────────────────────────────────────┤
│                    基础设施层 (Infrastructure Layer)                │
│      Persistence / Cache / Message Queue / AI / External Services  │
└─────────────────────────────────────────────────────────────────┘
```

### 2.1 目录结构

```
backend/
├── cmd/                            # 应用入口
│   └── server/
│       └── main.go                # 程序启动入口
│
├── internal/                       # 内部实现（不对外暴露）
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
│   │   ├── persistence/            # 数据持久化
│   │   │   ├── gorm/              # GORM 实现
│   │   │   │   ├── model/         # 数据库模型
│   │   │   │   └── repository/    # 仓储实现
│   │   │   └── migration/         # 数据库迁移
│   │   ├── ai/                    # AI 功能模块
│   │   │   ├── agent/             # 智能体
│   │   │   ├── llm/               # LLM 提供者
│   │   │   ├── rag/               # RAG 服务
│   │   │   ├── tool/              # 工具系统
│   │   │   ├── chat/              # 会话管理
│   │   │   ├── embedding/         # 向量嵌入
│   │   │   └── vectorstore/       # 向量存储
│   │   ├── cache/                 # 缓存模块
│   │   ├── mq/                    # 消息队列
│   │   ├── eventbus/              # 事件总线
│   │   ├── lock/                  # 分布式锁
│   │   ├── scheduler/             # 任务调度
│   │   ├── websocket/             # WebSocket
│   │   ├── sms/                   # 短信服务
│   │   ├── email/                 # 邮件服务
│   │   ├── storage/               # 文件存储
│   │   ├── snowflake/             # 雪花ID
│   │   ├── discovery/             # 服务发现
│   │   └── codegen/               # 代码生成
│   │
│   └── interfaces/                 # 接口层
│       └── http/
│           ├── handler/            # HTTP 处理器
│           ├── middleware/         # 中间件
│           ├── router/            # 路由定义
│           └── swagger/           # API 文档
│
├── pkg/                            # 公共工具包
│   ├── config/                     # 配置管理
│   ├── logger/                     # 日志模块
│   ├── auth/                       # 认证模块
│   ├── response/                   # 统一响应
│   └── utils/                      # 工具函数
│
├── configs/                        # 配置文件
├── docs/                           # 文档
├── scripts/                        # 脚本文件
└── test/                           # 测试文件
    ├── integration/                # 集成测试
    └── e2e/                        # 端到端测试
```

---

## 3. 领域层 (Domain Layer)

领域层是整个系统的核心，包含了业务逻辑和业务规则。该层不依赖任何外部框架或服务。

### 3.1 实体 (Entities)

#### 用户实体 (User)

```go
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
    Roles      []*Role   `json:"roles,omitempty" gorm:"many2many:t_user_bind_role;"`
}
```

#### 角色实体 (Role)

```go
type Role struct {
    ID          int64         `json:"id"`
    RoleName    string        `json:"roleName"`
    RoleCode    string        `json:"roleCode"`
    TenantID    int64         `json:"tenantId"`
    IsDelete    bool          `json:"isDelete"`
    CreateTime  time.Time     `json:"createTime"`
    Permissions []*Permission `json:"permissions,omitempty" gorm:"many2many:t_role_bind_permission;"`
}
```

#### AI 应用实体 (AIApp)

```go
type AIApp struct {
    ID           int64     `json:"id"`
    AppName      string    `json:"appName"`
    AppDesc      string    `json:"appDesc"`
    ModelID      int64     `json:"modelId"`
    SystemPrompt string    `json:"systemPrompt"`
    TenantID     int64     `json:"tenantId"`
    IsDelete     bool      `json:"isDelete"`
    CreateTime   time.Time `json:"createTime"`
    Model        *AIModel  `json:"model,omitempty"`
    Skills       []*Skill  `json:"skills,omitempty" gorm:"many2many:t_ai_apps_bind_skill;"`
}
```

#### 聊天会话实体 (ChatSession)

```go
type ChatSession struct {
    ID         int64     `json:"id"`
    UserID     int64     `json:"userId"`
    AppID      int64     `json:"appId"`
    Title      string    `json:"title"`
    IsDelete   bool      `json:"isDelete"`
    CreateTime time.Time `json:"createTime"`
    UpdateTime time.Time `json:"updateTime"`
}
```

#### 知识库实体 (KnowledgeBase)

```go
type KnowledgeBase struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    VectorModel string    `json:"vectorModel"`
    TenantID    int64     `json:"tenantId"`
    IsDelete    bool      `json:"isDelete"`
    CreateTime  time.Time `json:"createTime"`
}
```

### 3.2 仓储接口 (Repository Interfaces)

仓储接口定义了数据访问的契约，由基础设施层实现。

```go
// UserRepository 用户仓储接口
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id int64) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id int64) error
    GetByUserName(ctx context.Context, userName string) (*entity.User, error)
    GetByEmail(ctx context.Context, email string) (*entity.User, error)
    ListByTenantID(ctx context.Context, tenantID int64, page, pageSize int) ([]*entity.User, int64, error)
    GetWithRoles(ctx context.Context, id int64) (*entity.User, error)
    BindRoles(ctx context.Context, userID int64, roleIDs []int64) error
}
```

### 3.3 领域事件 (Domain Events)

领域事件用于实现模块间的解耦通信。

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
    Role      string
    Content   string
    Timestamp time.Time
}
```

---

## 4. 应用层 (Application Layer)

应用层负责协调领域对象完成业务用例，不包含业务逻辑。

### 4.1 命令处理器 (Command Handlers)

命令处理器处理写操作，遵循 CQRS 模式。

#### CreateUserHandler

```go
type CreateUserHandler struct {
    userRepo repository.UserRepository
    eventBus event.EventBus
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd *CreateUserCommand) (*CreateUserResult, error) {
    // 1. 验证用户名唯一性
    // 2. 密码加密
    // 3. 创建用户实体
    // 4. 持久化到数据库
    // 5. 发布领域事件
}
```

#### UpdateUserHandler

```go
type UpdateUserHandler struct {
    userRepo repository.UserRepository
}

func (h *UpdateUserHandler) Handle(ctx context.Context, cmd *UpdateUserCommand) error {
    // 1. 查询用户
    // 2. 更新用户信息
    // 3. 持久化到数据库
}
```

#### DeleteUserHandler

```go
type DeleteUserHandler struct {
    userRepo repository.UserRepository
}

func (h *DeleteUserHandler) Handle(ctx context.Context, cmd *DeleteUserCommand) error {
    // 1. 软删除用户
}
```

### 4.2 查询处理器 (Query Handlers)

查询处理器处理读操作，返回 DTO 对象。

#### GetUserHandler

```go
type GetUserHandler struct {
    userRepo repository.UserRepository
}

func (h *GetUserHandler) Handle(ctx context.Context, query *GetUserQuery) (*GetUserResult, error) {
    // 1. 查询用户及其角色
    // 2. 转换为 DTO
    // 3. 返回结果
}
```

#### ListUsersHandler

```go
type ListUsersHandler struct {
    userRepo repository.UserRepository
}

func (h *ListUsersHandler) Handle(ctx context.Context, query *ListUsersQuery) (*ListUsersResult, error) {
    // 1. 分页查询用户列表
    // 2. 转换为 DTO 列表
    // 3. 返回结果和总数
}
```

### 4.3 数据传输对象 (DTOs)

DTO 用于应用层和接口层之间的数据传输。

```go
// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
    UserName string `json:"userName" binding:"required,min=3,max=50"`
    Password string `json:"password" binding:"required,min=6"`
    RealName string `json:"realName"`
    Email    string `json:"email" binding:"omitempty,email"`
    Phone    string `json:"phone"`
}

// UserResponse 用户响应
type UserResponse struct {
    ID       int64          `json:"id"`
    UserName string         `json:"userName"`
    RealName string         `json:"realName"`
    Email    string         `json:"email"`
    Phone    string         `json:"phone"`
    TenantID int64          `json:"tenantId"`
    Roles    []RoleResponse `json:"roles,omitempty"`
}
```

---

## 5. 基础设施层 (Infrastructure Layer)

基础设施层实现了领域层定义的接口，提供技术能力。

### 5.1 数据持久化 (Persistence)

#### 数据库连接

```go
// database.go
func InitMySQL(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    // 配置连接池
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    return db, nil
}
```

#### GORM 模型

```go
// user.go
type UserGORM struct {
    ID         int64                 `gorm:"primaryKey;autoIncrement"`
    UserName   string                `gorm:"type:varchar(50);uniqueIndex;not null"`
    Password   string                `gorm:"type:varchar(100);not null"`
    RealName   string                `gorm:"type:varchar(50)"`
    Email      string                `gorm:"type:varchar(100)"`
    Phone      string                `gorm:"type:varchar(20)"`
    TenantID   int64                 `gorm:"index;not null"`
    IsDelete   soft_delete.DeletedAt `gorm:"softDelete:flag;index"`
    CreateTime time.Time             `gorm:"autoCreateTime"`
    UpdateTime time.Time             `gorm:"autoUpdateTime"`
}
```

#### 仓储实现

```go
// user_repo_impl.go
type UserRepositoryImpl struct {
    db *gorm.DB
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
    gormModel := r.toGORMModel(user)
    return r.db.WithContext(ctx).Create(gormModel).Error
}
```

### 5.2 AI 功能模块

#### LLM 提供者

```go
// provider.go
type Provider interface {
    Chat(ctx context.Context, messages []Message, opts ...Option) (*Response, error)
    ChatStream(ctx context.Context, messages []Message, opts ...Option) (<-chan *StreamChunk, error)
}

// openai.go
type OpenAIProvider struct {
    apiKey  string
    baseURL string
    model   string
}
```

#### 智能体 (Agent)

```go
// agent.go
type Agent struct {
    name         string
    description  string
    model        string
    systemPrompt string
    provider     llm.Provider
    tools        []tool.Tool
    maxTurns     int
}

func (a *Agent) Run(ctx context.Context, message string) (string, error) {
    // 1. 构建消息列表
    // 2. 调用 LLM
    // 3. 处理工具调用
    // 4. 返回结果
}
```

#### RAG 服务

```go
// rag.go
type RAGService interface {
    IndexDocument(ctx context.Context, doc *Document) error
    Query(ctx context.Context, query string, kbID int64, topK int) ([]*SearchResult, error)
}

// qdrant_rag.go
type QdrantRAGService struct {
    vectorStore vectorstore.VectorStore
    embedding   embedding.EmbeddingService
    llmProvider llm.Provider
}
```

### 5.3 缓存模块 (Cache)

```go
// cache.go
type Cache interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
}

// redis.go
type RedisCache struct {
    client *redis.Client
}
```

### 5.4 消息队列 (Message Queue)

```go
// mq.go
type MessageQueue interface {
    Publish(ctx context.Context, queue string, message interface{}) error
    Subscribe(ctx context.Context, queue string, handler func([]byte) error) error
    Close() error
}

// rabbitmq.go
type RabbitMQ struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}
```

### 5.5 事件总线 (Event Bus)

```go
// memory.go
type InMemoryEventBus struct {
    handlers map[string][]event.EventHandler
    mu       sync.RWMutex
}

func (b *InMemoryEventBus) Subscribe(eventName string, handler event.EventHandler) {
    b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func (b *InMemoryEventBus) Publish(event event.DomainEvent) {
    if handlers, ok := b.handlers[event.EventName()]; ok {
        for _, handler := range handlers {
            go handler(event) // 异步处理
        }
    }
}
```

### 5.6 分布式锁 (Distributed Lock)

```go
// lock.go
type Lock interface {
    Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error)
    Release(ctx context.Context, key string) error
}

// redis_lock.go
type RedisLock struct {
    client *redis.Client
}
```

### 5.7 任务调度 (Scheduler)

```go
// scheduler.go
type Scheduler struct {
    cron *cron.Cron
}

func (s *Scheduler) AddJob(spec string, cmd cron.Cmd) error {
    return s.cron.AddJob(spec, cmd)
}
```

### 5.8 WebSocket

```go
// hub.go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

// websocket.go
type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
}
```

### 5.9 短信服务 (SMS)

```go
// sms.go
type SMSProvider interface {
    Send(ctx context.Context, phone string, template string, params map[string]string) error
}

// aliyun.go
type AliyunSMS struct {
    accessKeyId  string
    accessKeySecret string
    signName     string
}

// tencent.go
type TencentSMS struct {
    appId     string
    secretId  string
    secretKey string
}
```

### 5.10 邮件服务 (Email)

```go
// email.go
type EmailService interface {
    Send(ctx context.Context, to []string, subject string, body string) error
}

// smtp.go
type SMTPService struct {
    host     string
    port     int
    username string
    password string
}
```

### 5.11 文件存储 (Storage)

```go
// storage.go
type Storage interface {
    Upload(ctx context.Context, file *File) (string, error)
    Download(ctx context.Context, fileID string) ([]byte, error)
    Delete(ctx context.Context, fileID string) error
}

// local.go
type LocalStorage struct {
    basePath string
}
```

### 5.12 雪花ID (Snowflake)

```go
// snowflake.go
type Snowflake struct {
    node *sf.Node
}

func (s *Snowflake) Generate() int64 {
    return s.node.Generate().Int64()
}
```

### 5.13 服务发现 (Service Discovery)

```go
// discovery.go
type ServiceDiscovery interface {
    Register(ctx context.Context, service *Service) error
    Deregister(ctx context.Context, serviceID string) error
    Discover(ctx context.Context, serviceName string) ([]*Service, error)
}

// consul.go
type ConsulDiscovery struct {
    client *consul.Client
}
```

### 5.14 代码生成器 (Code Generator)

```go
// codegen.go
type CodeGenerator struct {
    templates map[string]*template.Template
}

func (g *CodeGenerator) Generate(tplName string, data interface{}) ([]byte, error) {
    tpl, ok := g.templates[tplName]
    if !ok {
        return nil, fmt.Errorf("template %s not found", tplName)
    }
    var buf bytes.Buffer
    err := tpl.Execute(&buf, data)
    return buf.Bytes(), err
}
```

---

## 6. 接口层 (Interface Layer)

接口层负责处理外部请求，将请求转换为应用层的命令或查询。

### 6.1 HTTP 处理器 (Handlers)

#### AuthHandler

```go
type AuthHandler struct {
    userRepo   repository.UserRepository
    jwtSecret  string
    expireHour int
}

func (h *AuthHandler) Login(c *gin.Context) {
    // 1. 解析请求参数
    // 2. 验证用户名密码
    // 3. 生成 JWT Token
    // 4. 返回 Token
}

func (h *AuthHandler) Register(c *gin.Context) {
    // 1. 解析请求参数
    // 2. 调用 CreateUserHandler
    // 3. 返回创建结果
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
    // 1. 获取当前 Token
    // 2. 验证并刷新 Token
    // 3. 返回新 Token
}
```

#### UserHandler

```go
type UserHandler struct {
    createUser *command.CreateUserHandler
    getUser    *query.GetUserHandler
    listUsers  *query.ListUsersHandler
    updateUser *command.UpdateUserHandler
    deleteUser *command.DeleteUserHandler
}

func (h *UserHandler) Create(c *gin.Context) {
    // 1. 解析请求参数
    // 2. 调用 CreateUserHandler
    // 3. 返回创建结果
}

func (h *UserHandler) GetByID(c *gin.Context) {
    // 1. 解析用户 ID
    // 2. 调用 GetUserHandler
    // 3. 返回用户信息
}

func (h *UserHandler) List(c *gin.Context) {
    // 1. 解析分页参数
    // 2. 调用 ListUsersHandler
    // 3. 返回用户列表
}

func (h *UserHandler) Update(c *gin.Context) {
    // 1. 解析用户 ID 和更新参数
    // 2. 调用 UpdateUserHandler
    // 3. 返回更新结果
}

func (h *UserHandler) Delete(c *gin.Context) {
    // 1. 解析用户 ID
    // 2. 调用 DeleteUserHandler
    // 3. 返回删除结果
}
```

### 6.2 中间件 (Middleware)

#### JWT 认证中间件

```go
func JWTAuth(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从 Header 获取 Token
        // 2. 解析 Bearer Token
        // 3. 验证 Token 有效性
        // 4. 将用户信息存入上下文
        // 5. 继续处理请求
    }
}
```

#### 租户上下文中间件

```go
func TenantContext() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从上下文获取租户 ID
        // 2. 验证租户信息
        // 3. 继续处理请求
    }
}
```

#### CORS 跨域中间件

```go
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 设置跨域响应头
        // 2. 处理 OPTIONS 预检请求
        // 3. 继续处理请求
    }
}
```

### 6.3 路由定义 (Router)

```go
func RegisterRoutes(r *gin.Engine, userHandler *handler.UserHandler, authHandler *handler.AuthHandler) {
    v1 := r.Group("/api/v1")

    // 公开接口
    public := v1.Group("")
    {
        public.POST("/auth/login", authHandler.Login)
        public.POST("/auth/register", authHandler.Register)
        public.POST("/auth/refresh", authHandler.RefreshToken)
    }

    // 需要认证的接口
    protected := v1.Group("")
    protected.Use(middleware.JWTAuth("your-secret-key"))
    protected.Use(middleware.TenantContext())
    {
        // 用户管理
        users := protected.Group("/user")
        {
            users.GET("", userHandler.List)
            users.GET("/:id", userHandler.GetByID)
            users.POST("", userHandler.Create)
            users.PUT("/:id", userHandler.Update)
            users.DELETE("/:id", userHandler.Delete)
        }

        // AI 应用管理
        apps := protected.Group("/aiapps")
        {
            apps.GET("", aiHandler.ListApps)
            apps.GET("/:id", aiHandler.GetApp)
            apps.POST("", aiHandler.CreateApp)
            apps.PUT("/:id", aiHandler.UpdateApp)
            apps.DELETE("/:id", aiHandler.DeleteApp)
        }

        // AI 聊天
        chat := protected.Group("/aichat")
        {
            chat.POST("/sessions", chatHandler.CreateSession)
            chat.GET("/sessions", chatHandler.ListSessions)
            chat.GET("/sessions/:id", chatHandler.GetSession)
            chat.POST("/sessions/:id/messages", chatHandler.SendMessage)
            chat.GET("/sessions/:id/messages", chatHandler.GetMessages)
        }

        // 知识库管理
        kb := protected.Group("/aikmss")
        {
            kb.GET("", kbHandler.List)
            kb.POST("", kbHandler.Create)
            kb.GET("/:id", kbHandler.Get)
            kb.DELETE("/:id", kbHandler.Delete)
            kb.POST("/:id/documents", kbHandler.UploadDocument)
            kb.POST("/:id/query", kbHandler.Query)
        }
    }
}
```

---

## 7. 公共工具包 (pkg)

公共工具包提供了跨模块使用的通用功能。

### 7.1 配置管理 (config)

```go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    AI       AIConfig       `mapstructure:"ai"`
    RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
    Log      LogConfig      `mapstructure:"log"`
}

func Init(configPath string) error {
    viper.SetConfigFile(configPath)
    viper.AutomaticEnv()
    if err := viper.ReadInConfig(); err != nil {
        return err
    }
    cfg = &Config{}
    return viper.Unmarshal(cfg)
}
```

### 7.2 日志模块 (logger)

```go
func Init() {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.TimeKey = "timestamp"
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    writerSyncer := zapcore.AddSync(os.Stdout)
    level := zapcore.InfoLevel

    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(encoderConfig),
        writerSyncer,
        level,
    )

    log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func Info(msg string, fields ...zap.Field) {
    log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
    log.Error(msg, fields...)
}
```

### 7.3 认证模块 (auth)

```go
type JWTClaims struct {
    UserID   int64  `json:"userId"`
    UserName string `json:"userName"`
    TenantID int64  `json:"tenantId"`
    jwt.RegisteredClaims
}

func GenerateToken(secret string, expireHour int, userID int64, userName string, tenantID int64) (string, error) {
    claims := &JWTClaims{
        UserID:   userID,
        UserName: userName,
        TenantID: tenantID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHour) * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "go-kevin",
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func ParseToken(secret string, tokenString string) (*JWTClaims, error) {
    // 解析并验证 Token
}
```

### 7.4 统一响应 (response)

```go
type Response struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    Timestamp int64       `json:"timestamp"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:      200,
        Message:   "操作成功",
        Data:      data,
        Timestamp: time.Now().Unix(),
    })
}

func Error(c *gin.Context, code int, message string) {
    c.JSON(code, Response{
        Code:      code,
        Message:   message,
        Timestamp: time.Now().Unix(),
    })
}

func Page(c *gin.Context, list interface{}, total int64, page, pageSize int) {
    Success(c, PageResult{
        List:     list,
        Total:    total,
        Page:     page,
        PageSize: pageSize,
    })
}
```

---

## 8. 核心功能模块

### 8.1 用户权限管理

| 功能 | 说明 | 状态 |
|------|------|------|
| 用户注册 | 用户名、密码、邮箱注册 | ✅ |
| 用户登录 | JWT Token 认证 | ✅ |
| 用户管理 | CRUD 操作 | ✅ |
| 角色管理 | 角色的增删改查 | ✅ |
| 权限管理 | 权限的增删改查 | ✅ |
| 用户角色绑定 | 用户与角色关联 | ✅ |
| 角色权限绑定 | 角色与权限关联 | ✅ |
| RBAC 权限控制 | 基于角色的访问控制 | ✅ |
| 多租户支持 | 租户数据隔离 | ✅ |

### 8.2 AI 智能体功能

| 功能 | 说明 | 状态 |
|------|------|------|
| LLM 集成 | 支持 OpenAI 等多种大语言模型 | ✅ |
| 智能体核心 | Agent 核心，支持工具调用 | ✅ |
| 工具系统 | 天气、计算、搜索等工具 | ✅ |
| RAG 服务 | 知识库检索增强生成 | ✅ |
| 会话管理 | 多轮对话和上下文管理 | ✅ |
| 向量存储 | Qdrant 向量数据库集成 | ✅ |
| 向量嵌入 | OpenAI Embedding 集成 | ✅ |

### 8.3 基础设施功能

| 功能 | 说明 | 状态 |
|------|------|------|
| Redis 缓存 | 分布式缓存支持 | ✅ |
| 消息队列 | RabbitMQ 异步消息处理 | ✅ |
| 事件总线 | 内存事件总线 | ✅ |
| 分布式锁 | Redis 分布式锁 | ✅ |
| 任务调度 | Cron 定时任务 | ✅ |
| WebSocket | 实时通信支持 | ✅ |
| 短信服务 | 阿里云、腾讯云短信 | ✅ |
| 邮件服务 | SMTP 邮件发送 | ✅ |
| 文件存储 | 本地文件存储 | ✅ |
| 雪花ID | 分布式唯一ID生成 | ✅ |
| 服务发现 | Consul 服务注册与发现 | ✅ |
| 代码生成器 | CRUD 代码自动生成 | ✅ |

### 8.4 API 文档

| 功能 | 说明 | 状态 |
|------|------|------|
| Swagger UI | 自动生成 API 文档 | ✅ |
| API 版本管理 | v1 版本管理 | ✅ |
| 接口测试 | 在线接口测试 | ✅ |

---

## 9. 数据库设计

### 9.1 用户权限模块

#### t_user 用户表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| user_name | VARCHAR(50) | 用户名，唯一索引 |
| password | VARCHAR(100) | 密码（bcrypt 加密） |
| real_name | VARCHAR(50) | 真实姓名 |
| email | VARCHAR(100) | 邮箱 |
| phone | VARCHAR(20) | 手机号 |
| tenant_id | BIGINT | 租户ID，索引 |
| is_delete | TINYINT | 软删除标记 |
| create_time | DATETIME | 创建时间 |
| update_time | DATETIME | 更新时间 |

#### t_role 角色表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| role_name | VARCHAR(50) | 角色名称 |
| role_code | VARCHAR(50) | 角色编码，唯一索引 |
| tenant_id | BIGINT | 租户ID，索引 |
| is_delete | TINYINT | 软删除标记 |
| create_time | DATETIME | 创建时间 |

#### t_permission 权限表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| parent_id | BIGINT | 父权限ID，索引 |
| name | VARCHAR(50) | 权限名称 |
| code | VARCHAR(100) | 权限编码，唯一索引 |
| type | INT | 类型（1:菜单 2:按钮 3:API） |
| url | VARCHAR(200) | 资源URL |
| icon | VARCHAR(50) | 图标 |
| sort | INT | 排序号 |

#### t_user_bind_role 用户角色关联表

| 字段 | 类型 | 说明 |
|------|------|------|
| user_id | BIGINT | 用户ID |
| role_id | BIGINT | 角色ID |

#### t_role_bind_permission 角色权限关联表

| 字段 | 类型 | 说明 |
|------|------|------|
| role_id | BIGINT | 角色ID |
| permission_id | BIGINT | 权限ID |

### 9.2 AI 模块

#### t_ai_apps AI应用表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| app_name | VARCHAR(100) | 应用名称 |
| app_desc | VARCHAR(500) | 应用描述 |
| model_id | BIGINT | 模型ID，索引 |
| system_prompt | TEXT | 系统提示词 |
| tenant_id | BIGINT | 租户ID，索引 |
| is_delete | TINYINT | 软删除标记 |
| create_time | DATETIME | 创建时间 |

#### t_ai_models AI模型配置表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| model_name | VARCHAR(100) | 模型名称 |
| provider | VARCHAR(50) | 提供者（openai/claude等） |
| api_key | VARCHAR(200) | API密钥 |
| base_url | VARCHAR(200) | API基础URL |
| max_tokens | INT | 最大Token数 |
| is_delete | TINYINT | 软删除标记 |
| create_time | DATETIME | 创建时间 |

#### t_ai_skills 技能表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| skill_name | VARCHAR(100) | 技能名称 |
| skill_code | VARCHAR(50) | 技能编码，唯一索引 |
| description | VARCHAR(500) | 技能描述 |
| skill_type | VARCHAR(20) | 类型（tool/workflow/script） |
| config | JSON | 配置信息 |
| is_delete | TINYINT | 软删除标记 |
| create_time | DATETIME | 创建时间 |

#### t_ai_apps_bind_skill 应用技能关联表

| 字段 | 类型 | 说明 |
|------|------|------|
| app_id | BIGINT | 应用ID |
| skill_id | BIGINT | 技能ID |

#### t_ai_chat_sessions 聊天会话表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| user_id | BIGINT | 用户ID，索引 |
| app_id | BIGINT | 应用ID，索引 |
| title | VARCHAR(200) | 会话标题 |
| is_delete | TINYINT | 软删除标记 |
| create_time | DATETIME | 创建时间 |
| update_time | DATETIME | 更新时间 |

#### t_ai_chat_messages 聊天消息表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| session_id | BIGINT | 会话ID，索引 |
| role | VARCHAR(20) | 角色（user/assistant/system/tool） |
| content | TEXT | 消息内容 |
| tokens | INT | Token数量 |
| create_time | DATETIME | 创建时间 |

#### t_ai_knowledge_bases 知识库表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| name | VARCHAR(100) | 知识库名称 |
| description | VARCHAR(500) | 知识库描述 |
| vector_model | VARCHAR(50) | 向量模型 |
| tenant_id | BIGINT | 租户ID，索引 |
| is_delete | TINYINT | 软删除标记 |
| create_time | DATETIME | 创建时间 |

#### t_ai_documents 文档表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| kb_id | BIGINT | 知识库ID，索引 |
| file_name | VARCHAR(200) | 文件名 |
| file_url | VARCHAR(500) | 文件URL |
| status | VARCHAR(20) | 状态（pending/processing/completed/failed） |
| chunk_count | INT | 分块数量 |
| is_delete | TINYINT | 软删除标记 |
| create_time | DATETIME | 创建时间 |

---

## 10. API 接口设计

### 10.1 认证接口

#### POST /api/v1/auth/login

用户登录

**请求参数：**
```json
{
    "userName": "admin",
    "password": "123456"
}
```

**响应：**
```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expiresIn": 86400
    },
    "timestamp": 1686123456
}
```

#### POST /api/v1/auth/register

用户注册

**请求参数：**
```json
{
    "userName": "newuser",
    "password": "password123",
    "realName": "新用户",
    "email": "newuser@example.com",
    "phone": "13800138000"
}
```

**响应：**
```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        "id": 1,
        "userName": "newuser"
    },
    "timestamp": 1686123456
}
```

#### POST /api/v1/auth/refresh

刷新 Token

**请求头：**
```
Authorization: Bearer <token>
```

**响应：**
```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expiresIn": 86400
    },
    "timestamp": 1686123456
}
```

### 10.2 用户管理接口

#### GET /api/v1/user

获取用户列表

**请求参数：**
- `page`: 页码（默认 1）
- `pageSize`: 每页数量（默认 10）

**响应：**
```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        "list": [
            {
                "id": 1,
                "userName": "admin",
                "realName": "管理员",
                "email": "admin@example.com",
                "phone": "13800138000",
                "tenantId": 1000,
                "roles": [
                    {
                        "id": 1,
                        "roleName": "管理员",
                        "roleCode": "admin"
                    }
                ]
            }
        ],
        "total": 100,
        "page": 1,
        "pageSize": 10
    },
    "timestamp": 1686123456
}
```

#### GET /api/v1/user/:id

获取用户详情

**响应：**
```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        "id": 1,
        "userName": "admin",
        "realName": "管理员",
        "email": "admin@example.com",
        "phone": "13800138000",
        "tenantId": 1000,
        "roles": [
            {
                "id": 1,
                "roleName": "管理员",
                "roleCode": "admin"
            }
        ]
    },
    "timestamp": 1686123456
}
```

#### POST /api/v1/user

创建用户

**请求参数：**
```json
{
    "userName": "newuser",
    "password": "password123",
    "realName": "新用户",
    "email": "newuser@example.com",
    "phone": "13800138000"
}
```

**响应：**
```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        "id": 1,
        "userName": "newuser"
    },
    "timestamp": 1686123456
}
```

#### PUT /api/v1/user/:id

更新用户

**请求参数：**
```json
{
    "realName": "更新后的姓名",
    "email": "updated@example.com",
    "phone": "13900139000"
}
```

**响应：**
```json
{
    "code": 200,
    "message": "操作成功",
    "data": null,
    "timestamp": 1686123456
}
```

#### DELETE /api/v1/user/:id

删除用户

**响应：**
```json
{
    "code": 200,
    "message": "操作成功",
    "data": null,
    "timestamp": 1686123456
}
```

### 10.3 AI 应用管理接口

#### GET /api/v1/aiapps

获取 AI 应用列表

**请求参数：**
- `page`: 页码
- `pageSize`: 每页数量

**响应：**
```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        "list": [
            {
                "id": 1,
                "appName": "智能客服",
                "appDesc": "24小时在线客服机器人",
                "modelId": 1,
                "systemPrompt": "你是一个专业的客服助手",
                "tenantId": 1000,
                "model": {
                    "id": 1,
                    "modelName": "gpt-4",
                    "provider": "openai"
                },
                "skills": [
                    {
                        "id": 1,
                        "skillName": "天气查询",
                        "skillCode": "weather",
                        "skillType": "tool"
                    }
                ]
            }
        ],
        "total": 10,
        "page": 1,
        "pageSize": 10
    },
    "timestamp": 1686123456
}
```

#### POST /api/v1/aiapps

创建 AI 应用

**请求参数：**
```json
{
    "appName": "智能客服",
    "appDesc": "24小时在线客服机器人",
    "modelId": 1,
    "systemPrompt": "你是一个专业的客服助手"
}
```

### 10.4 AI 聊天接口

#### POST /api/v1/aichat/sessions

创建聊天会话

**请求参数：**
```json
{
    "appId": 1,
    "title": "新对话"
}
```

**响应：**
```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        "id": 1,
        "userId": 100,
        "appId": 1,
        "title": "新对话",
        "createTime": "2026-06-07T10:00:00Z"
    },
    "timestamp": 1686123456
}
```

#### POST /api/v1/aichat/sessions/:id/messages

发送聊天消息

**请求参数：**
```json
{
    "content": "你好，请问有什么可以帮助您的？"
}
```

**响应：**
```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        "id": 1,
        "sessionId": 1,
        "role": "assistant",
        "content": "您好！我是智能客服助手，很高兴为您服务。请问有什么可以帮助您的？",
        "tokens": 50,
        "createTime": "2026-06-07T10:00:01Z"
    },
    "timestamp": 1686123456
}
```

### 10.5 知识库管理接口

#### GET /api/v1/aikmss

获取知识库列表

**响应：**
```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        "list": [
            {
                "id": 1,
                "name": "产品知识库",
                "description": "包含所有产品文档",
                "vectorModel": "text-embedding-ada-002",
                "tenantId": 1000
            }
        ],
        "total": 5,
        "page": 1,
        "pageSize": 10
    },
    "timestamp": 1686123456
}
```

#### POST /api/v1/aikmss/:id/query

查询知识库

**请求参数：**
```json
{
    "query": "如何使用产品功能？",
    "topK": 5
}
```

**响应：**
```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        "results": [
            {
                "content": "产品功能使用说明...",
                "score": 0.95,
                "documentId": 1,
                "chunkIndex": 5
            }
        ]
    },
    "timestamp": 1686123456
}
```

---

## 附录

### A. 依赖注入流程

```go
// main.go 中的依赖注入流程
func main() {
    // 1. 加载配置
    config.Init(configPath)
    cfg := config.Get()

    // 2. 初始化日志
    logger.Init()

    // 3. 初始化数据库
    db, _ := persistence.InitMySQL(dsn)

    // 4. 执行数据库迁移
    migrator := migration.NewMigrator(db)
    migrator.Run()

    // 5. 初始化事件总线
    bus := eventbus.NewInMemoryEventBus()

    // 6. 初始化仓储
    userRepo := repository.NewUserRepositoryImpl(db)

    // 7. 初始化命令和查询处理器
    createUserHandler := command.NewCreateUserHandler(userRepo, bus)
    getUserHandler := query.NewGetUserHandler(userRepo)

    // 8. 初始化 HTTP 处理器
    authHandler := handler.NewAuthHandler(userRepo, cfg.JWT.Secret, cfg.JWT.ExpireHour)
    userHandler := handler.NewUserHandler(createUserHandler, getUserHandler, ...)

    // 9. 注册路由
    router.RegisterRoutes(r, userHandler, authHandler)

    // 10. 启动服务器
    srv.ListenAndServe()
}
```

### B. 测试策略

| 测试类型 | 工具 | 覆盖范围 | 状态 |
|----------|------|----------|------|
| 单元测试 | testify | 业务逻辑 | 部分完成 |
| 集成测试 | TestContainers | 数据库/缓存 | 基础完成 |
| API 测试 | httptest | HTTP 接口 | 部分完成 |
| 性能测试 | k6/vegeta | 关键接口 | 未开始 |
| 安全测试 | OWASP ZAP | 安全漏洞 | 未开始 |

### C. 错误码定义

| 错误码 | 说明 |
|--------|------|
| 200 | 操作成功 |
| 400 | 参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

**文档版本**: v1.0
**最后更新**: 2026/06/07
**维护者**: NetCoreKevin 开发团队
