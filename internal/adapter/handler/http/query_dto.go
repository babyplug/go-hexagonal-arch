package http

// listResourceRequest represents the request body for listing resources
type listResourceRequest struct {
	Page int64 `form:"page" binding:"min=0" example:"0"`
	Size int64 `form:"size" binding:"min=10,max=1000" example:"20"`
}
