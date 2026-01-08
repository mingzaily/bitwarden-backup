// Toast 全局插件
import { ref } from 'vue'

const toastInstance = ref(null)

export const useToast = () => {
  if (!toastInstance.value) {
    console.warn('Toast not initialized')
    return {
      success: () => {},
      error: () => {},
      warning: () => {},
      info: () => {}
    }
  }
  return toastInstance.value
}

export const setToastInstance = (instance) => {
  toastInstance.value = instance
}
