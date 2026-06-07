# NetCoreKevin 运维指南

本文档介绍了 NetCoreKevin 项目的部署架构、监控运维、安全规范、故障排查和系统调优。

---

## 目录

- [部署架构](#部署架构)
  - [Docker 部署](#docker-部署)
  - [Linux 系统服务部署](#linux-系统服务部署)
  - [环境变量配置](#环境变量配置)
- [监控运维](#监控运维)
  - [日志管理](#日志管理)
  - [性能监控](#性能监控)
  - [常见问题速查](#常见问题速查)
- [安全规范](#安全规范)
  - [输入验证](#输入验证)
  - [权限控制](#权限控制)
  - [数据加密](#数据加密)
  - [安全防护](#安全防护)
- [故障排查](#故障排查)
  - [启动失败](#启动失败)
  - [AI 工具调用失败](#ai-工具调用失败)
  - [知识库检索失败](#知识库检索失败)
- [系统调优](#系统调优)
  - [数据库优化](#数据库优化)
  - [缓存策略](#缓存策略)
  - [异步编程优化](#异步编程优化)
  - [日志优化](#日志优化)
  - [容器资源调优](#容器资源调优)

---

## 部署架构

### Docker 部署

#### Dockerfile 示例

```dockerfile
FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS base
WORKDIR /app
EXPOSE 80
EXPOSE 443

FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src
COPY ["App/WebApi/App.WebApi.csproj", "App/WebApi/"]
RUN dotnet restore "App/WebApi/App.WebApi.csproj"
COPY . .
WORKDIR "/src/App/WebApi"
RUN dotnet build "App.WebApi.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "App.WebApi.csproj" -c Release -o /app/publish

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "App.WebApi.dll"]
```

#### Docker Compose 示例

```yaml
version: '3.8'
services:
  webapi:
    build: .
    ports:
      - "9901:80"
    environment:
      - ASPNETCORE_ENVIRONMENT=Production
      - ConnectionStrings__dbConnection=server=mysql;port=3306;database=kevin_app;user id=root;password=admin123
    depends_on:
      - mysql
      - redis

  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=admin123
      - MYSQL_DATABASE=kevin_app

  redis:
    image: redis:7.0
    ports:
      - "6379:6379"
```

### Linux 系统服务部署

#### 创建 systemd 服务文件

创建文件 `/etc/systemd/system/kevin-webapi.service`：

```ini
[Unit]
Description=NetCoreKevin Web API
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/kevin
ExecStart=/usr/bin/dotnet /var/www/kevin/App.WebApi.dll
Restart=always
RestartSec=5
Environment=ASPNETCORE_ENVIRONMENT=Production

[Install]
WantedBy=multi-user.target
```

#### 启动服务

```bash
sudo systemctl daemon-reload
sudo systemctl start kevin-webapi
sudo systemctl enable kevin-webapi
```

#### 服务管理命令

```bash
# 查看服务状态
sudo systemctl status kevin-webapi

# 停止服务
sudo systemctl stop kevin-webapi

# 重启服务
sudo systemctl restart kevin-webapi

# 查看日志
sudo journalctl -u kevin-webapi -f
```

### 环境变量配置

#### Windows PowerShell

```powershell
$env:ASPNETCORE_ENVIRONMENT="Development"
```

#### Linux/macOS

```bash
export ASPNETCORE_ENVIRONMENT=Development
```

#### 生产环境数据库配置

**appsettings.json** 中的连接字符串配置：

```json
{
  "ConnectionStrings": {
    "dbConnection": "server=127.0.0.1;port=3306;database=kevin_app;user id=root;password=admin123;Convert Zero Datetime=True;TreatTinyAsBoolean=false;AllowLoadLocalInfile=true;Charset=utf8;Command Timeout=120;",
    "redisConnection": "127.0.0.1:6379,DefaultDatabase=0,password=123456"
  }
}
```

---

## 监控运维

### 日志管理

日志文件位于 `App/WebApi/Logs/` 目录：

| 日志文件 | 说明 |
|----------|------|
| `log_*.txt` | 应用日志 |
| `error_*.txt` | 错误日志 |
| `http_*.txt` | HTTP 请求日志 |

#### 日志使用规范

```csharp
// 使用 log4net 记录日志
LogHelper<MyService>.logger.Info("信息级别日志");
LogHelper<MyService>.logger.Warn("警告级别日志");
LogHelper<MyService>.logger.Error("错误级别日志", exception);
```

#### 日志级别说明

| 级别 | 用途 | 生产环境建议 |
|------|------|-------------|
| Debug | 调试信息 | 关闭 |
| Info | 一般信息 | 开启 |
| Warn | 警告信息 | 开启 |
| Error | 错误信息 | 开启 |
| Fatal | 致命错误 | 开启 |

### 性能监控

| 监控工具 | 地址 | 说明 |
|----------|------|------|
| **Hangfire Dashboard** | http://localhost:9901/pchangfire | 任务调度监控 |
| **Redis 监控** | 使用 Redis Insight | 缓存状态监控 |
| **MySQL 监控** | 使用 Prometheus + Grafana | 数据库性能监控 |

#### Hangfire 监控要点

- 查看任务执行历史和状态
- 监控失败任务并手动触发重试
- 观察任务队列积压情况
- 检查 recurring jobs 的调度频率

#### Redis 监控要点

- 监控内存使用率
- 观察连接数和命中率
- 检查慢查询日志
- 监控主从同步状态

### 常见问题速查

| 问题 | 解决方案 |
|------|----------|
| **数据库连接失败** | 检查 MySQL 服务是否启动，连接字符串配置是否正确 |
| **Redis 连接失败** | 检查 Redis 服务是否启动，密码配置是否正确 |
| **Qdrant 连接失败** | 检查 Qdrant 服务是否启动，配置的 URL 是否正确 |
| **AI 工具调用失败** | 检查工具注册配置，确保 `InitData` 正确传递参数 |
| **Hangfire 任务不执行** | 检查 Redis 连接，确保服务已启动 |
| **端口被占用** | 使用 `netstat -ano | findstr :9901` 查找占用进程并释放 |

---

## 安全规范

### 输入验证

- 所有用户输入必须进行验证
- 使用数据注解或 FluentValidation
- 防止 SQL 注入、XSS 攻击

### 权限控制

- 使用基于角色的访问控制（RBAC）
- 敏感接口需要权限验证
- 日志记录所有权限检查失败的请求

#### 权限配置项

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `IsOpenPermission` | 是否开启权限验证 | true |
| `TenantId` | 默认租户 ID | 1000 |
| `Jwt:AccessTokenExpirationMinutes` | Token 过期时间（分钟） | 60 |

### 数据加密

- 数据库密码使用 AES 加密存储
- API 密钥等敏感配置使用环境变量
- 传输层使用 HTTPS

### 安全防护

#### 危险命令拦截

系统会拦截以下危险命令：

- `rm -rf`
- `del /s /q`
- 格式化磁盘命令
- 其他破坏性操作

#### 域名白名单

HTTP/HTTPS 请求需要配置授权域名白名单：

```json
{
  "AuthorizedDomains": [
    "example.com",
    "api.example.com"
  ]
}
```

---

## 故障排查

### 启动失败

**症状：** 应用无法启动

**排查步骤：**

1. 检查端口是否被占用

```bash
netstat -ano | findstr :9901
```

2. 检查数据库连接

```bash
# 测试 MySQL 连接
mysql -h 127.0.0.1 -P 3306 -u root -p
```

3. 检查 Redis 连接

```bash
redis-cli -h 127.0.0.1 -p 6379 -a 123456 ping
```

4. 检查应用日志

```bash
# 查看最新日志
tail -f App/WebApi/Logs/log_*.txt
```

### AI 工具调用失败

**症状：** Error: Function failed

**排查步骤：**

1. 检查工具注册配置
2. 检查 `InitData` 参数传递
3. 查看日志文件 `Logs/log_*.txt`
4. 检查工具方法签名是否正确
5. 确认工具已正确注册到 DI 容器

### 知识库检索失败

**症状：** 问答返回空或错误

**排查步骤：**

1. 检查 Qdrant 服务状态

```bash
curl http://localhost:6333/health
```

2. 检查知识库文档是否已上传
3. 检查向量模型配置
4. 验证 Embedding 模型 API 是否可用

---

## 系统调优

### 数据库优化

| 优化项 | 说明 | 实施建议 |
|--------|------|----------|
| 索引优化 | 为常用查询字段创建索引 | 分析慢查询日志，针对高频查询建立复合索引 |
| 分页查询 | 使用 Skip/Take 避免全表扫描 | 对大数据量查询强制分页，避免一次加载全部数据 |
| 批量操作 | 使用 EF Core 的批量操作 API | 减少单条 SQL 执行次数，合并写入操作 |
| 读写分离 | 配置数据库读写分离 | 主库处理写操作，从库处理读查询，降低主库压力 |
| 连接池管理 | 合理配置数据库连接池 | 避免连接泄漏，设置合理的最大/最小连接数 |

#### 数据库字段索引配置

```json
{
  "DBDefaultHasIndexFields": "tableid,createtime,updatetime,deletetime,tenantid,createuserid,updateuserid,deleteuserid"
}
```

### 缓存策略

```csharp
// 使用 Redis 缓存
[CacheDataFilter(Duration = 60)] // 缓存60秒
public async Task<MyDto> GetData(int id)
{
    // 查询逻辑
}
```

#### 缓存最佳实践

- 热点数据设置较短的过期时间（30-60 秒）
- 变化频率低的数据设置较长的过期时间（10-30 分钟）
- 使用缓存穿透防护，对空结果也进行缓存
- 定期审查缓存命中率，调整缓存策略

### 异步编程优化

- 使用 `async/await` 模式
- 避免同步阻塞调用
- 使用 `ConfigureAwait(false)` 优化性能
- 避免在异步方法中使用 `.Result` 或 `.Wait()` 造成死锁

### 日志优化

- 生产环境关闭 Debug 级别日志
- 使用异步日志写入
- 定期清理日志文件
- 对高频日志使用结构化日志格式
- 配置日志文件滚动策略，防止单个文件过大

### 容器资源调优

#### Docker 资源限制

```yaml
services:
  webapi:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 2G
        reservations:
          cpus: '1.0'
          memory: 512M
```

#### 健康检查配置

```yaml
services:
  webapi:
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:80/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 15s
```

---

## 附录

### 配置项说明

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `IsOpenPermission` | 是否开启权限验证 | true |
| `TenantId` | 默认租户 ID | 1000 |
| `Jwt:AccessTokenExpirationMinutes` | Token 过期时间（分钟） | 60 |
| `HangfireSetting` | Hangfire 配置 | - |
| `CorsSetting` | CORS 跨域配置 | - |
| `CodeGeneratorSetting` | 代码生成器配置 | - |
| `AuthorizedDomains` | 授权域名白名单 | - |

### 服务端口一览

| 服务 | 默认端口 | 说明 |
|------|----------|------|
| Web API | 9901 | 应用主服务 |
| MySQL | 3306 | 数据库 |
| Redis | 6379 | 缓存 |
| Qdrant | 6333 | 向量数据库 |
| Ollama | 11434 | 本地大语言模型（可选） |

---

**文档版本**: v1.0
**最后更新**: 2026/06/07
**来源文档**: SYSTEM_DOCUMENTATION.md (第九、十、十一、十二章)
