package server

import (
	"github.com/gin-gonic/gin"

	"hostMgr/internal/host"
	"hostMgr/internal/opt"
)

// Handler bundles HTTP handlers for host and opt operations.
type Handler struct {
	svc    *host.Service
	optSvc *opt.Service
}

// NewHandler creates a Gin handler with the provided services.
func NewHandler(svc *host.Service, optSvc *opt.Service) *Handler {
	return &Handler{
		svc:    svc,
		optSvc: optSvc,
	}
}

// RegisterRoutes wires all endpoints into the given router.
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// host 相关路由
	r.GET("/host", h.getHost)
	r.POST("/host", h.createHost)
	r.DELETE("/host", h.deleteHost)
	r.GET("/host/list", h.listHosts)

	// opt 相关路由
	r.POST("/opt/report", h.reportOpt)
	r.GET("/opt", h.getCurrentOpt)
	r.GET("/opt/change", h.changeOpt)
}

// respondError 通用错误响应函数
func respondError(c *gin.Context, status int, err error) {
	c.JSON(status, host.MutationResponse{
		Code:    status,
		Message: err.Error(),
	})
}
