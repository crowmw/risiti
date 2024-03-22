package model

type User struct {
	ID    uint64 `json:"id"`
	Email string `json:"email" validate:"regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	// Min one small and up letters, one number, one special, min 8 chars
	Password string `json:"password" validate:"min=8"`
}
