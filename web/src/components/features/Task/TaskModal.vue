<template>
  <Teleport to="body">
    <div class="modal-backdrop" @click.self.stop>
      <div class="modal-panel modal-panel-wide" role="dialog" aria-modal="true" :aria-label="task ? '编辑任务' : '新建任务'">
        <div class="modal-header">
          <div>
            <h3 class="modal-title">{{ task ? '编辑任务' : '新建任务' }}</h3>
            <p class="modal-subtitle">{{ task ? '更新任务来源、目标和运行计划；已保存的启用状态会保留。' : '把 Bitwarden 源站、存储目标和执行计划组合成一个可追踪的任务。' }}</p>
          </div>
          <button class="icon-button" type="button" aria-label="关闭" @click="$emit('close')">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="m6 6 12 12M18 6 6 18" /></svg>
          </button>
        </div>

        <form id="task-form" class="modal-body grid gap-5" @submit.prevent="handleSubmit">
          <section class="form-section">
            <div class="form-section-heading">
              <h4 class="form-section-title">运行设置</h4>
              <p class="form-section-description">先命名任务，再选择只手动执行或按计划自动运行。</p>
            </div>
            <div class="field">
              <label class="field-label" for="task-name">任务名称</label>
              <input id="task-name" v-model.trim="formData.name" class="input" type="text" required placeholder="例如：每日备份" />
            </div>
            <div v-if="task" class="surface-muted flex items-center justify-between gap-4 p-3">
              <div>
                <p class="text-sm font-semibold text-main">任务状态</p>
                <p class="mt-1 text-xs text-muted">停用后不会自动调度，也不能从列表立即执行。</p>
              </div>
              <ToggleButton v-model="formData.enabled" :label="formData.enabled ? '已启用' : '已停用'" :aria-label="formData.enabled ? '任务已启用' : '任务已停用'" />
            </div>

            <div class="field">
              <div class="field-label-row">
                <span class="field-label">运行方式</span>
                <span class="field-label-note schedule-status-label">{{ scheduleMode === 'scheduled' ? '自动运行' : '仅手动' }}</span>
              </div>
              <div class="schedule-mode-list" role="radiogroup" aria-label="运行方式">
                <button
                  type="button"
                  :class="['schedule-mode-card selection-card', scheduleMode === 'manual' ? 'is-active' : '']"
                  role="radio"
                  :aria-checked="scheduleMode === 'manual'"
                  @click="setScheduleMode('manual')"
                >
                  <span class="schedule-mode-icon" aria-hidden="true">
                    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M8 5v14l11-7L8 5Z" /></svg>
                  </span>
                  <span class="schedule-mode-copy">
                    <span class="schedule-mode-title">仅手动执行</span>
                    <span class="schedule-mode-description">需要时从任务列表立即执行</span>
                  </span>
                  <svg v-if="scheduleMode === 'manual'" class="schedule-mode-check" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m5 12 4 4L19 6" /></svg>
                </button>
                <button
                  type="button"
                  :class="['schedule-mode-card selection-card', scheduleMode === 'scheduled' ? 'is-active' : '']"
                  role="radio"
                  :aria-checked="scheduleMode === 'scheduled'"
                  @click="setScheduleMode('scheduled')"
                >
                  <span class="schedule-mode-icon" aria-hidden="true">
                    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M7 4v3m10-3v3M4.5 9.5h15M6 5h12a2 2 0 0 1 2 2v11.5a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V7a2 2 0 0 1 2-2ZM8 13h3m2 0h3m-8 3h3m2 0h3" /></svg>
                  </span>
                  <span class="schedule-mode-copy">
                    <span class="schedule-mode-title">按计划自动运行</span>
                    <span class="schedule-mode-description">按照 Cron 表达式定时触发</span>
                  </span>
                  <svg v-if="scheduleMode === 'scheduled'" class="schedule-mode-check" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m5 12 4 4L19 6" /></svg>
                </button>
              </div>
            </div>

            <div v-if="scheduleMode === 'scheduled'" class="schedule-editor">
              <div class="field">
                <div class="field-label-row">
                  <label class="field-label" for="task-cron">Cron 表达式</label>
                  <span class="field-label-note">支持 5 / 6 位</span>
                </div>
                <input id="task-cron" v-model="formData.cron_expression" class="input schedule-input" type="text" aria-describedby="task-cron-hint" placeholder="0 0 2 * * *" />
                <p id="task-cron-hint" class="field-hint">例如 <code>0 0 2 * * *</code> 表示每天凌晨 2 点；留意服务所在时区。</p>
              </div>
              <div class="field">
                <span class="field-label">常用计划</span>
                <div class="schedule-preset-list">
                  <button v-for="preset in schedulePresets" :key="preset.value" :class="['schedule-preset', formData.cron_expression === preset.value ? 'is-active' : '']" type="button" :aria-pressed="formData.cron_expression === preset.value" @click="selectSchedulePreset(preset.value)">{{ preset.label }}</button>
                </div>
              </div>
              <div class="schedule-preview">
                <span class="schedule-preview-icon" aria-hidden="true">
                  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M7 4v3m10-3v3M4.5 9.5h15M6 5h12a2 2 0 0 1 2 2v11.5a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V7a2 2 0 0 1 2-2ZM8 13h3m2 0h3m-8 3h3" /></svg>
                </span>
                <span class="schedule-preview-copy"><span class="schedule-preview-label">执行摘要</span><strong>{{ scheduleSummary }}</strong></span>
                <span :class="['status-badge', scheduleIsValid ? 'status-success' : 'status-warning']">{{ scheduleIsValid ? '已配置' : '待完善' }}</span>
              </div>
            </div>
            <div v-else class="schedule-manual-card">
              <span class="schedule-preview-icon" aria-hidden="true">
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M8 5v14l11-7L8 5Z" /></svg>
              </span>
              <p class="schedule-manual-copy">不会自动调度，创建后可随时从任务列表手动执行。</p>
            </div>
          </section>

          <section class="form-section">
            <div class="form-section-heading">
              <h4 class="form-section-title">备份链路</h4>
              <p class="form-section-description">每次执行会从一个 Bitwarden 源站导出，再写入一个或多个存储目标。</p>
            </div>
            <CustomSelect
              v-model="formData.source_server_id"
              :options="serverOptions"
              label="Bitwarden 源站"
              placeholder="请选择 Bitwarden 源站"
              empty-text="暂无可用源站，请先创建源站"
            />
            <CheckboxGroup
              v-model="formData.destination_ids"
              :options="destinationOptions"
              label="存储目标（可多选）"
              empty-text="暂无可用存储目标，请先创建存储目标"
            />
          </section>
        </form>

        <div class="modal-footer">
          <button type="button" class="btn-secondary" @click="$emit('close')">取消</button>
          <button form="task-form" type="submit" class="btn-primary" :disabled="loading">
            <span v-if="loading" class="spinner"></span>{{ loading ? '保存中…' : task ? '保存更改' : '保存任务' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { tasksApi, serversApi, destinationsApi } from '@/api'
import { useToast } from '@/composables/useToast'
import CheckboxGroup from '@/components/ui/CheckboxGroup.vue'
import CustomSelect from '@/components/ui/CustomSelect.vue'
import ToggleButton from '@/components/ui/ToggleButton.vue'

const props = defineProps({ task: Object })
const emit = defineEmits(['close', 'saved'])
const toast = useToast()
const servers = ref([])
const destinations = ref([])
const emptyForm = () => ({ name: '', cron_expression: '', source_server_id: '', destination_ids: [], enabled: true })
const formData = ref(emptyForm())
const loading = ref(false)
const scheduleMode = ref('manual')
const schedulePresets = [
  { label: '每天 02:00', value: '0 0 2 * * *' },
  { label: '每 6 小时', value: '0 0 */6 * * *' },
  { label: '每周日 03:00', value: '0 0 3 * * 0' }
]

const getTypeLabel = (type) => ({ local: '本地存储', webdav: 'WebDAV', s3: 'S3', server: '服务器' }[type] || type)
const serverOptions = computed(() => {
  const currentID = Number(formData.value.source_server_id || 0)
  return servers.value
    .filter(server => server.enabled || Number(server.id) === currentID)
    .map(server => ({
      label: server.name,
      value: server.id,
      description: `${server.server_url}${server.enabled ? '' : ' · 已停用'}`
    }))
})
const destinationOptions = computed(() => {
  const currentIDs = new Set((formData.value.destination_ids || []).map(id => Number(id)))
  return destinations.value
    .filter(destination => destination.enabled || currentIDs.has(Number(destination.id)))
    .map(destination => ({
      label: destination.name,
      value: destination.id,
      description: `类型：${getTypeLabel(destination.type)}${destination.enabled ? '' : ' · 已停用'}`
    }))
})

watch(() => props.task, (newTask) => {
  if (newTask) {
    formData.value = {
      name: newTask.name || '',
      cron_expression: newTask.cron_expression || '',
      source_server_id: newTask.source_server?.id || newTask.source_server_id || '',
      destination_ids: Array.isArray(newTask.destinations) ? newTask.destinations.map(destination => destination.id) : (newTask.destination_ids || []),
      enabled: newTask.enabled ?? true
    }
    scheduleMode.value = newTask.cron_expression?.trim() ? 'scheduled' : 'manual'
  } else {
    formData.value = emptyForm()
    scheduleMode.value = 'manual'
  }
}, { immediate: true })

const loadServers = async () => {
  try {
    const res = await serversApi.getAll({ page: 1, page_size: 1000 })
    servers.value = res.data || []
  } catch (error) {
    console.error('Failed to load servers:', error)
  }
}
const loadDestinations = async () => {
  try {
    const res = await destinationsApi.getAll({ page: 1, page_size: 1000 })
    destinations.value = res.data || []
  } catch (error) {
    console.error('Failed to load destinations:', error)
  }
}
onMounted(() => { loadServers(); loadDestinations() })

const isValidCronExpression = (expression) => {
  if (!expression || expression.trim() === '') return true
  const parts = expression.trim().split(/\s+/)
  return parts.length === 5 || parts.length === 6
}

const scheduleIsValid = computed(() => {
  const expression = formData.value.cron_expression.trim()
  return scheduleMode.value !== 'scheduled' || (expression !== '' && isValidCronExpression(expression))
})

const scheduleSummary = computed(() => {
  const expression = formData.value.cron_expression.trim()
  if (!expression) return '输入 Cron 表达式后会显示计划'
  const preset = schedulePresets.find(item => item.value === expression)
  return preset ? preset.label : `自定义计划 · ${expression}`
})

const setScheduleMode = (mode) => {
  scheduleMode.value = mode
  if (mode === 'manual') {
    formData.value.cron_expression = ''
  } else if (!formData.value.cron_expression.trim()) {
    formData.value.cron_expression = schedulePresets[0].value
  }
}

const selectSchedulePreset = (value) => {
  scheduleMode.value = 'scheduled'
  formData.value.cron_expression = value
}

const handleSubmit = async () => {
  if (!formData.value.name.trim()) {
    toast.error('请输入任务名称')
    return
  }
  if (!formData.value.source_server_id) {
    toast.error('请选择 Bitwarden 源站')
    return
  }
  if (!formData.value.destination_ids || formData.value.destination_ids.length === 0) {
    toast.error('请至少选择一个存储目标')
    return
  }
  if (scheduleMode.value === 'scheduled' && !formData.value.cron_expression.trim()) {
    toast.error('请输入 Cron 表达式，或切换为仅手动执行')
    return
  }
  if (formData.value.cron_expression && !isValidCronExpression(formData.value.cron_expression)) {
    toast.error('Cron 表达式格式不正确，应为 5 位或 6 位格式')
    return
  }

  loading.value = true
  try {
    const data = { ...formData.value }
    if (!props.task?.id) delete data.enabled
    if (props.task?.id) {
      await tasksApi.update(props.task.id, data)
      toast.success('任务已更新')
    } else {
      await tasksApi.create(data)
      toast.success('任务已创建')
    }
    emit('saved')
  } catch (error) {
    console.error('Failed to save task:', error)
    toast.error(error.message || '保存失败')
  } finally {
    loading.value = false
  }
}
</script>
