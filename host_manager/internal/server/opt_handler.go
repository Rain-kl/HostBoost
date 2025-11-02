package server

import (
	"hostMgr/common/code"
	"net/http"

	"github.com/gin-gonic/gin"

	"hostMgr/internal/opt"
)

// reportOpt 处理优选上报请求
func (h *Handler) reportOpt(c *gin.Context) {
	var req opt.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, opt.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if err := h.optSvc.ReportOpt(req); err != nil {
		c.JSON(http.StatusBadRequest, opt.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, opt.BaseResponse{
		Code:    code.Success,
		Message: "success",
	})
}

// getCurrentOpt 获取指定类型的当前优选
func (h *Handler) getCurrentOpt(c *gin.Context) {
	// 从 query 参数获取 type
	optType := c.Query("type")
	if optType == "" {
		c.JSON(http.StatusBadRequest, opt.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "type parameter is required",
		})
		return
	}

	optType, optInfo, err := h.optSvc.GetCurrentOpt(optType)
	if err != nil {
		c.JSON(http.StatusNotFound, opt.BaseResponse{
			Code:    code.NotFound,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, opt.GetOptResponse{
		Code:    code.Success,
		Message: "success",
		Type:    optType,
		Data:    optInfo,
	})
}

// changeOpt 更换指定类型的当前优选
func (h *Handler) changeOpt(c *gin.Context) {
	// 从 query 参数获取 type
	optType := c.Query("type")
	if optType == "" {
		c.JSON(http.StatusBadRequest, opt.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "type parameter is required",
		})
		return
	}

	if err := h.optSvc.ChangeOpt(optType); err != nil {
		c.JSON(http.StatusNotFound, opt.BaseResponse{
			Code:    code.NotFound,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, opt.BaseResponse{
		Code:    code.Success,
		Message: "success",
	})
}
