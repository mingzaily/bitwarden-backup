# TODO - 代码重构计划

## 📋 重构目标

完整重构代码架构，采用**实用主义分层架构**（避免 DDD 复杂性）

---

## 🎯 目标架构

```
internal/
├── model/                  # 数据模型（纯数据结构）
│   ├── server.go          # 服务器模型
│   ├── destination.go     # 目标模型
│   ├── task.go            # 任务模型
│   └── log.go             # 日志模型
│
├── repository/             # 数据访问层（封装 GORM）
│   ├── server.go          # 服务器仓储
│   ├── destination.go     # 目标仓储
│   ├── task.go            # 任务仓储
│   └── log.go             # 日志仓储
│
├── service/                # 业务逻辑层
│   ├── server/
│   │   └── service.go     # 服务器业务逻辑
│   ├── destination/
│   │   └── service.go     # 目标业务逻辑
│   ├── task/
│   │   └── service.go     # 任务业务逻辑
│   └── backup/
│       └── service.go     # 备份业务逻辑
│
├── handler/                # HTTP 处理层
│   ├── server/
│   │   └── handler.go     # 服务器 API
│   ├── destination/
│   │   └── handler.go     # 目标 API
│   ├── task/
│   │   └── handler.go     # 任务 API
│   └── log/
│       └── handler.go     # 日志 API
│
├── scheduler/              # 调度层
│   ├── scheduler.go       # 调度器核心
│   ├── executor/          # 执行器
│   │   ├── executor.go
│   │   └── immediate.go
│   └── backup/            # 备份实现
│       ├── local.go
│       ├── webdav.go
│       └── server.go
│
├── database/               # 数据库基础设施
│   ├── db.go              # 数据库连接
│   └── migrate.go         # 迁移
│
├── crypto/                 # 加密（保持不变）
├── config/                 # 配置（保持不变）
├── bitwarden/              # Bitwarden 集成（保持不变）
└── webdav/                 # WebDAV 客户端（保持不变）
```

---

## 📐 各层职责

### 1. Model 层（数据模型）
- ✅ 纯粹的数据结构
- ✅ 包含 GORM 标签
- ✅ 包含 JSON 标签
- ❌ **不包含业务逻辑**

### 2. Repository 层（数据访问）
- ✅ 封装所有数据库操作
- ✅ 提供 CRUD 接口
- ✅ 处理查询逻辑
- ❌ **不包含业务逻辑**

### 3. Service 层（业务逻辑）
- ✅ 核心业务逻辑
- ✅ 调用 Repository
- ✅ 事务管理
- ✅ 业务规则验证

### 4. Handler 层（HTTP 接口）
- ✅ 处理 HTTP 请求/响应
- ✅ 参数验证
- ✅ 调用 Service
- ❌ **不包含业务逻辑**

### 5. Scheduler 层（任务调度）
- ✅ 定时任务管理
- ✅ 调用 Service 执行业务
- ✅ 备份实现

---

## 🚀 重构步骤（分 5 个阶段）

### ⏳ 阶段 1: 创建 Model 层（预计 30 分钟）

**任务：**
- [ ] 创建 `internal/model/` 目录
- [ ] 从 `database/models.go` 提取 `ServerConfig` 到 `model/server.go`
- [ ] 从 `database/models.go` 提取 `BackupDestination` 到 `model/destination.go`
- [ ] 从 `database/models.go` 提取 `BackupTask` 到 `model/task.go`
- [ ] 从 `database/models.go` 提取 `BackupLog` 到 `model/log.go`
- [ ] 更新所有导入路径
- [ ] 测试编译通过

**风险：** 低
**优先级：** 高

---

### ⏳ 阶段 2: 创建 Repository 层（预计 1 小时）

**任务：**
- [ ] 创建 `internal/repository/` 目录
- [ ] 创建 `repository/server.go` - 封装服务器数据访问
  - `Create(server *model.ServerConfig) error`
  - `Update(id uint, server *model.ServerConfig) error`
  - `Delete(id uint) error`
  - `FindByID(id uint) (*model.ServerConfig, error)`
  - `FindAll() ([]*model.ServerConfig, error)`
- [ ] 创建 `repository/destination.go` - 封装目标数据访问
- [ ] 创建 `repository/task.go` - 封装任务数据访问
- [ ] 创建 `repository/log.go` - 封装日志数据访问
- [ ] 从 Handler 中提取数据访问代码到 Repository
- [ ] 测试编译通过

**风险：** 中
**优先级：** 高

---

### ⏳ 阶段 3: 创建 Service 层（预计 1.5 小时）

**任务：**
- [ ] 创建 `internal/service/` 目录
- [ ] 创建 `service/server/service.go` - 服务器业务逻辑
  - 从 Handler 提取业务逻辑
  - 调用 Repository
  - 添加事务管理
- [ ] 创建 `service/destination/service.go` - 目标业务逻辑
- [ ] 创建 `service/task/service.go` - 任务业务逻辑
- [ ] 创建 `service/backup/service.go` - 备份业务逻辑
- [ ] 更新 Handler 调用 Service
- [ ] 测试编译通过

**风险：** 中
**优先级：** 高

---

### ⏳ 阶段 4: 重构 Handler 层（预计 1 小时）

**任务：**
- [ ] 创建 `internal/handler/server/` 目录
- [ ] 移动 `handler/server.go` 和 `handler/server_crud.go` 到 `handler/server/handler.go`
- [ ] 创建 `internal/handler/destination/` 目录
- [ ] 移动 `handler/destination.go` 和 `handler/destination_crud.go` 到 `handler/destination/handler.go`
- [ ] 创建 `internal/handler/task/` 目录
- [ ] 移动 `handler/task.go`、`handler/task_crud.go`、`handler/task_execute.go` 到 `handler/task/handler.go`
- [ ] 创建 `internal/handler/log/` 目录
- [ ] 移动 `handler/log.go` 到 `handler/log/handler.go`
- [ ] 简化 Handler 代码（仅保留 HTTP 处理逻辑）
- [ ] 更新路由注册
- [ ] 测试编译通过

**风险：** 中
**优先级：** 中

---

### ⏳ 阶段 5: 重构 Scheduler 层（预计 1 小时）

**任务：**
- [ ] 创建 `internal/scheduler/executor/` 目录
- [ ] 移动 `scheduler/execute.go` 到 `scheduler/executor/executor.go`
- [ ] 移动 `scheduler/execute_now.go` 到 `scheduler/executor/immediate.go`
- [ ] 创建 `internal/scheduler/backup/` 目录
- [ ] 移动 `scheduler/backup.go` 到 `scheduler/backup/backup.go`
- [ ] 移动 `scheduler/backup_local.go` 到 `scheduler/backup/local.go`
- [ ] 移动 `scheduler/backup_webdav.go` 到 `scheduler/backup/webdav.go`
- [ ] 移动 `scheduler/backup_server.go` 到 `scheduler/backup/server.go`
- [ ] 移动 `scheduler/multi_destination.go` 到 `scheduler/backup/multi.go`
- [ ] 移动 `scheduler/destination_handler.go` 到 `scheduler/backup/handler.go`
- [ ] 保留 `scheduler/scheduler.go` 和 `scheduler/task.go` 在根目录
- [ ] 删除 `scheduler/migrate.go`（已废弃）
- [ ] 测试编译通过

**风险：** 低
**优先级：** 低

---

## 📊 总体评估

**总工作量：** 约 5 小时
**总风险：** 中等
**建议时机：** 功能稳定后，有充足时间时进行

---

## ⚠️ 注意事项

1. **每个阶段独立测试**
   - 每完成一个阶段，必须编译通过
   - 运行基本功能测试
   - 确保没有破坏现有功能

2. **保持向后兼容**
   - 不改变 API 接口
   - 不改变数据库结构
   - 不改变配置格式

3. **渐进式重构**
   - 可以随时停止
   - 每个阶段都是可用状态
   - 不影响正常开发

4. **代码审查**
   - 每个阶段完成后进行代码审查
   - 确保符合架构设计
   - 检查是否有遗漏

---

## 📝 当前状态

- ✅ 架构设计完成
- ⏳ 等待合适时机开始重构
- 📅 建议在功能稳定后进行

---

## 🔗 相关文档

- [SECURITY.md](SECURITY.md) - 安全文档
- [README.md](README.md) - 项目说明
- [VUE_DEVELOPMENT.md](VUE_DEVELOPMENT.md) - 前端开发指南
