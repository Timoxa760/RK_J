package auth

import (
	"backend_project/internal/otp"
	"backend_project/internal/userstore"
)

const resetStubCode = "0000"

// Deps — зависимости auth handlers.
type Deps struct {
	Users userstore.Store
	OTP   otp.Store
}

func NewDeps(users userstore.Store, otpStore otp.Store) *Deps {
	return &Deps{
		Users: users,
		OTP:   otpStore,
	}
}

func resetKey(phone string) string {
	return "reset:" + phone
}
