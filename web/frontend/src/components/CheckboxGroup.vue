<template>
  <div class="space-y-2">
    <label v-if="label" class="block text-sm font-bold text-gray-900 mb-3">{{ label }}</label>
    <div class="space-y-2">
      <div
        v-for="option in options"
        :key="option.value"
        @click="toggleOption(option.value)"
        :class="[
          'flex items-center gap-3 px-4 py-3 rounded-lg border-2 border-black cursor-pointer transition-all',
          isSelected(option.value)
            ? 'bg-brutalist-blue/10 shadow-brutalist'
            : 'bg-white hover:bg-gray-50'
        ]"
      >
        <div
          :class="[
            'flex-shrink-0 h-5 w-5 rounded border-2 border-black flex items-center justify-center transition-all',
            isSelected(option.value) ? 'bg-brutalist-blue' : 'bg-white'
          ]"
        >
          <svg
            v-if="isSelected(option.value)"
            class="h-3 w-3 text-white"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"></path>
          </svg>
        </div>
        <div class="flex-1">
          <div class="text-sm font-bold text-gray-900">{{ option.label }}</div>
          <div v-if="option.description" class="text-xs text-gray-600 mt-0.5">{{ option.description }}</div>
        </div>
      </div>
    </div>
    <p v-if="options.length === 0" class="text-sm text-gray-600 py-2">
      {{ emptyText }}
    </p>
  </div>
</template>

<script setup>
const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  options: {
    type: Array,
    required: true,
    // options: [{ label: '本地备份', value: 1, description: '...' }, ...]
  },
  label: {
    type: String,
    default: ''
  },
  emptyText: {
    type: String,
    default: '暂无可选项'
  }
})

const emit = defineEmits(['update:modelValue'])

const isSelected = (value) => {
  return props.modelValue.includes(value)
}

const toggleOption = (value) => {
  const newValue = [...props.modelValue]
  const index = newValue.indexOf(value)

  if (index > -1) {
    newValue.splice(index, 1)
  } else {
    newValue.push(value)
  }

  emit('update:modelValue', newValue)
}
</script>
