import { useDashboardStore } from '~/store/dashboardStore'

export function useDashboard() {
  const store = useDashboardStore()

  return {
    categories: computed(() => store.categories),
    timemachine: computed(() => store.timemachine),
    loading: computed(() => store.loading),
    error: computed(() => store.error),
    loadAll: (options?: { silent?: boolean }) => store.loadAll(options),
    retry: () => store.loadAll()
  }
}
