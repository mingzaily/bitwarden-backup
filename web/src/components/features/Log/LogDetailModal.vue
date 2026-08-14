<template>
  <Teleport to="body">
    <div class="modal-backdrop" @click.self="$emit('close')">
      <div class="modal-panel modal-panel-wide" role="dialog" aria-modal="true" aria-label="执行日志详情">
        <div class="modal-header">
          <div>
            <h3 class="modal-title">执行日志详情</h3>
            <p class="modal-subtitle">查看这次备份执行的状态、产物和原始过程。</p>
          </div>
          <button class="icon-button" type="button" aria-label="关闭" @click="$emit('close')">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="m6 6 12 12M18 6 6 18" /></svg>
          </button>
        </div>

        <div class="modal-body grid gap-5">
          <div class="flex flex-wrap items-center gap-3">
            <span :class="statusClass">{{ statusLabel }}</span>
            <span class="text-sm font-semibold text-main">{{ log.task_name }}</span>
            <span class="text-xs text-muted">{{ formatTime(log.created_at) }}</span>
          </div>

          <div v-if="log.message" class="surface-muted rounded-xl px-4 py-3">
            <p class="text-xs font-semibold text-muted">执行摘要</p>
            <p class="mt-1 break-all text-sm leading-6 text-main">{{ log.message === 'Backup completed successfully' ? '备份成功' : log.message }}</p>
          </div>

          <div v-if="log.backup_file" class="log-file">
            <div class="mb-1 flex items-center gap-2 font-sans text-xs font-semibold text-accent"><svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M9 12h6m-6 4h6m2 5H7a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5.6a1 1 0 0 1 .7.3l5.4 5.4a1 1 0 0 1 .3.7V19a2 2 0 0 1-2 2Z" /></svg>备份文件</div>
            {{ log.backup_file }}
          </div>

          <section>
            <h4 class="mb-2 text-sm font-semibold text-main">执行过程</h4>
            <div v-if="executionLogs.length > 0" class="log-console">
              <div v-for="(entry, index) in executionLogs" :key="index" class="flex gap-3 py-0.5">
                <span class="flex-shrink-0 text-subtle">{{ entry.time }}</span>
                <span :class="getLogClass(entry.message)">{{ entry.message }}</span>
              </div>
            </div>
            <div v-else class="surface-muted py-8 text-center text-sm text-muted">暂无详细执行日志</div>
          </section>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn-secondary w-full sm:w-auto" @click="$emit('close')">关闭</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({ log: { type: Object, required: true } })
defineEmits(['close'])

const executionLogs = computed(() => {
  if (!props.log.execution_logs) return []
  try { return JSON.parse(props.log.execution_logs) } catch { return [] }
})
const statusLabel = computed(() => ({ success: '成功', failed: '失败', running: '运行中' }[props.log.status] || props.log.status))
const statusClass = computed(() => ({
  success: 'status-badge status-success',
  failed: 'status-badge status-danger',
  running: 'status-badge status-info'
}[props.log.status] || 'status-badge status-neutral'))
const formatTime = (time) => time ? new Date(time).toLocaleString('zh-CN') : 'N/A'
const getLogClass = (message) => {
  if (message.includes('exit=0')) return 'text-accent'
  if (message.includes('exit=1') || message.includes('stderr')) return 'text-danger'
  if (message.includes('Executing task')) return 'text-info'
  return 'text-muted'
}
</script>
