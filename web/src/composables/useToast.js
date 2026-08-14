// Toast 全局插件
import { ref } from 'vue'

const toastInstance = ref(null)

export const useToast = () => {
  // Return a lazy proxy so setup() can safely call useToast before the
  // Teleported ToastContainer has mounted and assigned its component ref.
  return {
    success: (...args) => toastInstance.value?.success?.(...args),
    error: (...args) => toastInstance.value?.error?.(...args),
    warning: (...args) => toastInstance.value?.warning?.(...args),
    info: (...args) => toastInstance.value?.info?.(...args)
  }
}

export const setToastInstance = (instance) => {
  toastInstance.value = instance
}
