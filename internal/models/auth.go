package models

type UserLogin struct {
	Email    string `json:"email"    binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required"       example:"aPassword1!"`
} //@name auth.UserLogin

type UserPasswordChange struct {
	Email              string `json:"email"                binding:"required,email" example:"user@example.com"`
	Password           string `json:"password"             binding:"required"       example:"aPassword1!"`
	NewPassword        string `json:"new_password"         binding:"required"       example:"newPassword2!"`
	NewPasswordConfirm string `json:"new_password_confirm" binding:"required"       example:"newPassword2!"`
} //@name auth.UserPasswordChange
