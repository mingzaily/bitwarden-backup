<template>
  <div ref="container" class="field relative">
    <label v-if="label" class="field-label">{{ label }}</label>

    <button
      type="button"
      class="input-select flex items-center justify-between gap-3 text-left"
      role="combobox"
      :aria-expanded="isOpen"
      :aria-label="label || placeholder"
      @click="toggle"
      @keydown.esc="isOpen = false"
    >
      <span v-if="selectedOption" class="truncate font-semibold text-main">{{ selectedOption.label }}</span>
      <span v-else class="truncate text-muted">{{ placeholder }}</span>
      <svg :class="['h-4 w-4 flex-shrink-0 text-subtle transition-transform', isOpen ? 'rotate-180' : '']" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="m6 9 6 6 6-6" />
      </svg>
    </button>

    <div v-if="isOpen" class="select-menu" role="listbox" :aria-label="label || placeholder">
      <div v-if="options.length === 0" class="px-3 py-3 text-center text-sm text-muted">{{ emptyText }}</div>
      <button
        v-for="option in options"
        :key="option.value"
        type="button"
        role="option"
        :aria-selected="option.value === modelValue"
        :class="['select-option', option.value === modelValue ? 'is-selected' : '']"
        @click="selectOption(option)"
      >
        <span class="min-w-0">
          <span class="block truncate text-sm font-semibold">{{ option.label }}</span>
          <span v-if="option.description" class="mt-0.5 block truncate text-xs text-muted">{{ option.description }}</span>
        </span>
        <svg v-if="option.value === modelValue" class="h-4 w-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m5 12 4 4L19 6" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'

const props = defineProps({
  modelValue: { type: [String, Number], default: '' },
  options: { type: Array, required: true },
  label: { type: String, default: '' },
  placeholder: { type: String, default: '请选择' },
  emptyText: { type: String, default: '暂无可选项' }
})

const emit = defineEmits(['update:modelValue'])
const container = ref(null)
const isOpen = ref(false)
const selectedOption = computed(() => props.options.find(option => option.value === props.modelValue))

const toggle = () => { isOpen.value = !isOpen.value }
const selectOption = (option) => {
  emit('update:modelValue', option.value)
  isOpen.value = false
}
const handleClickOutside = (event) => {
  if (container.value && !container.value.contains(event.target)) isOpen.value = false
}

onMounted(() => document.addEventListener('click', handleClickOutside))
onUnmounted(() => document.removeEventListener('click', handleClickOutside))
</script>
