# GoKevin - 企业级 AI 智能体 SaaS 平台

基于 Go 语言和 Vue.js 构建的企业级 AI 智能体 SaaS 平台，从 NetCoreKevin (.NET 9) 完整重构而来。

## 项目概述

GoKevin 是一个完整的全栈 AI 智能体平台，采用现代化的技术栈和架构设计：

- **后端**: Go 1.22 + Gin + GORM + tRPC-Agent-Go
- **前端**: Vue 3 + Element Plus + Vite
- **架构**: DDD + Event-Driven + CQRS
- **部署**: Docker + Docker Compose

## 项目结构

```
GoKevin/
├── backend/                # Go 后端服务
│   ├── cmd/               # 应用入口
│   ├── internal/          # 内部实现
│   │   ├── domain/        # 领域层
│   │   ├── application/   # 应用层
│   │   ├── infrastructure/# 基础设施层
│   │   └── interfaces/    # 接口层
│   ├── pkg/               # 公共工具包
│   ├── configs/           # 配置文件
│   ├── scripts/           # 脚本文件
│   └── test/              # 测试文件
│
├── frontend/              # Vue.js 前端应用
│   ├── src/               # 源代码
│   ├── public/            # 静态资源
│   └── package.json       # 依赖配置
│
├── docs/                  # 项目文档
│   ├── go-backend-development-tasks.md
│   ├── go-refactoring-implementation.md
│   ├── go-refactoring-design.md
│   └── go-remaining-modules.md
│
└── scripts/               # 部署脚本
```

## 技术栈

### 后端技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.22+ |
| Web 框架 | Gin | 1.12+ |
| ORM | GORM | 1.31+ |
| AI 框架 | tRPC-Agent-Go | latest |
| 认证 | JWT + RBAC | - |
| 缓存 | Redis | 7.0+ |
| 消息队列 | RabbitMQ | 3.x |
| 向量数据库 | Qdrant | 1.7+ |
| 日志 | Zap | 1.26+ |
| 配置 | Viper | 1.18+ |

### 前端技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| 框架 | Vue | 3.x |
| UI 组件库 | Element Plus | latest |
| 构建工具 | Vite | latest |
| HTTP 客户端 | Axios | latest |
| 状态管理 | Pinia | latest |
| 路由 | Vue Router | 4.x |

## 快速开始

### 环境要求

- Go 1.22+
- Node.js 18+
- MySQL 8.0+
- Redis 7.0+

### 后端启动

```bash
cd backend

# 安装依赖
go mod download

# 配置数据库
# 编辑 configs/config.yaml

# 运行
make run
# 或
go run cmd/server/main.go
```

### 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 开发模式
npm run serve

# 构建生产版本
npm run build
```

### Docker 部署

```bash
cd backend

# 构建并启动所有服务
docker-compose up -d

# 或使用 Makefile
make docker-build
make docker-run
```

## 核心功能

### AI 智能体功能

- **LLM 集成**: 支持 OpenAI、Claude 等多种大语言模型
- **Agent 系统**: 智能体核心，支持工具调用和任务分解
- **RAG 服务**: 知识库检索增强生成
- **向量存储**: Qdrant 向量数据库集成
- **会话管理**: 多轮对话和上下文管理

### 基础设施功能

- **用户权限**: 完整的 RBAC 权限管理系统
- **认证授权**: JWT Token 认证，支持刷新机制
- **缓存系统**: Redis 分布式缓存
- **消息队列**: RabbitMQ 异步消息处理
- **任务调度**: Cron 定时任务
- **文件存储**: 本地文件存储
- **API 文档**: Swagger UI 自动生成

## 开发指南

### 后端开发

1. 遵循 DDD 架构原则
2. 使用 CQRS 模式分离命令和查询
3. 通过领域事件解耦模块
4. 编写单元测试和集成测试

### 前端开发

1. 使用 Vue 3 Composition API
2. 遵循 Element Plus 设计规范
3. 使用 TypeScript 类型定义
4. 组件化开发

## 文档

项目提供了完整的文档体系，帮助开发者快速了解和参与项目开发。

### 📚 核心文档

| 文档 | 说明 | 适用人群 |
|------|------|----------|
| [后端架构功能说明](docs/backend-architecture.md) | 详细的后端架构设计、分层结构、模块功能说明 | 后端开发者、架构师 |
| [运行和部署指南](docs/deployment-guide.md) | 本地开发环境搭建、Docker 部署、生产环境部署 | 所有开发者、运维人员 |
| [API 接口文档](backend/docs/api/README.md) | RESTful API 接口详细说明 | 前后端开发者 |

### 📋 开发文档

| 文档 | 说明 | 适用人群 |
|------|------|----------|
| [Go 后端开发任务清单](docs/go-backend-development-tasks.md) | 后端开发任务列表、进度跟踪、待优化项 | 后端开发者 |
| [Go 重构设计方案](docs/go-refactoring-design.md) | 从 .NET 到 Go 的重构架构设计、技术选型 | 架构师、技术负责人 |
| [Go 重构实现计划](docs/go-refactoring-implementation.md) | 详细的重构实现步骤、代码示例 | 后端开发者 |
| [Go 剩余模块计划](docs/go-remaining-modules.md) | 待开发模块清单和计划 | 项目管理者 |

### 📖 文档使用指南

#### 新手入门

1. **首次接触项目**: 阅读 [README.md](README.md) 了解项目概况
2. **搭建开发环境**: 阅读 [运行和部署指南](docs/deployment-guide.md) 的"本地开发环境搭建"章节
3. **了解后端架构**: 阅读 [后端架构功能说明](docs/backend-architecture.md) 了解系统设计
4. **开始开发**: 参考 [Go 后端开发任务清单](docs/go-backend-development-tasks.md) 选择任务

#### 后端开发者

1. **架构理解**: [后端架构功能说明](docs/backend-architecture.md) - 了解 DDD + CQRS 架构
2. **API 开发**: [API 接口文档](backend/docs/api/README.md) - 查看接口规范
3. **任务认领**: [Go 后端开发任务清单](docs/go-backend-development-tasks.md) - 查看待开发任务
4. **重构参考**: [Go 重构实现计划](docs/go-refactoring-implementation.md) - 查看实现细节

#### 前端开发者

1. **环境搭建**: [运行和部署指南](docs/deployment-guide.md) - 搭建前端开发环境
2. **API 对接**: [API 接口文档](backend/docs/api/README.md) - 了解后端接口
3. **架构参考**: [后端架构功能说明](docs/backend-architecture.md) - 了解数据模型

#### 架构师/技术负责人

1. **架构设计**: [后端架构功能说明](docs/backend-architecture.md) - 了解整体架构
2. **重构方案**: [Go 重构设计方案](docs/go-refactoring-design.md) - 了解技术选型
3. **项目规划**: [Go 剩余模块计划](docs/go-remaining-modules.md) - 了解项目规划

#### 运维人员

1. **部署指南**: [运行和部署指南](docs/deployment-guide.md) - Docker 部署和生产环境配置
2. **架构理解**: [后端架构功能说明](docs/backend-architecture.md) - 了解系统组件

### 📁 文档目录结构

```text
docs/
├── backend-architecture.md           # 后端架构功能说明（新增）
├── deployment-guide.md               # 运行和部署指南（新增）
├── go-backend-development-tasks.md   # Go 后端开发任务清单
├── go-refactoring-design.md          # Go 重构设计方案
├── go-refactoring-implementation.md  # Go 重构实现计划
└── go-remaining-modules.md           # Go 剩余模块计划
```

### 🔍 快速查找

| 我想要... | 阅读文档 |
|-----------|----------|
| 了解项目整体架构 | [后端架构功能说明](docs/backend-architecture.md) |
| 搭建本地开发环境 | [运行和部署指南](docs/deployment-guide.md) |
| 使用 Docker 部署 | [运行和部署指南](docs/deployment-guide.md#5-docker-部署) |
| 查看 API 接口 | [API 接口文档](backend/docs/api/README.md) |
| 了解数据库设计 | [后端架构功能说明](docs/backend-architecture.md#9-数据库设计) |
| 查看待开发任务 | [Go 后端开发任务清单](docs/go-backend-development-tasks.md) |
| 了解重构背景 | [Go 重构设计方案](docs/go-refactoring-design.md) |
| 配置生产环境 | [运行和部署指南](docs/deployment-guide.md#6-生产环境部署) |

## 测试

### 后端测试

```bash
cd backend

# 运行所有测试
make test

# 运行单元测试
make test-unit

# 运行集成测试
make test-integration

# 生成覆盖率报告
make test-coverage
```

### 前端测试

```bash
cd frontend

# 运行单元测试
npm run test:unit

# 运行端到端测试
npm run test:e2e
```

## 部署

### 生产环境部署

1. 配置生产环境变量
2. 构建 Docker 镜像
3. 使用 Docker Compose 部署
4. 配置 Nginx 反向代理
5. 配置 SSL 证书

详细部署指南请参考 [部署文档](docs/deployment.md)。

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

MIT License

## 联系方式

- 项目负责人: NetCoreKevin 开发团队
- 技术支持: support@example.com

---

**项目版本**: v1.0.0
**最后更新**: 2026/06/07
**维护者**: NetCoreKevin 开发团队
