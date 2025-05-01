package handlers

import (
	"time"

	"github.com/cverish/go-jwt-rest-api/internal/config"
	"github.com/cverish/go-jwt-rest-api/internal/database"
	"github.com/cverish/go-jwt-rest-api/internal/models"
	"github.com/cverish/go-jwt-rest-api/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	db  *database.Database
	cfg *config.Config
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(db *database.Database, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:  db,
		cfg: cfg,
	}
}

// Login godoc
// @Summary login registered user
// @Description Logs in user with given email and password.
// @Description
// @Description On success, returns status code 200 and sends back HTTPOnly access and refresh tokens.
// @Tags auth
// @Param loginInfo body models.UserLogin true "Login info"
// @Success 200 {object} StatusOK "Successful Response"
// @Failure 400 {object} StatusError "Status Bad Request"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var login models.UserLogin

	if err := c.ShouldBindJSON(&login); err != nil {
		models.ResponseBadRequest(c, err.Error())
		return
	}

	user, err := h.db.GetUserByCredentials(&login)
	if err != nil {
		models.ResponseInternalServerError(c)
		return
	} else if user == nil {
		models.ResponseUnauthorized(c, "invalid user credentials")
		return
	}

	// generate tokens
	now := time.Now()
	accessToken, errAccessToken := utils.GenerateTokenString(
		[]byte(h.cfg.JWT.AccessTokenSecret),
		h.cfg.JWT.AccessTokenExpiry,
		now,
		user.ID.String(),
		user.Email,
		user.Role.String(),
	)
	refreshToken, errRefreshToken := utils.GenerateTokenString(
		[]byte(h.cfg.JWT.RefreshTokenSecret),
		h.cfg.JWT.RefreshTokenExpiry,
		now,
		user.ID.String(),
		user.Email,
		user.Role.String(),
	)
	if errAccessToken != nil || errRefreshToken != nil {
		models.ResponseInternalServerError(c)
		return
	}

	utils.SetToken(c, h.cfg, h.cfg.JWT.AccessTokenKey, accessToken, h.cfg.JWT.AccessTokenExpiry)
	utils.SetToken(c, h.cfg, h.cfg.JWT.RefreshTokenKey, refreshToken, h.cfg.JWT.RefreshTokenExpiry)
	models.ResponseOK(c, "access_token and refresh_token cookies returned")
}

// ChangePassword godoc
// @Summary change user password
// @Description Given correct login information, allow user to change their password.
// @Tags auth
// @Param changeLoginInfo body models.UserPasswordChange true "New login info"
// @Success 200 {object} StatusOK "Successful Response"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Security UserAccessCookie
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var passwordChange models.UserPasswordChange

	if err := c.ShouldBindJSON(&passwordChange); err != nil {
		models.ResponseBadRequest(c, err.Error())
		return
	}

	credentialsValid, changeValid, err := h.db.UpdateUserPassword(&passwordChange)
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

// RefreshToken godoc
// @Summary refresh access token from refresh token
// @Description Reads the user's refresh_token cookie and sends back a new access_token.
// @Description
// @Description Not allowed if tokens are refreshed on the backend (`JWT_BACKEND_REFRESH` environment variable).
// @Tags auth
// @Success 200 {object} StatusOK "Successful Response"
// @Failure 401 {object} StatusError "Status Unauthorized"
// @Failure 404 {object} StatusError "Status Not Found"
// @Failure 500 {object} StatusError "Status Internal Server Error"
// @Security UserAccessCookie
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	if h.cfg.JWT.BackendRefresh {
		models.ResponseNotFound(c, "404 page not found")
		return
	}

	cookie, err := c.Cookie(h.cfg.JWT.RefreshTokenKey)
	if err != nil {
		models.ResponseUnauthorized(c, "missing cookie")
		return
	}

	claims, err := utils.GetTokenClaims(cookie, []byte(h.cfg.JWT.RefreshTokenSecret))
	if err != nil {
		models.ResponseUnauthorized(c, "incorrect token")
		return
	}

	// check authorization
	userId, exists := claims["user_id"]
	if !exists {
		models.ResponseUnauthorized(c, "missing claims")
		return
	}

	user, _ := h.db.GetUserById(userId.(string))
	if user == nil {
		// user's access was revoked
		models.ResponseUnauthorized(c, "not authorized")
		utils.RevokeTokens(c, h.cfg)
		return
	}

	// generate new access token
	accessToken, err := utils.GenerateTokenString(
		[]byte(h.cfg.JWT.AccessTokenSecret),
		h.cfg.JWT.AccessTokenExpiry,
		time.Now(),
		user.ID.String(),
		user.Email,
		user.Role.String(),
	)
	if err != nil {
		models.ResponseInternalServerError(c)
		return
	}

	utils.SetToken(c, h.cfg, h.cfg.JWT.AccessTokenKey, accessToken, h.cfg.JWT.AccessTokenExpiry)
	models.ResponseOK(c, "access_token refreshed successfully")
}

// Logout godoc
// @Summary logout
// @Description Logs out user by removing access_token and refresh_token cookies
// @Tags auth
// @Success 200 {object} StatusOK "Successful Response"
// @Security UserAccessCookie
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	utils.RevokeTokens(c, h.cfg)
	models.ResponseOK(c, "logged out successfully")
}
