import { useDashboardStore } from '~/store/dashboardStore'

export function useDashboard() {
  const store = useDashboardStore()

  return {
    sankey: computed(() => store.sankey),
    stores: computed(() => store.stores),
    categories: computed(() => store.categories),
    compare: computed(() => store.compare),
    timemachine: computed(() => store.timemachine),
    loading: computed(() => store.loading),
    error: computed(() => store.error),
    loadAll: () => store.loadAll(),
    retry: () => store.loadAll()
  }
}
