package errcode

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

var (
	Success            = &AppError{Code: 0, Message: "ok"}
	ErrBadRequest      = &AppError{Code: 40000, Message: "bad request"}
	ErrUnauthorized    = &AppError{Code: 40001, Message: "unauthorized"}
	ErrForbidden       = &AppError{Code: 40003, Message: "forbidden"}
	ErrNotFound        = &AppError{Code: 40004, Message: "not found"}
	ErrMethodNotAllowed = &AppError{Code: 40005, Message: "method not allowed"}
	ErrConflict        = &AppError{Code: 40009, Message: "conflict"}
	ErrTooManyRequests = &AppError{Code: 42900, Message: "too many requests"}
	ErrInternal        = &AppError{Code: 50000, Message: "internal error"}

	// Auth
	ErrInvalidCredentials = &AppError{Code: 40101, Message: "invalid email or password"}
	ErrUserExists         = &AppError{Code: 40102, Message: "user already exists"}
	ErrTokenExpired       = &AppError{Code: 40103, Message: "token expired"}
	ErrTokenInvalid       = &AppError{Code: 40104, Message: "token invalid"}

	// Business
	ErrForumNotFound  = &AppError{Code: 40401, Message: "forum not found"}
	ErrThreadNotFound = &AppError{Code: 40402, Message: "thread not found"}
	ErrPostNotFound   = &AppError{Code: 40403, Message: "post not found"}
	ErrUserNotFound   = &AppError{Code: 40404, Message: "user not found"}
	ErrNoPermission   = &AppError{Code: 40301, Message: "insufficient permission"}
	ErrThreadClosed   = &AppError{Code: 40302, Message: "thread is closed"}
	ErrFileTooLarge   = &AppError{Code: 40013, Message: "file too large"}
	ErrFileTypeNotAllowed = &AppError{Code: 40014, Message: "file type not allowed"}
	ErrIPBanned        = &AppError{Code: 40310, Message: "ip is banned"}
)

// Response helpers

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}

func OKMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  msg,
		Data: nil,
	})
}

func Fail(c *gin.Context, httpStatus int, err *AppError) {
	c.JSON(httpStatus, Response{
		Code: err.Code,
		Msg:  err.Message,
		Data: nil,
	})
}

func FailMsg(c *gin.Context, httpStatus int, err *AppError, msg string) {
	c.JSON(httpStatus, Response{
		Code: err.Code,
		Msg:  msg,
		Data: nil,
	})
}

func FailValidation(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{
		Code: ErrBadRequest.Code,
		Msg:  msg,
		Data: nil,
	})
}

func OKData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}
