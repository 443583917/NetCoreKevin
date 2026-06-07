# NetCoreKevin 项目架构功能拆分文档

> **⚠️ 文档已迁移**
> 本文档内容已迁移到 [docs/architecture.md](docs/architecture.md)。
> 请参考新文档获取最新信息。

---

## 一、项目总体架构

### 1.1 架构概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              前端层 (Vue3 + AntDesign)                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                              API 网关层 (WebApi)                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                         Kevin 核心框架模块                             │  │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────────┐  │  │
│  │  │ Application │ │   Domain    │ │ EF Core     │ │  kevin.Module │  │  │
│  │  │   应用服务   │ │   领域层    │ │  数据访问    │ │   功能模块    │  │  │
│  │  └──────┬──────┘ └──────┬──────┘ └──────┬──────┘ └───────┬───────┘  │  │
│  └─────────┼───────────────┼───────────────┼────────────────┼──────────┘  │
│            │               │               │                │             │
│  ┌─────────▼───────────────▼───────────────▼────────────────▼──────────┐  │
│  │                         App 业务应用模块                             │  │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────────┐ │  │
│  │  │ Application │ │   Domain    │ │ Repositories│ │    WebApi     │ │  │
│  │  │   应用服务   │ │   领域层    │ │   仓储实现   │ │   API入口     │ │  │
│  │  └─────────────┘ └─────────────┘ └─────────────┘ └───────────────┘ │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                              数据存储层                                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  MySQL   │  │  Redis   │  │  Qdrant  │  │ RabbitMQ │  │ 文件存储  │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 项目目录结构

```
NetCoreKevin/
├── App/                            # 业务应用模块
│   ├── Application/                # 应用服务层
│   │   └── Services/
│   │       └── v1/                 # V1版本业务服务
│   ├── Domain/                     # 领域层
│   │   ├── Entities/               # 实体定义
│   │   └── Interfaces/             # 接口定义
│   │       ├── Repositorie/        # 仓储接口
│   │       └── Services/           # 服务接口
│   ├── RepositorieRps/             # 仓储实现层
│   │   └── Repositories/
│   ├── AppShare/                   # 共享组件
│   └── WebApi/                     # API入口
│       ├── Controllers/            # 控制器
│       └── Program.cs              # 程序入口
├── Kevin/                          # 核心框架模块
│   ├── Application/                # 核心应用服务
│   │   └── Services/
│   │       └── AI/                 # AI相关服务
│   ├── Domain/                     # 核心领域模型
│   │   ├── Attributes/             # 自定义属性
│   │   ├── Auth/                   # 认证授权
│   │   ├── BaseDatas/              # 基础数据
│   │   ├── Bases/                  # 基类定义
│   │   ├── Entities/               # 核心实体
│   │   ├── EventHandlers/          # 事件处理器
│   │   ├── Events/                 # 领域事件
│   │   └── Interfaces/             # 核心接口
│   ├── Kevin.EntityFrameworkCore/  # EF Core实现
│   │   ├── Configuration/          # 实体配置
│   │   ├── Database/               # 数据库上下文
│   │   └── Interceptors/           # 拦截器
│   ├── Kevin.Web.Basics/           # Web基础组件
│   └── kevin.Module/               # 功能模块集合
│       ├── kevin.AI.AgentFramework/# AI智能体框架
│       ├── Kevin.Authentication.Jwt/# JWT认证
│       ├── kevin.Cache/            # 缓存模块
│       ├── kevin.Cap/              # CAP消息队列
│       ├── kevin.CodeGenerator/    # 代码生成器
│       ├── Kevin.Common/           # 通用工具
│       ├── kevin.Consul/           # 服务发现
│       ├── kevin.FileStorage/      # 文件存储
│       ├── Kevin.Hangfire/         # 任务调度
│       ├── Kevin.HttpApiClients/   # HTTP客户端
│       ├── kevin.Ioc/              # 依赖注入
│       ├── Kevin.log4Net/          # 日志模块
│       ├── Kevin.MCP.Server/       # MCP服务
│       ├── kevin.Permission/       # 权限模块
│       ├── kevin.RabbitMQ/         # RabbitMQ
│       ├── Kevin.RAG/              # RAG检索
│       ├── Kevin.SignalR/          # 实时通信
│       ├── Kevin.SMS/              # 短信服务
│       ├── Kevin.SnowflakeId/      # 雪花ID
│       └── Kevin.Versioning.Swagger/# API版本管理
├── Doc/                            # 文档资源
├── InitData/                       # 初始化数据
├── Test/                           # 测试项目
└── vue/                            # 前端项目
```

---

## 二、核心功能模块拆分

### 2.1 AI 智能体模块 (kevin.AI.AgentFramework)

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **智能体服务** | 基于SemanticKernel的AI代理系统 | `AIAgentService` |
| **技能管理** | 动态加载和管理AI技能 | `IAIAgentToolSkillService` |
| **工具集成** | 集成各类外部工具 | `CommonToolsService` |
| **工作流** | 多步骤任务编排 | `WorkFlowsAndAIAgentsDemo` |
| **脚本执行** | Python脚本运行支持 | `PySubprocessScriptRunner` |

**目录结构:**
```
kevin.AI.AgentFramework/
├── Agent/                  # 智能体核心
├── Const/                  # 常量定义
│   └── SystemPrompt.cs     # 系统提示词
├── Dto/                    # 数据传输对象
├── Interfaces/             # 接口定义
├── ScriptRunners/          # 脚本执行器
├── SkillClass/             # 技能类
│   ├── GetWeatherSkill.cs
│   └── UnitConverterSkill.cs
├── Skills/                 # 技能集合
├── Tools/                  # 工具服务
│   ├── AgentHttpClientToolsService.cs
│   ├── CommonToolsService.cs
│   ├── PythonToolsService.cs
│   └── ShellToolsService.cs
└── WorkFlows/              # 工作流定义
```

---

### 2.2 AI 知识库模块 (Kevin.RAG)

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **向量化服务** | 文档向量化处理 | `IOllamaService` |
| **向量存储** | Qdrant向量数据库操作 | `IQdrantService` |
| **RAG检索** | 检索增强生成 | `IRAGService` |
| **知识库管理** | 知识库CRUD | `AIKmssService` |

**目录结构:**
```
Kevin.RAG/
├── Dto/                    # 数据传输对象
├── Interfaces/             # 接口定义
│   ├── IOllamaService.cs
│   ├── IQdrantService.cs
│   └── IRAGService.cs
├── Ollama/                 # Ollama集成
│   └── OllamaService.cs
└── Qdrant/                 # Qdrant集成
    └── QdrantService.cs
```

---

### 2.3 用户权限模块 (kevin.Permission)

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **用户管理** | 用户CRUD、角色绑定 | `UserService` |
| **角色管理** | 角色CRUD、权限配置 | `RoleService` |
| **权限管理** | 菜单权限、API权限 | `PermissionService` |
| **组织架构** | 部门、职位管理 | `DepartmentService` |

**核心实体:**
- `TUser` - 用户表
- `TRole` - 角色表
- `TPermission` - 权限表
- `TDepartment` - 部门表
- `TPosition` - 职位表
- `TTenant` - 租户表

---

### 2.4 任务调度模块 (Kevin.Hangfire)

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **定时任务** | Cron表达式任务调度 | `IModuleConfigTasks` |
| **任务监控** | 任务执行状态监控 | Hangfire Dashboard |
| **自动注册** | 启动时自动注册任务 | `ServiceCollectionExtensions` |

**配置示例:**
```csharp
public class AIKmssModuleConfigTasks : IModuleConfigTasks
{
    public Task<bool> ConfigTasks(IRecurringJobManager recurringJobManager)
    {
        recurringJobManager.AddOrUpdate<IAIKmssService>(
            recurringJobId: "AI文档知识库处理任务",
            (s) => s.ProcessKmssVectorData(default),
            "0 0/6 0/1 * * ?",
            new RecurringJobOptions { TimeZone = TimeZoneInfo.Local }
        );
        return Task.FromResult(true);
    }
}
```

---

### 2.5 消息队列模块 (kevin.Cap)

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **事件发布** | 分布式事件发布 | CAP Publisher |
| **事件订阅** | 事件订阅处理 | `DemoSubscribe` |
| **消息过滤** | 消息过滤器 | `CapSubscribeFilter` |

---

### 2.6 文件存储模块 (kevin.FileStorage)

| 功能 | 说明 | 支持平台 |
|------|------|----------|
| **文件上传** | 多文件上传 | - |
| **文件下载** | 文件下载服务 | - |
| **云存储** | 多云存储支持 | 腾讯云COS、阿里云OSS、七牛云 |

**目录结构:**
```
kevin.FileStorage/
├── AliCloud/               # 阿里云OSS
├── KevinStaticFiles/       # 本地文件存储
├── QiniuCloud/             # 七牛云
└── TencentCloud/           # 腾讯云COS
```

---

### 2.7 认证授权模块 (Kevin.Authentication.Jwt)

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **JWT生成** | Token生成服务 | `TokenService` |
| **Token刷新** | Token刷新机制 | `RefreshTokenDto` |
| **Claims管理** | 用户声明管理 | `JwtKeinClaimTypes` |

---

### 2.8 实时通信模块 (Kevin.SignalR)

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **消息推送** | 实时消息推送 | `SignalRService` |
| **连接管理** | 客户端连接管理 | Hub |
| **群组管理** | 消息群组管理 | - |

---

### 2.9 缓存模块 (kevin.Cache)

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **Redis缓存** | 分布式缓存服务 | `CacheService` |
| **缓存策略** | 多种缓存策略支持 | `ICacheService` |

---

### 2.10 日志模块 (Kevin.log4Net)

| 功能 | 说明 | 说明 |
|------|------|------|
| **日志记录** | 多级别日志记录 | Info/Warn/Error/Fatal |
| **日志配置** | 灵活的日志配置 | XML配置文件 |

---

### 2.11 服务发现模块 (kevin.Consul)

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **服务注册** | 服务自动注册 | Consul Client |
| **服务发现** | 服务地址发现 | - |
| **健康检查** | 服务健康检查 | - |

---

### 2.12 API版本管理 (Kevin.Versioning.Swagger)

| 功能 | 说明 | 说明 |
|------|------|------|
| **API版本控制** | URL版本管理 | `/api/v1/`, `/api/v2/` |
| **Swagger文档** | 自动生成API文档 | Swagger UI |

---

### 2.13 分布式锁模块 (kevin.DistributedLock)

| 功能 | 说明 | 说明 |
|------|------|------|
| **分布式锁** | Redis分布式锁 | 防止并发冲突 |

---

### 2.14 雪花ID模块 (Kevin.SnowflakeId)

| 功能 | 说明 | 说明 |
|------|------|------|
| **ID生成** | 分布式唯一ID生成 | 雪花算法 |

---

### 2.15 短信服务模块 (Kevin.SMS)

| 功能 | 说明 | 支持平台 |
|------|------|----------|
| **短信发送** | 多平台短信服务 | 阿里云、腾讯云 |

---

### 2.16 邮件服务模块 (Kevin.Email)

| 功能 | 说明 | 说明 |
|------|------|------|
| **邮件发送** | 邮件发送服务 | SMTP |

---

## 三、业务应用模块拆分 (App)

### 3.1 应用服务层 (Application)

```
App/Application/
├── Services/
│   └── v1/
│       ├── AppCodeTestService.cs      # 代码测试服务
│       ├── AppDemoService.cs          # 示例服务
│       └── AppInfoTest/
│           └── AppInfoTestService.cs  # 信息测试服务
├── ModuleInitializer.cs               # 模块初始化
└── GlobalUsings.cs                    # 全局引用
```

### 3.2 领域层 (Domain)

```
App/Domain/
├── Entities/
│   ├── TAppCodeTest.cs               # 代码测试实体
│   ├── TAppDemo.cs                   # 示例实体
│   └── AppInfoTest/
│       └── TAppInfoTest.cs           # 信息测试实体
├── Interfaces/
│   ├── Repositorie/
│   │   └── v1/
│   │       ├── IAppCodeTestRp.cs     # 代码测试仓储接口
│   │       ├── IAppDemoRp.cs         # 示例仓储接口
│   │       └── AppInfoTest/
│   │           └── IAppInfoTestRp.cs # 信息测试仓储接口
│   └── Services/
│       └── v1/
│           ├── IAppCodeTestService.cs
│           ├── IAppDemoService.cs
│           └── AppInfoTest/
│               └── IAppInfoTestService.cs
├── ModuleInitializer.cs
└── GlobalUsings.cs
```

### 3.3 仓储实现层 (RepositorieRps)

```
App/RepositorieRps/
├── Repositories/
│   └── v1/
│       ├── AppCodeTestRp.cs
│       ├── AppDemoRp.cs
│       └── AppInfoTest/
│           └── AppInfoTestRp.cs
├── ModuleInitializer.cs
└── GlobalUsings.cs
```

### 3.4 API层 (WebApi)

```
App/WebApi/
├── Controllers/
│   └── v1/
│       └── VersionController.cs
│   └── v2/
│       └── VersionController.cs
├── Program.cs
├── appsettings.json
└── Properties/
    └── launchSettings.json
```

---

## 四、核心领域实体拆分 (Kevin/Domain)

### 4.1 AI相关实体

| 实体 | 说明 | 对应表 |
|------|------|--------|
| `TAIApps` | AI应用配置 | AI应用表 |
| `TAIChatHistorys` | 聊天历史 | 聊天历史表 |
| `TAIChats` | 聊天会话 | 聊天会话表 |
| `TAIKmss` | 知识库管理 | 知识库表 |
| `TAIModels` | AI模型配置 | 模型配置表 |
| `TAIPrompts` | 提示词管理 | 提示词表 |
| `TAISkillToolManagement` | 技能工具管理 | 技能工具表 |
| `TAISkillToolBindId` | 技能工具绑定 | 技能绑定表 |
| `TAIAppsBindId` | 应用绑定 | 应用绑定表 |

### 4.2 权限相关实体

| 实体 | 说明 | 对应表 |
|------|------|--------|
| `TUser` | 用户表 | 用户表 |
| `TRole` | 角色表 | 角色表 |
| `TPermission` | 权限表 | 权限表 |
| `TUserBindRole` | 用户角色绑定 | 用户角色表 |
| `TUserInfo` | 用户信息 | 用户信息表 |
| `TTenant` | 租户表 | 租户表 |

### 4.3 组织架构实体

| 实体 | 说明 | 对应表 |
|------|------|--------|
| `TDepartment` | 部门表 | 部门表 |
| `TPosition` | 职位表 | 职位表 |

### 4.4 系统配置实体

| 实体 | 说明 | 对应表 |
|------|------|--------|
| `TDictionary` | 数据字典 | 字典表 |
| `THttpLog` | HTTP日志 | HTTP日志表 |
| `TOSLog` | 系统日志 | 系统日志表 |
| `TMessage` | 消息表 | 消息表 |

---

## 五、技术栈详解

### 5.1 后端技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| .NET | 9.0 | 运行时框架 |
| ASP.NET Core | 9.0 | Web API框架 |
| Entity Framework Core | 9.0 | ORM框架 |
| MySQL | 8.0+ | 关系型数据库 |
| Redis | 7.0+ | 缓存、分布式锁 |
| Qdrant | 1.7+ | 向量数据库 |
| Hangfire | 1.8+ | 任务调度 |
| CAP | 7.0+ | 分布式事件总线 |
| RabbitMQ | 3.x | 消息队列 |
| Consul | 1.x | 服务发现 |
| SemanticKernel | 1.0+ | AI框架 |
| Ollama | - | 本地AI模型 |
| log4net | 2.x | 日志框架 |

### 5.2 前端技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.x | 前端框架 |
| AntDesign Vue | 4.x | UI组件库 |
| Axios | 1.x | HTTP客户端 |
| Pinia | 2.x | 状态管理 |
| Vue Router | 4.x | 路由管理 |

---

## 六、数据库设计

### 6.1 核心表结构

```sql
-- 用户表
CREATE TABLE T_User (
    Id BIGINT PRIMARY KEY,
    UserName VARCHAR(50),
    Password VARCHAR(100),
    RealName VARCHAR(50),
    Email VARCHAR(100),
    Phone VARCHAR(20),
    TenantId BIGINT,
    IsDelete BIT,
    CreateTime DATETIME,
    UpdateTime DATETIME
);

-- 角色表
CREATE TABLE T_Role (
    Id BIGINT PRIMARY KEY,
    RoleName VARCHAR(50),
    RoleCode VARCHAR(50),
    TenantId BIGINT,
    IsDelete BIT,
    CreateTime DATETIME
);

-- 权限表
CREATE TABLE T_Permission (
    Id BIGINT PRIMARYKEY,
    ParentId BIGINT,
    Name VARCHAR(50),
    Code VARCHAR(100),
    Type INT, -- 1:菜单 2:按钮 3:API
    Url VARCHAR(200),
    Icon VARCHAR(50),
    Sort INT,
    IsDelete BIT
);

-- AI应用表
CREATE TABLE T_AI_Apps (
    Id BIGINT PRIMARY KEY,
    AppName VARCHAR(100),
    AppDesc VARCHAR(500),
    ModelId BIGINT,
    TenantId BIGINT,
    IsDelete BIT,
    CreateTime DATETIME
);

-- 知识库表
CREATE TABLE T_AI_Kmss (
    Id BIGINT PRIMARY KEY,
    KmssName VARCHAR(100),
    KmssDesc VARCHAR(500),
    VectorModel VARCHAR(50),
    TenantId BIGINT,
    IsDelete BIT,
    CreateTime DATETIME
);
```

---

## 七、API接口规范

### 7.1 接口版本管理

```
/api/v1/[controller]    # V1版本接口
/api/v2/[controller]    # V2版本接口
```

### 7.2 接口响应格式

```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        // 业务数据
    },
    "timestamp": 1234567890
}
```

### 7.3 核心API模块

| 模块 | 路径前缀 | 说明 |
|------|----------|------|
| 用户管理 | `/api/v1/user` | 用户CRUD |
| 角色管理 | `/api/v1/role` | 角色CRUD |
| 权限管理 | `/api/v1/permission` | 权限管理 |
| AI应用 | `/api/v1/aiapps` | AI应用管理 |
| 知识库 | `/api/v1/aikmss` | 知识库管理 |
| AI聊天 | `/api/v1/aichat` | 智能对话 |
| 文件管理 | `/api/v1/file` | 文件上传下载 |

---

## 八、部署架构

### 8.1 Docker部署

```yaml
# docker-compose.yml
version: '3.8'
services:
  webapi:
    build: .
    ports:
      - "9901:80"
    depends_on:
      - mysql
      - redis
      - qdrant

  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"

  redis:
    image: redis:7.0
    ports:
      - "6379:6379"

  qdrant:
    image: qdrant/qdrant
    ports:
      - "6333:6333"

  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"
      - "15672:15672"
```

### 8.2 微服务架构

```
┌─────────────────────────────────────────────────────────┐
│                    Consul 服务注册中心                    │
└─────────────────────────────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│   WebApi服务  │  │   AI服务      │  │   调度服务    │
│   (主服务)    │  │  (智能体)     │  │  (Hangfire)  │
└───────────────┘  └───────────────┘  └───────────────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
   ┌─────────┐       ┌─────────┐       ┌─────────┐
   │  MySQL  │       │  Redis  │       │ Qdrant  │
   └─────────┘       └─────────┘       └─────────┘
```

---

## 九、开发规范

### 9.1 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 实体类 | `T`前缀 + 大驼峰 | `TUser`, `TRole` |
| 接口 | `I`前缀 + 大驼峰 | `IUserService` |
| 服务类 | 大驼峰 + `Service` | `UserService` |
| 仓储类 | 大驼峰 + `Rp` | `UserRp` |
| 控制器 | 大驼峰 + `Controller` | `UserController` |

### 9.2 分层规范

- **Controller层**: 只负责接收请求和返回响应
- **Application层**: 业务逻辑编排，调用Domain服务
- **Domain层**: 核心业务逻辑，领域模型
- **Repository层**: 数据访问，只做CRUD操作

### 9.3 依赖注入规范

```csharp
// 使用特性注入
[AppService]
public class UserService : IUserService
{
    // 自动注入
}

// 手动注册
services.AddScoped<IUserService, UserService>();
```

---

## 十、附录

### 10.1 环境变量配置

```json
{
  "ConnectionStrings": {
    "dbConnection": "server=127.0.0.1;port=3306;database=kevin_app;user id=root;password=admin123",
    "redisConnection": "127.0.0.1:6379,DefaultDatabase=0,password=123456"
  },
  "OllamaApiSetting": {
    "Url": "http://localhost:11434/api/embeddings",
    "DefaultModel": "qwen3:4b"
  },
  "QdrantClientSetting": {
    "Url": "http://localhost:6333"
  }
}
```

### 10.2 常用命令

```bash
# 启动项目
cd App/WebApi
dotnet run --environment Development

# 数据库迁移
Add-Migration "迁移名称"
Update-Database

# Docker启动
docker-compose up -d
```

### 10.3 相关文档

- [系统文档](SYSTEM_DOCUMENTATION.md)
- [README](README.md)
- [CSDN专栏](https://blog.csdn.net/weixin_42629287/category_13037923.html)

---

**文档版本**: v1.0  
**生成日期**: 2026/06/07  
**维护者**: NetCoreKevin 开发团队
