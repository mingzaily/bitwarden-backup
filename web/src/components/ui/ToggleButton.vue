<template>
  <button
    type="button"
    @click="toggle"
    :class="[
      'relative inline-flex items-center px-4 py-2 text-sm font-bold rounded-lg border-2 border-black transition-all',
      modelValue
        ? 'bg-brutalist-blue text-white shadow-brutalist hover:shadow-brutalist-hover'
        : 'bg-white text-gray-700 hover:bg-gray-50'
    ]"
    :disabled="disabled"
  >
    <span class="flex items-center gap-2">
      <svg
        v-if="modelValue"
        class="h-4 w-4"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
      </svg>
      <svg
        v-else
        class="h-4 w-4"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
      </svg>
      <slot>{{ label }}</slot>
    </span>
  </button>
</template>

<script setup>
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  label: {
    type: String,
    default: ''
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const toggle = () => {
  if (!props.disabled) {
    emit('update:modelValue', !props.modelValue)
  }
}
</script>
