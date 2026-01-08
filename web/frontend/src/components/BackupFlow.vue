<template>
  <div class="space-y-2">
    <!-- 源服务器 -->
    <div v-if="sourceServer" class="flex items-center gap-2">
      <svg class="flex-shrink-0 h-4 w-4 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"></path>
      </svg>
      <span class="text-sm font-bold text-gray-900">{{ sourceServer.name }}</span>
      <ServerTag :is-official="isOfficialServer(sourceServer)" />
    </div>

    <!-- 流程箭头 -->
    <div class="flex items-center pl-2">
      <svg class="h-5 w-5 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3"></path>
      </svg>
    </div>

    <!-- 目标列表 -->
    <div v-if="destinations && destinations.length > 0" class="flex items-center gap-2">
      <svg class="flex-shrink-0 h-4 w-4 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"></path>
      </svg>
      <DestinationList :destinations="destinations" :max-visible="3" />
    </div>
  </div>
</template>

<script setup>
import ServerTag from './ServerTag.vue'
import DestinationList from './DestinationList.vue'

const props = defineProps({
  sourceServer: {
    type: Object,
    default: null
  },
  destinations: {
    type: Array,
    default: () => []
  }
})

const isOfficialServer = (server) => {
  if (!server) return false
  return server.server_type === 'official' || server.is_official === true
}
</script>
