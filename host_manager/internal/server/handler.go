package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"

	"hostMgr/internal/host"
	"hostMgr/internal/opt"
	"hostMgr/internal/tool"
)

// Handler bundles HTTP handlers for host and opt operations.
type Handler struct {
	svc     *host.Service
	optSvc  *opt.Service
	toolSvc *tool.ToolService
	cache   *cache.Cache
}

// NewHandler creates a Gin handler with the provided services.
func NewHandler(svc *host.Service, optSvc *opt.Service, toolSvc *tool.ToolService) *Handler {
	// 创建缓存实例：默认过期时间 5 分钟，清理周期 10 分钟
	c := cache.New(5*time.Minute, 10*time.Minute)

	return &Handler{
		svc:     svc,
		optSvc:  optSvc,
		toolSvc: toolSvc,
		cache:   c,
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

	// tool 相关路由
	r.GET("/tool/webDetails", h.getWebDetails)
}

// respondError 通用错误响应函数
func respondError(c *gin.Context, status int, err error) {
	c.JSON(200, host.MutationResponse{
		Code:    status,
		Message: err.Error(),
	})
}
