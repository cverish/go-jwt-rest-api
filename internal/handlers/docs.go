// @BasePath /api/v1
// @Title JWT Cookie Rest
// @Version 1.0
// @Description Rest API for JWT cookie authentication and user registration
// @Accept json
// @Produce json
// @Schemes http
// @Host localhost:8080

// @tag.name health check

// @tag.name auth
// @tag.description login, logout, and refresh token

// @tag.name registration
// @tag.description register as an invited user

// @tag.name users
// @tag.description view and update user information

// @tag.name admin
// @tag.description admin user capabilities

// @SecurityDefinitions.apikey BasicAuth
// @in cookie
// @name UserAccessCookie
// @description Requires login. Don't actually fill this out! Use the /login endpoint.

// @SecurityDefinitions.apikey ApiKeyAuth
// @in cookie
// @name AdminAccessInCookie
// @description Requires logged in user to be an admin. Don't actually fill this out! Use the /login endpoint with an admin user instead.
package handlers

type StatusOK struct {
	Message string `json:"message" example:"ok"`
} //@name http.StatusOK

type StatusOKList struct {
	Items []any `json:"items" swaggertype:"array,object"`
} //@name http.StatusOKList

type StatusCreated struct {
	Message string `json:"message" example:"created successfully"`
	ID      string `json:"id" example:"id-of-created-object"`
} //@name http.StatusCreated

type StatusError struct {
	Error string `json:"error" example:"error message"`
} //@name http.StatusError
