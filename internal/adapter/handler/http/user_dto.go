package http

import "clean-arch/internal/core/domain"

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse domain.User

type userListResponse struct {
	Meta *meta           `json:"meta"`
	Data []*userResponse `json:"users"`
}
