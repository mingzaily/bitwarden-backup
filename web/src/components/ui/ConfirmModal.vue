<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="confirm-overlay" @click.self="handleCancel">
        <div class="confirm-modal">
          <div class="p-6">
            <div class="flex items-start gap-4">
              <!-- Icon -->
              <div :class="iconClass">
                <svg v-if="type === 'danger'" class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
                <svg v-else class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <!-- Content -->
              <div class="flex-1">
                <h3 class="text-lg font-black text-gray-900">{{ title }}</h3>
                <p class="mt-2 text-sm text-gray-600">{{ message }}</p>
              </div>
            </div>
          </div>
          <!-- Actions -->
          <div class="flex justify-end gap-3 px-6 py-4 border-t-2 border-black bg-gray-50">
            <button
              type="button"
              @click="handleCancel"
              class="px-4 py-2 text-sm font-bold text-gray-700 bg-white border-2 border-black rounded-lg hover:bg-gray-100 transition-all"
            >
              {{ cancelText }}
            </button>
            <button
              type="button"
              @click="handleConfirm"
              :class="confirmButtonClass"
            >
              {{ confirmText }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  visible: { type: Boolean, default: false },
  title: { type: String, default: '确认操作' },
  message: { type: String, default: '确定要执行此操作吗？' },
  type: { type: String, default: 'warning' }, // 'warning' | 'danger'
  confirmText: { type: String, default: '确定' },
  cancelText: { type: String, default: '取消' }
})

const emit = defineEmits(['confirm', 'cancel'])

const iconClass = computed(() => {
  const base = 'flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center border-2 border-black'
  return props.type === 'danger'
    ? `${base} bg-red-100 text-red-600`
    : `${base} bg-yellow-100 text-yellow-600`
})

const confirmButtonClass = computed(() => {
  const base = 'px-4 py-2 text-sm font-bold text-white border-2 border-black rounded-lg transition-all'
  return props.type === 'danger'
    ? `${base} bg-red-500 hover:bg-red-600`
    : `${base} bg-brutalist-blue hover:bg-blue-600`
})

const handleConfirm = () => emit('confirm')
const handleCancel = () => emit('cancel')
</script>

<style scoped>
.confirm-overlay {
  position: fixed;
  inset: 0;
  z-index: 10000;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
}

.confirm-modal {
  background: white;
  border: 2px solid black;
  border-radius: 0.5rem;
  box-shadow: 4px 4px 0 0 rgba(0, 0, 0, 1);
  max-width: 28rem;
  width: calc(100% - 2rem);
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-active .confirm-modal,
.modal-leave-active .confirm-modal {
  transition: transform 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .confirm-modal,
.modal-leave-to .confirm-modal {
  transform: scale(0.95);
}
</style>
