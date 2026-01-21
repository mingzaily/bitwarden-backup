<template>
  <Teleport to="body">
    <div class="modal show" @click.self="$emit('close')">
      <div class="modal-content bg-white rounded-lg border-2 border-black shadow-brutalist max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
      <div class="px-6 py-4 border-b-2 border-black bg-brutalist-cream/20">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-black text-gray-900">{{ task ? '编辑任务' : '新建任务' }}</h3>
          <button @click="$emit('close')" class="text-gray-700 hover:text-gray-900 font-bold">
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
      </div>

      <form @submit.prevent="handleSubmit" class="p-6 space-y-4">
        <div>
          <label class="block text-sm font-bold text-gray-900 mb-2">任务名称</label>
          <input v-model="formData.name" type="text" required class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue" placeholder="例如：每日备份" />
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-900 mb-2">Cron 表达式（可选）</label>
          <input v-model="formData.cron_expression" type="text" class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue" placeholder="0 0 2 * * *（可选，留空则仅支持手动触发）" />
          <p class="mt-1 text-xs text-gray-600">💡 留空则创建手动触发任务，填写则自动定时执行。示例: 0 0 2 * * * (每天凌晨2点)</p>
        </div>

        <div>
          <CustomSelect
            v-model="formData.source_server_id"
            :options="servers.map(s => ({
              label: s.name,
              value: s.id,
              description: s.server_url
            }))"
            label="源服务器"
            placeholder="请选择源服务器"
            empty-text="⚠️ 暂无可用服务器，请先创建服务器"
          />
        </div>

        <!-- 备份目标多选 -->
        <CheckboxGroup
          v-model="formData.destination_ids"
          :options="destinations.map(d => ({
            label: d.name,
            value: d.id,
            description: `类型: ${getTypeLabel(d.type)}`
          }))"
          label="备份目标（可多选）"
          empty-text="暂无可用备份目标，请先创建备份目标"
        />

        <div class="flex justify-end gap-3 pt-4 border-t-2 border-black">
          <button type="button" @click="$emit('close')" class="px-4 py-2 text-sm font-bold text-gray-700 bg-white border-2 border-black rounded-lg hover:bg-gray-50 transition-all">取消</button>
          <button type="submit" :disabled="loading" class="btn-brutalist px-4 py-2 text-sm rounded-lg disabled:opacity-50">{{ loading ? '保存中...' : '保存' }}</button>
        </div>
      </form>
    </div>
  </div>
  </Teleport>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { tasksApi, serversApi, destinationsApi } from '@/api'
import { useToast } from '@/composables/useToast'
import CheckboxGroup from '@/components/ui/CheckboxGroup.vue'
import ToggleButton from '@/components/ui/ToggleButton.vue'
import CustomSelect from '@/components/ui/CustomSelect.vue'

const props = defineProps({ task: Object })
const emit = defineEmits(['close', 'saved'])
const toast = useToast()

const servers = ref([])
const destinations = ref([])
const formData = ref({
  name: '',
  cron_expression: '',
  source_server_id: '',
  destination_ids: [],
  enabled: true
})

const loading = ref(false)

const getTypeLabel = (type) => {
  const labels = {
    'local': '本地存储',
    'webdav': 'WebDAV',
    's3': 'S3',
    'server': '服务器'
  }
  return labels[type] || type
}

watch(() => props.task, (newTask) => {
  if (newTask) {
    const destinationIds = newTask.destinations?.map(d => d.id) || []
    const sourceServerId = newTask.source_server?.id || newTask.source_server_id || ''

    formData.value = {
      name: newTask.name || '',
      cron_expression: newTask.cron_expression || '',
      source_server_id: sourceServerId,
      destination_ids: destinationIds,
      enabled: newTask.enabled ?? true
    }
  } else {
    formData.value = {
      name: '',
      cron_expression: '',
      source_server_id: '',
      destination_ids: [],
      enabled: true
    }
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

onMounted(() => {
  loadServers()
  loadDestinations()
})

// 验证 Cron 表达式格式（简单校验）
const isValidCronExpression = (expr) => {
  if (!expr || expr.trim() === '') return true // 可选字段
  const parts = expr.trim().split(/\s+/)
  // 支持 5 位（分时日月周）或 6 位（秒分时日月周）
  return parts.length === 5 || parts.length === 6
}

const handleSubmit = async () => {
  // 必填校验：源服务器和备份目标
  if (!formData.value.source_server_id) {
    toast.error('请选择源服务器')
    return
  }
  if (!formData.value.destination_ids || formData.value.destination_ids.length === 0) {
    toast.error('请至少选择一个备份目标')
    return
  }

  // Cron 表达式格式校验
  if (formData.value.cron_expression && !isValidCronExpression(formData.value.cron_expression)) {
    toast.error('Cron 表达式格式不正确，应为 5 位或 6 位格式')
    return
  }

  loading.value = true
  try {
    const data = { ...formData.value }

    // 新增时不传 enabled 参数（后端默认为 true）
    if (!props.task?.id) {
      delete data.enabled
    }

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

<style scoped>
.modal {
  display: flex;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  background-color: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  align-items: center;
  justify-content: center;
}
</style>
