<template>
  <Teleport to="body">
    <div class="modal show" @click.self="$emit('close')">
      <div class="modal-content bg-white rounded-lg border-2 border-black shadow-brutalist max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
      <!-- Modal Header -->
      <div class="px-6 py-4 border-b-2 border-black bg-brutalist-cream/20">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-black text-gray-900">
            {{ server ? '编辑服务器' : '新建服务器' }}
          </h3>
          <button
            @click="$emit('close')"
            class="text-gray-700 hover:text-gray-900 font-bold"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
      </div>

      <!-- Modal Body -->
      <form @submit.prevent="handleSubmit" class="p-6 space-y-4">
        <div>
          <label class="block text-sm font-bold text-gray-900 mb-2">服务器名称</label>
          <input
            v-model="formData.name"
            type="text"
            required
            class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
            placeholder="例如：我的 Bitwarden"
          />
        </div>

        <!-- 服务器类型选择 -->
        <TabSelector
          v-model="formData.server_type"
          :options="serverTypes"
          label="服务器类型"
        />

        <!-- 服务器区域选择 -->
        <div v-if="formData.server_type === 'official'">
          <TabSelector
            v-model="formData.region"
            :options="officialRegions"
            label="服务器区域"
          />
        </div>

        <!-- 自建服务器地址输入 -->
        <div v-if="formData.server_type === 'self-hosted'">
          <label class="block text-sm font-bold text-gray-900 mb-2">服务器地址</label>
          <input
            v-model="formData.url"
            type="url"
            required
            class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
            placeholder="https://vault.example.com"
          />
          <p class="mt-1 text-xs text-gray-600">
            💡 请填写自建 Bitwarden 服务器的完整地址
          </p>
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-900 mb-2">Client ID</label>
          <input
            v-model="formData.client_id"
            type="text"
            required
            class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
            placeholder="user.xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
          />
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-900 mb-2">Client Secret</label>
          <input
            v-model="formData.client_secret"
            type="password"
            :required="!server"
            class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
            placeholder="••••••••••••••••"
          />
          <p v-if="server" class="mt-1 text-xs text-gray-600">
            💡 留空表示不修改
          </p>
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-900 mb-2">
            Master Password
            <span class="text-xs text-gray-600 font-normal ml-2">
              (用于解锁 Bitwarden 保险库)
            </span>
          </label>
          <input
            v-model="formData.master_password"
            type="password"
            :required="!server"
            class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
            placeholder="••••••••••••••••"
          />
          <p class="mt-1 text-xs text-gray-600">
            <template v-if="server">💡 留空表示不修改</template>
            <template v-else>💡 备份时需要使用 Bitwarden CLI 解锁保险库</template>
          </p>
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
import { ref, watch } from 'vue'
import { serversApi } from '@/api'
import { useToast } from '@/composables/useToast'
import TabSelector from '@/components/ui/TabSelector.vue'

const props = defineProps({
  server: Object
})

const emit = defineEmits(['close', 'saved'])
const toast = useToast()

const serverTypes = [
  { label: '官方服务器', value: 'official' },
  { label: '自建服务器', value: 'self-hosted' }
]

const officialRegions = [
  { label: 'bitwarden.com', value: 'https://vault.bitwarden.com' },
  { label: 'bitwarden.eu', value: 'https://vault.bitwarden.eu' }
]

const formData = ref({
  name: '',
  server_type: 'official',
  region: 'https://vault.bitwarden.com',
  url: '',
  client_id: '',
  client_secret: '',
  master_password: ''
})

const loading = ref(false)

// 根据服务器类型和区域自动设置 URL
watch([() => formData.value.server_type, () => formData.value.region], ([type, region]) => {
  if (type === 'official') {
    formData.value.url = region
  } else if (type === 'self-hosted') {
    // 切换到自建服务器时，清空 URL 让用户输入
    formData.value.url = ''
  }
})

watch(() => props.server, (newServer) => {
  if (newServer) {
    // 修复：从后端的 server_url 字段读取 URL
    const serverUrl = newServer.server_url || newServer.url || ''
    const isOfficial = serverUrl === 'https://vault.bitwarden.com' || serverUrl === 'https://vault.bitwarden.eu'
    formData.value = {
      ...newServer,
      url: serverUrl,  // 映射 server_url 到前端的 url 字段
      server_type: isOfficial ? 'official' : 'self-hosted',
      region: isOfficial ? serverUrl : 'https://vault.bitwarden.com'
    }
  } else {
    formData.value = {
      name: '',
      server_type: 'official',
      region: 'https://vault.bitwarden.com',
      url: 'https://vault.bitwarden.com',
      client_id: '',
      client_secret: '',
      master_password: ''
    }
  }
}, { immediate: true })

const handleSubmit = async () => {
  loading.value = true
  try {
    const submitData = {
      name: formData.value.name,
      server_url: formData.value.url,  // 修复：使用 server_url 匹配后端字段
      client_id: formData.value.client_id,
      client_secret: formData.value.client_secret,
      master_password: formData.value.master_password
    }

    if (props.server?.id) {
      await serversApi.update(props.server.id, submitData)
      toast.success('服务器已更新')
    } else {
      await serversApi.create(submitData)
      toast.success('服务器已创建')
    }
    emit('saved')
  } catch (error) {
    console.error('Failed to save server:', error)
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
