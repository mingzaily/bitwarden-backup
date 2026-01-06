<template>
  <Teleport to="body">
    <div class="modal show" @click.self="$emit('close')">
      <div class="modal-content bg-white rounded-lg border-2 border-black shadow-brutalist max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
      <!-- Modal Header -->
      <div class="px-6 py-4 border-b-2 border-black bg-brutalist-cream/20">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-black text-gray-900">
            {{ destination ? '编辑备份目标' : '新建备份目标' }}
          </h3>
          <button @click="$emit('close')" class="text-gray-700 hover:text-gray-900 font-bold">
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
      </div>

      <!-- Modal Body -->
      <form @submit.prevent="handleSubmit" class="p-6 space-y-4">
        <div>
          <label class="block text-sm font-bold text-gray-900 mb-2">目标名称</label>
          <input
            v-model="formData.name"
            type="text"
            required
            class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
            placeholder="例如：本地备份"
          />
        </div>

        <TabSelector
          v-model="formData.type"
          :options="storageTypes"
          label="存储类型"
        />

        <!-- Local Config -->
        <div v-if="formData.type === 'local'" class="space-y-4">
          <div>
            <label class="block text-sm font-bold text-gray-900 mb-2">本地路径</label>
            <input
              v-model="formData.local_path"
              type="text"
              required
              class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
              placeholder="/data/backups"
            />
            <p class="mt-1 text-xs text-gray-600">
              💡 请填写绝对路径，例如：<code class="px-1 py-0.5 bg-gray-100 rounded">/app/backups</code> 或 <code class="px-1 py-0.5 bg-gray-100 rounded">D:/backups</code>
            </p>
          </div>
        </div>

        <!-- WebDAV Config -->
        <div v-if="formData.type === 'webdav'" class="space-y-4">
          <div>
            <label class="block text-sm font-bold text-gray-900 mb-2">WebDAV URL</label>
            <input
              v-model="formData.webdav_url"
              type="text"
              required
              class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
              placeholder="https://dav.example.com"
            />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-bold text-gray-900 mb-2">用户名</label>
              <input
                v-model="formData.webdav_username"
                type="text"
                class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
              />
            </div>
            <div>
              <label class="block text-sm font-bold text-gray-900 mb-2">密码</label>
              <input
                v-model="formData.webdav_password"
                type="password"
                class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
              />
            </div>
          </div>
          <div>
            <label class="block text-sm font-bold text-gray-900 mb-2">存储路径</label>
            <input
              v-model="formData.webdav_path"
              type="text"
              class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
              placeholder="/backups"
            />
          </div>
        </div>

        <!-- S3 Config -->
        <div v-if="formData.type === 's3'" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-bold text-gray-900 mb-2">Endpoint</label>
              <input
                v-model="formData.s3_endpoint"
                type="text"
                required
                class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
                placeholder="s3.amazonaws.com"
              />
            </div>
            <div>
              <label class="block text-sm font-bold text-gray-900 mb-2">区域 (Region)</label>
              <input
                v-model="formData.s3_region"
                type="text"
                required
                class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
                placeholder="us-east-1"
              />
            </div>
          </div>
          <div>
            <label class="block text-sm font-bold text-gray-900 mb-2">Bucket 名称</label>
            <input
              v-model="formData.s3_bucket"
              type="text"
              required
              class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
            />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-bold text-gray-900 mb-2">Access Key</label>
              <input
                v-model="formData.s3_access_key"
                type="text"
                class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
              />
            </div>
            <div>
              <label class="block text-sm font-bold text-gray-900 mb-2">Secret Key</label>
              <input
                v-model="formData.s3_secret_key"
                type="password"
                class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
              />
            </div>
          </div>
          <div>
            <label class="block text-sm font-bold text-gray-900 mb-2">存储路径</label>
            <input
              v-model="formData.s3_path"
              type="text"
              class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
              placeholder="/backups"
            />
          </div>
        </div>

        <!-- Server Config -->
        <div v-if="formData.type === 'server'" class="space-y-4">
          <CustomSelect
            v-model="formData.target_server_id"
            :options="servers.map(s => ({
              label: s.name,
              value: s.id,
              description: s.server_url
            }))"
            label="目标服务器"
            placeholder="请选择目标服务器"
            empty-text="⚠️ 暂无可用服务器，请先创建服务器"
          />
        </div>

        <!-- 加密选项（仅本地和 WebDAV 和 S3 显示） -->
        <div v-if="['local', 'webdav', 's3'].includes(formData.type)" class="space-y-2">
          <label class="block text-sm font-bold text-gray-900">备份文件加密</label>
          <ToggleButton v-model="formData.encrypted" label="加密备份文件" />
          <p class="text-xs text-gray-600 mt-1">
            💡 启用后将使用 Bitwarden CLI 的 <code class="px-1 py-0.5 bg-gray-100 rounded">--format encrypted_json</code> 导出加密文件
          </p>
        </div>

        <div class="space-y-2">
          <label class="block text-sm font-bold text-gray-900">启用状态</label>
          <ToggleButton v-model="formData.enabled" label="启用此备份目标" />
        </div>

        <!-- Modal Footer -->
        <div class="flex justify-end gap-3 pt-4 border-t-2 border-black">
          <button
            type="button"
            @click="$emit('close')"
            class="px-4 py-2 text-sm font-bold text-gray-700 bg-white border-2 border-black rounded-lg hover:bg-gray-50 transition-all"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="btn-brutalist px-4 py-2 text-sm rounded-lg disabled:opacity-50"
          >
            {{ loading ? '保存中...' : '保存' }}
          </button>
        </div>
      </form>
    </div>
  </div>
  </Teleport>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { destinationsApi, serversApi } from '@/api'
import { useToast } from '@/composables/useToast'
import TabSelector from './TabSelector.vue'
import ToggleButton from './ToggleButton.vue'
import CustomSelect from './CustomSelect.vue'

const props = defineProps({
  destination: Object
})

const emit = defineEmits(['close', 'saved'])
const toast = useToast()

const storageTypes = [
  { label: '本地存储', value: 'local' },
  { label: 'WebDAV', value: 'webdav' },
  { label: 'S3', value: 's3' },
  { label: '官方服务器', value: 'server' }
]

const servers = ref([])

const formData = ref({
  name: '',
  type: 'local',
  local_path: '',
  webdav_url: '',
  webdav_username: '',
  webdav_password: '',
  webdav_path: '',
  s3_endpoint: '',
  s3_region: '',
  s3_bucket: '',
  s3_access_key: '',
  s3_secret_key: '',
  s3_path: '',
  target_server_id: '',
  enabled: true,
  encrypted: false
})

const loading = ref(false)

watch(() => props.destination, (newDest) => {
  if (newDest) {
    formData.value = {
      ...newDest,
      // Ensure local_path is set from either local_path or path if old data exists
      local_path: newDest.local_path || (newDest.type === 'local' ? newDest.path : ''),
      target_server_id: newDest.target_server_id || '',
      encrypted: newDest.encrypted || false
    }
  } else {
    formData.value = {
      name: '',
      type: 'local',
      local_path: '',
      webdav_url: '',
      webdav_username: '',
      webdav_password: '',
      webdav_path: '',
      s3_endpoint: '',
      s3_region: '',
      s3_bucket: '',
      s3_access_key: '',
      s3_secret_key: '',
      s3_path: '',
      target_server_id: '',
      enabled: true,
      encrypted: false
    }
  }
}, { immediate: true })

const loadServers = async () => {
  try {
    servers.value = await serversApi.getAll({ enabled: true })
  } catch (error) {
    console.error('Failed to load servers:', error)
  }
}

onMounted(() => {
  loadServers()
})

const handleSubmit = async () => {
  loading.value = true
  try {
    const data = { ...formData.value }

    if (props.destination?.id) {
      await destinationsApi.update(props.destination.id, data)
      toast.success('备份目标已更新')
    } else {
      await destinationsApi.create(data)
      toast.success('备份目标已创建')
    }
    emit('saved')
  } catch (error) {
    console.error('Failed to save destination:', error)
    toast.error('保存失败')
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