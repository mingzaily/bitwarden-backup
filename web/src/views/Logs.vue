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
          @update:modelValue="handleTaskChange"
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
      <div
        v-for="log in logs"
        :key="log.id"
        :class="[
          'bg-white rounded-lg border-2 border-black shadow-brutalist-sm',
          'border-l-4',
          log.status === 'success' ? 'border-l-brutalist-green' : log.status === 'failed' ? 'border-l-red-500' : 'border-l-brutalist-blue'
        ]"
      >
        <div class="px-6 py-4">
          <div class="flex items-center justify-between">
            <!-- 左侧：日志信息 -->
            <div class="flex-1 space-y-2">
              <div class="flex items-center gap-2 mb-2">
                <span :class="['px-2 py-1 text-xs font-bold rounded border-2 border-black', getStatusClass(log.status)]">
                  {{ getStatusLabel(log.status) }}
                </span>
                <span class="text-sm font-bold text-gray-900">{{ log.task_name }}</span>
              </div>
              <p v-if="log.message" class="text-sm text-gray-700">{{ formatMessage(log.message) }}</p>
              <!-- 成功时显示备份文件路径 -->
              <div v-if="log.status === 'success' && log.backup_file" class="flex items-start gap-1.5 text-xs text-gray-500">
                <svg class="h-3.5 w-3.5 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                </svg>
                <span class="font-mono break-all">{{ log.backup_file }}</span>
              </div>
              <!-- 执行时间 -->
              <div class="flex items-center gap-1.5 text-xs text-gray-500">
                <svg class="h-3.5 w-3.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                <span>{{ formatTime(log.created_at) }}</span>
              </div>
              <!-- 只有失败日志且消息被优化过才显示展开按钮 -->
              <template v-if="log.status === 'failed' && formatMessage(log.message) !== log.message">
                <button
                  @click="toggleDetail(log.id)"
                  class="text-xs text-gray-500 hover:text-gray-700 flex items-center gap-1"
                >
                  <svg :class="['h-3 w-3 transition-transform', expandedLogs.has(log.id) ? 'rotate-90' : '']" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                  </svg>
                  {{ expandedLogs.has(log.id) ? '收起' : '查看错误' }}
                </button>
                <div v-if="expandedLogs.has(log.id)" class="mt-2 border-l-2 border-red-300 pl-3">
                  <code class="text-xs text-red-600 bg-red-50 px-2 py-1.5 rounded block font-mono break-all">
                    {{ log.message }}
                  </code>
                </div>
              </template>
            </div>

            <!-- 右侧：操作按钮 -->
            <div class="flex items-center ml-4">
              <button
                @click="showLogDetail(log)"
                class="px-3 py-1.5 text-sm font-bold text-brutalist-blue hover:bg-blue-50 rounded border-2 border-black transition-all"
              >
                查看
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页组件 -->
      <Pagination
        :page="pagination.page"
        :page-size="pagination.pageSize"
        :total="pagination.total"
        :total-page="pagination.totalPage"
        @change="handlePageChange"
      />
    </div>

    <!-- 日志详情 Modal -->
    <LogDetailModal
      v-if="selectedLog"
      :log="selectedLog"
      @close="selectedLog = null"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
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

// 分页状态
const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0,
  totalPage: 0
})

const toggleDetail = (logId) => {
  if (expandedLogs.value.has(logId)) {
    expandedLogs.value.delete(logId)
  } else {
    expandedLogs.value.add(logId)
  }
  expandedLogs.value = new Set(expandedLogs.value)
}

const showLogDetail = (log) => {
  selectedLog.value = log
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
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    }
    if (selectedTaskId.value) {
      params.task_id = selectedTaskId.value
    }
    const res = await logsApi.getAll(params)
    logs.value = res.data || []
    pagination.value = {
      page: res.pagination?.page || 1,
      pageSize: res.pagination?.page_size || 10,
      total: res.pagination?.total || 0,
      totalPage: res.pagination?.total_page || 0
    }
  } catch (error) {
    console.error('Failed to load logs:', error)
    toast.error('加载日志失败')
  } finally {
    loading.value = false
  }
}

// 切换任务筛选时重置到第一页
const handleTaskChange = () => {
  pagination.value.page = 1
  loadLogs()
}

// 分页切换
const handlePageChange = (page) => {
  pagination.value.page = page
  loadLogs()
}

const loadTasks = async () => {
  try {
    const res = await tasksApi.getAll()
    tasks.value = res.data || []
  } catch (error) {
    console.error('Failed to load tasks:', error)
  }
}

onMounted(() => {
  loadTasks()
  loadLogs()
})
</script>
