<template>
  <div id="app" class="min-h-screen bg-brutalist-cream p-6 flex items-start justify-center">
    <!-- 整个应用的大卡片容器 -->
    <div class="w-full max-w-7xl bg-white rounded-lg border-2 border-black shadow-brutalist overflow-hidden">
      <!-- Header -->
      <header class="bg-white border-b-2 border-black px-6 py-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="h-10 w-10 bg-brutalist-blue rounded border-2 border-black flex items-center justify-center">
            <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
            </svg>
          </div>
          <h1 class="text-xl font-black text-gray-900">Bitwarden 备份助手</h1>
        </div>
      </div>

      <!-- Tab Navigation -->
      <nav class="flex gap-2 mt-6 border-b-2 border-black" role="tablist" aria-label="主导航标签">
        <router-link
          v-for="tab in tabs"
          :key="tab.path"
          :to="tab.path"
          class="tab-btn px-4 pb-3 pt-2 text-sm font-bold text-gray-700 hover:text-brutalist-blue rounded-t-lg"
          active-class="active"
        >
          {{ tab.label }}
        </router-link>
      </nav>
    </header>

    <!-- Main Content -->
    <main class="flex-1 overflow-y-auto bg-brutalist-cream/30 p-6">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- Toast Container -->
    <ToastContainer ref="toastRef" />

    <!-- Confirm Modal -->
    <ConfirmModal
      :visible="confirmState.visible"
      :title="confirmState.title"
      :message="confirmState.message"
      :type="confirmState.type"
      :confirm-text="confirmState.confirmText"
      :cancel-text="confirmState.cancelText"
      @confirm="handleConfirm"
      @cancel="handleCancel"
    />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import ToastContainer from '@/components/ui/ToastContainer.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import { setToastInstance } from '@/composables/useToast'
import { useConfirm } from '@/composables/useConfirm'

const { state: confirmState, handleConfirm, handleCancel } = useConfirm()

const tabs = ref([
  { path: '/servers', label: '服务器列表' },
  { path: '/destinations', label: '备份目标' },
  { path: '/tasks', label: '备份任务' },
  { path: '/logs', label: '运行日志' }
])

const toastRef = ref(null)

onMounted(() => {
  if (toastRef.value) {
    setToastInstance(toastRef.value)
  }
})
</script>

<style scoped>
/* Tab Button Styles - 修复抖动问题 */
.tab-btn {
  position: relative;
  /* 只对 color 和 transform 添加过渡，避免 all 导致的布局重排 */
  transition: color 0.15s ease, transform 0.15s ease;
}

.tab-btn.active {
  color: #2563EB;
  background-color: white;
  border-left: 2px solid black;
  border-right: 2px solid black;
  border-top: 2px solid black;
  border-bottom: 2px solid white;
  margin-bottom: -2px;
}

/* Fade Transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
