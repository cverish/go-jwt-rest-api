package handlers

import (
	"github.com/cverish/go-jwt-rest-api/internal/database"
	"github.com/cverish/go-jwt-rest-api/internal/models"
	"github.com/cverish/go-jwt-rest-api/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	db *database.Database
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *database.Database) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

// Register godoc
// @Summary register new user
// @Tags registration
// @Param registrationInfo body models.UserRegister true "Registration Information"
// @Success 201 {object} StatusCreated "Successful Response"
// @Failure 400 {object} StatusError "Status Bad Request"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var registrationInfo models.UserRegister

	// validate input JSON
	if err := c.ShouldBindJSON(&registrationInfo); err != nil {
		models.ResponseBadRequest(c, err.Error())
		return
	}

	credentials := &models.UserLogin{
		Email:    registrationInfo.Email,
		Password: registrationInfo.Password,
	}

	// find invited user that matches user register information
	invitedUser, err := h.db.GetInvitedUserByCredentials(credentials)
	if err != nil {
		models.ResponseInternalServerError(c)
		return
	} else if invitedUser == nil {
		models.ResponseUnauthorized(c, "invalid credentials")
		return
	}

	// check if username is unique
	if unique, err := h.db.ValidateUniqueUsername(registrationInfo.Username); err != nil {
		models.ResponseInternalServerError(c)
		return
	} else if !unique {
		models.ResponseBadRequest(c, "username has already been taken")
		return
	}

	// verify that new passwords match
	if registrationInfo.NewPassword != registrationInfo.NewPasswordConfirm {
		models.ResponseBadRequest(c, "passwords must match")
		return
	}

	// check that password meets requirements
	if err := utils.ValidatePassword(registrationInfo.NewPassword); err != nil {
		models.ResponseBadRequest(c, err.Error())
		return
	}

	// hash password
	hashedPassword, err := utils.HashPassword(registrationInfo.NewPassword)
	if err != nil {
		models.ResponseInternalServerError(c)
		return
	}

	createdUser := &models.User{
		Username:     registrationInfo.Username,
		Email:        registrationInfo.Email,
		PasswordHash: hashedPassword,
		FirstName:    registrationInfo.FirstName,
		LastName:     registrationInfo.LastName,
	}

	// open transaction to create new user and delete existing invited user
	if err := h.db.RegisterUser(invitedUser, createdUser); err != nil {
		models.ResponseInternalServerError(c)
		return
	}

	models.ResponseCreated(c, "user registered successfully", createdUser.ID.String())
}

// GetUser godoc
// @Summary get user information
// @Tags users
// @Param user_id path string true "User ID (uuid)"
// @Success 200 {object} StatusOK "Successful Response"
// @Failure 400 {object} StatusError "Status Bad Request"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 403 {object} StatusError "Status Forbidden"
// @Failure 404 {object} StatusError "Status Not Found"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Security UserAccessCookie
// @Router /users/{user_id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userId := c.Param("user_id")
	// only allow a user to update themself, or an admin to view any user, by checking JWT
	if c.GetString("user_id") != userId && c.GetString("role") != models.RoleAdmin.String() {
		models.ResponseForbidden(c, "not authorized to view other user profiles")
		return
	}

	user, err := h.db.GetUserById(userId)
	if err != nil {
		models.ResponseInternalServerError(c)
		return
	} else if user == nil {
		models.ResponseNotFound(c, "user not found")
		return
	}

	models.ResponseOKItem(c, user)
}

// UpdateUser godoc
// @Summary update user information
// @Tags users
// @Param user_id path string true "User ID (uuid)"
// @Param user body models.User true "Updated user info"
// @Success 200 {object} StatusOK "Successful Response"
// @Failure 400 {object} StatusError "Status Bad Request"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 403 {object} StatusError "Status Forbidden"
// @Failure 404 {object} StatusError "Status Not Found"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Security UserAccessCookie
// @Router /users/{user_id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userId := c.Param("user_id")

	// only allow a user to update themself, or an admin to update any user, by checking JWT
	if c.GetString("user_id") != userId && c.GetString("role") != models.RoleAdmin.String() {
		models.ResponseForbidden(c, "Not authorized to change other user profiles")
		return
	}

	user, err := h.db.GetUserById(userId)
	if err != nil {
		models.ResponseInternalServerError(c)
		return
	} else if user == nil {
		models.ResponseNotFound(c, "user not found")
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		models.ResponseBadRequest(c, err.Error())
		return
	}

	// validate that username is still unique
	if user.Username != updatedUser.Username {
		if unique, err := h.db.ValidateUniqueUsername(updatedUser.Username); err != nil {
			models.ResponseInternalServerError(c)
			return
		} else if !unique {
			models.ResponseBadRequest(c, "username already in use")
			return
		}
	}

	// validate that email is still unique
	if user.Email != updatedUser.Email {
		if unique, err := h.db.ValidateUniqueEmail(updatedUser.Email); err != nil {
			models.ResponseInternalServerError(c)
			return
		} else if !unique {
			models.ResponseBadRequest(c, "email already in use")
			return
		}
	}

	user.Username = updatedUser.Username
	user.FirstName = updatedUser.FirstName
	user.LastName = updatedUser.LastName
	user.Email = updatedUser.Email

	if err := h.db.UpdateUser(user); err != nil {
		models.ResponseInternalServerError(c)
		return
	}

	models.ResponseOK(c, "updated user successfully")
}
