package utils

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
    c.JSON(statusCode, Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err string) {
    c.JSON(statusCode, Response{
        Success: false,
        Message: message,
        Error:   err,
    })
}

func InternalServerError(c *gin.Context, err error) {
    ErrorResponse(c, http.StatusInternalServerError, "Internal server error", err.Error())
}

func BadRequestError(c *gin.Context, message string) {
    ErrorResponse(c, http.StatusBadRequest, message, "Bad request")
}

func UnauthorizedError(c *gin.Context) {
    ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid credentials")
}

func NotFoundError(c *gin.Context, resource string) {
    ErrorResponse(c, http.StatusNotFound, resource+" not found", "Resource not found")
}
