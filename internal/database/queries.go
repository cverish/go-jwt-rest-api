package database

import (
	"github.com/cverish/go-jwt-rest-api/internal/models"
	"github.com/cverish/go-jwt-rest-api/internal/utils"
	"gorm.io/gorm"
)

// get all users
func (db *Database) GetAllUsers() ([]models.User, error) {
	var users []models.User

	if q := db.DB.Find(&users); q.Error != nil {
		return nil, q.Error
	}

	return users, nil
}

// get all invited users
func (db *Database) GetAllInvitedUsers() ([]models.InvitedUser, error) {
	var invitedUsers []models.InvitedUser

	if q := db.DB.Find(&invitedUsers); q.Error != nil {
		return nil, q.Error
	}

	return invitedUsers, nil
}

// find user by given id
// returns (nil, error) if server error, (nil, nil) if not found
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

// find invited user by given ID
// returns (nil, error) if server error, (nil, nil) if not found
func (db *Database) GetInvitedUserById(userId string) (*models.InvitedUser, error) {
	var invitedUser models.InvitedUser

	q := db.DB.Where("id = ?", userId).First(&invitedUser)

	if q.Error == gorm.ErrRecordNotFound {
		return nil, nil
	} else if q.Error != nil {
		return nil, q.Error
	}

	return &invitedUser, nil
}

// find user by given email
// returns (nil, error) if server error, (nil, nil) if not found
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

// find invited user by given email
// returns (nil, error) if server error, (nil, nil) if not found
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

// get a user by their credentials
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

// get an invited user by their credentials
func (db *Database) GetInvitedUserByCredentials(login *models.UserLogin) (*models.InvitedUser, error) {
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

// get a count of users
func (db *Database) GetUserCount() (int, error) {
	q := db.DB.Find(&models.User{})
	return int(q.RowsAffected), q.Error
}

// get a count of invited users
func (db *Database) GetInvitedGetUserCount() (int, error) {
	q := db.DB.Find(&models.InvitedUser{})
	return int(q.RowsAffected), q.Error
}

// create a new user and delete the invited user
func (db *Database) RegisterUser(invitedUser *models.InvitedUser, user *models.User) error {
	tx_err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := db.DB.Create(user).Error; err != nil {
			return err
		}
		if err := db.DB.Delete(invitedUser).Error; err != nil {
			return err
		}
		return nil
	})

	return tx_err
}

// create a new user
func (db *Database) CreateUser(user *models.User) error {
	if q := db.DB.Create(user); q.Error != nil {
		return q.Error
	}
	return nil
}

// create a new invited user
func (db *Database) CreateInvitedUser(invitedUser *models.InvitedUser) error {
	if q := db.DB.Create(invitedUser); q.Error != nil {
		return q.Error
	}
	return nil
}

// update user
func (db *Database) UpdateUser(user *models.User) error {
	q := db.DB.Save(&user)
	return q.Error
}

// update invited user
func (db *Database) UpdateInvitedUser(invitedUser *models.InvitedUser) error {
	q := db.DB.Save(&invitedUser)
	return q.Error
}

// update user's password
// returns bool (valid credentials), bool (successful change), err (server error)
func (db *Database) UpdateUserPassword(passwordChange *models.UserPasswordChange) (bool, bool, error) {
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

// reset invited user's password
// returns bool (valid invited user), bool (successful change), err
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

// reset user's password
// returns bool (valid user), bool (successful change), err
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

// delete user
func (db *Database) DeleteUser(user *models.User) error {
	q := db.DB.Delete(&user)
	return q.Error
}

// delete invited user
func (db *Database) DeleteInvitedUser(invitedUser *models.InvitedUser) error {
	q := db.DB.Delete(&invitedUser)
	return q.Error
}

// check that the given email address has not been used by a User or InvitedUser
// returns boolean for validation and error (in the case of other database error)
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

// check that a given username has not been used by a User or InvitedUser
// returns boolean for validation and error (in the case of other database error)
func (db *Database) ValidateUniqueUsername(username string) (bool, error) {
	q := db.DB.Where("username = ?", username).First(&models.User{})

	if q.Error != gorm.ErrRecordNotFound {
		if q.Error != nil {
			return false, q.Error
		} else {
			return false, nil
		}
	}

	return true, nil
}
