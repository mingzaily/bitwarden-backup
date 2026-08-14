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
            <p class="mt-1 break-all text-sm leading-6 text-main">{{ formatSummary(log.message) }}</p>
          </div>

          <div v-if="log.backup_file" class="log-file">
            <div class="mb-1 flex items-center gap-2 font-sans text-xs font-semibold text-accent"><svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M9 12h6m-6 4h6m2 5H7a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5.6a1 1 0 0 1 .7.3l5.4 5.4a1 1 0 0 1 .3.7V19a2 2 0 0 1-2 2Z" /></svg>备份文件</div>
            {{ log.backup_file }}
          </div>

          <section>
            <h4 class="mb-2 text-sm font-semibold text-main">执行过程</h4>
            <div v-if="executionLogs.length > 0" class="log-console">
              <div v-for="(entry, index) in executionLogs" :key="index" class="log-entry">
                <span class="log-entry-time text-subtle">{{ entry.time }}</span>
                <span :class="['log-source-badge', sourceClass(entry.source)]">{{ sourceLabel(entry.source) }}</span>
                <span :class="['log-entry-message', getLogClass(entry.message, entry.level)]">{{ entry.message }}</span>
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

const formatSummary = (message) => {
  if (!message) return ''
  const text = String(message)
  if (text === 'Backup completed successfully') return '备份成功'
  const partialFailure = text.match(/^Backup completed with destination errors:\s*(.*)$/i)
  if (partialFailure) return `部分目标失败：${partialFailure[1]}`
  const allFailed = text.match(/^all \d+ backup destinations failed:\s*(.*)$/i)
  if (allFailed) return `备份目标全部失败：${allFailed[1]}`
  const errorMappings = [
    { pattern: /source server is disabled/i, text: '源服务器已停用，请先启用源站' },
    { pattern: /no enabled backup destinations/i, text: '没有可用的已启用备份目标' },
    { pattern: /target server is disabled/i, text: '目标服务器已停用，请先启用目标服务器' },
    { pattern: /failed to clean up old backups/i, text: '备份已生成，但清理旧备份失败' }
  ]
  for (const { pattern: errorPattern, text: mappedText } of errorMappings) {
    if (errorPattern.test(text)) return mappedText
  }
  return text
}

const formatLogMessage = (message) => {
  const text = String(message || '')
  const statusPrefix = 'bw status stdout:'
  if (!text.toLowerCase().startsWith(statusPrefix)) return text

  try {
    const parsed = JSON.parse(text.slice(statusPrefix.length).trim())
    if (parsed?.status) {
      return parsed.serverUrl ? `bw status: ${parsed.status} (${parsed.serverUrl})` : `bw status: ${parsed.status}`
    }
  } catch {
    // Older logs may contain truncated JSON. Do not expose the raw payload.
  }
  return 'bw status: 状态输出已收到'
}

const executionLogs = computed(() => {
  if (!props.log.execution_logs) return []
  try {
    const parsed = JSON.parse(props.log.execution_logs)
    if (!Array.isArray(parsed)) return []
    return parsed.filter(Boolean).map((entry) => ({
      time: String(entry.time || '—'),
      source: String(entry.source || 'bitwarden'),
      level: String(entry.level || ''),
      message: formatLogMessage(entry.message)
    }))
  } catch { return [] }
})
const statusLabel = computed(() => ({ success: '成功', failed: '失败', running: '运行中' }[props.log.status] || props.log.status))
const statusClass = computed(() => ({
  success: 'status-badge status-success',
  failed: 'status-badge status-danger',
  running: 'status-badge status-info'
}[props.log.status] || 'status-badge status-neutral'))
const formatTime = (time) => time ? new Date(time).toLocaleString('zh-CN') : 'N/A'
const sourceLabel = (source) => ({
  bitwarden: 'Bitwarden',
  local: '本地 CP',
  webdav: 'WebDAV',
  s3: 'OSS',
  server: '服务器'
}[source || 'bitwarden'] || source || '系统')
const sourceClass = (source) => ({
  bitwarden: 'log-source-bitwarden',
  local: 'log-source-local',
  webdav: 'log-source-webdav',
  s3: 'log-source-s3',
  server: 'log-source-server'
}[source || 'bitwarden'] || 'log-source-system')
const getLogClass = (message, level) => {
  const text = String(message || '')
  if (level === 'error') return 'text-danger'
  if (/^bw logout \(exit=1/.test(text) || /already logged out/i.test(text)) return 'text-muted'
  if (text.includes('exit=1') || text.includes('stderr')) return 'text-danger'
  if (text.includes('Executing task')) return 'text-info'
  if (text.includes('exit=0')) return 'text-accent'
  if (level === 'info') return 'text-muted'
  return 'text-muted'
}
</script>
