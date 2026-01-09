<template>
  <div v-if="totalPage > 1" class="flex items-center justify-between mt-4 pt-4 border-t-2 border-black">
    <!-- 左侧：分页信息 -->
    <div class="text-sm font-bold text-gray-600">
      共 {{ total }} 条，第 {{ page }} / {{ totalPage }} 页
    </div>

    <!-- 右侧：分页按钮 -->
    <div class="flex items-center gap-1">
      <!-- 首页 -->
      <button
        @click="goToPage(1)"
        :disabled="page <= 1"
        :class="buttonClass(page <= 1)"
        title="首页"
      >
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
        </svg>
      </button>

      <!-- 上一页 -->
      <button
        @click="goToPage(page - 1)"
        :disabled="page <= 1"
        :class="buttonClass(page <= 1)"
        title="上一页"
      >
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>

      <!-- 页码按钮 -->
      <template v-for="p in visiblePages" :key="p">
        <span v-if="p === '...'" class="px-2 text-gray-500">...</span>
        <button
          v-else
          @click="goToPage(p)"
          :class="pageButtonClass(p === page)"
        >
          {{ p }}
        </button>
      </template>

      <!-- 下一页 -->
      <button
        @click="goToPage(page + 1)"
        :disabled="page >= totalPage"
        :class="buttonClass(page >= totalPage)"
        title="下一页"
      >
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>

      <!-- 末页 -->
      <button
        @click="goToPage(totalPage)"
        :disabled="page >= totalPage"
        :class="buttonClass(page >= totalPage)"
        title="末页"
      >
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 5l7 7-7 7M5 5l7 7-7 7" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  page: { type: Number, default: 1 },
  pageSize: { type: Number, default: 10 },
  total: { type: Number, default: 0 },
  totalPage: { type: Number, default: 0 }
})

const emit = defineEmits(['change'])

const goToPage = (p) => {
  if (p >= 1 && p <= props.totalPage && p !== props.page) {
    emit('change', p)
  }
}

const visiblePages = computed(() => {
  const pages = []
  const total = props.totalPage
  const current = props.page

  if (total <= 7) {
    for (let i = 1; i <= total; i++) pages.push(i)
  } else {
    pages.push(1)
    if (current > 3) pages.push('...')

    const start = Math.max(2, current - 1)
    const end = Math.min(total - 1, current + 1)
    for (let i = start; i <= end; i++) pages.push(i)

    if (current < total - 2) pages.push('...')
    pages.push(total)
  }
  return pages
})

const buttonClass = (disabled) => [
  'p-2 rounded border-2 border-black transition-all',
  disabled
    ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
    : 'bg-white text-gray-700 hover:bg-gray-50 hover:shadow-brutalist-sm'
]

const pageButtonClass = (active) => [
  'min-w-[36px] h-9 px-2 rounded border-2 border-black font-bold transition-all',
  active
    ? 'bg-brutalist-blue text-white shadow-brutalist-sm'
    : 'bg-white text-gray-700 hover:bg-gray-50'
]
</script>
