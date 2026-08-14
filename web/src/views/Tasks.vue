<template>
  <div class="space-y-6">
    <div class="page-header">
      <div>
        <p class="eyebrow">BACKUP TASKS</p>
        <h2 class="page-title">备份任务</h2>
        <p class="page-subtitle">把源站、存储目标和执行计划编排成可追踪的备份任务</p>
      </div>
      <button class="btn-primary" type="button" @click="showModal = true">
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M12 5v14m7-7H5" /></svg>
        新建任务
      </button>
    </div>

    <div v-if="loading" class="loading-state"><div><div class="spinner mx-auto text-accent"></div><p class="mt-3 text-sm">正在读取任务…</p></div></div>
    <div v-else-if="tasks.length === 0" class="empty-state"><div><div class="empty-state-icon"><svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M7 4h10a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2Zm1 5h8M8 13h5M8 17h3" /></svg></div><p class="text-sm font-semibold text-main">暂无备份任务</p><p class="mt-1 text-xs text-muted">创建任务后，可以从这里执行和追踪备份。</p></div></div>
    <div v-else class="grid gap-3">
      <article v-for="task in tasks" :key="task.id" :class="['resource-card', !task.enabled ? 'is-disabled' : '']">
        <div :class="['resource-leading', !task.enabled ? 'is-muted' : '']" aria-hidden="true">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M7 4h10a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2Zm1 5h8M8 13h5M8 17h3" /></svg>
        </div>
        <div class="resource-content">
          <div class="resource-title-row">
            <h3 class="resource-title" :title="task.name">{{ task.name }}</h3>
            <span :class="['type-badge', task.cron_expression ? 'type-webdav' : 'type-server']">{{ task.cron_expression ? '定时' : '手动' }}</span>
            <span :class="['status-badge', task.enabled ? 'status-success' : 'status-neutral']">{{ task.enabled ? '已启用' : '已停用' }}</span>
          </div>
          <div class="resource-meta"><svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M12 8v4l3 3m6-3a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" /></svg><span>{{ task.cron_expression || '手动触发' }}</span></div>
          <div class="resource-meta"><svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2H5a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2Z" /></svg><span>创建于 {{ formatDateTime(task.created_at) }}</span></div>
          <div class="mt-4 border-t border-theme pt-3"><BackupFlow :source-server="task.source_server" :destinations="task.destinations" /></div>
        </div>
        <div class="resource-actions">
          <button class="btn-secondary" type="button" :disabled="!task.enabled" :title="task.enabled ? '立即执行任务' : '请先启用任务'" @click="executeTask(task.id)">立即执行</button>
          <button class="btn-ghost" type="button" @click="toggleTask(task.id, !task.enabled)">{{ task.enabled ? '禁用' : '启用' }}</button>
          <button class="btn-secondary" type="button" @click="editTask(task)">编辑</button>
          <button class="btn-danger" type="button" @click="deleteTask(task.id)">删除</button>
        </div>
      </article>
    </div>

    <Pagination :page="pagination.page" :page-size="pagination.page_size" :total="pagination.total" :total-page="pagination.total_page" @page-change="handlePageChange" />
    <TaskModal v-if="showModal" :task="editingTask" @close="closeModal" @saved="handleSaved" />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { tasksApi } from '@/api'
import { useToast } from '@/composables/useToast'
import { useConfirm } from '@/composables/useConfirm'
import TaskModal from '@/components/features/Task/TaskModal.vue'
import BackupFlow from '@/components/features/Task/BackupFlow.vue'
import Pagination from '@/components/ui/Pagination.vue'

const toast = useToast()
const { confirm } = useConfirm()
const tasks = ref([])
const loading = ref(false)
const showModal = ref(false)
const editingTask = ref(null)
const pagination = ref({ page: 1, page_size: 10, total: 0, total_page: 0 })

const formatDateTime = (dateStr) => {
  if (!dateStr) return 'N/A'
  const date = new Date(dateStr)
  const pad = (value) => String(value).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}
const loadTasks = async () => {
  loading.value = true
  try {
    const res = await tasksApi.getAll({ page: pagination.value.page, page_size: pagination.value.page_size })
    tasks.value = res.data
    pagination.value = res.pagination
  } catch (error) {
    console.error('Failed to load tasks:', error)
    toast.error('加载任务列表失败')
  } finally { loading.value = false }
}
const handlePageChange = (page) => { pagination.value.page = page; loadTasks() }
const editTask = (task) => { editingTask.value = task; showModal.value = true }
const toggleTask = async (id, enabled) => {
  const taskIndex = tasks.value.findIndex(task => task.id === id)
  const originalEnabled = tasks.value[taskIndex]?.enabled
  try {
    if (taskIndex !== -1) tasks.value[taskIndex].enabled = enabled
    await tasksApi.setEnabled(id, enabled)
    toast.success(enabled ? '任务已启用' : '任务已禁用')
    loadTasks()
  } catch (error) {
    console.error('Failed to toggle task:', error)
    if (taskIndex !== -1 && originalEnabled !== undefined) tasks.value[taskIndex].enabled = originalEnabled
    toast.error('操作失败')
  }
}
const executeTask = async (id) => {
  const confirmed = await confirm({ title: '执行备份任务', message: '确定要立即执行此备份任务吗？', type: 'warning', confirmText: '执行' })
  if (!confirmed) return
  try { await tasksApi.execute(id); toast.success('任务已启动，请查看日志') }
  catch (error) { console.error('Failed to execute task:', error); toast.error('任务启动失败') }
}
const deleteTask = async (id) => {
  const confirmed = await confirm({ title: '删除任务', message: '确定要删除这个备份任务吗？此操作不可恢复。', type: 'danger', confirmText: '删除' })
  if (!confirmed) return
  try { await tasksApi.delete(id); toast.success('任务已删除'); loadTasks() }
  catch (error) { console.error('Failed to delete task:', error); toast.error('删除失败') }
}
const closeModal = () => { showModal.value = false; editingTask.value = null }
const handleSaved = () => { closeModal(); loadTasks() }
onMounted(loadTasks)
</script>
