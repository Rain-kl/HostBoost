package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"hostMgr/internal/tool"
)

// getWebDetails 获取指定域名的详细信息
// 1. 解析域名获取 IP 地址
// 2. 查询第一个 IP 的地理位置信息
func (h *Handler) getWebDetails(c *gin.Context) {
	domain := c.Query("domain")

	// 如果域名为空，返回 400
	if domain == "" {
		c.JSON(http.StatusBadRequest, tool.DetailResponse{
			Code:    "400",
			Message: "domain parameter is required",
			Data:    nil,
		})
		return
	}

	// 步骤 1: 解析域名获取 IP 地址
	ips, err := h.toolSvc.ResolveDomain(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, tool.DetailResponse{
			Code:    "500",
			Message: fmt.Sprintf("failed to resolve domain: %v", err),
			Data:    nil,
		})
		return
	}

	if len(ips) == 0 {
		c.JSON(http.StatusNotFound, tool.DetailResponse{
			Code:    "404",
			Message: "no IP addresses found for domain",
			Data:    nil,
		})
		return
	}

	// 步骤 2: 查询第一个 IP 的详细信息
	ipInfo, err := h.toolSvc.GetIPInfo(ips[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, tool.DetailResponse{
			Code:    "500",
			Message: fmt.Sprintf("failed to get IP info: %v", err),
			Data:    nil,
		})
		return
	}

	// 构造响应数据
	data := &tool.DomainDetail{
		Organization: ipInfo.Organization,
		IP:           ipInfo.IP,
		ISP:          ipInfo.ISP,
	}

	c.JSON(http.StatusOK, tool.DetailResponse{
		Code:    "200",
		Message: "success",
		Data:    data,
	})
}
