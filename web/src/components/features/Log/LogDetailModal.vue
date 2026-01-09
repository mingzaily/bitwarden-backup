<template>
  <div class="fixed inset-0 z-50 overflow-y-auto">
    <div class="flex min-h-screen items-center justify-center p-4">
      <!-- Backdrop -->
      <div class="fixed inset-0 bg-black/50" @click="$emit('close')"></div>

      <!-- Modal -->
      <div class="relative bg-white rounded-lg border-2 border-black shadow-brutalist w-full max-w-2xl max-h-[80vh] flex flex-col">
        <!-- Header -->
        <div class="flex items-center justify-between px-6 py-4 border-b-2 border-black">
          <h3 class="text-lg font-black text-gray-900">执行日志详情</h3>
          <button
            @click="$emit('close')"
            class="p-1 hover:bg-gray-100 rounded border-2 border-black"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Content -->
        <div class="flex-1 overflow-y-auto p-6">
          <!-- 基本信息 -->
          <div class="mb-4 flex items-center gap-3">
            <span :class="['px-2 py-1 text-xs font-bold rounded border-2 border-black', statusClass]">
              {{ statusLabel }}
            </span>
            <span class="font-bold text-gray-900">{{ log.task_name }}</span>
            <span class="text-sm text-gray-500">{{ formatTime(log.created_at) }}</span>
          </div>

          <!-- 备份文件 -->
          <div v-if="log.backup_file" class="mb-4 p-3 bg-green-50 rounded-lg border-2 border-black">
            <div class="flex items-start gap-2">
              <svg class="h-4 w-4 text-green-600 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <span class="text-sm font-mono text-green-800 break-all">{{ log.backup_file }}</span>
            </div>
          </div>

          <!-- 执行日志 -->
          <div v-if="executionLogs.length > 0" class="space-y-1">
            <h4 class="text-sm font-bold text-gray-700 mb-2">执行过程</h4>
            <div class="bg-gray-900 rounded-lg border-2 border-black p-4 font-mono text-sm overflow-x-auto">
              <div
                v-for="(entry, index) in executionLogs"
                :key="index"
                class="flex gap-2 py-0.5"
              >
                <span class="text-gray-500 flex-shrink-0">{{ entry.time }}</span>
                <span :class="getLogClass(entry.message)">{{ entry.message }}</span>
              </div>
            </div>
          </div>

          <!-- 无日志提示 -->
          <div v-else class="text-center py-8 text-gray-500">
            <p>暂无详细执行日志</p>
          </div>
        </div>

        <!-- Footer -->
        <div class="px-6 py-4 border-t-2 border-black">
          <button
            @click="$emit('close')"
            class="w-full px-4 py-2 bg-gray-100 hover:bg-gray-200 rounded-lg border-2 border-black font-bold transition-colors"
          >
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  log: {
    type: Object,
    required: true
  }
})

defineEmits(['close'])

const executionLogs = computed(() => {
  if (!props.log.execution_logs) return []
  try {
    return JSON.parse(props.log.execution_logs)
  } catch {
    return []
  }
})

const statusLabel = computed(() => {
  const labels = { success: '成功', failed: '失败', running: '运行中' }
  return labels[props.log.status] || props.log.status
})

const statusClass = computed(() => {
  const classes = {
    success: 'bg-green-100 text-green-800',
    failed: 'bg-red-100 text-red-800',
    running: 'bg-blue-100 text-blue-800'
  }
  return classes[props.log.status] || 'bg-gray-100 text-gray-800'
})

const formatTime = (time) => {
  if (!time) return 'N/A'
  return new Date(time).toLocaleString('zh-CN')
}

const getLogClass = (message) => {
  if (message.includes('exit=0')) return 'text-green-400'
  if (message.includes('exit=1') || message.includes('stderr')) return 'text-red-400'
  if (message.includes('Executing task')) return 'text-blue-400'
  return 'text-gray-300'
}
</script>
