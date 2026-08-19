package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prasetyoary48/wms/internal/domain"
	"github.com/prasetyoary48/wms/internal/repository"
	"github.com/prasetyoary48/wms/pkg/response"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

type productRequest struct {
	SKU         string `json:"sku" binding:"required"`
	Name        string `json:"name" binding:"required"`
	CategoryID  uint   `json:"category_id"`
	Unit        string `json:"unit" binding:"required"`
	Description string `json:"description"`
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req productRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid")
		return
	}
	product := &domain.Product{
		SKU: req.SKU, Name: req.Name, CategoryID: req.CategoryID,
		Unit: req.Unit, Description: req.Description, IsActive: true,
	}
	if err := h.repo.Create(product); err != nil {
		response.Error(c, http.StatusBadRequest, "gagal membuat produk, SKU mungkin sudah dipakai")
		return
	}
	response.Success(c, http.StatusCreated, "produk berhasil dibuat", product)
}

func (h *ProductHandler) List(c *gin.Context) {
	products, err := h.repo.FindAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "berhasil mengambil daftar produk", products)
}

func (h *ProductHandler) Detail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	product, err := h.repo.FindByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "produk tidak ditemukan")
		return
	}
	response.Success(c, http.StatusOK, "berhasil mengambil detail produk", product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	product, err := h.repo.FindByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "produk tidak ditemukan")
		return
	}
	var req productRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid")
		return
	}
	product.SKU, product.Name, product.CategoryID = req.SKU, req.Name, req.CategoryID
	product.Unit, product.Description = req.Unit, req.Description

	if err := h.repo.Update(product); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "produk berhasil diperbarui", product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	if err := h.repo.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "produk berhasil dihapus", nil)
}