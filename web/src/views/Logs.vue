<template>
  <div class="space-y-6">
    <div class="page-header">
      <div>
        <p class="eyebrow">ACTIVITY</p>
        <h2 class="page-title">运行记录</h2>
        <p class="page-subtitle">查看备份任务的状态、产物和执行过程</p>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <button class="btn-secondary" type="button" :disabled="loading" title="刷新运行记录" @click="refreshLogs">
          <svg class="h-4 w-4" :class="loading ? 'animate-spin' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M4 4v5h5M20 20v-5h-5M6.7 9A7 7 0 0 1 19 6.5L20 9M4 15l1 2.5A7 7 0 0 0 17.3 15" /></svg>
          {{ loading ? '刷新中…' : '刷新' }}
        </button>
        <CustomSelect v-model="selectedTaskId" :options="taskOptions" placeholder="全部任务" class="w-52" @update:modelValue="handleTaskChange" />
      </div>
    </div>

    <div v-if="loading" class="loading-state"><div><div class="spinner mx-auto text-accent"></div><p class="mt-3 text-sm">正在读取日志…</p></div></div>
    <div v-else-if="logs.length === 0" class="empty-state"><div><div class="empty-state-icon"><svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M5 5h14v14H5zM8 9h8M8 13h5" /></svg></div><p class="text-sm font-semibold text-main">暂无运行记录</p><p class="mt-1 text-xs text-muted">任务执行后，状态和备份产物会出现在这里。</p></div></div>
    <div v-else class="space-y-3">
      <div class="log-selection-toolbar">
        <label class="log-select-all">
          <input
            ref="selectAllCheckbox"
            class="log-checkbox"
            type="checkbox"
            :checked="isPageSelected"
            :aria-checked="isPagePartiallySelected ? 'mixed' : String(isPageSelected)"
            :disabled="selectableLogs.length === 0 || deleting"
            aria-label="选择当前页可删除的运行记录"
            @change="togglePageSelection"
          />
          <span>全选当前页</span>
          <span v-if="selectedLogIds.size" class="log-selection-count">已选 {{ selectedLogIds.size }} 条</span>
        </label>
        <div class="flex flex-wrap items-center gap-2">
          <span v-if="selectableLogs.length < logs.length" class="text-xs text-subtle">运行中的记录不能删除</span>
          <button class="btn-danger" type="button" :disabled="selectedLogIds.size === 0 || deleting" @click="deleteSelectedLogs">
            <svg v-if="!deleting" class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M5 7h14m-9 4v5m4-5v5M9 7V5h6v2m-8 0 1 13h8l1-13" /></svg>
            <span v-else class="spinner"></span>
            {{ deleting ? '删除中…' : '删除选中' }}
          </button>
          <button v-if="selectedLogIds.size" class="btn-ghost" type="button" :disabled="deleting" @click="clearSelection">取消选择</button>
        </div>
      </div>
      <div class="grid gap-3">
        <article v-for="log in logs" :key="log.id" :class="['resource-card', 'log-resource-card', selectedLogIds.has(log.id) ? 'is-selected' : '']">
        <div class="log-status-column">
          <label class="log-select-control" :title="log.status === 'running' ? '运行中的记录不能删除' : '选择运行记录'">
            <input
              class="log-checkbox"
              type="checkbox"
              :checked="selectedLogIds.has(log.id)"
              :disabled="log.status === 'running' || deleting"
              :aria-label="`选择 ${log.task_name} 的运行记录`"
              @change="toggleLogSelection(log.id)"
            />
          </label>
          <div :class="['resource-leading', log.status === 'failed' ? 'is-danger' : log.status === 'running' ? 'is-info' : '']" aria-hidden="true">
          <svg v-if="log.status === 'success'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="m5 12.5 4 4L19 7" /></svg>
          <svg v-else-if="log.status === 'failed'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M6 18 18 6M6 6l12 12" /></svg>
          <svg v-else fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M12 7v5l3 2m6-2a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" /></svg>
          </div>
          <span :class="getStatusClass(log.status)">{{ getStatusLabel(log.status) }}</span>
        </div>
        <div class="resource-content log-card-content">
          <div class="log-card-main">
            <h3 class="log-card-title" :title="log.task_name">{{ log.task_name }}</h3>
            <p v-if="log.message" :class="['log-summary', isErrorSummary(log) ? 'is-danger' : '']">{{ formatMessage(log.message) }}</p>
            <template v-if="log.status === 'failed' && log.message">
              <button class="log-error-toggle" type="button" @click="toggleDetail(log.id)"><svg :class="['h-3 w-3 transition-transform', expandedLogs.has(log.id) ? 'rotate-90' : '']" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="m9 5 7 7-7 7" /></svg>{{ expandedLogs.has(log.id) ? '收起错误' : '查看完整错误' }}</button>
              <div v-if="expandedLogs.has(log.id)" class="log-error-detail"><code>{{ log.message }}</code></div>
            </template>
          </div>
          <div class="log-card-meta">
            <div class="log-card-meta-item">
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M12 8v4l3 3m6-3a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" /></svg>
              <time class="log-card-meta-value" :datetime="log.created_at">{{ formatTime(log.created_at) }}</time>
            </div>
            <div v-if="log.backup_file" class="log-card-meta-item" :title="log.backup_file">
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M9 12h6m-6 4h6m2 5H7a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5.6a1 1 0 0 1 .7.3l5.4 5.4a1 1 0 0 1 .3.7V19a2 2 0 0 1-2 2Z" /></svg>
              <span class="log-card-meta-value mono">{{ formatBackupFileName(log.backup_file) }}</span>
            </div>
          </div>
        </div>
        <div class="resource-actions log-card-actions"><button class="btn-secondary" type="button" @click="showLogDetail(log)">查看详情</button></div>
        </article>
        <Pagination :page="pagination.page" :page-size="pagination.pageSize" :total="pagination.total" :total-page="pagination.totalPage" @change="handlePageChange" />
      </div>
    </div>

    <LogDetailModal v-if="selectedLog" :log="selectedLog" @close="selectedLog = null" />
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { logsApi, tasksApi } from '@/api'
import { useToast } from '@/composables/useToast'
import { useConfirm } from '@/composables/useConfirm'
import CustomSelect from '@/components/ui/CustomSelect.vue'
import Pagination from '@/components/ui/Pagination.vue'
import LogDetailModal from '@/components/features/Log/LogDetailModal.vue'

const toast = useToast()
const { confirm } = useConfirm()
const logs = ref([])
const tasks = ref([])
const loading = ref(false)
const selectedTaskId = ref('')
const expandedLogs = ref(new Set())
const selectedLog = ref(null)
const selectedLogIds = ref(new Set())
const deleting = ref(false)
const selectAllCheckbox = ref(null)
const pagination = ref({ page: 1, pageSize: 10, total: 0, totalPage: 0 })

const toggleDetail = (logId) => {
  if (expandedLogs.value.has(logId)) expandedLogs.value.delete(logId)
  else expandedLogs.value.add(logId)
  expandedLogs.value = new Set(expandedLogs.value)
}
const showLogDetail = (log) => { selectedLog.value = log }
const taskOptions = computed(() => [{ label: '全部任务', value: '' }, ...tasks.value.map(task => ({ label: task.name, value: task.id }))])
const selectableLogs = computed(() => logs.value.filter(log => log.status !== 'running'))
const isPageSelected = computed(() => selectableLogs.value.length > 0 && selectableLogs.value.every(log => selectedLogIds.value.has(log.id)))
const isPagePartiallySelected = computed(() => {
  const selectedCount = selectableLogs.value.filter(log => selectedLogIds.value.has(log.id)).length
  return selectedCount > 0 && !isPageSelected.value
})
const getStatusLabel = (status) => ({ success: '成功', failed: '失败', running: '运行中' }[status] || status)
const getStatusClass = (status) => ({ success: 'status-badge status-success', failed: 'status-badge status-danger', running: 'status-badge status-info' }[status] || 'status-badge status-neutral')
const formatTime = (time) => {
  if (!time) return 'N/A'
  const date = new Date(time)
  if (Number.isNaN(date.getTime())) return 'N/A'
  return date.toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false
  }).replace(/\//g, '-')
}
const formatBackupFileName = (file) => {
  const normalized = String(file || '').replace(/\\/g, '/')
  const parts = normalized.split('/').filter(Boolean)
  return parts[parts.length - 1] || normalized
}
const isErrorSummary = (log) => log.status === 'failed' || /destination errors/i.test(String(log.message || ''))
const formatMessage = (message) => {
  if (!message) return ''
  if (message === 'Backup completed successfully') return '备份成功'
  const partialFailure = message.match(/^Backup completed with destination errors:\s*(.*)$/i)
  if (partialFailure) return `部分目标失败：${partialFailure[1]}`
  const allFailed = message.match(/^all \d+ backup destinations failed:\s*(.*)$/i)
  if (allFailed) return `备份目标全部失败：${allFailed[1]}`
  const errorMappings = [
    { pattern: /source server is disabled/i, text: '源服务器已停用，请先启用源站' },
    { pattern: /no enabled backup destinations/i, text: '没有可用的已启用备份目标' },
    { pattern: /target server is disabled/i, text: '目标服务器已停用，请先启用目标服务器' },
    { pattern: /failed to clean up old backups/i, text: '备份已生成，但清理旧备份失败' },
    { pattern: /unlock returned empty session token/i, text: '解锁失败，请检查主密码是否正确' },
    { pattern: /unlock failed/i, text: '解锁失败，请检查主密码' },
    { pattern: /login failed/i, text: '登录失败，请检查 Client ID 和 Secret' },
    { pattern: /unauthenticated/i, text: '未登录，请检查凭证配置' },
    { pattern: /config server failed/i, text: '服务器配置失败，请检查服务器地址' },
    { pattern: /export failed/i, text: '导出失败' },
    { pattern: /import failed/i, text: '导入失败' },
    { pattern: /vault is not unlocked/i, text: '保险库未解锁' },
    { pattern: /failed to create.*directory/i, text: '创建目录失败，请检查路径权限' },
    { pattern: /bw status failed/i, text: 'Bitwarden CLI 状态检查失败' }
  ]
  for (const { pattern, text } of errorMappings) if (pattern.test(message)) return text
  return message
}

const clearSelection = () => { selectedLogIds.value = new Set() }
const toggleLogSelection = (logId) => {
  const next = new Set(selectedLogIds.value)
  if (next.has(logId)) next.delete(logId)
  else next.add(logId)
  selectedLogIds.value = next
}
const togglePageSelection = () => {
  const next = new Set(selectedLogIds.value)
  if (isPageSelected.value) selectableLogs.value.forEach(log => next.delete(log.id))
  else selectableLogs.value.forEach(log => next.add(log.id))
  selectedLogIds.value = next
}
const deleteSelectedLogs = async () => {
  if (selectedLogIds.value.size === 0 || deleting.value) return
  const selectedCount = selectedLogIds.value.size
  const confirmed = await confirm({
    title: '删除运行记录',
    message: `确定删除选中的 ${selectedCount} 条运行记录吗？只会删除记录，不会删除备份文件。`,
    type: 'danger',
    confirmText: '删除'
  })
  if (!confirmed) return

  deleting.value = true
  try {
    const res = await logsApi.deleteMany([...selectedLogIds.value])
    const deletedCount = Number(res?.deleted || 0)
    clearSelection()
    if (deletedCount > 0) toast.success(`已删除 ${deletedCount} 条运行记录`)
    else toast.success('选中的运行记录已不存在或正在运行')
    await loadLogs()
  } catch (error) {
    console.error('Failed to delete logs:', error)
    toast.error(error.message || '删除运行记录失败')
  } finally {
    deleting.value = false
  }
}

const loadLogs = async () => {
  loading.value = true
  try {
    const params = { page: pagination.value.page, page_size: pagination.value.pageSize }
    if (selectedTaskId.value) params.task_id = selectedTaskId.value
    const res = await logsApi.getAll(params)
    const responsePagination = res.pagination || {}
    const totalPage = responsePagination.total_page || 0
    const requestedPage = pagination.value.page
    if ((totalPage === 0 && requestedPage !== 1) || (totalPage > 0 && requestedPage > totalPage)) {
      pagination.value.page = totalPage > 0 ? totalPage : 1
      return loadLogs()
    }
    logs.value = res.data || []
    pagination.value = { page: responsePagination.page || 1, pageSize: responsePagination.page_size || 10, total: responsePagination.total || 0, totalPage }
    selectedLogIds.value = new Set([...selectedLogIds.value].filter(id => logs.value.some(log => log.id === id)))
  } catch (error) {
    console.error('Failed to load logs:', error)
    toast.error('加载日志失败')
  } finally { loading.value = false }
}
const handleTaskChange = () => { clearSelection(); pagination.value.page = 1; loadLogs() }
const handlePageChange = (page) => { clearSelection(); pagination.value.page = page; loadLogs() }
const refreshLogs = () => { clearSelection(); loadLogs() }
const loadTasks = async () => {
  try { const res = await tasksApi.getAll(); tasks.value = res.data || [] }
  catch (error) { console.error('Failed to load tasks:', error) }
}
onMounted(() => { loadTasks(); loadLogs() })
watch([isPageSelected, isPagePartiallySelected], () => {
  if (selectAllCheckbox.value) selectAllCheckbox.value.indeterminate = isPagePartiallySelected.value
}, { immediate: true })
</script>
