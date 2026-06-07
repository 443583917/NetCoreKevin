# NetCoreKevin 使用指南

本文档介绍了 NetCoreKevin 项目的环境要求、安装配置、AI 智能体配置、API 接口使用以及常见问题排查。

---

## 目录

- [环境要求](#环境要求)
  - [必备依赖](#必备依赖)
  - [可选依赖](#可选依赖)
- [安装配置](#安装配置)
  - [数据库配置](#数据库配置)
  - [数据库迁移](#数据库迁移)
  - [环境变量配置](#环境变量配置)
  - [启动运行](#启动运行)
  - [服务访问](#服务访问)
  - [默认账户](#默认账户)
- [AI 智能体配置](#ai-智能体配置)
  - [Qdrant 配置](#qdrant-配置)
  - [模型配置](#模型配置)
  - [AI 智能体使用流程](#ai-智能体使用流程)
  - [自定义 AI 工具开发](#自定义-ai-工具开发)
- [API 接口说明](#api-接口说明)
  - [认证接口](#认证接口)
  - [AI 接口](#ai-接口)
  - [任务调度接口](#任务调度接口)
  - [代码生成器接口](#代码生成器接口)
- [常见问题](#常见问题)
  - [数据库连接问题](#数据库连接问题)
  - [缓存服务问题](#缓存服务问题)
  - [AI 功能问题](#ai-功能问题)
  - [任务调度问题](#任务调度问题)
  - [启动失败问题](#启动失败问题)

---

## 环境要求

### 必备依赖

| 依赖 | 版本要求 | 说明 |
|------|----------|------|
| .NET SDK | 9.0+ | 开发环境 |
| MySQL | 8.0+ | 数据库 |
| Redis | 7.0+ | 缓存、Hangfire |
| Qdrant | 1.7+ | 向量数据库（AI 功能） |
| Node.js | 18+ | 前端开发（可选） |

### 可选依赖

| 依赖 | 说明 |
|------|------|
| RabbitMQ | CAP 消息队列 |
| Consul | 服务发现 |
| Ollama | 本地大语言模型 |

---

## 安装配置

### 数据库配置

**appsettings.json** 中的连接字符串配置：

```json
{
  "ConnectionStrings": {
    "dbConnection": "server=127.0.0.1;port=3306;database=kevin_app;user id=root;password=admin123;Convert Zero Datetime=True;TreatTinyAsBoolean=false;AllowLoadLocalInfile=true;Charset=utf8;Command Timeout=120;",
    "redisConnection": "127.0.0.1:6379,DefaultDatabase=0,password=123456"
  }
}
```

### 数据库迁移

在 **Package Manager Console** 中执行：

```powershell
# 1. 选择默认项目为 Kevin.EntityFrameworkCore
# 2. 创建迁移
Add-Migration "初始化数据库"

# 3. 应用迁移
Update-Database
```

### 环境变量配置

**Windows PowerShell：**
```powershell
$env:ASPNETCORE_ENVIRONMENT="Development"
```

**Linux/macOS：**
```bash
export ASPNETCORE_ENVIRONMENT=Development
```

### 启动运行

**方式一：使用 Visual Studio**

1. 打开解决方案 `kevin.abp.core.sln`
2. 设置 `App.WebApi` 为启动项目
3. 按 F5 启动调试

**方式二：使用命令行**
```bash
cd App/WebApi
dotnet run --environment Development
```

### 服务访问

| 服务 | 地址 |
|------|------|
| API 接口 | http://localhost:9901 |
| Swagger 文档 | http://localhost:9901/swagger |
| Hangfire 面板 | http://localhost:9901/pchangfire |

### 默认账户

| 账户 | 密码 | 租户 |
|------|------|------|
| admin | 123456 | 1000 |

---

## AI 智能体配置

### Qdrant 配置

```json
{
  "QdrantClientSetting": {
    "Url": "localhost"
  }
}
```

### 模型配置

**智谱 AI 配置：**
```json
{
  "OllamaApiSetting": {
    "Url": "https://open.bigmodel.cn/api/paas/v4/embeddings",
    "DefaultModel": "embedding-3",
    "ApiKey": "your-api-key"
  }
}
```

**Ollama 本地模型配置：**
```json
{
  "OllamaApiSetting": {
    "Url": "http://localhost:11434/api/embeddings",
    "DefaultModel": "qwen3:4b",
    "ApiKey": ""
  }
}
```

### AI 智能体使用流程

1. **注册 AI 账户** -- 获取 API Key
2. **配置模型** -- 设置向量模型和对话模型
3. **新建知识库** -- 上传文档，选择向量模型
4. **配置智能体** -- 绑定技能工具
5. **开始对话** -- 与 AI 智能体交互

### 自定义 AI 工具开发

#### 开发步骤

**1. 创建工具接口**

```csharp
public interface IMyCustomToolService
{
    [Description("自定义工具描述")]
    Task<string> MyToolMethod(
        [Description("参数1说明")] string param1,
        [Description("参数2说明")] int param2
    );
}
```

**2. 实现工具服务**

```csharp
public class MyCustomToolService : IMyCustomToolService
{
    public async Task<string> MyToolMethod(string param1, int param2)
    {
        // 实现逻辑
        return "执行结果";
    }
}
```

**3. 注册工具到容器**

在模块初始化中注册服务：
```csharp
services.AddTransient<IMyCustomToolService, MyCustomToolService>();
```

**4. 配置工具到智能体**

在 `AIAgentToolSkillService.cs` 中添加工具注册：
```csharp
aiTools.Add(
    AIFunctionFactory.Create(_myCustomToolService.MyToolMethod,
    new AIFunctionFactoryOptions
    {
        Name = "MyToolMethod",
        Description = "自定义工具描述"
    }
));
```

#### 工具开发规范

- 使用 `[Description]` 属性添加参数说明
- 返回类型推荐使用 `Task<string>`
- 错误返回以 `X` 开头的字符串
- 支持的参数类型：`string`, `int`, `bool`, `List<string>`

---

## API 接口说明

### 认证接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/Authorize/Login` | POST | 用户登录 |
| `/api/Authorize/Logout` | POST | 用户登出 |
| `/api/Authorize/RefreshToken` | POST | 刷新 Token |

**登录请求示例：**
```json
{
  "userName": "admin",
  "password": "123456",
  "tenantId": 1000
}
```

### AI 接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/AIChats/Chat` | POST | 发送聊天消息 |
| `/api/AIKmss/Create` | POST | 创建知识库 |
| `/api/AIKmss/UploadFile` | POST | 上传知识库文档 |
| `/api/AIModels/GetAll` | GET | 获取模型列表 |

### 任务调度接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/AITasks/AddOrUpdateCronTask` | POST | 创建/更新定时任务 |
| `/api/AITasks/RemoveCronTask` | POST | 删除定时任务 |
| `/api/AITasks/TriggerCronTask` | POST | 立即触发任务 |

### 代码生成器接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/CodeGenerator/GetAreaNames` | GET | 获取区域列表 |
| `/api/CodeGenerator/GetAreaNameEntityItems` | GET | 获取实体列表 |
| `/api/CodeGenerator/BulidCode` | POST | 生成代码 |

**代码生成器使用注意事项：**

- 仅超级管理员有权限使用代码生成器
- 生成的代码会覆盖现有文件，请谨慎操作
- 建议在生成前备份相关文件

**代码生成器配置（appsettings.json）：**

```json
{
  "CodeGeneratorSetting": {
    "CodeGeneratorItems": [
      {
        "AreaName": "App.WebApi.v1",
        "AreaPath": "App.Domain.Entities",
        "IRpBulidPath": "App.Domain.Interfaces.Repositorie.v1",
        "RpBulidPath": "App.RepositorieRps.Repositories.v1",
        "IServiceBulidPath": "App.Domain.Interfaces.Services.v1",
        "ServiceBulidPath": "App.Application.Services.v1"
      }
    ]
  }
}
```

---

## 常见问题

### 数据库连接问题

| 问题 | 排查步骤 | 解决方案 |
|------|----------|----------|
| 数据库连接失败 | 检查 MySQL 服务状态 | 确认 MySQL 已启动且连接字符串配置正确 |
| 数据库连接超时 | 测试网络连通性 | 检查 `Command Timeout` 参数和网络延迟 |

```bash
# 测试 MySQL 连接
mysql -h 127.0.0.1 -P 3306 -u root -p
```

### 缓存服务问题

| 问题 | 排查步骤 | 解决方案 |
|------|----------|----------|
| Redis 连接失败 | 检查 Redis 服务状态 | 确认 Redis 已启动且密码配置正确 |
| Hangfire 任务不执行 | 检查 Redis 连接 | 确保 Redis 连接正常且服务已启动 |

```bash
# 测试 Redis 连接
redis-cli -h 127.0.0.1 -p 6379 -a 123456 ping
```

### AI 功能问题

| 问题 | 排查步骤 | 解决方案 |
|------|----------|----------|
| Qdrant 连接失败 | 检查 Qdrant 服务状态 | 确认 Qdrant 已启动且 URL 配置正确 |
| AI 工具调用失败 | 检查工具注册配置 | 确保 `InitData` 正确传递参数 |
| 知识库检索失败 | 检查向量模型配置 | 确认文档已上传且向量模型正常 |
| 问答返回空或错误 | 检查 Qdrant 服务状态 | 查看日志文件 `Logs/log_*.txt` |

```bash
# 测试 Qdrant 健康状态
curl http://localhost:6333/health
```

**AI 工具调用失败排查步骤：**

1. 检查工具注册配置
2. 检查 `InitData` 参数传递
3. 查看日志文件 `Logs/log_*.txt`
4. 检查工具方法签名是否正确

**知识库检索失败排查步骤：**

1. 检查 Qdrant 服务状态
2. 检查知识库文档是否已上传
3. 检查向量模型配置

### 任务调度问题

| 问题 | 排查步骤 | 解决方案 |
|------|----------|----------|
| 定时任务不触发 | 检查 Redis 连接 | 确认 Hangfire 服务已正确启动 |
| 任务执行失败 | 查看 Hangfire 面板 | 检查任务日志中的错误信息 |

### 启动失败问题

**症状：** 应用无法启动

**排查步骤：**

1. **检查端口是否被占用**
```bash
netstat -ano | findstr :9901
```

2. **检查数据库连接**
```bash
mysql -h 127.0.0.1 -P 3306 -u root -p
```

3. **检查 Redis 连接**
```bash
redis-cli -h 127.0.0.1 -p 6379 -a 123456 ping
```

4. **检查环境变量配置**
```powershell
# Windows
$env:ASPNETCORE_ENVIRONMENT

# Linux/macOS
echo $ASPNETCORE_ENVIRONMENT
```

---

**文档版本**: v1.0
**最后更新**: 2026/06/07
**来源文档**: SYSTEM_DOCUMENTATION.md
