# NetCoreKevin 术语表

本文档定义了NetCoreKevin项目中使用的统一术语，确保文档一致性。

---

## 技术术语

| 术语 | 英文 | 说明 |
|------|------|------|
| 智能体 | Agent | 基于AI的代理系统，支持多步推理与任务自动化 |
| 知识库 | Knowledge Base | 使用RAG检索增强的文档管理系统 |
| 向量化 | Vectorization | 将文本转换为向量表示的过程 |
| 仓储 | Repository | 数据访问层，负责CRUD操作 |
| 领域 | Domain | 核心业务逻辑层 |
| 应用服务 | Application Service | 业务逻辑编排层 |
| 实体 | Entity | 领域模型中的核心对象 |
| 值对象 | Value Object | 没有唯一标识的领域对象 |
| 聚合 | Aggregate | 一组相关对象的集合 |
| 领域事件 | Domain Event | 领域中发生的有意义的事件 |

## 业务术语

| 术语 | 英文 | 说明 |
|------|------|------|
| 租户 | Tenant | 多租户架构中的独立业务单元 |
| 权限 | Permission | 用户对资源的访问控制 |
| 角色 | Role | 用户权限的集合 |
| 技能 | Skill | AI智能体可调用的功能模块 |
| 工具 | Tool | AI智能体可调用的外部服务 |

## 模块术语

| 术语 | 英文 | 说明 |
|------|------|------|
| CAP | CAP | 分布式事件总线框架 |
| Hangfire | Hangfire | 任务调度框架 |
| SemanticKernel | SemanticKernel | AI框架 |
| Qdrant | Qdrant | 向量数据库 |
| Ollama | Ollama | 本地AI模型运行框架 |

---

**文档版本**: v1.0  
**最后更新**: 2026/06/07