# Vue SPA 开发指南

## 项目结构

```
web/
├── frontend/          # Vue 3 前端项目
│   ├── src/
│   │   ├── App.vue           # 主应用组件
│   │   ├── main.js           # 入口文件
│   │   ├── router/           # Vue Router 配置
│   │   ├── views/            # 页面组件
│   │   ├── components/       # 可复用组件
│   │   ├── api/              # API 服务层
│   │   └── composables/      # 组合式函数
│   ├── package.json
│   └── vite.config.js
└── dist/              # 构建产物（自动生成）
```

## 开发模式

### 1. 安装依赖

```bash
cd web/frontend
npm install
```

### 2. 启动开发服务器

```bash
# 终端 1: 启动 Vue 开发服务器（带热重载）
cd web/frontend
npm run dev

# 终端 2: 启动 Go 后端
cd ../..
go run cmd/server/main.go cmd/server/router.go
```

Vue 开发服务器会运行在 `http://localhost:5173`，并自动代理 `/api` 请求到 Go 后端。

## 生产构建

### 方式 1: 使用构建脚本（推荐）

```bash
./build-and-run.sh
```

这个脚本会：
1. 安装前端依赖（如果需要）
2. 构建 Vue 前端到 `web/dist`
3. 启动 Go 服务器

### 方式 2: 手动构建

```bash
# 1. 构建前端
cd web/frontend
npm run build

# 2. 启动后端
cd ../..
go run cmd/server/main.go cmd/server/router.go
```

## 已修复的问题

### ✅ 问题 1: Tab 切换抖动
**原因**: CSS `transition: all` 导致所有属性变化都有动画，引起布局重排。

**解决方案**:
- 只对 `color` 和 `transform` 添加过渡效果
- 见 `App.vue:61` 的 CSS 修改

### ✅ 问题 2: Toast 错误提示双 X 按钮
**原因**: 原 JS 版本 DOM 结构问题。

**解决方案**:
- 使用 `flex items-center`（水平布局）替代 `items-start`
- 见 `ToastContainer.vue:29` 的布局修改

## 路由配置

Vue Router 使用 History 模式，路由映射：

- `/` → 重定向到 `/servers`
- `/servers` → 服务器列表
- `/destinations` → 备份位置
- `/tasks` → 定时任务
- `/logs` → 运行日志

Go 后端已配置 SPA fallback（`router.go:51-70`），支持：
- 浏览器刷新
- 直接访问子路由
- 前进/后退按钮

## 设计风格

项目保持**柔和粗野主义**（Soft Brutalism）风格：

- 边框: `border-2 border-black`
- 阴影: `shadow-brutalist` (4px 4px 0px 0px rgba(0, 0, 0, 1))
- 字体: Inter, font-bold/font-black
- 颜色: brutalist-blue, brutalist-green, brutalist-red

## API 集成

所有 API 调用通过 `src/api/index.js` 统一管理：

```javascript
import { serversApi } from '@/api'

// 获取服务器列表
const servers = await serversApi.getAll()

// 创建服务器
await serversApi.create(data)
```

## Toast 通知

使用全局 Toast 组件：

```javascript
import { useToast } from '@/composables/useToast'

const toast = useToast()
toast.success('操作成功')
toast.error('操作失败')
toast.warning('警告信息')
toast.info('提示信息')
```

## 下一步开发

当前已完成：
- ✅ Vue 3 + Vite 项目结构
- ✅ Vue Router 配置
- ✅ Toast 组件（修复布局问题）
- ✅ Servers 视图完整功能
- ✅ Go 后端 SPA 支持

待完善：
- [ ] Destinations 视图完整功能
- [ ] Tasks 视图完整功能
- [ ] Logs 视图完整功能
- [ ] Modal 组件复用优化
- [ ] 表单验证增强
