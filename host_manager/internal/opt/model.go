package opt

// OptInfo 优选信息模型
type OptInfo struct {
	IP    string `json:"ip" binding:"required"`
	Delay string `json:"delay" binding:"required"`
	Rate  string `json:"rate" binding:"required"`
}

// OptData 优选数据(包含类型和 IP 列表)
type OptData struct {
	Type    string    `json:"type"`
	Data    []OptInfo `json:"data"`
	Current int       `json:"current"` // 当前使用的优选索引
}

// OptStore 用于 JSON 文件存储的结构
type OptStore struct {
	Opts map[string]*OptData `json:"opts"` // key 为 type
}

// ReportRequest 优选上报请求
type ReportRequest struct {
	Type string    `json:"type" binding:"required"`
	Data []OptInfo `json:"data" binding:"required"`
}

// BaseResponse 基础响应
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetOptResponse 获取当前优选响应
type GetOptResponse struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Type    string  `json:"type"`
	Data    OptInfo `json:"data"`
}
