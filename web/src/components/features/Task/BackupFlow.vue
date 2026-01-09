<template>
  <div class="flex flex-wrap items-center gap-2 sm:gap-3">
    <!-- 源服务器 -->
    <div v-if="sourceServer" class="flex items-center gap-2 px-2 py-1 bg-white rounded border border-gray-200 shadow-sm">
      <svg class="flex-shrink-0 h-4 w-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"></path>
      </svg>
      <span class="text-sm font-bold text-gray-900">{{ sourceServer.name }}</span>
      <ServerTag :is-official="isOfficialServer(sourceServer)" class="scale-90 origin-left" />
    </div>

    <!-- 流程箭头 -->
    <div class="flex items-center text-gray-400">
      <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path>
      </svg>
    </div>

    <!-- 目标列表 -->
    <template v-if="destinations && destinations.length > 0">
      <div
        v-for="(dest, index) in visibleDestinations"
        :key="dest.id"
        class="flex items-center gap-2 px-2 py-1 bg-white rounded border border-gray-200 shadow-sm"
      >
        <svg class="flex-shrink-0 h-4 w-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"></path>
        </svg>
        <span class="text-sm font-bold text-gray-900">{{ dest.name }}</span>
        <span :class="getTypeBadgeClass(dest.type)">{{ getTypeLabel(dest.type) }}</span>
      </div>
      <div
        v-if="remainingCount > 0"
        class="px-2 py-1 text-sm font-bold text-gray-500 bg-gray-100 rounded border border-gray-200"
      >
        +{{ remainingCount }}
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import ServerTag from '@/components/features/Server/ServerTag.vue'

const props = defineProps({
  sourceServer: {
    type: Object,
    default: null
  },
  destinations: {
    type: Array,
    default: () => []
  },
  maxVisible: {
    type: Number,
    default: 2
  }
})

const visibleDestinations = computed(() =>
  props.destinations.slice(0, props.maxVisible)
)

const remainingCount = computed(() =>
  Math.max(0, props.destinations.length - props.maxVisible)
)

const isOfficialServer = (server) => {
  if (!server) return false
  return server.server_type === 'official' || server.is_official === true
}

const getTypeLabel = (type) => {
  const labels = {
    'local': '本地',
    'webdav': 'WebDAV',
    's3': 'S3',
    'server': '服务器'
  }
  return labels[type] || type
}

const getTypeBadgeClass = (type) => {
  const base = 'px-1.5 py-0.5 text-xs font-bold rounded border border-black'
  const colors = {
    'local': 'bg-brutalist-green text-white',
    'webdav': 'bg-brutalist-blue text-white',
    's3': 'bg-yellow-500 text-white',
    'server': 'bg-gray-800 text-white'
  }
  return `${base} ${colors[type] || 'bg-gray-100 text-gray-800'}`
}
</script>
