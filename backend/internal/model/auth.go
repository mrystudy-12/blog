package model

import (
	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义 JWT 声明结构体
type Claims struct {
	ID                   int64  `json:"ID,omitempty"`
	Username             string `json:"Username,omitempty"`
	Role                 string `json:"Role,omitempty"`
	jwt.RegisteredClaims `json:"Jwt.RegisteredClaims"`
}

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
	Message string      `json:"message"`
	Token   string      `json:"token"`
	Data    interface{} `json:"data"`
}

// RegisterResponse 注册响应结构体
type RegisterResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
