<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-lg font-black text-gray-900">运行日志</h2>
      <p class="text-sm text-gray-700 font-bold">查看备份任务执行记录</p>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-brutalist-blue border-t-transparent"></div>
      <p class="mt-2 text-sm text-gray-600">加载中...</p>
    </div>

    <div v-else-if="logs.length === 0" class="text-center py-12 bg-brutalist-cream/20 rounded-lg border-2 border-black">
      <p class="text-gray-600">暂无运行日志</p>
    </div>

    <div v-else class="space-y-3">
      <div v-for="log in logs" :key="log.id" class="bg-white rounded-lg border-2 border-black shadow-brutalist-sm p-4">
        <div class="flex items-start justify-between mb-2">
          <div class="flex items-center gap-2">
            <span :class="['px-2 py-1 text-xs font-bold rounded border-2 border-black', getStatusClass(log.status)]">
              {{ getStatusLabel(log.status) }}
            </span>
            <span class="text-sm font-bold text-gray-900">{{ log.task_name }}</span>
          </div>
          <span class="text-xs text-gray-600">{{ formatTime(log.created_at) }}</span>
        </div>
        <p v-if="log.message" class="text-sm text-gray-700">{{ log.message }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { logsApi } from '@/api'
import { useToast } from '@/composables/useToast'

const toast = useToast()
const logs = ref([])
const loading = ref(false)

const getStatusLabel = (status) => {
  const labels = {
    'success': '成功',
    'failed': '失败',
    'running': '运行中'
  }
  return labels[status] || status
}

const getStatusClass = (status) => {
  const classes = {
    'success': 'bg-green-100 text-green-800',
    'failed': 'bg-red-100 text-red-800',
    'running': 'bg-blue-100 text-blue-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const formatTime = (time) => {
  if (!time) return 'N/A'
  return new Date(time).toLocaleString('zh-CN')
}

const loadLogs = async () => {
  loading.value = true
  try {
    logs.value = await logsApi.getAll()
  } catch (error) {
    console.error('Failed to load logs:', error)
    toast.error('加载日志失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadLogs()
})
</script>
