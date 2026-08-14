<template>
  <div class="field">
    <div v-if="label" class="field-label">{{ label }}</div>
    <div class="check-list" role="group" :aria-label="label || '选项列表'">
      <button
        v-for="option in options"
        :key="option.value"
        type="button"
        role="checkbox"
        :aria-checked="isSelected(option.value)"
        :class="['check-option text-left', isSelected(option.value) ? 'is-selected' : '']"
        @click="toggleOption(option.value)"
      >
        <span class="check-box" aria-hidden="true">
          <svg v-if="isSelected(option.value)" class="h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="m5 12.5 4 4L19 7" />
          </svg>
        </span>
        <span class="min-w-0 flex-1">
          <span class="block text-sm font-semibold text-main">{{ option.label }}</span>
          <span v-if="option.description" class="mt-0.5 block text-xs text-muted">{{ option.description }}</span>
        </span>
      </button>
    </div>
    <p v-if="options.length === 0" class="field-hint py-2">{{ emptyText }}</p>
  </div>
</template>

<script setup>
const props = defineProps({
  modelValue: { type: Array, default: () => [] },
  options: { type: Array, required: true },
  label: { type: String, default: '' },
  emptyText: { type: String, default: '暂无可选项' }
})

const emit = defineEmits(['update:modelValue'])

const isSelected = (value) => props.modelValue.includes(value)

const toggleOption = (value) => {
  const nextValue = [...props.modelValue]
  const index = nextValue.indexOf(value)
  if (index > -1) nextValue.splice(index, 1)
  else nextValue.push(value)
  emit('update:modelValue', nextValue)
}
</script>
