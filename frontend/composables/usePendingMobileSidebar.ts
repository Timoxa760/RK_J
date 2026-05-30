/** Открыть мобильный сайдбар без вызова useSidebar вне дерева SidebarProvider. */
export function usePendingMobileSidebar() {
  const pending = useState('mm-sidebar-open-mobile-pending', () => false)

  function requestOpen() {
    pending.value = true
  }

  function clear() {
    pending.value = false
  }

  return { pending, requestOpen, clear }
}
