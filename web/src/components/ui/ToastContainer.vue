<template>
  <Teleport to="body">
    <div
      class="fixed top-5 right-5 z-[10000] flex flex-col gap-3 pointer-events-none"
      aria-live="polite"
      aria-atomic="true"
    >
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          :class="[
            'pointer-events-auto w-full max-w-sm overflow-hidden rounded-lg transform transition-all duration-200',
            'border-2 border-black',
            toast.type === 'success' ? 'bg-green-500 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]' : '',
            toast.type === 'error' ? 'bg-red-500 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]' : '',
            toast.type === 'warning' ? 'bg-yellow-400 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]' : '',
            toast.type === 'info' ? 'bg-blue-500 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]' : ''
          ]"
        >
          <div class="p-4">
            <!-- 修复：使用 flex items-center（水平布局）而非 items-start（竖向布局） -->
            <div class="flex items-center">
              <!-- Icon -->
              <div class="flex-shrink-0">
                <component :is="getIcon(toast.type)" class="w-5 h-5" :class="toast.type === 'warning' ? 'text-gray-900' : 'text-white'" />
              </div>

              <!-- Message -->
              <div class="ml-3 flex-1">
                <p :class="['text-sm font-bold', toast.type === 'warning' ? 'text-gray-900' : 'text-white']">
                  {{ toast.message }}
                </p>
              </div>

              <!-- Close Button - 修复：只显示一个 X -->
              <div class="ml-4 flex-shrink-0">
                <button
                  @click="removeToast(toast.id)"
                  class="bg-white rounded-md p-1 inline-flex items-center justify-center text-gray-900 hover:bg-gray-100 focus:outline-none border-2 border-black shadow-brutalist-sm hover:shadow-none transition-all"
                  aria-label="关闭通知"
                >
                  <svg class="h-4 w-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, h } from 'vue'

const toasts = ref([])
let toastId = 0

const getIcon = (type) => {
  const icons = {
    success: () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M5 13l4 4L19 7' })
    ]),
    error: () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M6 18L18 6M6 6l12 12' })
    ]),
    warning: () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z' })
    ]),
    info: () => h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z' })
    ])
  }
  return icons[type] || icons.info
}

const addToast = (message, type = 'info', duration = 3000) => {
  const id = toastId++

  // 防止重复
  if (toasts.value.some(t => t.message === message)) return

  toasts.value.push({ id, message, type })

  if (duration > 0) {
    setTimeout(() => {
      removeToast(id)
    }, duration)
  }
}

const removeToast = (id) => {
  const index = toasts.value.findIndex(t => t.id === id)
  if (index > -1) {
    toasts.value.splice(index, 1)
  }
}

// 暴露方法供全局使用
defineExpose({
  success: (msg, duration) => addToast(msg, 'success', duration),
  error: (msg, duration = 5000) => addToast(msg, 'error', duration),
  warning: (msg, duration = 4000) => addToast(msg, 'warning', duration),
  info: (msg, duration) => addToast(msg, 'info', duration)
})
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.2s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(2.5rem);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(2.5rem);
}
</style>
