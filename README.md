# NetCoreKevin

> 基于 .NET 9 的企业级 AI 智能体 SaaS 平台，集成知识库管理、技能调度、本地模型支持、多租户架构和微服务治理，支持 RAG 检索增强和智能对话。

---

## 技术亮点

| 核心技术 | 说明 | 优势 |
| --- | --- | --- |
| **.NET 9** | 最新 LTS 版本 | 高性能、新特性支持 |
| **DDD 架构** | 领域驱动设计 | 模块化、易维护扩展 |
| **微服务** | Consul + CAP + Hangfire | 服务解耦、分布式事务 |
| **AI 集成** | SemanticKernel + MCP | 多模型支持、协议标准化 |
| **RAG 检索** | Qdrant 向量数据库 | 知识库问答、语义搜索 |
| **多租户** | 一库多租户架构 | 数据隔离、资源复用 |

---

## 快速开始

### 环境要求

- .NET SDK 9.0+
- MySQL 8.0+
- Redis 7.0+
- Qdrant 1.7+（AI 功能）

### 配置步骤

**1. 配置数据库连接**

编辑 `App/WebApi/appsettings.json`：

```json
{
  "ConnectionStrings": {
    "dbConnection": "server=127.0.0.1;port=3306;database=kevin_app;user id=root;password=admin123",
    "redisConnection": "127.0.0.1:6379,DefaultDatabase=0,password=123456"
  }
}
```

**2. 初始化数据库**

在 **程序包管理控制台** 执行：

```powershell
# 选择 Kevin.EntityFrameworkCore 项目
Add-Migration "初始化数据库"
Update-Database
```

**3. 启动应用**

```bash
cd App/WebApi
dotnet run --environment Development
```

**4. 验证服务**

| 服务 | 地址 | 验证方法 |
| --- | --- | --- |
| API | http://localhost:9901 | 访问根路径返回 200 |
| Swagger | http://localhost:9901/swagger | 查看 API 文档 |
| Hangfire | http://localhost:9901/pchangfire | 查看任务面板 |

**5. 默认账户**

| 用户名 | 密码 | 租户 | 说明 |
| --- | --- | --- | --- |
| admin | 123456 | 1000 | 管理员账户 |

---

## 功能模块状态

| 模块 | 功能 | 状态 | 说明 |
| --- | --- | --- | --- |
| 用户管理 | 用户 CRUD、权限绑定 | ✅ | 已完成 |
| 角色管理 | 角色 CRUD、权限配置 | ✅ | 已完成 |
| 权限管理 | 菜单权限、API 权限 | ✅ | 已完成 |
| AI 智能体 | 智能对话、工具调用 | ✅ | 已完成 |
| 知识库 | 文档管理、RAG 检索 | ✅ | 已完成 |
| 任务调度 | Hangfire 定时任务 | ✅ | 已完成 |
| 消息服务 | 钉钉消息推送 | ✅ | 已完成 |
| 文件存储 | 多云存储支持 | ✅ | 已完成 |
| 代码生成器 | CRUD 代码生成 | 🔄 | 开发中 |
| 多语言支持 | 国际化 | ❌ | 待开发 |

---

## 文档导航

### 核心文档

| 文档 | 说明 | 适用人群 |
| --- | --- | --- |
| [架构设计](docs/architecture.md) | 系统架构、技术栈、模块设计 | 架构师、高级开发 |
| [使用指南](docs/usage.md) | 环境配置、安装部署、API 使用 | 运维人员、测试人员 |
| [开发指南](docs/development.md) | 开发规范、代码生成器、性能优化 | 开发人员 |
| [运维指南](docs/operations.md) | 部署架构、监控运维、故障排查 | 运维人员 |
| [术语表](docs/glossary.md) | 项目术语统一定义 | 所有人员 |

### 外部资源

- **系统文档**: [SYSTEM_DOCUMENTATION.md](SYSTEM_DOCUMENTATION.md)
- **项目架构**: [PROJECT_ARCHITECTURE.md](PROJECT_ARCHITECTURE.md)
- **教学文档**: [CSDN 专栏](https://blog.csdn.net/weixin_42629287/category_13037923.html)
- **新项目教程**: [基于 NetCoreKevin 二次开发](https://gitee.com/netkevin-li/ainet)

---

## 项目结构

```text
NetCoreKevin/
├── App/                          # 业务应用模块
│   ├── Application/              # 应用服务层
│   ├── Domain/                   # 领域层
│   ├── RepositorieRps/           # 仓储实现
│   └── WebApi/                   # API 入口
├── Kevin/                        # 核心框架模块
│   ├── Application/              # 核心服务
│   │   └── Services/AI/          # AI 相关服务
│   ├── Domain/                   # 核心领域模型
│   ├── Kevin.EntityFrameworkCore/  # EF Core 实现
│   └── Kevin.Web.Basics/         # Web 基础组件
├── docs/                         # 项目文档
│   ├── architecture.md           # 架构设计
│   ├── usage.md                  # 使用指南
│   ├── development.md            # 开发指南
│   ├── operations.md             # 运维指南
│   └── glossary.md               # 术语表
├── Doc/                          # 文档资源
└── InitData/                     # 初始化数据
```

---

## 功能效果图

### AI 智能体技能工具管理

| 动态管理 | 智能体配置 | 对话交互 |
| --- | --- | --- |
| ![动态管理](Doc/Img/list1/1.png) | ![智能体配置](Doc/Img/list1/2.png) | ![对话交互](Doc/Img/list1/3.png) |

### AI 知识库系统

| 知识库管理 | 文档上传 | 智能问答 |
| --- | --- | --- |
| ![知识库管理](Doc/Img/list2/1.png) | ![文档上传](Doc/Img/list2/2.png) | ![智能问答](Doc/Img/list2/3.png) |

### AI 智能体技能

| 技能列表 | 技能配置 | 技能执行 |
| --- | --- | --- |
| ![技能列表](Doc/Img/list3/1.png) | ![技能配置](Doc/Img/list3/2.png) | ![技能执行](Doc/Img/list3/3.png) |

### Hangfire 任务调度

| 任务列表 | 任务配置 | 执行记录 |
| --- | --- | --- |
| ![任务列表](Doc/Img/list4/1.png) | ![任务配置](Doc/Img/list4/2.png) | ![执行记录](Doc/Img/list4/3.png) |

### 后台管理系统 (Vue3 + AntDesign)

| 用户管理 | 角色管理 | 权限管理 | 系统配置 |
| --- | --- | --- | --- |
| ![用户管理](Doc/Img/list5/1.png) | ![角色管理](Doc/Img/list5/2.png) | ![权限管理](Doc/Img/list5/3.png) | ![系统关键日志](Doc/Img/list5/4.png) |

---

## 交流社区

| 微信 | 交流群 |
| --- | --- |
| ![微信](Doc/wx.jpeg) | ![交流群](Doc/wx_jiaoliuqun.JPG) |

---

## Star History

<a href="https://www.star-history.com/?repos=junkai-li/NetCoreKevin&type=timeline">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/chart?repos=junkai-li/NetCoreKevin&type=timeline&theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/chart?repos=junkai-li/NetCoreKevin&type=timeline" />
    <img alt="Star History Chart" src="https://api.star-history.com/chart?repos=junkai-li/NetCoreKevin&type=timeline" />
  </picture>
</a>

---

**版本**: v1.0  
**License**: MIT  
**维护者**: NetCoreKevin 开发团队
