<template>
  <div class="field">
    <div v-if="label" class="field-label">{{ label }}</div>
    <div class="tabs" role="tablist" :aria-label="label || '选项'">
      <button
        v-for="option in options"
        :key="option.value"
        type="button"
        role="tab"
        :aria-selected="modelValue === option.value"
        :class="['tab-button', modelValue === option.value ? 'is-active' : '']"
        :disabled="disabled"
        @click="selectOption(option.value)"
      >
        {{ option.label }}
      </button>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  modelValue: { type: String, required: true },
  options: { type: Array, required: true },
  label: { type: String, default: '' },
  disabled: { type: Boolean, default: false }
})

const emit = defineEmits(['update:modelValue'])
const selectOption = (value) => {
  if (!props.disabled) emit('update:modelValue', value)
}
</script>
