# 备份保留策略功能实施计划

## 需求概述

为 local、webdav、s3 三种备份目标类型添加备份保留策略功能：
- 新增 `max_backup_count` 参数指定最大保存数量
- 实现循环覆盖机制（达到上限时删除最旧备份）
- 参数可选，未设置（0）时保留所有备份
- 作用域：按 destination 全局

## 技术方案

**后端**：存储端列举+删除
**前端**：Toggle + 条件输入

---

## 实施步骤

### 阶段 1：数据模型修改

**文件**: `internal/model/destination.go`

1. `BackupDestination` 结构体新增字段：
```go
MaxBackupCount int `gorm:"default:0" json:"max_backup_count"`
```

2. `DestinationResponse` 新增字段：
```go
MaxBackupCount int `json:"max_backup_count"`
```

3. `ToResponse()` 方法添加字段映射

---

### 阶段 2：Provider 接口扩展

**文件**: `internal/provider/interface.go`

新增 `RetentionProvider` 接口：
```go
type RetentionProvider interface {
    // Cleanup 清理超出保留数量的旧备份
    // maxCount: 最大保留数量，0 表示不限制
    // 返回删除的文件数量
    Cleanup(ctx BackupContext, maxCount int) (int, error)
}
```

---

### 阶段 3：Local Provider 实现

**文件**: `internal/provider/local.go`

实现 `Cleanup` 方法：
1. 读取目录下所有 `backup_*.json` 文件
2. 按修改时间排序（最新在前）
3. 删除超出 maxCount 的旧文件

---

### 阶段 4：WebDAV Provider 实现

**文件**: `internal/provider/webdav.go`
**文件**: `internal/webdav/client.go`

1. WebDAV Client 新增方法：
   - `ListFiles(path string)` - PROPFIND 列举文件
   - `Delete(path string)` - DELETE 删除文件

2. WebDAV Provider 实现 `Cleanup` 方法

---

### 阶段 5：S3 Provider 实现

**文件**: `internal/provider/s3.go`

实现 `Cleanup` 方法：
1. 使用 `ListObjectsV2` 列举对象
2. 按 `LastModified` 排序
3. 使用 `DeleteObjects` 批量删除

---

### 阶段 6：调度层集成

**文件**: `internal/scheduler/multi_destination.go`

备份成功后调用清理逻辑：
```go
if dest.MaxBackupCount > 0 {
    if rp, ok := provider.(RetentionProvider); ok {
        deleted, err := rp.Cleanup(ctx, dest.MaxBackupCount)
        // 记录日志，错误不影响备份结果
    }
}
```

---

### 阶段 7：前端实现

**文件**: `web/src/components/features/Destination/DestinationModal.vue`

1. `formData` 新增 `max_backup_count` 字段
2. 新增 `retentionEnabled` 状态控制
3. UI：Toggle + 条件数字输入框
4. 表单提交处理逻辑

---

## 文件修改清单

| 序号 | 文件 | 修改类型 | 说明 |
|------|------|----------|------|
| 1 | `internal/model/destination.go` | 修改 | 新增 MaxBackupCount 字段 |
| 2 | `internal/provider/interface.go` | 修改 | 新增 RetentionProvider 接口 |
| 3 | `internal/provider/local.go` | 修改 | 实现 Cleanup 方法 |
| 4 | `internal/webdav/client.go` | 修改 | 新增 ListFiles/Delete 方法 |
| 5 | `internal/provider/webdav.go` | 修改 | 实现 Cleanup 方法 |
| 6 | `internal/provider/s3.go` | 修改 | 实现 Cleanup 方法 |
| 7 | `internal/scheduler/multi_destination.go` | 修改 | 集成清理逻辑 |
| 8 | `web/src/.../DestinationModal.vue` | 修改 | 前端表单 |

---

## 风险与注意事项

1. **WebDAV 兼容性**：不同服务器对 PROPFIND 支持不一致，需容错处理
2. **并发安全**：同一 destination 并发备份时需考虑竞态
3. **错误处理**：清理失败不应影响备份成功状态，仅记录警告日志
4. **向后兼容**：max_backup_count=0 表示不限制，保持现有行为
