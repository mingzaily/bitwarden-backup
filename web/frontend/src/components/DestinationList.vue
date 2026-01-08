<template>
  <div class="flex items-center gap-2 flex-wrap">
    <span
      v-for="(dest, index) in visibleDestinations"
      :key="dest.id"
      :class="getDestinationBadgeClass(dest.type)"
    >
      {{ dest.name }}
    </span>
    <span
      v-if="remainingCount > 0"
      class="px-2 py-0.5 text-xs font-bold rounded border-2 border-black bg-gray-200 text-gray-700"
    >
      +{{ remainingCount }}
    </span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  destinations: {
    type: Array,
    default: () => []
  },
  maxVisible: {
    type: Number,
    default: 3
  }
})

const visibleDestinations = computed(() =>
  props.destinations.slice(0, props.maxVisible)
)

const remainingCount = computed(() =>
  Math.max(0, props.destinations.length - props.maxVisible)
)

const getDestinationBadgeClass = (type) => {
  const baseClasses = 'px-2 py-0.5 text-xs font-bold rounded border-2 border-black'
  const typeClasses = {
    'local': 'bg-brutalist-green text-white',
    'webdav': 'bg-brutalist-blue text-white',
    's3': 'bg-yellow-500 text-white',
    'server': 'bg-gray-800 text-white'
  }
  return `${baseClasses} ${typeClasses[type] || 'bg-gray-100 text-gray-800'}`
}
</script>
