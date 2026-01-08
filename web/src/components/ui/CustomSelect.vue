<template>
  <div class="relative" ref="container">
    <label v-if="label" class="block text-sm font-bold text-gray-900 mb-2">{{ label }}</label>
    
    <!-- Trigger Button -->
    <button
      type="button"
      @click="toggle"
      :class="[
        'w-full px-3 py-2 border-2 border-black rounded-lg bg-white flex items-center justify-between focus:outline-none transition-all',
        isOpen ? 'ring-2 ring-brutalist-blue' : ''
      ]"
    >
      <span v-if="selectedOption" class="text-sm font-bold text-gray-900 truncate">
        {{ selectedOption.label }}
      </span>
      <span v-else class="text-sm text-gray-500 truncate">
        {{ placeholder }}
      </span>
      <svg
        :class="['h-5 w-5 text-gray-700 transition-transform', isOpen ? 'rotate-180' : '']"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
      </svg>
    </button>

    <!-- Dropdown Menu -->
    <div
      v-if="isOpen"
      class="absolute z-10 w-full mt-1 bg-white border-2 border-black rounded-lg shadow-brutalist max-h-60 overflow-y-auto"
    >
      <div v-if="options.length === 0" class="px-4 py-3 text-sm text-gray-500 text-center">
        {{ emptyText }}
      </div>
      <div
        v-for="option in options"
        :key="option.value"
        @click="selectOption(option)"
        :class="[
          'px-4 py-2.5 text-sm font-bold cursor-pointer transition-colors hover:bg-gray-50',
          option.value === modelValue ? 'bg-brutalist-blue/10 text-brutalist-blue' : 'text-gray-900'
        ]"
      >
        <div class="flex items-center justify-between">
          <span>{{ option.label }}</span>
          <svg v-if="option.value === modelValue" class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
          </svg>
        </div>
        <div v-if="option.description" class="text-xs text-gray-500 font-normal mt-0.5">
          {{ option.description }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  options: {
    type: Array,
    required: true,
    // options: [{ label: 'Option 1', value: 1, description: 'Desc' }, ...]
  },
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: '请选择'
  },
  emptyText: {
    type: String,
    default: '暂无可选项'
  }
})

const emit = defineEmits(['update:modelValue'])
const container = ref(null)
const isOpen = ref(false)

const selectedOption = computed(() => {
  return props.options.find(opt => opt.value === props.modelValue)
})

const toggle = () => {
  isOpen.value = !isOpen.value
}

const selectOption = (option) => {
  emit('update:modelValue', option.value)
  isOpen.value = false
}

const handleClickOutside = (event) => {
  if (container.value && !container.value.contains(event.target)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
