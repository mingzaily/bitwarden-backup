<template>
  <div class="space-y-6">
    <div class="page-header">
      <div>
        <p class="eyebrow">BITWARDEN SOURCES</p>
        <h2 class="page-title">Bitwarden 源站</h2>
        <p class="page-subtitle">管理用于导出和迁移的 Bitwarden 连接</p>
      </div>
      <button class="btn-primary" type="button" @click="showModal = true">
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M12 5v14m7-7H5" /></svg>
        新建源站
      </button>
    </div>

    <div v-if="loading" class="loading-state">
      <div><div class="spinner mx-auto text-accent"></div><p class="mt-3 text-sm">正在读取源站…</p></div>
    </div>
    <div v-else-if="servers.length === 0" class="empty-state">
      <div><div class="empty-state-icon"><svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M5 5.5h14v5H5zM5 13.5h14v5H5zM8 8h.01M8 16h.01" /></svg></div><p class="text-sm font-semibold text-main">暂无 Bitwarden 源站</p><p class="mt-1 text-xs text-muted">点击右上角添加第一个 Bitwarden 源站。</p></div>
    </div>
    <div v-else class="grid gap-3">
      <article v-for="server in servers" :key="server.id" :class="['resource-card', !server.enabled ? 'is-disabled' : '']">
        <div :class="['resource-leading', !server.enabled ? 'is-muted' : '']" aria-hidden="true">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M5 5.5h14v5H5zM5 13.5h14v5H5zM8 8h.01M8 16h.01" /></svg>
        </div>
        <div class="resource-content">
          <div class="resource-title-row">
            <h3 class="resource-title" :title="server.name">{{ server.name }}</h3>
            <ServerTag :is-official="isOfficialServer(server)" />
            <span :class="['status-badge', server.enabled ? 'status-success' : 'status-neutral']">{{ server.enabled ? '已启用' : '已停用' }}</span>
          </div>
          <div class="resource-meta"><svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M21 12a9 9 0 0 1-9 9m9-9a9 9 0 0 0-9-9m9 9H3m9 9a9 9 0 0 1-9-9m9 9c1.66 0 3-4.03 3-9s-1.34-9-3-9m0 18c-1.66 0-3-4.03-3-9s1.34-9 3-9" /></svg><span class="mono">{{ server.server_url || server.url }}</span></div>
          <div class="resource-meta"><svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M16 7a4 4 0 1 1-8 0 4 4 0 0 1 8 0ZM5 21a7 7 0 0 1 14 0" /></svg><span class="mono">Client ID · {{ server.client_id || 'N/A' }}</span></div>
          <div class="resource-meta"><svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M12 8v4l3 3m6-3a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" /></svg><span>创建于 {{ formatDateTime(server.created_at) }}</span></div>
        </div>
        <div class="resource-actions">
          <button :class="server.enabled ? 'btn-ghost' : 'btn-secondary'" type="button" @click="toggleServer(server.id, !server.enabled)">{{ server.enabled ? '禁用' : '启用' }}</button>
          <button class="btn-secondary" type="button" @click="editServer(server)">编辑</button>
          <button class="btn-danger" type="button" @click="deleteServer(server.id)">删除</button>
        </div>
      </article>
    </div>

    <Pagination :page="pagination.page" :page-size="pagination.page_size" :total="pagination.total" :total-page="pagination.total_page" @page-change="handlePageChange" />
    <ServerModal v-if="showModal" :server="editingServer" @close="closeModal" @saved="handleSaved" />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { serversApi } from '@/api'
import { useToast } from '@/composables/useToast'
import { useConfirm } from '@/composables/useConfirm'
import ServerModal from '@/components/features/Server/ServerModal.vue'
import ServerTag from '@/components/features/Server/ServerTag.vue'
import Pagination from '@/components/ui/Pagination.vue'

const toast = useToast()
const { confirm } = useConfirm()
const servers = ref([])
const loading = ref(false)
const showModal = ref(false)
const editingServer = ref(null)
const pagination = ref({ page: 1, page_size: 10, total: 0, total_page: 0 })

const formatDateTime = (dateStr) => {
  if (!dateStr) return 'N/A'
  const date = new Date(dateStr)
  const pad = (value) => String(value).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

const loadServers = async () => {
  loading.value = true
  try {
    const res = await serversApi.getAll({ page: pagination.value.page, page_size: pagination.value.page_size })
    servers.value = res.data
    pagination.value = res.pagination
  } catch (error) {
    console.error('Failed to load servers:', error)
    toast.error('加载源站列表失败')
  } finally {
    loading.value = false
  }
}
const handlePageChange = (page) => { pagination.value.page = page; loadServers() }
const editServer = (server) => { editingServer.value = server; showModal.value = true }
const isOfficialServer = (server) => ['https://vault.bitwarden.com', 'https://vault.bitwarden.eu'].includes(server.server_url || server.url)
const toggleServer = async (id, enabled) => {
  try { await serversApi.setEnabled(id, enabled); toast.success(enabled ? '已启用' : '已禁用'); loadServers() }
  catch (error) { console.error('Failed to toggle server:', error); toast.error('操作失败') }
}
const deleteServer = async (id) => {
  const confirmed = await confirm({ title: '删除源站', message: '确定要删除这个源站吗？此操作不可恢复。', type: 'danger', confirmText: '删除' })
  if (!confirmed) return
  try { await serversApi.delete(id); toast.success('源站已删除'); loadServers() }
  catch (error) { console.error('Failed to delete server:', error); toast.error('删除失败') }
}
const closeModal = () => { showModal.value = false; editingServer.value = null }
const handleSaved = () => { closeModal(); loadServers() }
onMounted(loadServers)
</script>
