<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h2 class="text-lg font-black text-gray-900">Bitwarden 服务器</h2>
        <p class="text-sm text-gray-700 font-bold">管理 Bitwarden 服务器连接信息</p>
      </div>
      <button
        @click="showModal = true"
        class="btn-brutalist inline-flex items-center px-5 py-2.5 text-sm font-bold rounded-lg text-white"
      >
        <svg class="mr-2 -ml-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
        </svg>
        新建服务器
      </button>
    </div>

    <!-- Server List -->
    <div v-if="loading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-brutalist-blue border-t-transparent"></div>
      <p class="mt-2 text-sm text-gray-600">加载中...</p>
    </div>

    <div v-else-if="servers.length === 0" class="text-center py-12 bg-white rounded-lg border-2 border-black">
      <p class="text-gray-600">暂无服务器，点击上方按钮添加</p>
    </div>

    <div v-else class="grid gap-4">
      <div
        v-for="server in servers"
        :key="server.id"
        class="bg-white overflow-hidden rounded-lg border-2 border-black shadow-brutalist hover:shadow-brutalist-hover transition-all"
      >
        <div class="px-4 py-3 bg-brutalist-cream/20">
          <div class="flex items-center justify-between">
            <!-- 左侧：服务器信息 -->
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-2">
                <h3 class="text-base font-black text-gray-900">{{ server.name }}</h3>
                <!-- 服务器类型标签 -->
                <span
                  :class="[
                    'px-2 py-0.5 text-xs font-bold rounded border-2 border-black',
                    isOfficialServer(server)
                      ? 'bg-brutalist-blue text-white'
                      : 'bg-brutalist-green text-white'
                  ]"
                >
                  {{ isOfficialServer(server) ? '官方' : '自建' }}
                </span>
              </div>
              <div class="flex items-center text-sm mb-1">
                <svg class="flex-shrink-0 mr-2 h-4 w-4 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"></path>
                </svg>
                <span class="text-gray-700 font-bold break-all">{{ server.server_url || server.url }}</span>
              </div>
              <div class="flex items-center text-sm">
                <svg class="flex-shrink-0 mr-2 h-4 w-4 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                </svg>
                <span class="text-gray-700 font-bold break-all">ID: {{ server.client_id || 'N/A' }}</span>
              </div>
            </div>

            <!-- 右侧：操作按钮 -->
            <div class="flex items-center gap-2 ml-4">
              <button
                @click="toggleServer(server.id, !server.enabled)"
                :class="[
                  'px-3 py-1 text-sm font-bold rounded border-2 border-black transition-all',
                  server.enabled
                    ? 'text-gray-700 hover:bg-gray-50'
                    : 'text-brutalist-green hover:bg-green-50'
                ]"
              >
                {{ server.enabled ? '禁用' : '启用' }}
              </button>
              <button
                @click="editServer(server)"
                class="px-3 py-1 text-sm font-bold text-brutalist-blue hover:bg-blue-50 rounded border-2 border-black transition-all"
              >
                编辑
              </button>
              <button
                @click="deleteServer(server.id)"
                class="px-3 py-1 text-sm font-bold text-brutalist-red hover:bg-red-50 rounded border-2 border-black transition-all"
              >
                删除
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Server Modal -->
    <ServerModal
      v-if="showModal"
      :server="editingServer"
      @close="closeModal"
      @saved="handleSaved"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { serversApi } from '@/api'
import { useToast } from '@/composables/useToast'
import ServerModal from '@/components/ServerModal.vue'

const toast = useToast()
const servers = ref([])
const loading = ref(false)
const showModal = ref(false)
const editingServer = ref(null)

const loadServers = async () => {
  loading.value = true
  try {
    servers.value = await serversApi.getAll()
  } catch (error) {
    console.error('Failed to load servers:', error)
    toast.error('加载服务器列表失败')
  } finally {
    loading.value = false
  }
}

const editServer = (server) => {
  editingServer.value = server
  showModal.value = true
}

const isOfficialServer = (server) => {
  const officialUrls = ['https://vault.bitwarden.com', 'https://vault.bitwarden.eu']
  return officialUrls.includes(server.server_url || server.url)
}

const toggleServer = async (id, enabled) => {
  try {
    await serversApi.update(id, { enabled })
    toast.success(enabled ? '已启用' : '已禁用')
    loadServers()
  } catch (error) {
    console.error('Failed to toggle server:', error)
    toast.error('操作失败')
  }
}

const deleteServer = async (id) => {
  if (!confirm('确定要删除这个服务器吗？')) return

  try {
    await serversApi.delete(id)
    toast.success('服务器已删除')
    loadServers()
  } catch (error) {
    console.error('Failed to delete server:', error)
    toast.error('删除失败')
  }
}

const closeModal = () => {
  showModal.value = false
  editingServer.value = null
}

const handleSaved = () => {
  closeModal()
  loadServers()
}

onMounted(() => {
  loadServers()
})
</script>
