<template>
  <div v-if="totalPage > 1" class="pagination">
    <div class="pagination-info">共 {{ total }} 条 · 第 {{ page }} / {{ totalPage }} 页</div>
    <div class="pagination-controls" aria-label="分页">
      <button class="pagination-button" :disabled="page <= 1" title="首页" aria-label="首页" @click="goToPage(1)">
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="m11 19-7-7 7-7m8 14-7-7 7-7" /></svg>
      </button>
      <button class="pagination-button" :disabled="page <= 1" title="上一页" aria-label="上一页" @click="goToPage(page - 1)">
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="m15 19-7-7 7-7" /></svg>
      </button>

      <template v-for="(p, index) in visiblePages" :key="`${p}-${index}`">
        <span v-if="p === '...'" class="px-1 text-subtle">…</span>
        <button v-else :class="['pagination-button', p === page ? 'is-active' : '']" :aria-current="p === page ? 'page' : undefined" @click="goToPage(p)">{{ p }}</button>
      </template>

      <button class="pagination-button" :disabled="page >= totalPage" title="下一页" aria-label="下一页" @click="goToPage(page + 1)">
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="m9 5 7 7-7 7" /></svg>
      </button>
      <button class="pagination-button" :disabled="page >= totalPage" title="末页" aria-label="末页" @click="goToPage(totalPage)">
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="m13 5 7 7-7 7M5 5l7 7-7 7" /></svg>
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

const emit = defineEmits(['change', 'page-change'])

const goToPage = (nextPage) => {
  if (nextPage >= 1 && nextPage <= props.totalPage && nextPage !== props.page) {
    emit('change', nextPage)
    emit('page-change', nextPage)
  }
}

const visiblePages = computed(() => {
  const pages = []
  const total = props.totalPage
  const current = props.page
  if (total <= 7) {
    for (let index = 1; index <= total; index += 1) pages.push(index)
  } else {
    pages.push(1)
    if (current > 3) pages.push('...')
    for (let index = Math.max(2, current - 1); index <= Math.min(total - 1, current + 1); index += 1) pages.push(index)
    if (current < total - 2) pages.push('...')
    pages.push(total)
  }
  return pages
})
</script>
