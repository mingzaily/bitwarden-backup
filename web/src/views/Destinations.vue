<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h2 class="text-lg font-black text-gray-900">备份目标</h2>
        <p class="text-sm text-gray-700 font-bold">管理备份文件存储目标</p>
      </div>
      <button
        @click="showModal = true"
        class="btn-brutalist inline-flex items-center px-5 py-2.5 text-sm font-bold rounded-lg text-white"
      >
        <svg class="mr-2 -ml-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
        </svg>
        新建目标
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-brutalist-blue border-t-transparent"></div>
      <p class="mt-2 text-sm text-gray-600">加载中...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="destinations.length === 0" class="text-center py-12 bg-brutalist-cream/20 rounded-lg border-2 border-black">
      <p class="text-gray-600">暂无备份目标，点击上方按钮添加</p>
    </div>

    <!-- Destinations List -->
    <div v-else class="grid gap-4">
      <div
        v-for="destination in destinations"
        :key="destination.id"
        :class="[
          'bg-white overflow-hidden rounded-lg border-2 border-black shadow-brutalist hover:shadow-brutalist-hover transition-all',
          'border-l-4',
          destination.enabled ? 'border-l-brutalist-green' : 'border-l-gray-400',
          !destination.enabled && 'opacity-50'
        ]"
      >
        <div class="px-6 py-4 bg-brutalist-cream/20 min-h-[90px]">
          <div class="flex items-center justify-between">
            <!-- 左侧：目标信息 -->
            <div class="flex-1 space-y-2">
              <div class="flex items-center gap-2 mb-3">
                <h3 class="text-base font-black text-gray-900 leading-6">{{ destination.name }}</h3>
                <!-- 类型徽章 -->
                <span :class="getTypeBadgeClass(destination.type)">
                  {{ destination.type_label || getTypeLabel(destination.type) }}
                </span>
                <!-- 加密状态徽章（仅非服务器类型显示） -->
                <span
                  v-if="destination.type !== 'server' && destination.encrypted"
                  class="px-2 py-0.5 text-xs font-bold rounded border-2 border-black bg-purple-600 text-white cursor-help"
                  title="已启用 AES-256-GCM 加密"
                >
                  加密
                </span>
              </div>
              <div class="flex items-center text-sm mb-2">
                <svg class="flex-shrink-0 mr-2 h-4 w-4 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                </svg>
                <span class="text-gray-700 font-bold break-all leading-4">{{ destination.display_path || getDestinationPath(destination) }}</span>
              </div>
              <!-- 创建时间 -->
              <div class="flex items-center text-sm">
                <svg class="flex-shrink-0 mr-2 h-4 w-4 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                <span class="text-gray-600 font-bold leading-4">创建于 {{ formatDateTime(destination.created_at) }}</span>
              </div>
            </div>

            <!-- 右侧：操作按钮 -->
            <div class="flex items-center gap-2 ml-4">
              <button
                @click="toggleDestination(destination.id, !destination.enabled)"
                :class="[
                  'px-3 py-1 text-sm font-bold rounded border-2 border-black transition-all',
                  destination.enabled ? 'text-gray-700 hover:bg-gray-50' : 'text-brutalist-green hover:bg-green-50'
                ]"
              >
                {{ destination.enabled ? '禁用' : '启用' }}
              </button>
              <button
                @click="editDestination(destination)"
                class="px-3 py-1 text-sm font-bold text-brutalist-blue hover:bg-blue-50 rounded border-2 border-black transition-all"
              >
                编辑
              </button>
              <button
                @click="deleteDestination(destination.id)"
                class="px-3 py-1 text-sm font-bold text-brutalist-red hover:bg-red-50 rounded border-2 border-black transition-all"
              >
                删除
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <Pagination
      :page="pagination.page"
      :page-size="pagination.page_size"
      :total="pagination.total"
      :total-page="pagination.total_page"
      @page-change="handlePageChange"
    />

    <!-- Destination Modal -->
    <DestinationModal
      v-if="showModal"
      :destination="editingDestination"
      @close="closeModal"
      @saved="handleSaved"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { destinationsApi } from '@/api'
import { useToast } from '@/composables/useToast'
import { useConfirm } from '@/composables/useConfirm'
import DestinationModal from '@/components/features/Destination/DestinationModal.vue'
import Pagination from '@/components/ui/Pagination.vue'

const toast = useToast()
const { confirm } = useConfirm()
const destinations = ref([])
const loading = ref(false)
const showModal = ref(false)
const editingDestination = ref(null)
const pagination = ref({
  page: 1,
  page_size: 10,
  total: 0,
  total_page: 0
})

const formatDateTime = (dateStr) => {
  if (!dateStr) return 'N/A'
  const date = new Date(dateStr)
  const pad = (n) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

const getTypeLabel = (type) => {
  const labels = {
    'local': '本地存储',
    'webdav': 'WebDAV',
    's3': 'S3',
    'server': '服务器'
  }
  return labels[type] || type
}

const getDestinationPath = (dest) => {
  switch (dest.type) {
    case 'local':
      return dest.path || dest.local_path || 'N/A'
    case 'webdav':
      return dest.webdav_url ? `${dest.webdav_url}${dest.webdav_path || ''}` : (dest.path || 'N/A')
    case 's3':
      return dest.s3_bucket ? `s3://${dest.s3_bucket}${dest.s3_path || ''}` : (dest.path || 'N/A')
    default:
      return dest.path || 'N/A'
  }
}

const getTypeBadgeClass = (type) => {
  const baseClasses = 'px-2 py-0.5 text-xs font-bold rounded border-2 border-black'
  const typeClasses = {
    'local': 'bg-brutalist-green text-white',
    'webdav': 'bg-brutalist-blue text-white',
    's3': 'bg-brutalist-orange text-white',
    'server': 'bg-gray-800 text-white'
  }
  return `${baseClasses} ${typeClasses[type] || 'bg-gray-100 text-gray-800'}`
}

const loadDestinations = async () => {
  loading.value = true
  try {
    const res = await destinationsApi.getAll({
      page: pagination.value.page,
      page_size: pagination.value.page_size
    })
    destinations.value = res.data
    pagination.value = res.pagination
  } catch (error) {
    console.error('Failed to load destinations:', error)
    toast.error('加载备份目标失败')
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  pagination.value.page = page
  loadDestinations()
}

const editDestination = (destination) => {
  editingDestination.value = destination
  showModal.value = true
}

const toggleDestination = async (id, enabled) => {
  try {
    await destinationsApi.update(id, { enabled })
    toast.success(enabled ? '已启用' : '已禁用')
    loadDestinations()
  } catch (error) {
    console.error('Failed to toggle destination:', error)
    toast.error('操作失败')
  }
}

const deleteDestination = async (id) => {
  const confirmed = await confirm({
    title: '删除备份目标',
    message: '确定要删除这个备份目标吗？此操作不可恢复。',
    type: 'danger',
    confirmText: '删除'
  })
  if (!confirmed) return

  try {
    await destinationsApi.delete(id)
    toast.success('备份目标已删除')
    loadDestinations()
  } catch (error) {
    console.error('Failed to delete destination:', error)
    toast.error('删除失败')
  }
}

const closeModal = () => {
  showModal.value = false
  editingDestination.value = null
}

const handleSaved = () => {
  closeModal()
  loadDestinations()
}

onMounted(() => {
  loadDestinations()
})
</script>
