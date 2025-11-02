package server

import (
	"hostMgr/common/code"
	"net/http"

	"github.com/gin-gonic/gin"

	"hostMgr/internal/host"
)

// Handler bundles HTTP handlers for host operations.
type Handler struct {
	svc *host.Service
}

// NewHandler creates a Gin handler with the provided service.
func NewHandler(svc *host.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes wires the host endpoints into the given router.
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/host", h.getHost)
	r.POST("/host", h.createHost)
	r.DELETE("/host", h.deleteHost)
	r.GET("/host/list", h.listHosts)
}

func (h *Handler) getHost(c *gin.Context) {
	domain := c.Query("domain")
	hostEntry, err := h.svc.GetHost(domain)
	if err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, host.QueryHostResponse{
		Code:    code.Success,
		Message: "success",
		Data:    hostEntry,
	})
}

func (h *Handler) listHosts(c *gin.Context) {
	hosts := h.svc.ListHosts()

	c.JSON(http.StatusOK, host.QueryHostListResponse{
		Code:    code.Success,
		Message: "success",
		Data: host.QueryHostListResult{
			Total: len(hosts),
			List:  hosts,
		},
	})
}

func (h *Handler) createHost(c *gin.Context) {
	var req host.AddHostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.CreateHost(req); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, host.MutationResponse{
		Code:    code.Success,
		Message: "created",
	})
}

func (h *Handler) deleteHost(c *gin.Context) {
	var req host.DeleteHostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.DeleteHost(req.Domain); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, host.MutationResponse{
		Code:    code.Success,
		Message: "deleted",
	})
}

func respondError(c *gin.Context, status int, err error) {
	c.JSON(status, host.MutationResponse{
		Code:    status,
		Message: err.Error(),
	})
}
