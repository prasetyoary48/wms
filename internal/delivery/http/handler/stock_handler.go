package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prasetyoary48/wms/internal/usecase"
	"github.com/prasetyoary48/wms/pkg/response"
)

type StockHandler struct {
	stockUC usecase.StockUsecase
}

func NewStockHandler(stockUC usecase.StockUsecase) *StockHandler {
	return &StockHandler{stockUC: stockUC}
}

type stockInOutRequest struct {
	ProductID  uint   `json:"product_id" binding:"required"`
	LocationID uint   `json:"location_id" binding:"required"`
	Qty        int    `json:"qty" binding:"required"`
	Note       string `json:"note"`
}

func (h *StockHandler) StockIn(c *gin.Context) {
	var req stockInOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid")
		return
	}
	userID := c.GetUint("user_id")

	movement, err := h.stockUC.StockIn(req.ProductID, req.LocationID, req.Qty, req.Note, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "stock in berhasil dicatat", movement)
}

func (h *StockHandler) StockOut(c *gin.Context) {
	var req stockInOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid")
		return
	}
	userID := c.GetUint("user_id")

	movement, err := h.stockUC.StockOut(req.ProductID, req.LocationID, req.Qty, req.Note, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "stock out berhasil dicatat", movement)
}

// LocationsOfProduct menjawab "produk X ada di rak mana saja".
func (h *StockHandler) LocationsOfProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "product id tidak valid")
		return
	}

	stocks, err := h.stockUC.LocationsOfProduct(uint(productID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "berhasil mengambil lokasi produk", stocks)
}

func (h *StockHandler) StockInLocation(c *gin.Context) {
	locationID, err := strconv.Atoi(c.Param("locationId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "location id tidak valid")
		return
	}

	stocks, err := h.stockUC.StockInLocation(uint(locationID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "berhasil mengambil stok pada lokasi", stocks)
}

func (h *StockHandler) AllStocks(c *gin.Context) {
	stocks, err := h.stockUC.AllStocks()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "berhasil mengambil semua stok", stocks)
}