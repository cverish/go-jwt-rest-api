package handlers

import (
	"github.com/cverish/go-jwt-rest-api/internal/database"
	"github.com/cverish/go-jwt-rest-api/internal/models"
	"github.com/cverish/go-jwt-rest-api/internal/utils"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	db *database.Database
}

func NewAdminHandler(db *database.Database) *AdminHandler {
	return &AdminHandler{
		db: db,
	}
}

// GetInvitedUsers godoc
// @Summary get all invited users
// @Description Get a list of all invited, but not registered, users
// @Tags admin
// @Success 200 {object} StatusOKList "Successful Response"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 403 {object} StatusError "Status Forbidden"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Security AdminAccessInCookie
// @Router /admin/invites [get]
func (h *AdminHandler) GetInvitedUsers(c *gin.Context) {
	invitedUsers, err := h.db.GetAllInvitedUsers()
	if err != nil {
		models.ResponseInternalServerError(c)
		return
	}

	models.ResponseOKList(c, invitedUsers)
}

// CreateInvite godoc
// @Summary create an invitation for a new user
// @Description Creates a InvitedUser with an email address and temporary password. Invited user must use this email and password to register for an account.
// @Tags admin
// @Param invitationInfo body models.UserLogin true "Invitation email and password"
// @Success 201 {object} StatusCreated "Successful Response"
// @Failure 400 {object} StatusError "Status Bad Request"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 403 {object} StatusError "Status Forbidden"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Security AdminAccessInCookie
// @Router /admin/invites/create [post]
func (h *AdminHandler) CreateInvite(c *gin.Context) {
	var invitedUser models.UserLogin

	// validate input JSON
	if err := c.ShouldBindJSON(&invitedUser); err != nil {
		models.ResponseBadRequest(c, err.Error())
		return
	}

	// validate email
	if err := utils.ValidateEmail(invitedUser.Email); err != nil {
		models.ResponseBadRequest(c, err.Error())
		return
	} else if unique, err := h.db.ValidateUniqueEmail(invitedUser.Email); err != nil {
		models.ResponseInternalServerError(c)
		return
	} else if !unique {
		models.ResponseBadRequest(c, "email already in use")
		return
	}

	// hash password
	hashedPassword, err := utils.HashPassword(invitedUser.Password)
	if err != nil {
		models.ResponseInternalServerError(c)
		return
	}

	// create invitation
	createdInvitedUser := &models.InvitedUser{
		Email:        invitedUser.Email,
		PasswordHash: hashedPassword,
	}

	if err := h.db.CreateInvitedUser(createdInvitedUser); err != nil {
		models.ResponseInternalServerError(c)
		return
	}

	models.ResponseCreated(c, "user registered successfully", createdInvitedUser.ID.String())
}

// ResetInvitedUserPassword godoc
// @Summary reset the password of an invited user
// @Description Reset the password of an invited user in case password was forgotten.
// @Tags admin
// @Param loginInfo body models.UserLogin true "Login info"
// @Success 200 {object} StatusOK "Successful Response"
// @Failure 400 {object} StatusError "Status Bad Request"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 403 {object} StatusError "Status Forbidden"
// @Failure 404 {object} StatusError "Status Not Found"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Security AdminAccessInCookie
// @Router /admin/invites/reset-password [post]
func (h *AdminHandler) ResetInvitedUserPassword(c *gin.Context) {
	var loginInfo models.UserLogin

	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		models.ResponseBadRequest(c, err.Error())
		return
	}

	found, ok, err := h.db.ResetInvitedUserPassword(loginInfo)
	if err != nil {
		models.ResponseInternalServerError(c)
		return
	} else if !found {
		models.ResponseNotFound(c, "user not found")
		return
	} else if !ok {
		models.ResponseNotAcceptable(c, "password is not valid")
		return
	}

	models.ResponseOK(c, "password changed successfully")
}

// DeleteInvitedUser godoc
// @Summary delete invited user
// @Description Delete an invited user
// @Tags admin
// @Param user_id path string true "User ID (uuid)"
// @Success 200 {object} StatusOK "Successful Response"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 403 {object} StatusError "Status Forbidden"
// @Failure 404 {object} StatusError "Status Not Found"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Security AdminAccessInCookie
// @Router /admin/invites/delete/{user_id} [delete]
func (h *AdminHandler) DeleteInvitedUser(c *gin.Context) {
	userId := c.Param("user_id")

	user, err := h.db.GetInvitedUserById(userId)
	if err != nil {
		models.ResponseInternalServerError(c)
		return
	} else if user == nil {
		models.ResponseNotFound(c, "user not found")
		return
	}

	if err := h.db.DeleteInvitedUser(user); err != nil {
		models.ResponseInternalServerError(c)
		return
	}

	models.ResponseOK(c, "deleted invited user successfully")
}

// GetUsers godoc
// @Summary get all users
// @Description Get a list of all registered users
// @Tags admin
// @Success 200 {object} StatusOKList "Successful Response"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 403 {object} StatusError "Status Forbidden"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Security AdminAccessInCookie
// @Router /admin/users [get]
func (h *AdminHandler) GetUsers(c *gin.Context) {
	users, err := h.db.GetAllUsers()
	if err != nil {
		models.ResponseInternalServerError(c)
		return
	}

	models.ResponseOKList(c, users)
}

// ResetUserPassword godoc
// @Summary reset the password of a user
// @Description Reset the password of a user in case password was forgotten.
// @Tags admin
// @Param loginInfo body models.UserLogin true "New login info"
// @Success 200 {object} StatusOK "Successful Response"
// @Failure 400 {object} StatusError "Status Bad Request"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 403 {object} StatusError "Status Forbidden"
// @Failure 404 {object} StatusError "Status Not Found"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Security AdminAccessInCookie
// @Router /admin/users/reset-password [post]
func (h *AdminHandler) ResetUserPassword(c *gin.Context) {
	var loginInfo models.UserLogin

	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		models.ResponseBadRequest(c, err.Error())
		return
	}

	credentialsValid, changeValid, err := h.db.ResetUserPassword(loginInfo)
	if err != nil {
		models.ResponseInternalServerError(c)
		return
	} else if !credentialsValid {
		models.ResponseUnauthorized(c, "invalid user credentials")
		return
	} else if !changeValid {
		models.ResponseBadRequest(c, "new password is not valid")
		return
	}

	models.ResponseOK(c, "password updated successfully")
}

// DeleteUser godoc
// @Summary delete user
// @Description Delete a user
// @Tags admin
// @Param user_id path string true "User ID (uuid)"
// @Success 200 {object} StatusOK "Successful Response"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 403 {object} StatusError "Status Forbidden"
// @Failure 404 {object} StatusError "Status Not Found"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Security AdminAccessInCookie
// @Router /admin/users/delete/{user_id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userId := c.Param("user_id")

	user, err := h.db.GetUserById(userId)
	if err != nil {
		models.ResponseInternalServerError(c)
		return
	} else if user == nil {
		models.ResponseNotFound(c, "user not found")
		return
	}

	if err := h.db.DeleteUser(user); err != nil {
		models.ResponseInternalServerError(c)
		return
	}

	models.ResponseOK(c, "deleted user successfully")
}

// CreateInitialAdmin godoc
// @Summary create first user as admin
// @Description When no users exist in the database, create an admin user with the given information.
// @Tags admin
// @Param registrationInfo body models.AdminRegister true "Admin registration information"
// @Success 201 {object} StatusCreated "Successful Response"
// @Failure 400 {object} StatusError "Status Bad Request"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Security AdminAccessInCookie
// @Router /admin/manage/create-initial-admin [post]
func (h *AdminHandler) CreateInitialAdmin(c *gin.Context) {
	// validate that no users currently exist
	if numUsers, err := h.db.GetUserCount(); err != nil {
		models.ResponseInternalServerError(c)
		return
	} else if numUsers != 0 {
		models.ResponseUnauthorized(c, "an admin user cannot be created")
		return
	}

	var user models.AdminRegister
	// validate input JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		models.ResponseBadRequest(c, err.Error())
		return
	}

	// check that password meets requirements
	if err := utils.ValidatePassword(user.Password); err != nil {
		models.ResponseBadRequest(c, err.Error())
		return
	}

	// hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		models.ResponseInternalServerError(c)
		return
	}

	createdAdmin := &models.User{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: hashedPassword,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         models.RoleAdmin,
	}

	if err := h.db.CreateUser(createdAdmin); err != nil {
		models.ResponseInternalServerError(c)
		return
	}

	models.ResponseCreated(c, "admin user registered successfully", createdAdmin.ID.String())
}

// AddAdmin godoc
// @Summary upgrade a user to admin
// @Description Elevate the given user's role to `admin`.
// @Tags admin
// @Param user_id path string true "ID of user to elevate to admin"
// @Success 201 {object} StatusCreated "Successful Response"
// @Failure 400 {object} StatusError "Status Bad Request"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 403 {object} StatusError "Status Forbidden"
// @Failure 404 {object} StatusError "Status Not Found"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Security AdminAccessInCookie
// @Router /admin/manage/add-admin/{user_id} [post]
func (h *AdminHandler) AddAdmin(c *gin.Context) {
	// get user to elevate to admin
	userId := c.Param("user_id")
	InvitedUserAdmin, err := h.db.GetUserById(userId)
	if err != nil {
		models.ResponseNotFound(c, err.Error())
		return
	}

	InvitedUserAdmin.Role = models.RoleAdmin
	if err := h.db.UpdateUser(InvitedUserAdmin); err != nil {
		models.ResponseInternalServerError(c)
		return
	}

	models.ResponseOK(c, "user updated to admin successfully")
}

// RemoveAdmin godoc
// @Summary downgrade an admin to user
// @Description Remove admin permissions from user. POST body requires current admin user's email and password for verification.
// @Tags admin
// @Param user_id path string true "ID of admin to downgrade to user"
// @Success 201 {object} StatusCreated "Successful Response"
// @Failure 400 {object} StatusError "Status Bad Request"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 403 {object} StatusError "Status Forbidden"
// @Failure 404 {object} StatusError "Status Not Found"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Security AdminAccessInCookie
// @Router /admin/manage/remove-admin/{user_id} [post]
func (h *AdminHandler) RemoveAdmin(c *gin.Context) {
	// get admin to downgrade to user
	userId := c.Param("user_id")
	adminToDowngrade, err := h.db.GetUserById(userId)
	if err != nil {
		models.ResponseNotFound(c, err.Error())
		return
	}

	adminToDowngrade.Role = models.RoleUser
	if err := h.db.UpdateUser(adminToDowngrade); err != nil {
		models.ResponseInternalServerError(c)
		return
	}

	models.ResponseOK(c, "admin updated to user successfully")
}
