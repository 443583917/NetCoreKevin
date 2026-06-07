# NetCoreKevin 开发指南

本文档介绍了NetCoreKevin项目的开发规范、代码生成器使用、自定义AI工具开发和性能优化。

---

## 目录

- [开发规范](#开发规范)
- [代码生成器](#代码生成器)
- [自定义AI工具](#自定义ai工具)
- [性能优化](#性能优化)
- [测试规范](#测试规范)

---

## 开发规范

### 代码风格

- 使用 **PascalCase** 命名类、接口、方法
- 使用 **camelCase** 命名参数、局部变量
- 文件编码使用 **UTF-8 无 BOM**
- 每行代码不超过 120 个字符

### 目录结构规范

```text
Application/Services/
├── [ModuleName]/
│   ├── [ServiceName]Service.cs      # 服务实现
│   └── Dto/
│       ├── [EntityName]Dto.cs       # 数据传输对象
│       └── [EntityName]Input.cs     # 输入参数

Domain/
├── Entities/
│   └── [EntityName].cs              # 实体定义
├── Interfaces/
│   ├── IRepositories/
│   │   └── I[EntityName]Rp.cs       # 仓储接口
│   └── IServices/
│       └── I[ServiceName]Service.cs # 服务接口
```

### 异常处理

- 使用统一的异常处理中间件
- 捕获异常时记录日志（使用 log4net）
- API 返回统一格式的错误响应

### 日志规范

```csharp
// 使用 log4net 记录日志
LogHelper<MyService>.logger.Info("信息级别日志");
LogHelper<MyService>.logger.Warn("警告级别日志");
LogHelper<MyService>.logger.Error("错误级别日志", exception);
```

## 代码生成器

### 功能概述

代码生成器可以根据数据库实体自动生成：

- 仓储接口（IRepository）
- 仓储实现（Repository）
- 服务接口（IService）
- 服务实现（Service）

### 配置说明

编辑 `appsettings.json`：

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

### 使用步骤

1. **获取区域列表**

```bash
GET /api/CodeGenerator/GetAreaNames
```

2. **获取实体列表**

```bash
GET /api/CodeGenerator/GetAreaNameEntityItems?name=App.WebApi.v1
```

3. **生成代码**

```bash
POST /api/CodeGenerator/BulidCode
Content-Type: application/json

[
  {
    "EntityName": "MyEntity",
    "AreaName": "App.WebApi.v1"
  }
]
```

### 注意事项

- 仅超级管理员有权限使用代码生成器
- 生成的代码会覆盖现有文件，请谨慎操作
- 建议在生成前备份相关文件

## 自定义AI工具

### 开发步骤

1. **创建工具接口**

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

2. **实现工具服务**

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

3. **注册工具到容器**

在模块初始化中注册服务：

```csharp
services.AddTransient<IMyCustomToolService, MyCustomToolService>();
```

4. **配置工具到智能体**

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

### 工具开发规范

- 使用 `[Description]` 属性添加参数说明
- 返回类型推荐使用 `Task<string>`
- 错误返回以错误标识符开头的字符串
- 支持的参数类型：string, int, bool, List&lt;string&gt;

## 性能优化

### 数据库优化

| 优化项 | 说明 |
|--------|------|
| 索引优化 | 为常用查询字段创建索引 |
| 分页查询 | 使用 Skip/Take 避免全表扫描 |
| 批量操作 | 使用 EF Core 的批量操作 API |
| 读写分离 | 配置数据库读写分离 |

### 缓存策略

```csharp
// 使用 Redis 缓存
[CacheDataFilter(Duration = 60)] // 缓存60秒
public async Task<MyDto> GetData(int id)
{
    // 查询逻辑
}
```

### 异步编程

- 使用 `async/await` 模式
- 避免同步阻塞调用
- 使用 `ConfigureAwait(false)` 优化性能

### 日志优化

- 生产环境关闭 Debug 级别日志
- 使用异步日志写入
- 定期清理日志文件

## 测试规范

### 单元测试

- 使用 xUnit 作为测试框架
- 使用 Moq 进行依赖注入模拟
- 测试方法命名规范：`方法名_场景_预期结果`
- 每个测试方法只测试一个功能点

### 集成测试

- 使用 WebApplicationFactory 进行 API 测试
- 测试数据库使用内存数据库或测试专用数据库
- 测试前后清理测试数据

### 测试覆盖率

- 核心业务逻辑测试覆盖率要求 ≥ 80%
- API 接口测试覆盖率要求 ≥ 90%
- 使用 coverlet 进行代码覆盖率统计

### 调试技巧

#### 日志调试

```csharp
// 在关键位置添加日志
LogHelper<MyService>.logger.Debug("调试信息: {0}", variable);
```

#### 断点调试

- 在 Visual Studio 中使用条件断点
- 使用日志点（Logpoints）替代断点
- 使用异常设置捕获特定异常

#### 性能调试

- 使用 dotnet-trace 进行性能跟踪
- 使用 dotnet-counters 监控性能计数器
- 使用 MiniProfiler 进行 SQL 和代码性能分析

---

**文档版本**: v1.0  
**最后更新**: 2026/06/07