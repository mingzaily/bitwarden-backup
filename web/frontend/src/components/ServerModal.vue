<template>
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

        <div>
          <label class="block text-sm font-bold text-gray-900 mb-2">服务器地址</label>
          <input
            v-model="formData.url"
            type="url"
            required
            class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
            placeholder="https://vault.bitwarden.com"
          />
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-900 mb-2">Client ID</label>
          <input
            v-model="formData.client_id"
            type="text"
            required
            class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
            placeholder="organization.xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
          />
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-900 mb-2">Client Secret</label>
          <input
            v-model="formData.client_secret"
            type="password"
            required
            class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
            placeholder="••••••••••••••••"
          />
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
            required
            class="w-full px-3 py-2 border-2 border-black rounded-lg focus:outline-none focus:ring-2 focus:ring-brutalist-blue"
            placeholder="••••••••••••••••"
          />
          <p class="mt-1 text-xs text-gray-600">
            💡 为什么需要主密码？备份时需要使用 Bitwarden CLI 的 <code class="px-1 py-0.5 bg-gray-100 rounded text-xs">bw unlock</code> 命令解锁保险库，才能导出数据进行备份。
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
</template>

<script setup>
import { ref, watch } from 'vue'
import { serversApi } from '@/api'
import { useToast } from '@/composables/useToast'

const props = defineProps({
  server: Object
})

const emit = defineEmits(['close', 'saved'])
const toast = useToast()

const formData = ref({
  name: '',
  url: '',
  client_id: '',
  client_secret: '',
  master_password: ''
})

const loading = ref(false)

watch(() => props.server, (newServer) => {
  if (newServer) {
    formData.value = { ...newServer }
  } else {
    formData.value = {
      name: '',
      url: '',
      client_id: '',
      client_secret: '',
      master_password: ''
    }
  }
}, { immediate: true })

const handleSubmit = async () => {
  loading.value = true
  try {
    if (props.server?.id) {
      await serversApi.update(props.server.id, formData.value)
      toast.success('服务器已更新')
    } else {
      await serversApi.create(formData.value)
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
  inset: 0;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 50;
  background-color: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  align-items: center;
  justify-content: center;
}
</style>
