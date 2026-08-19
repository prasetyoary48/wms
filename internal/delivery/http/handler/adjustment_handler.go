package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prasetyoary48/wms/internal/usecase"
	"github.com/prasetyoary48/wms/pkg/response"
)

type AdjustmentHandler struct {
	adjUC usecase.AdjustmentUsecase
}

func NewAdjustmentHandler(adjUC usecase.AdjustmentUsecase) *AdjustmentHandler {
	return &AdjustmentHandler{adjUC: adjUC}
}

type requestAdjustmentRequest struct {
	ProductID    uint   `json:"product_id" binding:"required"`
	LocationID   uint   `json:"location_id" binding:"required"`
	RequestedQty int    `json:"requested_qty"`
	Reason       string `json:"reason"`
}

func (h *AdjustmentHandler) Request(c *gin.Context) {
	var req requestAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid")
		return
	}
	userID := c.GetUint("user_id")

	adj, err := h.adjUC.Request(req.ProductID, req.LocationID, req.RequestedQty, req.Reason, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "pengajuan adjustment berhasil dibuat", adj)
}

func (h *AdjustmentHandler) Approve(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	approverID := c.GetUint("user_id")

	if err := h.adjUC.Approve(uint(id), approverID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "adjustment disetujui", nil)
}

func (h *AdjustmentHandler) Reject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	var req rejectRequest
	_ = c.ShouldBindJSON(&req)
	approverID := c.GetUint("user_id")

	if err := h.adjUC.Reject(uint(id), approverID, req.Reason); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "adjustment ditolak", nil)
}

func (h *AdjustmentHandler) PendingList(c *gin.Context) {
	list, err := h.adjUC.PendingList()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "berhasil mengambil daftar pengajuan adjustment", list)
}

func (h *AdjustmentHandler) MyRequests(c *gin.Context) {
	userID := c.GetUint("user_id")
	list, err := h.adjUC.MyRequests(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "berhasil mengambil pengajuan anda", list)
}