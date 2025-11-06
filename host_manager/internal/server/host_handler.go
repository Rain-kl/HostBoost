package server

import (
	"hostMgr/common/code"
	"net/http"

	"github.com/gin-gonic/gin"

	"hostMgr/internal/host"
)

// getHost 获取指定域名的 host 配置
func (h *Handler) getHost(c *gin.Context) {
	domain := c.Query("domain")
	hostEntry, err := h.svc.GetHost(domain)
	if err != nil {
		respondError(c, http.StatusNoContent, err)
		return
	}

	c.JSON(http.StatusOK, host.QueryHostResponse{
		Code:    code.Success,
		Message: "success",
		Data:    hostEntry,
	})
}

// listHosts 列出所有 host 配置
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

// createHost 创建新的 host 配置
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

// deleteHost 删除指定的 host 配置
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
