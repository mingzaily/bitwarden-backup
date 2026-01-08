[根目录](../../CLAUDE.md) > [web](../) > **frontend**

---

# web/frontend - Vue 3 前端应用

## 变更记录 (Changelog)

### 2026-01-07 01:54:20
- 初始化模块文档

---

## 模块职责

Vue 3 单页应用 (SPA)，负责：
- 提供 Web 管理界面
- 服务器配置管理
- 备份目标配置管理
- 备份任务管理
- 备份日志查看
- 与后端 API 交互

---

## 入口与启动

### 主入口
**文件**: `src/main.js`

**启动流程**:
1. 创建 Vue 应用实例
2. 注册路由 (Vue Router)
3. 挂载到 `#app` DOM 节点

**开发命令**:
```bash
npm install
npm run dev      # 开发服务器 (Vite)
npm run build    # 生产构建
npm run preview  # 预览构建产物
```

---

## 对外接口

### 路由配置
**文件**: `src/router/index.js`

| 路径 | 组件 | 功能 |
|------|------|------|
| `/` | 重定向到 `/servers` | - |
| `/servers` | `Servers.vue` | 服务器配置管理 |
| `/destinations` | `Destinations.vue` | 备份目标管理 |
| `/tasks` | `Tasks.vue` | 备份任务管理 |
| `/logs` | `Logs.vue` | 备份日志查看 |

### API 服务层
**文件**: `src/api/index.js`

**服务模块**:
- `serversApi` - 服务器配置 API
  - `getAll()`, `getById(id)`, `create(data)`, `update(id, data)`, `delete(id)`
- `destinationsApi` - 备份目标 API
  - `getAll()`, `getById(id)`, `create(data)`, `update(id, data)`, `delete(id)`
- `tasksApi` - 备份任务 API
  - `getAll()`, `getById(id)`, `create(data)`, `update(id, data)`, `delete(id)`, `execute(id)`
- `logsApi` - 日志 API
  - `getAll()`

---

## 关键依赖与配置

### 核心依赖
- `vue@^3.4.0` - Vue 3 框架
- `vue-router@^4.2.0` - 路由管理
- `vite@^5.0.0` - 构建工具
- `tailwindcss@^3.4.0` - CSS 框架

### 构建配置
**文件**: `vite.config.js`

- 开发服务器端口: 默认 5173
- 构建输出目录: `dist/`
- API 代理: `/api` → `http://localhost:8080/api`

---

## 数据模型

### 组件结构

```
src/
├── main.js                    # 应用入口
├── App.vue                    # 根组件
├── router/
│   └── index.js               # 路由配置
├── api/
│   └── index.js               # API 服务层
├── composables/
│   └── useToast.js            # Toast 通知 Composable
├── components/
│   ├── ServerModal.vue        # 服务器配置弹窗
│   ├── DestinationModal.vue   # 备份目标弹窗
│   ├── TaskModal.vue          # 任务配置弹窗
│   ├── ToastContainer.vue     # Toast 通知容器
│   ├── ToggleButton.vue       # 开关按钮
│   ├── TabSelector.vue        # 标签选择器
│   ├── CheckboxGroup.vue      # 复选框组
│   └── CustomSelect.vue       # 自定义下拉选择
└── views/
    ├── Servers.vue            # 服务器管理页面
    ├── Destinations.vue       # 备份目标页面
    ├── Tasks.vue              # 任务管理页面
    └── Logs.vue               # 日志查看页面
```

---

## 测试与质量

### 当前状态
- 无单元测试文件
- 无 E2E 测试

### 建议补充
1. 添加 Vitest 单元测试配置
2. 为组件添加测试 (`*.spec.js`)
3. 添加 API 服务层测试 (Mock fetch)
4. 添加 Cypress E2E 测试

---

## 常见问题 (FAQ)

**Q: 如何添加新页面？**
A:
1. 在 `src/views/` 创建新组件
2. 在 `src/router/index.js` 添加路由
3. 在 `App.vue` 的 `tabs` 数组添加导航项

**Q: 如何调用后端 API？**
A: 使用 `src/api/index.js` 中的服务方法：
```javascript
import { serversApi } from '@/api'
const servers = await serversApi.getAll()
```

**Q: 如何显示 Toast 通知？**
A:
```javascript
import { useToast } from '@/composables/useToast'
const toast = useToast()
toast.success('操作成功')
toast.error('操作失败')
```

**Q: 开发时如何解决跨域问题？**
A: Vite 配置中已设置 API 代理，确保后端运行在 8080 端口。

---

## 相关文件清单

```
web/frontend/
├── package.json              # 依赖配置
├── vite.config.js            # Vite 构建配置
├── tailwind.config.js        # Tailwind CSS 配置
├── postcss.config.js         # PostCSS 配置
├── index.html                # HTML 模板
└── src/
    ├── main.js               # 应用入口
    ├── App.vue               # 根组件
    ├── style.css             # 全局样式
    ├── router/               # 路由配置
    ├── api/                  # API 服务层
    ├── composables/          # 组合式函数
    ├── components/           # 可复用组件
    └── views/                # 页面组件
```

---

## 下一步建议

1. 添加表单验证 (Vuelidate / Yup)
2. 添加加载状态指示器
3. 添加分页功能 (日志列表)
4. 添加搜索/过滤功能
5. 添加国际化支持 (i18n)
6. 优化移动端响应式布局
