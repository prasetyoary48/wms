package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prasetyoary48/wms/internal/usecase"
	"github.com/prasetyoary48/wms/pkg/response"
)

type TransferHandler struct {
	transferUC usecase.TransferUsecase
}

func NewTransferHandler(transferUC usecase.TransferUsecase) *TransferHandler {
	return &TransferHandler{transferUC: transferUC}
}

type requestTransferRequest struct {
	ProductID      uint   `json:"product_id" binding:"required"`
	FromLocationID uint   `json:"from_location_id" binding:"required"`
	ToLocationID   uint   `json:"to_location_id" binding:"required"`
	Qty            int    `json:"qty" binding:"required"`
	Note           string `json:"note"`
}

// RequestTransfer — dipakai Staff untuk mengajukan pemindahan barang
// (bisa antar rak dalam 1 gudang, maupun antar gudang berbeda).
func (h *TransferHandler) RequestTransfer(c *gin.Context) {
	var req requestTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid")
		return
	}
	userID := c.GetUint("user_id")

	movement, err := h.transferUC.RequestTransfer(req.ProductID, req.FromLocationID, req.ToLocationID, req.Qty, req.Note, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "pengajuan transfer berhasil dibuat, menunggu approval SPV", movement)
}

// Approve — hanya SPV.
func (h *TransferHandler) Approve(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	approverID := c.GetUint("user_id")

	if err := h.transferUC.ApproveTransfer(uint(id), approverID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "transfer barang disetujui", nil)
}

type rejectRequest struct {
	Reason string `json:"reason"`
}

func (h *TransferHandler) Reject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	var req rejectRequest
	_ = c.ShouldBindJSON(&req)
	approverID := c.GetUint("user_id")

	if err := h.transferUC.RejectTransfer(uint(id), approverID, req.Reason); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "transfer barang ditolak", nil)
}

func (h *TransferHandler) PendingList(c *gin.Context) {
	list, err := h.transferUC.PendingTransfers()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "berhasil mengambil daftar pengajuan transfer", list)
}