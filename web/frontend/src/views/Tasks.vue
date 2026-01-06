<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h2 class="text-lg font-black text-gray-900">备份任务</h2>
        <p class="text-sm text-gray-700 font-bold">设置自动备份计划和规则</p>
      </div>
      <button @click="showModal = true" class="btn-brutalist inline-flex items-center px-5 py-2.5 text-sm font-bold rounded-lg text-white">
        <svg class="mr-2 -ml-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
        </svg>
        新建任务
      </button>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-brutalist-blue border-t-transparent"></div>
      <p class="mt-2 text-sm text-gray-600">加载中...</p>
    </div>

    <div v-else-if="tasks.length === 0" class="text-center py-12 bg-brutalist-cream/20 rounded-lg border-2 border-black">
      <p class="text-gray-600">暂无备份任务，点击上方按钮添加</p>
    </div>

    <div v-else class="grid gap-4">
      <div v-for="task in tasks" :key="task.id" class="bg-white overflow-hidden rounded-lg border-2 border-black shadow-brutalist hover:shadow-brutalist-hover transition-all">
        <div class="px-4 py-3 flex justify-between items-center border-b-2 border-black">
          <div class="flex items-center gap-2">
            <div class="h-8 w-8 rounded bg-brutalist-orange flex items-center justify-center text-white font-black text-sm border-2 border-black">
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
            <h3 class="text-base font-black text-gray-900">{{ task.name }}</h3>
          </div>
          <button
            @click="toggleTask(task.id, !task.enabled)"
            :class="[
              'px-3 py-1 text-xs font-bold rounded border-2 border-black transition-all',
              task.enabled
                ? 'bg-green-100 text-green-800 hover:bg-green-200'
                : 'bg-gray-100 text-gray-800 hover:bg-gray-200'
            ]"
          >
            {{ task.enabled ? '已启用' : '已禁用' }}
          </button>
        </div>
        <div class="px-4 py-3 bg-brutalist-cream/20">
          <div class="flex items-center text-sm mb-2">
            <svg class="flex-shrink-0 mr-2 h-4 w-4 text-gray-700 -mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <span class="text-gray-700 font-bold">Cron: {{ task.cron_expression }}</span>
          </div>
        </div>
        <div class="bg-white px-4 py-2 border-t-2 border-black flex justify-end gap-2">
          <button @click="executeTask(task.id)" class="px-3 py-1 text-sm font-bold text-brutalist-green hover:bg-green-50 rounded border-2 border-black transition-all">
            立即执行
          </button>
          <button @click="editTask(task)" class="px-3 py-1 text-sm font-bold text-brutalist-blue hover:bg-blue-50 rounded border-2 border-black transition-all">
            编辑
          </button>
          <button @click="deleteTask(task.id)" class="px-3 py-1 text-sm font-bold text-brutalist-red hover:bg-red-50 rounded border-2 border-black transition-all">
            删除
          </button>
        </div>
      </div>
    </div>

    <TaskModal v-if="showModal" :task="editingTask" @close="closeModal" @saved="handleSaved" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { tasksApi } from '@/api'
import { useToast } from '@/composables/useToast'
import TaskModal from '@/components/TaskModal.vue'

const toast = useToast()
const tasks = ref([])
const loading = ref(false)
const showModal = ref(false)
const editingTask = ref(null)

const loadTasks = async () => {
  loading.value = true
  try {
    tasks.value = await tasksApi.getAll()
  } catch (error) {
    console.error('Failed to load tasks:', error)
    toast.error('加载任务列表失败')
  } finally {
    loading.value = false
  }
}

const editTask = (task) => {
  editingTask.value = task
  showModal.value = true
}

const toggleTask = async (id, enabled) => {
  try {
    await tasksApi.update(id, { enabled })
    toast.success(enabled ? '任务已启用' : '任务已禁用')
    loadTasks()
  } catch (error) {
    console.error('Failed to toggle task:', error)
    toast.error('操作失败')
  }
}

const executeTask = async (id) => {
  if (!confirm('确定要立即执行此备份任务吗？')) return

  try {
    await tasksApi.execute(id)
    toast.success('任务已启动，请查看日志')
  } catch (error) {
    console.error('Failed to execute task:', error)
    toast.error('任务启动失败')
  }
}

const deleteTask = async (id) => {
  if (!confirm('确定要删除这个备份任务吗？')) return

  try {
    await tasksApi.delete(id)
    toast.success('任务已删除')
    loadTasks()
  } catch (error) {
    console.error('Failed to delete task:', error)
    toast.error('删除失败')
  }
}

const closeModal = () => {
  showModal.value = false
  editingTask.value = null
}

const handleSaved = () => {
  closeModal()
  loadTasks()
}

onMounted(() => {
  loadTasks()
})
</script>
