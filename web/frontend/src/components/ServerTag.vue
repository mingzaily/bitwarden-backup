<template>
  <span
    :class="[
      'px-2 py-0.5 text-xs font-bold rounded border-2 border-black',
      isOfficial ? 'bg-brutalist-blue text-white' : 'bg-gray-300 text-gray-700'
    ]"
  >
    {{ isOfficial ? '官方' : '自建' }}
  </span>
</template>

<script setup>
defineProps({
  isOfficial: {
    type: Boolean,
    default: false
  }
})
</script>
