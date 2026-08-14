<template>
  <div class="space-y-6">
    <div class="page-header">
      <div>
        <p class="eyebrow">ACTIVITY</p>
        <h2 class="page-title">运行记录</h2>
        <p class="page-subtitle">查看备份任务的状态、产物和执行过程</p>
      </div>
      <CustomSelect v-model="selectedTaskId" :options="taskOptions" placeholder="全部任务" class="w-52" @update:modelValue="handleTaskChange" />
    </div>

    <div v-if="loading" class="loading-state"><div><div class="spinner mx-auto text-accent"></div><p class="mt-3 text-sm">正在读取日志…</p></div></div>
    <div v-else-if="logs.length === 0" class="empty-state"><div><div class="empty-state-icon"><svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M5 5h14v14H5zM8 9h8M8 13h5" /></svg></div><p class="text-sm font-semibold text-main">暂无运行记录</p><p class="mt-1 text-xs text-muted">任务执行后，状态和备份产物会出现在这里。</p></div></div>
    <div v-else class="grid gap-3">
      <article v-for="log in logs" :key="log.id" class="resource-card">
        <div :class="['resource-leading', log.status === 'failed' ? 'is-danger' : log.status === 'running' ? 'is-info' : '']" aria-hidden="true">
          <svg v-if="log.status === 'success'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="m5 12.5 4 4L19 7" /></svg>
          <svg v-else-if="log.status === 'failed'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M6 18 18 6M6 6l12 12" /></svg>
          <svg v-else fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M12 7v5l3 2m6-2a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" /></svg>
        </div>
        <div class="resource-content">
          <div class="resource-title-row"><span :class="getStatusClass(log.status)">{{ getStatusLabel(log.status) }}</span><span class="resource-title">{{ log.task_name }}</span></div>
          <p v-if="log.message" class="mt-2 text-sm leading-6 text-muted">{{ formatMessage(log.message) }}</p>
          <div v-if="log.status === 'success' && log.backup_file" class="resource-meta"><svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M9 12h6m-6 4h6m2 5H7a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5.6a1 1 0 0 1 .7.3l5.4 5.4a1 1 0 0 1 .3.7V19a2 2 0 0 1-2 2Z" /></svg><span class="mono">{{ log.backup_file }}</span></div>
          <div class="resource-meta"><svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M12 8v4l3 3m6-3a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" /></svg><span>{{ formatTime(log.created_at) }}</span></div>
          <template v-if="log.status === 'failed' && formatMessage(log.message) !== log.message">
            <button class="mt-3 inline-flex items-center gap-1 text-xs font-semibold text-muted hover:text-main" type="button" @click="toggleDetail(log.id)"><svg :class="['h-3 w-3 transition-transform', expandedLogs.has(log.id) ? 'rotate-90' : '']" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="m9 5 7 7-7 7" /></svg>{{ expandedLogs.has(log.id) ? '收起错误' : '查看错误' }}</button>
            <div v-if="expandedLogs.has(log.id)" class="mt-2 border-l-2 border-danger pl-3"><code class="block break-all rounded-lg bg-danger/10 px-3 py-2 text-xs text-danger">{{ log.message }}</code></div>
          </template>
        </div>
        <div class="resource-actions"><button class="btn-secondary" type="button" @click="showLogDetail(log)">查看详情</button></div>
      </article>
      <Pagination :page="pagination.page" :page-size="pagination.pageSize" :total="pagination.total" :total-page="pagination.totalPage" @change="handlePageChange" />
    </div>

    <LogDetailModal v-if="selectedLog" :log="selectedLog" @close="selectedLog = null" />
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { logsApi, tasksApi } from '@/api'
import { useToast } from '@/composables/useToast'
import CustomSelect from '@/components/ui/CustomSelect.vue'
import Pagination from '@/components/ui/Pagination.vue'
import LogDetailModal from '@/components/features/Log/LogDetailModal.vue'

const toast = useToast()
const logs = ref([])
const tasks = ref([])
const loading = ref(false)
const selectedTaskId = ref('')
const expandedLogs = ref(new Set())
const selectedLog = ref(null)
const pagination = ref({ page: 1, pageSize: 10, total: 0, totalPage: 0 })

const toggleDetail = (logId) => {
  if (expandedLogs.value.has(logId)) expandedLogs.value.delete(logId)
  else expandedLogs.value.add(logId)
  expandedLogs.value = new Set(expandedLogs.value)
}
const showLogDetail = (log) => { selectedLog.value = log }
const taskOptions = computed(() => [{ label: '全部任务', value: '' }, ...tasks.value.map(task => ({ label: task.name, value: task.id }))])
const getStatusLabel = (status) => ({ success: '成功', failed: '失败', running: '运行中' }[status] || status)
const getStatusClass = (status) => ({ success: 'status-badge status-success', failed: 'status-badge status-danger', running: 'status-badge status-info' }[status] || 'status-badge status-neutral')
const formatTime = (time) => time ? new Date(time).toLocaleString('zh-CN') : 'N/A'
const formatMessage = (message) => {
  if (!message) return ''
  if (message === 'Backup completed successfully') return '备份成功'
  const errorMappings = [
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

const loadLogs = async () => {
  loading.value = true
  try {
    const params = { page: pagination.value.page, page_size: pagination.value.pageSize }
    if (selectedTaskId.value) params.task_id = selectedTaskId.value
    const res = await logsApi.getAll(params)
    logs.value = res.data || []
    pagination.value = { page: res.pagination?.page || 1, pageSize: res.pagination?.page_size || 10, total: res.pagination?.total || 0, totalPage: res.pagination?.total_page || 0 }
  } catch (error) {
    console.error('Failed to load logs:', error)
    toast.error('加载日志失败')
  } finally { loading.value = false }
}
const handleTaskChange = () => { pagination.value.page = 1; loadLogs() }
const handlePageChange = (page) => { pagination.value.page = page; loadLogs() }
const loadTasks = async () => {
  try { const res = await tasksApi.getAll(); tasks.value = res.data || [] }
  catch (error) { console.error('Failed to load tasks:', error) }
}
onMounted(() => { loadTasks(); loadLogs() })
</script>
