package http

import (
	"clean-arch/internal/core/port"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service   port.AuthService
	jwtSecret string
}

// loginRequest represents the request body for logging in a user
type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func NewAuthHandler(service port.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Login godoc
//
//	@Summary		Login and get an access token
//	@Description	Logs in a registered user and returns an access token if the credentials are valid.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loginRequest	true	"Login request body"
//	@Success		200		{object}	authResponse	"Succesfully logged in"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		401		{object}	errorResponse	"Unauthorized error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var creds loginRequest

	if err := ctx.ShouldBindJSON(&creds); err != nil {
		validationError(ctx, err)
		return
	}

	token, err := h.service.Login(ctx, creds.Email, creds.Password)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newAuthResponse(token)

	handleSuccess(ctx, rsp)
}
