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
      <div
        v-for="task in tasks"
        :key="task.id"
        :class="[
          'bg-white overflow-hidden rounded-lg border-2 border-black shadow-brutalist hover:shadow-brutalist-hover transition-all',
          'border-l-4',
          task.enabled ? 'border-l-brutalist-green' : 'border-l-gray-400',
          !task.enabled && 'opacity-50'
        ]"
      >
        <div class="px-6 py-4 bg-brutalist-cream/20 min-h-[90px]">
          <div class="flex items-center justify-between">
            <!-- 左侧：任务信息 -->
            <div class="flex-1 space-y-2">
              <div class="flex items-center gap-2 mb-3">
                <h3 class="text-base font-black text-gray-900 leading-6">{{ task.name }}</h3>
                <!-- 任务类型标签 -->
                <span
                  :class="[
                    'px-2 py-0.5 text-xs font-bold rounded border-2 border-black',
                    task.cron_expression ? 'bg-brutalist-blue text-white' : 'bg-gray-300 text-gray-700'
                  ]"
                >
                  {{ task.cron_expression ? '定时' : '手动' }}
                </span>
              </div>

              <!-- 备份流程 -->
              <BackupFlow
                :source-server="task.source_server"
                :destinations="task.destinations"
                class="mb-3"
              />

              <div class="flex items-center text-sm">
                <svg class="flex-shrink-0 mr-2 h-4 w-4 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                <span class="text-gray-700 font-bold leading-4">
                  {{ task.cron_expression ? `Cron: ${task.cron_expression}` : '手动触发' }}
                </span>
              </div>
            </div>

            <!-- 右侧：操作按钮 -->
            <div class="ml-4 shrink-0">
              <div class="inline-flex items-center overflow-hidden rounded border-2 border-black divide-x-2 divide-black">
                <button
                  @click="executeTask(task.id)"
                  class="px-3 py-1 text-sm font-bold text-brutalist-green hover:bg-green-50 transition-colors whitespace-nowrap"
                >
                  立即执行
                </button>
                <button
                  @click="toggleTask(task.id, !task.enabled)"
                  :class="[
                    'px-3 py-1 text-sm font-bold transition-colors whitespace-nowrap',
                    task.enabled
                      ? 'text-gray-700 hover:bg-gray-50'
                      : 'text-brutalist-green hover:bg-green-50'
                  ]"
                >
                  {{ task.enabled ? '禁用' : '启用' }}
                </button>
                <button
                  @click="editTask(task)"
                  class="px-3 py-1 text-sm font-bold text-brutalist-blue hover:bg-blue-50 transition-colors whitespace-nowrap"
                >
                  编辑
                </button>
                <button
                  @click="deleteTask(task.id)"
                  class="px-3 py-1 text-sm font-bold text-brutalist-red hover:bg-red-50 transition-colors whitespace-nowrap"
                >
                  删除
                </button>
              </div>
            </div>
          </div>
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
import BackupFlow from '@/components/BackupFlow.vue'

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
  // 保存原始状态用于回滚
  const taskIndex = tasks.value.findIndex(t => t.id === id)
  const originalEnabled = tasks.value[taskIndex]?.enabled

  try {
    // 立即更新本地状态（乐观更新）
    if (taskIndex !== -1) {
      tasks.value[taskIndex].enabled = enabled
    }

    // 调用 API 更新
    await tasksApi.update(id, { enabled })
    toast.success(enabled ? '任务已启用' : '任务已禁用')

    // 后台同步数据（确保数据一致性）
    loadTasks()
  } catch (error) {
    console.error('Failed to toggle task:', error)

    // 错误时回滚本地状态
    if (taskIndex !== -1 && originalEnabled !== undefined) {
      tasks.value[taskIndex].enabled = originalEnabled
    }

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
