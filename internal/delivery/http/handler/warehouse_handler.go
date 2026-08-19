package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prasetyoary48/wms/internal/domain"
	"github.com/prasetyoary48/wms/internal/repository"
	"github.com/prasetyoary48/wms/pkg/response"
)

type WarehouseHandler struct {
	whRepo  repository.WarehouseRepository
	locRepo repository.LocationRepository
}

func NewWarehouseHandler(whRepo repository.WarehouseRepository, locRepo repository.LocationRepository) *WarehouseHandler {
	return &WarehouseHandler{whRepo: whRepo, locRepo: locRepo}
}

type warehouseRequest struct {
	Code    string `json:"code" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}

func (h *WarehouseHandler) Create(c *gin.Context) {
	var req warehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid")
		return
	}
	wh := &domain.Warehouse{Code: req.Code, Name: req.Name, Address: req.Address, IsActive: true}
	if err := h.whRepo.Create(wh); err != nil {
		response.Error(c, http.StatusBadRequest, "gagal membuat gudang, kode mungkin sudah dipakai")
		return
	}
	response.Success(c, http.StatusCreated, "gudang berhasil dibuat", wh)
}

func (h *WarehouseHandler) List(c *gin.Context) {
	list, err := h.whRepo.FindAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "berhasil mengambil daftar gudang", list)
}

// --- Location / Rack ---

type locationRequest struct {
	WarehouseID uint   `json:"warehouse_id" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Zone        string `json:"zone"`
	Capacity    int    `json:"capacity"`
}

func (h *WarehouseHandler) CreateLocation(c *gin.Context) {
	var req locationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid")
		return
	}
	loc := &domain.Location{
		WarehouseID: req.WarehouseID, Code: req.Code, Zone: req.Zone,
		Capacity: req.Capacity, IsActive: true,
	}
	if err := h.locRepo.Create(loc); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "lokasi/rak berhasil dibuat", loc)
}

func (h *WarehouseHandler) ListLocations(c *gin.Context) {
	warehouseID, err := strconv.Atoi(c.Query("warehouse_id"))
	if err == nil && warehouseID > 0 {
		list, err := h.locRepo.FindByWarehouse(uint(warehouseID))
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		response.Success(c, http.StatusOK, "berhasil mengambil daftar lokasi", list)
		return
	}

	list, err := h.locRepo.FindAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "berhasil mengambil daftar lokasi", list)
}