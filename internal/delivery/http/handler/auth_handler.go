package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prasetyoary48/wms/internal/usecase"
	"github.com/prasetyoary48/wms/pkg/response"
)

type AuthHandler struct {
	authUC usecase.AuthUsecase
}

func NewAuthHandler(authUC usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid")
		return
	}

	token, user, err := h.authUC.Login(req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "login berhasil", gin.H{
		"token": token,
		"user":  user,
	})
}

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

// Register hanya bisa diakses Admin (lihat routing) untuk membuat akun user baru.
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid")
		return
	}

	user, err := h.authUC.Register(req.Name, req.Email, req.Password, req.RoleID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "user berhasil dibuat", user)
}