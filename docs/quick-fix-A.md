# 快速修复方案 A（规划 + 落地改动）

## 目标

1) 修复 `bw unlock` 返回空 session token 的问题，且不在命令行暴露密码。
2) 提升可观测性：补齐 Bitwarden CLI 调用的 stdout/stderr 日志。
3) 前端按钮区更紧凑：增加分隔线、减少右侧留白。
4) 简化 GORM 关系标注：移除 `foreignKey` 标签，避免误配/冗余。

## 文件清单

- `internal/database/models.go`
- `internal/bitwarden/client.go`
- `internal/bitwarden/export.go`
- `internal/bitwarden/import.go`
- `web/frontend/src/views/Tasks.vue`

## 关键架构决策

### 1) Bitwarden CLI 调用收敛到统一执行器

- 在 `internal/bitwarden/client.go` 增加 `Client.runBW(...)`：
  - 分离采集 stdout/stderr（不再使用 `CombinedOutput()`）
  - 对日志中的敏感参数做脱敏（`--session` / `--password`）
- 统一的执行器为后续加超时、重试、告警埋点预留扩展点。

### 2) Unlock 改造策略（安全 + 兼容）

- **密码通过 stdin 输入**：`bw unlock --raw` + `stdin=masterPassword\n`
  - 避免出现在 `ps`/命令行历史。
- **`bw status` 预检**：
  - `unauthenticated`：直接报错（需要先 login）
  - 其他状态继续尝试 unlock
- **空 token 容错**：
  - unlock stdout 为空时再次 `bw status`，若已 `unlocked`，允许继续（此时后续命令不携带 `--session`）。
  - 目的：兼容 "CLI 已解锁但不返回 raw token" 的情况，避免阻塞备份。

### 3) 前端按钮组（紧凑 + 分隔线）

- 用 `inline-flex + border + divide-x` 形成按钮组：
  - 视觉上有明确分隔线
  - 去掉 `gap` 与单按钮外边框，减少右侧留白且更整齐

### 4) GORM foreignKey 标签移除

- `BackupTask.SourceServer` 依赖字段命名惯例自动推断 `SourceServerID`，移除显式 `foreignKey` 标签。

## 修改步骤（应用补丁顺序建议）

1) 应用 `internal/bitwarden/*` 改动，验证：
   - 锁定态：unlock 后 `sessionToken` 非空
   - 已解锁态：unlock 输出空 token 但 `bw status` 为 `unlocked` 时，export/import 仍可工作（不带 `--session`）
2) 应用 `web/frontend/src/views/Tasks.vue`，确认按钮分隔线与布局无异常。
3) 应用 `internal/database/models.go`，确认 GORM 预加载/关联查询正常。
