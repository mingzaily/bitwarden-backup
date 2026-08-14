<template>
  <div class="flex flex-wrap items-center gap-2">
    <span v-for="destination in visibleDestinations" :key="destination.id" :class="getDestinationBadgeClass(destination.type)">
      {{ destination.name }}
    </span>
    <span v-if="remainingCount > 0" class="type-badge type-server">+{{ remainingCount }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  destinations: { type: Array, default: () => [] },
  maxVisible: { type: Number, default: 3 }
})

const visibleDestinations = computed(() => props.destinations.slice(0, props.maxVisible))
const remainingCount = computed(() => Math.max(0, props.destinations.length - props.maxVisible))

const getDestinationBadgeClass = (type) => {
  const classes = {
    local: 'type-badge type-local',
    webdav: 'type-badge type-webdav',
    s3: 'type-badge type-s3',
    server: 'type-badge type-server'
  }
  return classes[type] || 'type-badge type-server'
}
</script>
