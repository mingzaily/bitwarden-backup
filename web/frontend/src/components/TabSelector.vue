<template>
  <div class="space-y-2">
    <label v-if="label" class="block text-sm font-bold text-gray-900">{{ label }}</label>
    <div class="inline-flex rounded-lg border-2 border-black bg-white p-1 gap-1">
      <button
        v-for="option in options"
        :key="option.value"
        type="button"
        @click="selectOption(option.value)"
        :class="[
          'px-4 py-2 text-sm font-bold rounded transition-all',
          modelValue === option.value
            ? 'bg-brutalist-blue text-white shadow-brutalist'
            : 'text-gray-700 hover:bg-gray-50'
        ]"
        :disabled="disabled"
      >
        {{ option.label }}
      </button>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  modelValue: {
    type: String,
    required: true
  },
  options: {
    type: Array,
    required: true,
    // options: [{ label: '本地存储', value: 'local' }, ...]
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

const selectOption = (value) => {
  if (!props.disabled) {
    emit('update:modelValue', value)
  }
}
</script>
