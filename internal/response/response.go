package response

import "github.com/gin-gonic/gin"

// successBody is the envelope for successful responses per PRD section 5.1.
type successBody struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

// errorDetail is the envelope for error responses per PRD section 5.1.
type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type errorBody struct {
	Success bool        `json:"success"`
	Error   errorDetail `json:"error"`
}

// Success sends a standard success JSON response.
func Success(c *gin.Context, status int, data any, message string) {
	c.JSON(status, successBody{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// Error sends a standard error JSON response.
func Error(c *gin.Context, status int, code string, message string, details any) {
	c.JSON(status, errorBody{
		Success: false,
		Error: errorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
