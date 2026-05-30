export function useAddExpenseSheet() {
  const open = useState('add-expense-open', () => false)
  const addedVersion = useState('add-expense-added-version', () => 0)

  function show() {
    open.value = true
  }

  function notifyAdded() {
    addedVersion.value += 1
  }

  return {
    open,
    addedVersion,
    show,
    notifyAdded
  }
}
