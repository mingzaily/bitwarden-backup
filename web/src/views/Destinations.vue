<template>
  <div class="space-y-6">
    <div class="page-header">
      <div>
        <p class="eyebrow">STORAGE TARGETS</p>
        <h2 class="page-title">存储目标</h2>
        <p class="page-subtitle">管理备份文件的落点与保留策略</p>
      </div>
      <button class="btn-primary" type="button" @click="showModal = true">
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M12 5v14m7-7H5" /></svg>
        新建存储目标
      </button>
    </div>

    <div v-if="loading" class="loading-state"><div><div class="spinner mx-auto text-accent"></div><p class="mt-3 text-sm">正在读取存储目标…</p></div></div>
    <div v-else-if="destinations.length === 0" class="empty-state"><div><div class="empty-state-icon"><svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M4 7h16v12H4zM8 7V5a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" /></svg></div><p class="text-sm font-semibold text-main">暂无存储目标</p><p class="mt-1 text-xs text-muted">添加一个本地、WebDAV、S3 或服务器目标。</p></div></div>
    <div v-else class="grid gap-3">
      <article v-for="destination in destinations" :key="destination.id" :class="['resource-card', !destination.enabled ? 'is-disabled' : '']">
        <div :class="['resource-leading', !destination.enabled ? 'is-muted' : '']" aria-hidden="true">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M4 7h16v12H4zM8 7V5a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" /></svg>
        </div>
        <div class="resource-content">
          <div class="resource-title-row">
            <h3 class="resource-title" :title="destination.name">{{ destination.name }}</h3>
            <span :class="getTypeBadgeClass(destination.type)">{{ destination.type_label || getTypeLabel(destination.type) }}</span>
            <span v-if="destination.type !== 'server' && destination.encrypted" class="status-badge status-info" title="已启用 AES-256-GCM 加密">已加密</span>
            <span :class="['status-badge', destination.enabled ? 'status-success' : 'status-neutral']">{{ destination.enabled ? '已启用' : '已停用' }}</span>
          </div>
          <div class="resource-meta"><svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M9 12h6m-6 4h6m2 5H7a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5.6a1 1 0 0 1 .7.3l5.4 5.4a1 1 0 0 1 .3.7V19a2 2 0 0 1-2 2Z" /></svg><span class="mono">{{ destination.display_path || getDestinationPath(destination) }}</span></div>
          <div class="resource-meta"><svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M12 8v4l3 3m6-3a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" /></svg><span>创建于 {{ formatDateTime(destination.created_at) }}</span></div>
        </div>
        <div class="resource-actions">
          <button v-if="destination.type === 'webdav'" class="btn-secondary" type="button" :disabled="testingDestinationId === destination.id" @click="testDestination(destination)">
            <span v-if="testingDestinationId === destination.id" class="spinner"></span>{{ testingDestinationId === destination.id ? '测试中…' : '测试连接' }}
          </button>
          <button :class="destination.enabled ? 'btn-ghost' : 'btn-secondary'" type="button" @click="toggleDestination(destination.id, !destination.enabled)">{{ destination.enabled ? '禁用' : '启用' }}</button>
          <button class="btn-secondary" type="button" @click="editDestination(destination)">编辑</button>
          <button class="btn-danger" type="button" @click="deleteDestination(destination.id)">删除</button>
        </div>
      </article>
    </div>

    <Pagination :page="pagination.page" :page-size="pagination.page_size" :total="pagination.total" :total-page="pagination.total_page" @page-change="handlePageChange" />
    <DestinationModal v-if="showModal" :destination="editingDestination" @close="closeModal" @saved="handleSaved" />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { destinationsApi } from '@/api'
import { useToast } from '@/composables/useToast'
import { useConfirm } from '@/composables/useConfirm'
import DestinationModal from '@/components/features/Destination/DestinationModal.vue'
import Pagination from '@/components/ui/Pagination.vue'

const toast = useToast()
const { confirm } = useConfirm()
const destinations = ref([])
const loading = ref(false)
const testingDestinationId = ref(null)
const showModal = ref(false)
const editingDestination = ref(null)
const pagination = ref({ page: 1, page_size: 10, total: 0, total_page: 0 })

const formatDateTime = (dateStr) => {
  if (!dateStr) return 'N/A'
  const date = new Date(dateStr)
  const pad = (value) => String(value).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}
const getTypeLabel = (type) => ({ local: '本地存储', webdav: 'WebDAV', s3: 'S3', server: '服务器' }[type] || type)
const getDestinationPath = (destination) => {
  switch (destination.type) {
    case 'local': return destination.path || destination.local_path || 'N/A'
    case 'webdav': return destination.webdav_url ? `${destination.webdav_url}${destination.webdav_path || ''}` : (destination.path || 'N/A')
    case 's3': return destination.s3_bucket ? `s3://${destination.s3_bucket}${destination.s3_path || ''}` : (destination.path || 'N/A')
    default: return destination.path || 'N/A'
  }
}
const getTypeBadgeClass = (type) => ({
  local: 'type-badge type-local',
  webdav: 'type-badge type-webdav',
  s3: 'type-badge type-s3',
  server: 'type-badge type-server'
}[type] || 'type-badge type-server')

const loadDestinations = async () => {
  loading.value = true
  try {
    const res = await destinationsApi.getAll({ page: pagination.value.page, page_size: pagination.value.page_size })
    destinations.value = res.data
    pagination.value = res.pagination
  } catch (error) {
    console.error('Failed to load destinations:', error)
    toast.error('加载存储目标失败')
  } finally {
    loading.value = false
  }
}
const handlePageChange = (page) => { pagination.value.page = page; loadDestinations() }
const editDestination = (destination) => { editingDestination.value = destination; showModal.value = true }
const toggleDestination = async (id, enabled) => {
  try { await destinationsApi.setEnabled(id, enabled); toast.success(enabled ? '已启用' : '已禁用'); loadDestinations() }
  catch (error) { console.error('Failed to toggle destination:', error); toast.error('操作失败') }
}
const testDestination = async (destination) => {
  testingDestinationId.value = destination.id
  try {
    await destinationsApi.test(destination.id)
    toast.success('WebDAV 连接正常')
  } catch (error) {
    console.error('Failed to test destination:', error)
    toast.error(error.message || 'WebDAV 连接测试失败')
  } finally {
    testingDestinationId.value = null
  }
}
const deleteDestination = async (id) => {
  const confirmed = await confirm({ title: '删除存储目标', message: '确定要删除这个存储目标吗？此操作不可恢复。', type: 'danger', confirmText: '删除' })
  if (!confirmed) return
  try { await destinationsApi.delete(id); toast.success('存储目标已删除'); loadDestinations() }
  catch (error) { console.error('Failed to delete destination:', error); toast.error('删除失败') }
}
const closeModal = () => { showModal.value = false; editingDestination.value = null }
const handleSaved = () => { closeModal(); loadDestinations() }
onMounted(loadDestinations)
</script>
