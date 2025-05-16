package http

import (
	"errors"
	"log/slog"
	"net/http"

	"clean-arch/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// response represents a response body format
type response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

// errorResponse represents an error response body format
type errorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

// authResponse represents an authentication response body
type authResponse struct {
	AccessToken string `json:"token" example:"v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..."`
}

// meta represents metadata for a paginated response
type meta struct {
	Page      int64 `json:"page" example:"10"`
	Size      int64 `json:"size" example:"0"`
	Total     int64 `json:"total" example:"100"`
	TotalPage int64 `json:"totalPage" example:"10"`
}

// handleSuccess sends a success response with the specified status code and optional data
func handleSuccess(ctx *gin.Context, data any) {
	rsp := newResponse(true, "Success", data)
	ctx.JSON(http.StatusOK, rsp)
}

// newResponse is a helper function to create a response body
func newResponse(success bool, message string, data any) response {
	return response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

// handleError determines the status code of an error and returns a JSON response with the error message and status code
func handleError(ctx *gin.Context, err error, code ...int) {
	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)

	// Default to 500 Internal Server Error if no code is provided
	statusCode := http.StatusInternalServerError
	if len(code) > 0 {
		statusCode = code[0]
	}

	if errors.Is(err, &domain.Error{}) {
		domainErr := err.(*domain.Error)
		statusCode = domainErr.Code
	}

	slog.Error("Error occurred", slog.String("error", err.Error()), slog.Int("status_code", statusCode))

	ctx.JSON(statusCode, errRsp)
}

// parseError parses error messages from the error object and returns a slice of error messages
func parseError(err error) []string {
	var errMsgs []string

	if errors.As(err, &validator.ValidationErrors{}) {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, err.Error())
		}
	} else {
		errMsgs = append(errMsgs, err.Error())
	}

	return errMsgs
}

// newErrorResponse is a helper function to create an error response body
func newErrorResponse(errMsgs []string) errorResponse {
	return errorResponse{
		Success:  false,
		Messages: errMsgs,
	}
}

// validationError sends an error response for some specific request validation error
func validationError(ctx *gin.Context, err error) {
	errMsgs := parseError(err)
	errRsp := newErrorResponse(errMsgs)
	ctx.JSON(http.StatusBadRequest, errRsp)
}

// newAuthResponse is a helper function to create a response body for handling authentication data
func newAuthResponse(token string) authResponse {
	return authResponse{
		AccessToken: token,
	}
}

// newMeta is a helper function to create metadata for a paginated response
func newMeta(total, page, size int64) meta {
	return meta{
		Total: total,
		TotalPage: total / size,
		Size:  size,
		Page:  page,
	}
}
