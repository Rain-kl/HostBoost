package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"cf_opt/utils"
)

// OptVo 优选IP数据结构
type OptVo struct {
	IP    string `json:"ip"`
	Delay string `json:"delay"`
	Rate  string `json:"rate"`
}

// OptRequest 上报请求结构
type OptRequest struct {
	Type string  `json:"type"`
	Data []OptVo `json:"data"`
}

// BaseResponse 响应结构
type BaseResponse struct {
	Code    interface{} `json:"code"`
	Message string      `json:"message"`
}

// ReportConfig 上报配置
type ReportConfig struct {
	ServerURL string // 服务器地址，默认 http://127.0.0.1:15920
	Type      string // 类型，如 "cloudflare"
	Timeout   int    // 超时时间（秒），默认 10
}

var DefaultReportConfig = ReportConfig{
	ServerURL: "http://127.0.0.1:15920",
	Type:      "cloudflare",
	Timeout:   10,
}

// Report 将优选的 IP 数据上报到远程服务器
// speedData: 速度测试结果数据（已按速度排序）
// config: 上报配置，如果为 nil 则使用默认配置
// 返回值: error 如果上报失败则返回错误
// 注意：只上报前 5 个最好的 IP
func Report(speedData []utils.CloudflareIPData, config *ReportConfig) error {
	if config == nil {
		config = &DefaultReportConfig
	}

	// 只上报前 5 个最好的 IP
	reportCount := len(speedData)
	if reportCount > 5 {
		reportCount = 5
	}

	if reportCount == 0 {
		return fmt.Errorf("没有可上报的 IP 数据")
	}

	// 转换数据格式
	optData := make([]OptVo, 0, reportCount)
	for i := 0; i < reportCount; i++ {
		data := speedData[i]
		// 计算延迟（毫秒）
		delayMs := data.Delay.Seconds() * 1000
		// 计算下载速度（MB/s）
		speedMBps := data.DownloadSpeed / 1024 / 1024

		optData = append(optData, OptVo{
			IP:    data.IP.String(),
			Delay: strconv.FormatFloat(delayMs, 'f', 0, 64),
			Rate:  strconv.FormatFloat(speedMBps, 'f', 2, 64),
		})
	}

	// 构建请求
	request := OptRequest{
		Type: config.Type,
		Data: optData,
	}

	// 序列化为 JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	// 发送 POST 请求
	url := config.ServerURL + "/opt/report"
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析响应
	var response BaseResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("服务器返回错误状态: %d, 消息: %s", resp.StatusCode, response.Message)
	}
	// 兼容 code 为数字或字符串
	var codeStr string
	switch v := response.Code.(type) {
	case string:
		codeStr = v
	case float64:
		codeStr = strconv.Itoa(int(v))
	default:
		codeStr = fmt.Sprintf("%v", v)
	}

	fmt.Printf("上报成功: 已上报 %d 个最优 IP, code=%s, message=%s\n", reportCount, codeStr, response.Message)
	return nil
}
