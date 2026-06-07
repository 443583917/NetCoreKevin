# NetCoreKevin 架构设计文档

本文档详细介绍了NetCoreKevin项目的系统架构、技术栈、模块设计和数据库设计。

---

## 目录

- [系统架构](#系统架构)
  - [架构概览](#架构概览)
  - [项目目录结构](#项目目录结构)
- [技术栈](#技术栈)
  - [后端技术栈](#后端技术栈)
  - [前端技术栈](#前端技术栈)
- [模块设计](#模块设计)
  - [AI 智能体模块](#ai-智能体模块)
  - [AI 知识库模块](#ai-知识库模块)
  - [用户权限模块](#用户权限模块)
  - [任务调度模块](#任务调度模块)
  - [消息队列模块](#消息队列模块)
  - [文件存储模块](#文件存储模块)
  - [认证授权模块](#认证授权模块)
  - [实时通信模块](#实时通信模块)
  - [缓存模块](#缓存模块)
  - [日志模块](#日志模块)
  - [服务发现模块](#服务发现模块)
  - [API版本管理模块](#api版本管理模块)
  - [分布式锁模块](#分布式锁模块)
  - [雪花ID模块](#雪花id模块)
  - [短信服务模块](#短信服务模块)
  - [邮件服务模块](#邮件服务模块)
  - [业务应用模块](#业务应用模块)
  - [核心领域实体](#核心领域实体)
- [数据库设计](#数据库设计)
  - [核心表结构](#核心表结构)
- [架构原则](#架构原则)
  - [分层架构原则](#分层架构原则)
  - [领域驱动设计](#领域驱动设计)
  - [模块化设计原则](#模块化设计原则)
  - [编码规范](#编码规范)
  - [依赖注入规范](#依赖注入规范)

---

## 系统架构

### 架构概览

系统采用分层架构设计，整体分为前端层、API网关层、核心框架层、业务应用层和数据存储层五大层次。

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

**各层职责：**

- **前端层**：基于 Vue3 + AntDesign 的单页面应用，负责用户界面交互
- **API 网关层**：ASP.NET Core WebApi，统一的接口入口，负责请求路由和响应处理
- **核心框架层**：Kevin核心模块，提供应用服务、领域模型、数据访问和通用功能模块
- **业务应用层**：App业务模块，包含具体的业务逻辑、领域实体、仓储实现和API控制器
- **数据存储层**：MySQL关系数据库、Redis缓存、Qdrant向量数据库、RabbitMQ消息队列、文件存储

### 项目目录结构

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

## 技术栈

### 后端技术栈

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

### 前端技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.x | 前端框架 |
| AntDesign Vue | 4.x | UI组件库 |
| Axios | 1.x | HTTP客户端 |
| Pinia | 2.x | 状态管理 |
| Vue Router | 4.x | 路由管理 |

---

## 模块设计

### AI 智能体模块

**模块路径：** `kevin.AI.AgentFramework`

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **智能体服务** | 基于SemanticKernel的AI代理系统 | `AIAgentService` |
| **技能管理** | 动态加载和管理AI技能 | `IAIAgentToolSkillService` |
| **工具集成** | 集成各类外部工具 | `CommonToolsService` |
| **工作流** | 多步骤任务编排 | `WorkFlowsAndAIAgentsDemo` |
| **脚本执行** | Python脚本运行支持 | `PySubprocessScriptRunner` |

**目录结构：**

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

### AI 知识库模块

**模块路径：** `Kevin.RAG`

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **向量化服务** | 文档向量化处理 | `IOllamaService` |
| **向量存储** | Qdrant向量数据库操作 | `IQdrantService` |
| **RAG检索** | 检索增强生成 | `IRAGService` |
| **知识库管理** | 知识库CRUD | `AIKmssService` |

**目录结构：**

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

### 用户权限模块

**模块路径：** `kevin.Permission`

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **用户管理** | 用户CRUD、角色绑定 | `UserService` |
| **角色管理** | 角色CRUD、权限配置 | `RoleService` |
| **权限管理** | 菜单权限、API权限 | `PermissionService` |
| **组织架构** | 部门、职位管理 | `DepartmentService` |

**核心实体：**

- `TUser` - 用户表
- `TRole` - 角色表
- `TPermission` - 权限表
- `TDepartment` - 部门表
- `TPosition` - 职位表
- `TTenant` - 租户表

---

### 任务调度模块

**模块路径：** `Kevin.Hangfire`

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **定时任务** | Cron表达式任务调度 | `IModuleConfigTasks` |
| **任务监控** | 任务执行状态监控 | Hangfire Dashboard |
| **自动注册** | 启动时自动注册任务 | `ServiceCollectionExtensions` |

**配置示例：**

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

### 消息队列模块

**模块路径：** `kevin.Cap`

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **事件发布** | 分布式事件发布 | CAP Publisher |
| **事件订阅** | 事件订阅处理 | `DemoSubscribe` |
| **消息过滤** | 消息过滤器 | `CapSubscribeFilter` |

---

### 文件存储模块

**模块路径：** `kevin.FileStorage`

| 功能 | 说明 | 支持平台 |
|------|------|----------|
| **文件上传** | 多文件上传 | - |
| **文件下载** | 文件下载服务 | - |
| **云存储** | 多云存储支持 | 腾讯云COS、阿里云OSS、七牛云 |

**目录结构：**

```
kevin.FileStorage/
├── AliCloud/               # 阿里云OSS
├── KevinStaticFiles/       # 本地文件存储
├── QiniuCloud/             # 七牛云
└── TencentCloud/           # 腾讯云COS
```

---

### 认证授权模块

**模块路径：** `Kevin.Authentication.Jwt`

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **JWT生成** | Token生成服务 | `TokenService` |
| **Token刷新** | Token刷新机制 | `RefreshTokenDto` |
| **Claims管理** | 用户声明管理 | `JwtKeinClaimTypes` |

---

### 实时通信模块

**模块路径：** `Kevin.SignalR`

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **消息推送** | 实时消息推送 | `SignalRService` |
| **连接管理** | 客户端连接管理 | Hub |
| **群组管理** | 消息群组管理 | - |

---

### 缓存模块

**模块路径：** `kevin.Cache`

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **Redis缓存** | 分布式缓存服务 | `CacheService` |
| **缓存策略** | 多种缓存策略支持 | `ICacheService` |

---

### 日志模块

**模块路径：** `Kevin.log4Net`

| 功能 | 说明 | 说明 |
|------|------|------|
| **日志记录** | 多级别日志记录 | Info/Warn/Error/Fatal |
| **日志配置** | 灵活的日志配置 | XML配置文件 |

---

### 服务发现模块

**模块路径：** `kevin.Consul`

| 功能 | 说明 | 核心类 |
|------|------|--------|
| **服务注册** | 服务自动注册 | Consul Client |
| **服务发现** | 服务地址发现 | - |
| **健康检查** | 服务健康检查 | - |

---

### API版本管理模块

**模块路径：** `Kevin.Versioning.Swagger`

| 功能 | 说明 | 说明 |
|------|------|------|
| **API版本控制** | URL版本管理 | `/api/v1/`, `/api/v2/` |
| **Swagger文档** | 自动生成API文档 | Swagger UI |

---

### 分布式锁模块

**模块路径：** `kevin.DistributedLock`

| 功能 | 说明 | 说明 |
|------|------|------|
| **分布式锁** | Redis分布式锁 | 防止并发冲突 |

---

### 雪花ID模块

**模块路径：** `Kevin.SnowflakeId`

| 功能 | 说明 | 说明 |
|------|------|------|
| **ID生成** | 分布式唯一ID生成 | 雪花算法 |

---

### 短信服务模块

**模块路径：** `Kevin.SMS`

| 功能 | 说明 | 支持平台 |
|------|------|----------|
| **短信发送** | 多平台短信服务 | 阿里云、腾讯云 |

---

### 邮件服务模块

**模块路径：** `Kevin.Email`

| 功能 | 说明 | 说明 |
|------|------|------|
| **邮件发送** | 邮件发送服务 | SMTP |

---

### 业务应用模块

业务应用模块（App）是项目的业务实现层，包含以下子模块：

#### 应用服务层 (Application)

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

#### 领域层 (Domain)

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

#### 仓储实现层 (RepositorieRps)

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

#### API层 (WebApi)

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

### 核心领域实体

#### AI相关实体

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

#### 权限相关实体

| 实体 | 说明 | 对应表 |
|------|------|--------|
| `TUser` | 用户表 | 用户表 |
| `TRole` | 角色表 | 角色表 |
| `TPermission` | 权限表 | 权限表 |
| `TUserBindRole` | 用户角色绑定 | 用户角色表 |
| `TUserInfo` | 用户信息 | 用户信息表 |
| `TTenant` | 租户表 | 租户表 |

#### 组织架构实体

| 实体 | 说明 | 对应表 |
|------|------|--------|
| `TDepartment` | 部门表 | 部门表 |
| `TPosition` | 职位表 | 职位表 |

#### 系统配置实体

| 实体 | 说明 | 对应表 |
|------|------|--------|
| `TDictionary` | 数据字典 | 字典表 |
| `THttpLog` | HTTP日志 | HTTP日志表 |
| `TOSLog` | 系统日志 | 系统日志表 |
| `TMessage` | 消息表 | 消息表 |

---

## 数据库设计

### 核心表结构

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

## 架构原则

### 分层架构原则

项目采用严格的分层架构，各层职责明确：

| 层次 | 职责 | 规则 |
|------|------|------|
| **Controller层** | 接收请求和返回响应 | 不包含业务逻辑，只做参数校验和结果封装 |
| **Application层** | 业务逻辑编排 | 调用Domain层服务，不直接访问数据层 |
| **Domain层** | 核心业务逻辑 | 领域模型和业务规则，不依赖基础设施 |
| **Repository层** | 数据访问 | 只做CRUD操作，不包含业务逻辑 |

**依赖方向：** Controller -> Application -> Domain <- Repository

### 领域驱动设计

项目遵循DDD（Domain-Driven Design）核心概念：

- **实体（Entity）**：具有唯一标识的领域对象，如 `TUser`、`TRole`
- **值对象（Value Object）**：没有唯一标识的不可变对象
- **聚合根（Aggregate Root）**：聚合的入口点，保证事务一致性
- **仓储（Repository）**：领域对象的持久化接口，隔离数据访问细节
- **领域事件（Domain Event）**：解耦模块间通信，实现事件驱动架构
- **应用服务（Application Service）**：编排业务用例，协调多个领域服务

### 模块化设计原则

项目采用模块化架构，每个功能模块独立封装：

- **单一职责**：每个模块只负责一个特定的业务领域
- **接口隔离**：模块间通过接口通信，不依赖具体实现
- **松耦合**：模块之间尽量减少直接依赖
- **高内聚**：模块内部各组件紧密协作完成特定功能
- **可插拔**：模块可以独立启用或禁用

**模块列表：**

| 模块 | 职责 | 关键依赖 |
|------|------|----------|
| `kevin.Permission` | 用户权限管理 | EF Core, JWT |
| `Kevin.Authentication.Jwt` | JWT认证授权 | Microsoft.Identity |
| `kevin.AI.AgentFramework` | AI智能体 | SemanticKernel |
| `Kevin.RAG` | RAG检索增强 | Ollama, Qdrant |
| `Kevin.Hangfire` | 任务调度 | Hangfire |
| `kevin.Cap` | 分布式事件 | CAP, RabbitMQ |
| `kevin.Cache` | 分布式缓存 | Redis |
| `kevin.FileStorage` | 文件存储 | 腾讯云/阿里云/七牛云 |
| `Kevin.SignalR` | 实时通信 | SignalR |
| `kevin.Consul` | 服务发现 | Consul |

### 编码规范

#### 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 实体类 | `T`前缀 + 大驼峰 | `TUser`, `TRole` |
| 接口 | `I`前缀 + 大驼峰 | `IUserService` |
| 服务类 | 大驼峰 + `Service` | `UserService` |
| 仓储类 | 大驼峰 + `Rp` | `UserRp` |
| 控制器 | 大驼峰 + `Controller` | `UserController` |
| 数据传输对象 | 大驼峰 + `Dto` | `UserDto` |

#### API接口规范

接口采用URL版本管理：

```
/api/v1/[controller]    # V1版本接口
/api/v2/[controller]    # V2版本接口
```

**统一响应格式：**

```json
{
    "code": 200,
    "message": "操作成功",
    "data": {},
    "timestamp": 1234567890
}
```

**核心API模块：**

| 模块 | 路径前缀 | 说明 |
|------|----------|------|
| 用户管理 | `/api/v1/user` | 用户CRUD |
| 角色管理 | `/api/v1/role` | 角色CRUD |
| 权限管理 | `/api/v1/permission` | 权限管理 |
| AI应用 | `/api/v1/aiapps` | AI应用管理 |
| 知识库 | `/api/v1/aikmss` | 知识库管理 |
| AI聊天 | `/api/v1/aichat` | 智能对话 |
| 文件管理 | `/api/v1/file` | 文件上传下载 |

### 依赖注入规范

项目使用 `[AppService]` 特性进行自动依赖注入：

```csharp
// 特性注入（推荐）
[AppService]
public class UserService : IUserService
{
    // 自动注册为 Scoped 生命周期
}

// 手动注册（特殊场景）
services.AddScoped<IUserService, UserService>();
```

---

**文档版本**: v1.0
**最后更新**: 2026/06/07
**维护者**: NetCoreKevin 开发团队
