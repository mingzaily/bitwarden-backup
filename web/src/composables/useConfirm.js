// Confirm 全局插件
import { ref } from 'vue'

const confirmState = ref({
  visible: false,
  title: '确认操作',
  message: '确定要执行此操作吗？',
  type: 'warning',
  confirmText: '确定',
  cancelText: '取消',
  resolve: null
})

export const useConfirm = () => {
  const confirm = (options = {}) => {
    return new Promise((resolve) => {
      confirmState.value = {
        visible: true,
        title: options.title || '确认操作',
        message: options.message || '确定要执行此操作吗？',
        type: options.type || 'warning',
        confirmText: options.confirmText || '确定',
        cancelText: options.cancelText || '取消',
        resolve
      }
    })
  }

  const handleConfirm = () => {
    if (confirmState.value.resolve) {
      confirmState.value.resolve(true)
    }
    confirmState.value.visible = false
  }

  const handleCancel = () => {
    if (confirmState.value.resolve) {
      confirmState.value.resolve(false)
    }
    confirmState.value.visible = false
  }

  return {
    state: confirmState,
    confirm,
    handleConfirm,
    handleCancel
  }
}
