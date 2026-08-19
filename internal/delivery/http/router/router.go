package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prasetyoary48/wms/internal/delivery/http/handler"
	"github.com/prasetyoary48/wms/internal/delivery/http/middleware"
	"github.com/prasetyoary48/wms/internal/domain"
)

type Handlers struct {
	Auth       *handler.AuthHandler
	Product    *handler.ProductHandler
	Warehouse  *handler.WarehouseHandler
	Stock      *handler.StockHandler
	Transfer   *handler.TransferHandler
	Adjustment *handler.AdjustmentHandler
}

func Setup(r *gin.Engine, h *Handlers, jwtSecret string) {
	api := r.Group("/api/v1")

	// --- Public ---
	api.POST("/auth/login", h.Auth.Login)

	// --- Authenticated (semua role) ---
	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Lihat stok & lokasi rak — semua role boleh lihat
		auth.GET("/stocks", h.Stock.AllStocks)
		auth.GET("/stocks/product/:productId", h.Stock.LocationsOfProduct)
		auth.GET("/stocks/location/:locationId", h.Stock.StockInLocation)
		auth.GET("/warehouses", h.Warehouse.List)
		auth.GET("/locations", h.Warehouse.ListLocations)
		auth.GET("/products", h.Product.List)
		auth.GET("/products/:id", h.Product.Detail)

		// Staff: lihat kegiatan/pengajuan sendiri
		auth.GET("/adjustments/mine", h.Adjustment.MyRequests)
	}

	// --- Admin only ---
	admin := api.Group("")
	admin.Use(middleware.AuthMiddleware(jwtSecret), middleware.RequireRole(domain.RoleAdmin))
	{
		admin.POST("/auth/register", h.Auth.Register) // manage user

		admin.POST("/products", h.Product.Create)
		admin.PUT("/products/:id", h.Product.Update)
		admin.DELETE("/products/:id", h.Product.Delete)

		admin.POST("/warehouses", h.Warehouse.Create)
		admin.POST("/locations", h.Warehouse.CreateLocation)
		// TODO: manage supplier, manage role -> pola sama seperti product/warehouse handler
	}

	// --- Staff only (operasional harian) ---
	staff := api.Group("")
	staff.Use(middleware.AuthMiddleware(jwtSecret), middleware.RequireRole(domain.RoleStaff, domain.RoleAdmin))
	{
		staff.POST("/stocks/in", h.Stock.StockIn)
		staff.POST("/stocks/out", h.Stock.StockOut)
		staff.POST("/transfers", h.Transfer.RequestTransfer)
		staff.POST("/adjustments", h.Adjustment.Request)
		// TODO: stock opname input -> pola sama, tambahkan opname handler di sini
	}

	// --- SPV only (approval & monitoring) ---
	spv := api.Group("")
	spv.Use(middleware.AuthMiddleware(jwtSecret), middleware.RequireRole(domain.RoleSPV, domain.RoleAdmin))
	{
		spv.GET("/transfers/pending", h.Transfer.PendingList)
		spv.PATCH("/transfers/:id/approve", h.Transfer.Approve)
		spv.PATCH("/transfers/:id/reject", h.Transfer.Reject)

		spv.GET("/adjustments/pending", h.Adjustment.PendingList)
		spv.PATCH("/adjustments/:id/approve", h.Adjustment.Approve)
		spv.PATCH("/adjustments/:id/reject", h.Adjustment.Reject)
		// TODO: approve stock opname, endpoint laporan & monitoring staff
	}
}