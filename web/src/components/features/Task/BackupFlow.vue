<template>
  <div class="flex flex-wrap items-center gap-2 sm:gap-3">
    <div v-if="sourceServer" class="surface-muted flex items-center gap-2 px-2.5 py-2">
      <svg class="h-4 w-4 flex-shrink-0 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M5 12h14M5 12a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v4M5 12a2 2 0 0 0-2 2v4a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-4a2 2 0 0 0-2-2" /></svg>
      <span class="text-sm font-semibold text-main">{{ sourceServer.name }}</span>
      <ServerTag :is-official="isOfficialServer(sourceServer)" class="scale-90 origin-left" />
    </div>

    <div class="flex items-center text-subtle" aria-hidden="true"><svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="m14 5 7 7-7 7m7-7H3" /></svg></div>

    <template v-if="destinations && destinations.length > 0">
      <div v-for="dest in visibleDestinations" :key="dest.id" class="surface-muted flex items-center gap-2 px-2.5 py-2">
        <svg class="h-4 w-4 flex-shrink-0 text-info" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M3 7v10a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-6l-2-2H5a2 2 0 0 0-2 2Z" /></svg>
        <span class="text-sm font-semibold text-main">{{ dest.name }}</span>
        <span :class="getTypeBadgeClass(dest.type)">{{ getTypeLabel(dest.type) }}</span>
      </div>
      <div v-if="remainingCount > 0" class="type-badge type-server">+{{ remainingCount }}</div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import ServerTag from '@/components/features/Server/ServerTag.vue'

const props = defineProps({
  sourceServer: { type: Object, default: null },
  destinations: { type: Array, default: () => [] },
  maxVisible: { type: Number, default: 2 }
})

const visibleDestinations = computed(() => props.destinations.slice(0, props.maxVisible))
const remainingCount = computed(() => Math.max(0, props.destinations.length - props.maxVisible))
const isOfficialServer = (server) => Boolean(server && (
  server.server_type === 'official'
  || server.is_official === true
  || ['https://vault.bitwarden.com', 'https://vault.bitwarden.eu'].includes(server.server_url || server.url)
))
const getTypeLabel = (type) => ({ local: '本地', webdav: 'WebDAV', s3: 'S3', server: '服务器' }[type] || type)
const getTypeBadgeClass = (type) => ({
  local: 'type-badge type-local',
  webdav: 'type-badge type-webdav',
  s3: 'type-badge type-s3',
  server: 'type-badge type-server'
}[type] || 'type-badge type-server')
</script>
