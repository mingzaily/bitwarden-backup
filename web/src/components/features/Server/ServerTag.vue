<template>
  <span :class="['server-badge', isOfficial ? 'server-official' : 'server-self']">
    {{ isOfficial ? '官方' : '自建' }}
  </span>
</template>

<script setup>
defineProps({
  isOfficial: { type: Boolean, default: false }
})
</script>
