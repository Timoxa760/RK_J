import { toast } from 'vue-sonner'

export type ToastType = 'error' | 'success' | 'info'

export function useToast() {
  function show(message: string, type: ToastType = 'info') {
    if (type === 'success') toast.success(message)
    else if (type === 'error') toast.error(message)
    else toast.info(message)
  }

  return { show }
}
