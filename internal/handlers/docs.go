package handlers

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

// StatusOK is the type associated with a 200 server response
// with a message body used by Swagger docs.
type StatusOK struct {
	Message string `json:"message" example:"ok"`
} //@name http.StatusOK

// StatusOKItem is the type associated with a 200 server response
// with an item body used by Swagger docs.
type StatusOKItem struct {
	Item any `json:"item" swaggertype:"object"`
} //@name http.StatusOKItem

// StatusOKList is the type associated with a 200 server response
// with an items body used by Swagger docs.
type StatusOKList struct {
	Items []any `json:"items" swaggertype:"array,object"`
} //@name http.StatusOKList

// StatusCreated is the type associated with a 201 server response
// with a message and id of object created, used by Swagger docs
type StatusCreated struct {
	Message string `json:"message" example:"created successfully"`
	ID      string `json:"id"      example:"id-of-created-object"`
} //@name http.StatusCreated

// StatusError is the type associated with any 4** or 5** responses
// with an error string, used by Swagger docs
type StatusError struct {
	Error string `json:"error" example:"error message"`
} //@name http.StatusError
