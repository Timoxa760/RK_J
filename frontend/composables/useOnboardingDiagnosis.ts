export function useOnboardingDiagnosis() {
  const { diagnosis, loading, error, fetchDiagnosis } = useDiagnosis()

  async function loadDiagnosis() {
    await fetchDiagnosis()
  }

  return {
    diagnosis,
    loading,
    error,
    loadDiagnosis
  }
}
