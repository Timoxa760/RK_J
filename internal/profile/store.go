package profile

// Store — чтение/запись финансового профиля.
type Store interface {
	Get(userID string) (FinancialProfile, error)
	Save(userID string, p FinancialProfile) error
}
