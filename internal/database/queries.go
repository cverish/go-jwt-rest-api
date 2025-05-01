package database

import (
	"github.com/cverish/go-jwt-rest-api/internal/models"
	"github.com/cverish/go-jwt-rest-api/internal/utils"
	"gorm.io/gorm"
)

// GetAllUsers queries the database and returns a list of Users.
//
// Returns:
//
//	([]User, nil): success
//	(nil, error): error
func (db *Database) GetAllUsers() ([]models.User, error) {
	var users []models.User

	if q := db.DB.Find(&users); q.Error != nil {
		return nil, q.Error
	}

	return users, nil
}

// GetAllInvitedUsers queries the database and returns a list of InvitedUsers.
//
// Returns:
//
//	([]InvitedUser, nil): success
//	(nil, error): error
func (db *Database) GetAllInvitedUsers() ([]models.InvitedUser, error) {
	var invitedUsers []models.InvitedUser

	if q := db.DB.Find(&invitedUsers); q.Error != nil {
		return nil, q.Error
	}

	return invitedUsers, nil
}

// GetUserById queries the database for a User with the given userId.
//
// Returns:
//
//	(*User, nil): found
//	(nil, nil): not found
//	(nil, error): server error
func (db *Database) GetUserById(userId string) (*models.User, error) {
	var user models.User

	q := db.DB.Where("id = ?", userId).First(&user)
	if q.Error == gorm.ErrRecordNotFound {
		return nil, nil
	} else if q.Error != nil {
		return nil, q.Error
	}

	return &user, nil
}

// GetInvitedUserById queries the database for an InvitedUser with the given invitedUserId.
//
// Returns:
//
//	(*InvitedUser, nil): found
//	(nil, nil): not found
//	(nil, error): server error
func (db *Database) GetInvitedUserById(invitedUserId string) (*models.InvitedUser, error) {
	var invitedUser models.InvitedUser

	q := db.DB.Where("id = ?", invitedUserId).First(&invitedUser)
	if q.Error == gorm.ErrRecordNotFound {
		return nil, nil
	} else if q.Error != nil {
		return nil, q.Error
	}

	return &invitedUser, nil
}

// GetUserByEmail queries the database for a User with the given email.
//
// Returns:
//
//	(*User, nil): found
//	(nil, nil): not found
//	(nil, error): server error
func (db *Database) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	q := db.DB.Where("email = ?", email).First(&user)
	if q.Error == gorm.ErrRecordNotFound {
		return nil, nil
	} else if q.Error != nil {
		return nil, q.Error
	}

	return &user, nil
}

// GetInvitedUserByEmail queries the database for an InvitedUser with the given email.
//
// Returns:
//
//	(*InvitedUser, nil): found
//	(nil, nil): not found
//	(nil, error): server error
func (db *Database) GetInvitedUserByEmail(email string) (*models.InvitedUser, error) {
	var invitedUser models.InvitedUser

	q := db.DB.Where("email = ?", email).First(&invitedUser)
	if q.Error == gorm.ErrRecordNotFound {
		return nil, nil
	} else if q.Error != nil {
		return nil, q.Error
	}

	return &invitedUser, nil
}

// GetUserByCredentials queries the database for a User with the email in UserLogin
// and validates the given password in UserLogin against the user's PasswordHash.
//
// Returns:
//
//	(*User, nil): found and valid credentials
//	(nil, error): server error,
//	(nil, nil): not found or user credentials are incorrect.
func (db *Database) GetUserByCredentials(login *models.UserLogin) (*models.User, error) {
	user, err := db.GetUserByEmail(login.Email)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, nil
	}

	if ok := utils.CheckPasswordHash(login.Password, user.PasswordHash); !ok {
		return nil, nil
	}

	return user, nil
}

// GetInvitedUserByCredentials queries the database for an InvitedUser with the email in UserLogin
// and validates the given password in UserLogin against the invited user's PasswordHash.
//
// Returns:
//
//	(*InvitedUser, nil): found and valid credentials
//	(nil, error): server error,
//	(nil, nil): not found or user credentials are incorrect.
func (db *Database) GetInvitedUserByCredentials(
	login *models.UserLogin,
) (*models.InvitedUser, error) {
	invitedUser, err := db.GetInvitedUserByEmail(login.Email)
	if err != nil {
		return nil, err
	} else if invitedUser == nil {
		return nil, nil
	}

	if ok := utils.CheckPasswordHash(login.Password, invitedUser.PasswordHash); !ok {
		return nil, nil
	}

	return invitedUser, nil
}

// GetUserCount queries the database and returns the number of Users.
//
// Returns:
//
//	(int, nil): found
//	(0, error): server error.
func (db *Database) GetUserCount() (int, error) {
	q := db.DB.Find(&models.User{})
	return int(q.RowsAffected), q.Error
}

// GetInvitedUserCount queries the database and returns the number of InvitedUsers.
//
// Returns:
//
//	(int, nil): count found
//	(0, error): server error
func (db *Database) GetInvitedGetUserCount() (int, error) {
	q := db.DB.Find(&models.InvitedUser{})
	return int(q.RowsAffected), q.Error
}

// RegisterUser creates a User and deletes an InvitedUser within a single transaction.
// If either fails, the transaction is rolled back.
//
// Returns:
//
//	nil: success
//	error: server error
func (db *Database) RegisterUser(invitedUser *models.InvitedUser, user *models.User) error {
	txErr := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := db.DB.Create(user).Error; err != nil {
			return err
		}
		if err := db.DB.Delete(invitedUser).Error; err != nil {
			return err
		}
		return nil
	})

	return txErr
}

// CreateUser creates a User.
//
// Returns:
//
//	nil: success
//	error: server error
func (db *Database) CreateUser(user *models.User) error {
	q := db.DB.Create(user)
	return q.Error
}

// CreateInvitedUser creates an InvitedUser.
// Returns nil if successful, error if server error.
func (db *Database) CreateInvitedUser(invitedUser *models.InvitedUser) error {
	q := db.DB.Create(invitedUser)
	return q.Error
}

// UpdateUser saves a given User.
// Returns nil if successful, error if server error.
func (db *Database) UpdateUser(user *models.User) error {
	q := db.DB.Save(&user)
	return q.Error
}

// UpdateInvitedUser saves a given InvitedUser.
// Returns nil if successful, error if server error.
func (db *Database) UpdateInvitedUser(invitedUser *models.InvitedUser) error {
	q := db.DB.Save(&invitedUser)
	return q.Error
}

// UpdateUserPassword gets a User by their credentials, validates their new password,
// hashes the password, and updates the user's PasswordHash.
//
// Returns:
//
//	(true, true, nil): valid credentials and changed successfully
//	(false, false, nil): invalid credentials
//	(true, false, error): valid credentials but server error
func (db *Database) UpdateUserPassword(
	passwordChange *models.UserPasswordChange,
) (bool, bool, error) {
	login := &models.UserLogin{Email: passwordChange.Email, Password: passwordChange.Password}
	user, err := db.GetUserByCredentials(login)
	if err != nil {
		return false, false, err
	} else if user == nil {
		return false, false, nil
	}

	// verify that new passwords match
	if passwordChange.NewPassword != passwordChange.NewPasswordConfirm {
		return true, false, nil
	}
	// check that password meets requirements
	if err := utils.ValidatePassword(passwordChange.NewPassword); err != nil {
		return true, false, nil
	}

	newPasswordHash, err := utils.HashPassword(passwordChange.NewPassword)
	if err != nil {
		return true, false, err
	}

	user.PasswordHash = newPasswordHash
	if err := db.UpdateUser(user); err != nil {
		return true, false, err
	}

	return true, true, nil
}

// ResetInvitedUserPassword gets an InvitedUser by the given email,
// checks the validity of the password, hashes the password,
// and updates the invited user's PasswordHash.
//
// Returns:
//
//	(true, true, nil): valid credentials and changed successfully
//	(false, false, nil): invited user not found
//	(true, false, error): valid credentials but server error
func (db *Database) ResetInvitedUserPassword(newCredentials models.UserLogin) (bool, bool, error) {
	invitedUser, err := db.GetInvitedUserByEmail(newCredentials.Email)

	if err != nil {
		return false, false, err
	} else if invitedUser == nil {
		return false, false, nil
	}

	// check that password meets requirements
	if err := utils.ValidatePassword(newCredentials.Password); err != nil {
		return true, false, nil
	}

	newPasswordHash, err := utils.HashPassword(newCredentials.Password)
	if err != nil {
		return true, false, err
	}

	invitedUser.PasswordHash = newPasswordHash
	if err := db.UpdateInvitedUser(invitedUser); err != nil {
		return true, false, err
	}

	return true, true, nil
}

// ResetUserPassword gets an User by the given email,
// checks the validity of the password, hashes the password,
// and updates the user's PasswordHash.
//
// Returns:
//
//	(true, true, nil): valid credentials and changed successfully
//	(false, false, nil): user not found
//	(true, false, error): valid credentials but server error
func (db *Database) ResetUserPassword(newCredentials models.UserLogin) (bool, bool, error) {
	user, err := db.GetUserByEmail(newCredentials.Email)
	if err != nil {
		return false, false, err
	} else if user == nil {
		return false, false, nil
	}

	// check that password meets requirements
	if err := utils.ValidatePassword(newCredentials.Password); err != nil {
		return true, false, nil
	}

	newPasswordHash, err := utils.HashPassword(newCredentials.Password)
	if err != nil {
		return true, false, err
	}

	user.PasswordHash = newPasswordHash
	if err := db.UpdateUser(user); err != nil {
		return true, false, err
	}

	return true, true, nil
}

// DeleteUser deletes a given User from the database.
//
// Returns:
//
//	nil: success
//	error: server error
func (db *Database) DeleteUser(user *models.User) error {
	q := db.DB.Delete(&user)
	return q.Error
}

// DeleteInvitedUser deletes a given InvitedUser from the database.
//
// Returns:
//
//	nil: success
//	error: server error
func (db *Database) DeleteInvitedUser(invitedUser *models.InvitedUser) error {
	q := db.DB.Delete(&invitedUser)
	return q.Error
}

// ValidateUniqueEmail checks that the given email has not been used by an
// existing User or InvitedUser.
//
// Returns:
//
//	(true, nil): email is unique
//	(false, nil): email is not unique
//	(false, error): server error
func (db *Database) ValidateUniqueEmail(email string) (bool, error) {
	q1 := db.DB.Where("email = ?", email).First(&models.User{})
	if q1.Error != gorm.ErrRecordNotFound {
		if q1.Error != nil {
			return false, q1.Error
		} else {
			return false, nil
		}
	}

	q2 := db.DB.Where("email = ?", email).First(&models.InvitedUser{})
	if q2.Error != gorm.ErrRecordNotFound {
		if q2.Error != nil {
			return false, q2.Error
		} else {
			return false, nil
		}
	}

	return true, nil
}

// ValidateUniqueUsername checks that the given username has not been used by an
// existing User.
//
// Returns:
//
//	(true, nil): username is unique
//	(false, nil): username is not unique
//	(false, error): server error
func (db *Database) ValidateUniqueUsername(username string) (bool, error) {
	q := db.DB.Where("username = ?", username).First(&models.User{})
	if q.Error == gorm.ErrRecordNotFound {
		return true, nil
	} else {
		return false, q.Error
	}
}
