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

---

# 方案 B：深度优化项（2026-01-08）

## 🎯 优化目标

解决系统稳定性和可观测性的深层问题，提升整体质量。

---

## 📋 优化清单

### 🔒 优先级 1：Bitwarden CLI 隔离（预计 2 小时）

**问题**：多任务共享 `~/.config/Bitwarden CLI` 全局状态，导致并发竞争和状态污染。

**任务**：
- [ ] **方案 A：全局互斥锁**（推荐，简单）
  - 在 `internal/scheduler/` 添加全局 `sync.Mutex`
  - 所有 `bw` 命令调用前加锁，完成后解锁
  - 确保任务串行执行，避免状态冲突

- [ ] **方案 B：独立数据目录**（复杂，彻底隔离）
  - 为每个任务创建独立的 CLI 数据目录
  - 使用环境变量 `BITWARDENCLI_APPDATA_DIR` 指定目录
  - 任务完成后清理临时目录

**风险**：中等
**优先级**：高
**预期收益**：彻底解决并发竞争问题

---

### 📊 优先级 2：日志系统重构（预计 2 小时）

**问题**：
1. `BackupLog.Message` 只存 `err.Error()`，缺少完整输出
2. 前端日志展示无换行，可读性差
3. 缺少分步骤日志（login/unlock/export/upload）

**任务**：
- [ ] **后端：扩展 BackupLog 模型**
  - 添加 `Details` 字段（TEXT 类型）存储 JSON 格式的详细日志
  - JSON 结构：`{"steps": [{"name": "login", "stdout": "...", "stderr": "...", "exit_code": 0}]}`
  - 敏感信息脱敏（密码、token 等）

- [ ] **后端：修改日志记录逻辑**
  - 在 `internal/scheduler/multi_destination.go` 记录每个步骤
  - 捕获所有 `bw` 命令的 stdout/stderr
  - 失败时记录完整错误上下文

- [ ] **前端：优化日志展示**
  - 修改 `web/frontend/src/views/Logs.vue`
  - 使用 `<pre class="whitespace-pre-wrap break-words">` 展示日志
  - 添加可折叠的详细日志区域
  - 修复 `task_name` 字段缺失问题（后端返回或前端映射）

**风险**：低
**优先级**：高
**预期收益**：大幅提升问题诊断效率

---

### 🗑️ 优先级 3：删除链路完善（预计 1.5 小时）

**问题**：删除 `ServerConfig/BackupTask/BackupDestination` 时未清理关联数据，导致孤儿记录。

**任务**：
- [ ] **删除服务器时清理关联**
  - 检查是否有任务引用该服务器
  - 如有引用，提示用户或级联删除
  - 清理相关日志记录

- [ ] **删除任务时清理关联**
  - 清理 `task_destinations` 多对多关联表
  - 清理相关日志记录

- [ ] **删除目标时清理关联**
  - 清理 `task_destinations` 多对多关联表
  - 检查是否有任务引用该目标

- [ ] **添加数据巡检工具**
  - 创建 `cmd/tools/cleanup-orphans/` 工具
  - 扫描并清理孤儿记录
  - 生成清理报告

**风险**：中等（需要仔细测试删除逻辑）
**优先级**：中
**预期收益**：保持数据库整洁，避免数据膨胀

---

### 🎨 优先级 4：前端日志优化（预计 30 分钟）

**问题**：日志列表展示简陋，缺少交互性。

**任务**：
- [ ] **添加日志筛选功能**
  - 按任务名称筛选
  - 按状态筛选（成功/失败）
  - 按时间范围筛选

- [ ] **添加日志详情展开**
  - 点击日志条目展开详细信息
  - 显示分步骤日志
  - 支持复制日志内容

- [ ] **优化日志列表样式**
  - 成功日志绿色边框
  - 失败日志红色边框
  - 添加时间戳格式化

**风险**：低
**优先级**：低
**预期收益**：提升用户体验

---

## 📊 总体评估

**总工作量**：约 6 小时
**总风险**：中等
**建议时机**：方案 A 完成并稳定运行后

---

## ⚠️ 实施建议

1. **优先级顺序**：
   - 先做优先级 1（CLI 隔离）- 解决核心稳定性问题
   - 再做优先级 2（日志系统）- 提升可观测性
   - 最后做优先级 3-4 - 完善细节

2. **测试策略**：
   - 每个优化项独立测试
   - 重点测试并发场景
   - 验证日志完整性

3. **回滚方案**：
   - 保留原有代码备份
   - 可随时回退到方案 A

---

## 📝 当前状态

- ✅ 方案 A（快速修复）规划完成
- ✅ 方案 B（深度优化）规划完成
- ⏳ 等待方案 A 实施完成
- 📅 方案 B 建议在方案 A 稳定后进行

