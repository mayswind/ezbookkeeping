package ext

import (
	"github.com/gin-gonic/gin"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	extapi "github.com/mayswind/ezbookkeeping/pkg/ext/api"
	extmw "github.com/mayswind/ezbookkeeping/pkg/ext/middleware"
)

// ApiBinder wraps an API handler into a gin handler (upstream's bindApi)
type ApiBinder func(core.ApiHandlerFunc) gin.HandlerFunc

// MiddlewareBinder wraps a middleware into a gin handler (upstream's bindMiddleware)
type MiddlewareBinder func(core.MiddlewareHandlerFunc) gin.HandlerFunc

// DelegationMiddleware returns the middleware that lets managers and staff work on an owner's data.
// Register it on the authenticated API group right after authorization, before any route.
func DelegationMiddleware(bind MiddlewareBinder) gin.HandlerFunc {
	return bind(extmw.Delegation())
}

// RegisterRoutes registers all ext routes on the authenticated "/api/v1" group.
// Who may call what is decided by pkg/ext/permissions, not here.
func RegisterRoutes(group *gin.RouterGroup, bind ApiBinder) {
	h := extapi.API

	// Business membership: always about the logged-in person
	group.GET("/ext/me/businesses.json", bind(h.MyBusinessesHandler))
	group.GET("/ext/me/settings.json", bind(h.MySettingsHandler))
	group.POST("/ext/me/settings/update.json", bind(h.MySettingsUpdateHandler))
	group.POST("/ext/staff/invite.json", bind(h.StaffInviteHandler))
	group.GET("/ext/staff/list.json", bind(h.StaffListHandler))
	group.GET("/ext/staff/people.json", bind(h.MyPeopleListHandler))
	group.POST("/ext/staff/set_role.json", bind(h.StaffSetRoleHandler))
	group.POST("/ext/staff/remove.json", bind(h.StaffRemoveHandler))
	group.POST("/ext/staff/respond.json", bind(h.StaffRespondHandler))
	group.POST("/ext/staff/leave.json", bind(h.StaffLeaveHandler))
	group.GET("/ext/audit/list.json", bind(h.AuditListHandler))

	// People (names for "recorded by") and reports
	group.GET("/ext/people/list.json", bind(h.PeopleListHandler))
	group.GET("/ext/reports/stock_value.json", bind(h.StockValueReportHandler))
	group.GET("/ext/reports/low_stock.json", bind(h.LowStockReportHandler))
	group.GET("/ext/reports/receivables.json", bind(h.ReceivablesReportHandler))

	// Locations
	group.GET("/ext/locations/list.json", bind(h.LocationListHandler))
	group.POST("/ext/locations/add.json", bind(h.LocationAddHandler))
	group.POST("/ext/locations/modify.json", bind(h.LocationModifyHandler))
	group.POST("/ext/locations/delete.json", bind(h.LocationDeleteHandler))

	// Items and stock
	group.GET("/ext/items/list.json", bind(h.ItemListHandler))
	group.GET("/ext/items/stock.json", bind(h.StockLevelsHandler))
	group.POST("/ext/items/add.json", bind(h.ItemAddHandler))
	group.POST("/ext/items/modify.json", bind(h.ItemModifyHandler))
	group.POST("/ext/items/delete.json", bind(h.ItemDeleteHandler))
	group.POST("/ext/stock/receive.json", bind(h.StockReceiveHandler))
	group.POST("/ext/stock/adjust.json", bind(h.StockAdjustHandler))
	group.POST("/ext/stock/transfer.json", bind(h.StockTransferHandler))
	group.GET("/ext/stock/movements/list.json", bind(h.StockMovementsHandler))

	// Customers, sales and credit repayments
	group.GET("/ext/customers/list.json", bind(h.CustomerListHandler))
	group.GET("/ext/customers/balances.json", bind(h.CustomerBalancesHandler))
	group.POST("/ext/customers/add.json", bind(h.CustomerAddHandler))
	group.POST("/ext/customers/modify.json", bind(h.CustomerModifyHandler))
	group.POST("/ext/customers/delete.json", bind(h.CustomerDeleteHandler))
	group.POST("/ext/sales/add.json", bind(h.SaleCreateHandler))
	group.GET("/ext/sales/list.json", bind(h.SaleListHandler))
	group.GET("/ext/sales/get.json", bind(h.SaleGetHandler))
	group.POST("/ext/sales/void.json", bind(h.SaleVoidHandler))
	group.POST("/ext/repayments/add.json", bind(h.RepaymentCreateHandler))
	group.GET("/ext/repayments/list.json", bind(h.RepaymentListHandler))
}
