<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-backdrop" @click.self="handleCancel">
        <div class="modal-panel max-w-md" role="dialog" aria-modal="true" :aria-label="title">
          <div class="modal-body">
            <div class="flex items-start gap-4">
              <div :class="iconClass" aria-hidden="true">
                <svg v-if="type === 'danger'" class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>
                <svg v-else class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M8.23 9A4 4 0 0 1 12 7c2.21 0 4 1.34 4 3 0 1.4-1.28 2.58-3.01 2.91-.54.1-.99.54-.99 1.09m0 3h.01M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" /></svg>
              </div>
              <div class="min-w-0 flex-1">
                <h3 class="modal-title">{{ title }}</h3>
                <p class="modal-subtitle">{{ message }}</p>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn-secondary" @click="handleCancel">{{ cancelText }}</button>
            <button type="button" :class="confirmButtonClass" @click="handleConfirm">{{ confirmText }}</button>
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
  type: { type: String, default: 'warning' },
  confirmText: { type: String, default: '确定' },
  cancelText: { type: String, default: '取消' }
})

const emit = defineEmits(['confirm', 'cancel'])
const iconClass = computed(() => props.type === 'danger' ? 'grid h-10 w-10 flex-shrink-0 place-items-center rounded-xl bg-danger/10 text-danger' : 'grid h-10 w-10 flex-shrink-0 place-items-center rounded-xl bg-warning/10 text-warning')
const confirmButtonClass = computed(() => props.type === 'danger' ? 'btn-danger' : 'btn-primary')
const handleConfirm = () => emit('confirm')
const handleCancel = () => emit('cancel')
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active { transition: opacity 180ms ease; }
.modal-enter-active .modal-panel,
.modal-leave-active .modal-panel { transition: transform 180ms ease; }
.modal-enter-from,
.modal-leave-to { opacity: 0; }
.modal-enter-from .modal-panel,
.modal-leave-to .modal-panel { transform: scale(0.97) translateY(0.35rem); }
</style>
