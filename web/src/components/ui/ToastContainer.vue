<template>
  <Teleport to="body">
    <div class="pointer-events-none fixed right-4 top-4 z-[10000] flex flex-col gap-3" aria-live="polite" aria-atomic="true">
      <TransitionGroup name="toast">
        <div v-for="toast in toasts" :key="toast.id" :class="['toast pointer-events-auto', `toast-${toast.type}`]">
          <div class="toast-icon" aria-hidden="true"><component :is="getIcon(toast.type)" class="h-4 w-4" /></div>
          <p class="toast-message">{{ toast.message }}</p>
          <button class="toast-close" type="button" aria-label="关闭通知" @click="removeToast(toast.id)">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="m6 6 12 12M18 6 6 18" /></svg>
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup>
import { h, ref } from 'vue'

const toasts = ref([])
let toastId = 0

const iconPath = {
  success: 'm5 12.5 4 4L19 7',
  error: 'M6 18 18 6M6 6l12 12',
  warning: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z',
  info: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0z'
}

const getIcon = (type) => () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: iconPath[type] || iconPath.info })
])

const addToast = (message, type = 'info', duration = 3000) => {
  if (toasts.value.some(toast => toast.message === message)) return
  const id = toastId++
  toasts.value.push({ id, message, type })
  if (duration > 0) window.setTimeout(() => removeToast(id), duration)
}

const removeToast = (id) => {
  const index = toasts.value.findIndex(toast => toast.id === id)
  if (index > -1) toasts.value.splice(index, 1)
}

defineExpose({
  success: (message, duration) => addToast(message, 'success', duration),
  error: (message, duration = 5000) => addToast(message, 'error', duration),
  warning: (message, duration = 4000) => addToast(message, 'warning', duration),
  info: (message, duration) => addToast(message, 'info', duration)
})
</script>
