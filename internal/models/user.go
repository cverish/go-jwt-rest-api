package models

import "database/sql/driver"

// User contains the information associated with a user.
type User struct {
	Base
	// User's username
	Username string `json:"username"   gorm:"uniqueIndex"                binding:"required"       example:"username"`
	// User's email address
	Email string `json:"email"      gorm:"uniqueIndex"                binding:"required,email" example:"username@example.com"`
	// User's role -- oneof "user", "admin"
	Role         Role   `json:"role"       gorm:"type:string;default:'user'" binding:"-"              example:"admin"`
	PasswordHash string `json:"-"          gorm:"type:varchar(255);unique"   binding:"-"                                             swaggerignore:"true"`
	// User's first name
	FirstName string `json:"first_name"                                   binding:"required"       example:"FirstName"`
	// User's last name
	LastName string `json:"last_name"                                    binding:"required"       example:"LastName"`
} //@name user.User

// InvitedUser contains the information associated with an invited user.
type InvitedUser struct {
	Base
	Email        string `json:"email" gorm:"uniqueIndex"                                binding:"required,email" example:"invited_user@example.com"`
	PasswordHash string `json:"-"     gorm:"type:varchar(255);unique binding:omitempty"                                                             swaggerignore:"true"`
} //@name user.InvitedUser

// UserRegister contains the information associated with a user registration request.
type UserRegister struct {
	// Email of invited user, given when invited
	Email string `json:"email"                binding:"required,email" example:"user@example.com"`
	// Temporary password of invited user, given when invited
	Password string `json:"password"             binding:"required"       example:"password-given-to-invited-user"`
	// invited user's chosen username
	Username string `json:"username"             binding:"required"       example:"user"`
	// invited user's chosen new password
	// Must be 8>= len >= 64 characters, contain at least one uppercase, one lowercase, one digit, one special character
	NewPassword string `json:"new_password"         binding:"required"       example:"user-chosen-password"`
	// Confirmation of chosen password
	// Must be 8>= len >= 64 characters, contain at least one uppercase, one lowercase, one digit, one special character
	NewPasswordConfirm string `json:"new_password_confirm" binding:"required"       example:"user-chosen-password"`
	// invited user's chosen first name
	FirstName string `json:"first_name"           binding:"required"       example:"FirstName"`
	// invited user's chosen last name
	LastName string `json:"last_name"            binding:"required"       example:"LastName"`
} //@name user.UserRegister

// AdminRegister contains the information associated with an admin registration request.
type AdminRegister struct {
	// Chosen username
	Username string `json:"username"   binding:"required"       example:"admin"`
	// Chosen password.
	// Must be 8>= len >= 64 characters, contain at least one uppercase, one lowercase, one digit, one special character
	Password string `json:"password"   binding:"required"       example:"adminPassword123!"`
	// Chosen email address
	Email string `json:"email"      binding:"required,email" example:"admin@example.com"`
	// Chosen first name
	FirstName string `json:"first_name" binding:"required"       example:"FirstName"`
	// Chosen last name
	LastName string `json:"last_name"  binding:"required"       example:"LastName"`
} //@name user.AdminRegister

// Role is a user role enum. Available values: RoleAdmin ("admin") and RoleUser ("user")
type Role string //@name user.RoleEnum

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (r *Role) Scan(value interface{}) error {
	*r = Role(value.(string))
	return nil
}

// String returns the string representation of a given Role.
func (r Role) String() string {
	return string(r)
}

func (r Role) Value() (driver.Value, error) {
	return r.String(), nil
}
