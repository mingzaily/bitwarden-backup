<template>
  <span class="inline-flex items-center gap-2">
    <button
      type="button"
      role="switch"
      :aria-checked="modelValue"
      :aria-label="ariaLabel || label || undefined"
      :class="['switch', modelValue ? 'is-on' : '']"
      :disabled="disabled"
      @click="toggle"
    ></button>
    <span v-if="label || $slots.default" class="text-sm font-semibold text-main">
      <slot>{{ label }}</slot>
    </span>
  </span>
</template>

<script setup>
const props = defineProps({
  modelValue: { type: Boolean, default: false },
  label: { type: String, default: '' },
  ariaLabel: { type: String, default: '' },
  disabled: { type: Boolean, default: false }
})

const emit = defineEmits(['update:modelValue'])
const toggle = () => {
  if (!props.disabled) emit('update:modelValue', !props.modelValue)
}
</script>
