package http

import (
	"clean-arch/internal/core/domain"
	"clean-arch/internal/core/port"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc port.UserService
}

func NewUserHandler(svc port.UserService) *UserHandler {
	return &UserHandler{svc}
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	create a new user account
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			registerRequest	body		registerRequest	true	"Register request"
//	@Success		200				{object}	userResponse	"User created"
//	@Failure		400				{object}	errorResponse	"Validation error"
//	@Failure		401				{object}	errorResponse	"Unauthorized error"
//	@Failure		404				{object}	errorResponse	"Data not found error"
//	@Failure		409				{object}	errorResponse	"Data conflict error"
//	@Failure		500				{object}	errorResponse	"Internal server error"
func (h *UserHandler) Register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.svc.Create(ctx, &user); err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, user)
}

// List godoc
//
//	@Summary		List users
//	@Description	List users with pagination
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			page	query		uint64				true	"Page"
//	@Param			size	query		uint64				true	"Size"
//	@Success		200		{object}	userListResponse	"Users displayed"
//	@Failure		400		{object}	errorResponse		"Validation error"
//	@Failure		500		{object}	errorResponse		"Internal server error"
//	@Router			/users [get]
//	@Security		BearerAuth
func (h *UserHandler) List(ctx *gin.Context) {
	var req listResourceRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	users, err := h.svc.List(ctx, req.Page, req.Size)
	if err != nil {
		handleError(ctx, err)
		return
	}

	total := int64(len(users))
	meta := newMeta(total, req.Page, req.Size)
	rsp := toMap(meta, users)

	handleSuccess(ctx, rsp)
}

func (h *UserHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := h.svc.GetByID(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, user)
}

func (h *UserHandler) Update(ctx *gin.Context) {
	user := &domain.User{}

	err := h.svc.Update(ctx, user)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, user)
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.svc.Delete(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}
