<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h2 class="text-lg font-black text-gray-900">运行日志</h2>
        <p class="text-sm text-gray-700 font-bold">查看备份任务执行记录</p>
      </div>
      <!-- 任务筛选 -->
      <div class="flex items-center gap-2">
        <CustomSelect
          v-model="selectedTaskId"
          :options="taskOptions"
          placeholder="全部任务"
          class="w-48"
          @update:modelValue="loadLogs"
        />
      </div>
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
        <div v-if="log.message" class="space-y-2">
          <p class="text-sm text-gray-700">{{ formatMessage(log.message) }}</p>
          <!-- 只有失败日志且消息被优化过才显示展开按钮 -->
          <template v-if="log.status === 'failed' && formatMessage(log.message) !== log.message">
            <button
              @click="toggleDetail(log.id)"
              class="text-xs text-gray-500 hover:text-gray-700 flex items-center gap-1"
            >
              <svg :class="['h-3 w-3 transition-transform', expandedLogs.has(log.id) ? 'rotate-90' : '']" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
              </svg>
              {{ expandedLogs.has(log.id) ? '收起' : '查看原始错误' }}
            </button>
            <div v-if="expandedLogs.has(log.id)" class="mt-2 border-l-2 border-red-300 pl-3">
              <code class="text-xs text-red-600 bg-red-50 px-2 py-1.5 rounded block font-mono break-all">
                {{ log.message }}
              </code>
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { logsApi, tasksApi } from '@/api'
import { useToast } from '@/composables/useToast'
import CustomSelect from '@/components/ui/CustomSelect.vue'

const toast = useToast()
const logs = ref([])
const tasks = ref([])
const loading = ref(false)
const selectedTaskId = ref('')
const expandedLogs = ref(new Set())

const toggleDetail = (logId) => {
  if (expandedLogs.value.has(logId)) {
    expandedLogs.value.delete(logId)
  } else {
    expandedLogs.value.add(logId)
  }
  expandedLogs.value = new Set(expandedLogs.value)
}

const taskOptions = computed(() => {
  return [
    { label: '全部任务', value: '' },
    ...tasks.value.map(t => ({ label: t.name, value: t.id }))
  ]
})

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

const formatMessage = (message) => {
  if (!message) return ''

  // 成功消息
  if (message === 'Backup completed successfully') {
    return '备份成功'
  }

  // 错误消息映射
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
    { pattern: /bw status failed/i, text: 'Bitwarden CLI 状态检查失败' },
  ]

  for (const { pattern, text } of errorMappings) {
    if (pattern.test(message)) {
      return text
    }
  }

  return message
}

const loadLogs = async () => {
  loading.value = true
  try {
    const params = {}
    if (selectedTaskId.value) {
      params.task_id = selectedTaskId.value
    }
    logs.value = await logsApi.getAll(params)
  } catch (error) {
    console.error('Failed to load logs:', error)
    toast.error('加载日志失败')
  } finally {
    loading.value = false
  }
}

const loadTasks = async () => {
  try {
    tasks.value = await tasksApi.getAll()
  } catch (error) {
    console.error('Failed to load tasks:', error)
  }
}

onMounted(() => {
  loadTasks()
  loadLogs()
})
</script>
