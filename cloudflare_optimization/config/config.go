package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 配置文件结构
type Config struct {
	// 延迟测速相关
	Routines      int    `yaml:"routines"`       // 延迟测速线程数
	PingTimes     int    `yaml:"ping_times"`     // 延迟测速次数
	TCPPort       int    `yaml:"tcp_port"`       // 测速端口
	Httping       bool   `yaml:"httping"`        // 是否使用HTTPing模式
	HttpingCode   int    `yaml:"httping_code"`   // HTTPing有效状态码
	HttpingCFColo string `yaml:"httping_cfcolo"` // 匹配指定地区
	URL           string `yaml:"url"`            // 测速地址

	// 过滤条件
	MaxDelay    int     `yaml:"max_delay"`     // 平均延迟上限(ms)
	MinDelay    int     `yaml:"min_delay"`     // 平均延迟下限(ms)
	MaxLossRate float64 `yaml:"max_loss_rate"` // 丢包率上限(0.00-1.00)

	// 下载测速相关
	TestCount       int     `yaml:"test_count"`       // 下载测速数量
	DownloadTime    int     `yaml:"download_time"`    // 下载测速时间(秒)
	MinSpeed        float64 `yaml:"min_speed"`        // 下载速度下限(MB/s)
	DisableDownload bool    `yaml:"disable_download"` // 禁用下载测速

	// IP相关
	IPFile  string `yaml:"ip_file"`  // IP段数据文件
	IPText  string `yaml:"ip_text"`  // 指定IP段数据
	TestAll bool   `yaml:"test_all"` // 测速全部IP

	// 输出相关
	PrintNum int    `yaml:"print_num"` // 显示结果数量
	Output   string `yaml:"output"`    // 输出结果文件

	// 其他
	Debug bool `yaml:"debug"` // 调试模式

	// 定时任务相关
	EnableSchedule bool   `yaml:"enable_schedule"` // 是否启用定时任务
	CronExpression string `yaml:"cron_expression"` // Cron 表达式

	// 上报相关
	EnableReport    bool   `yaml:"enable_report"`     // 是否启用自动上报
	ReportServerURL string `yaml:"report_server_url"` // 上报服务器地址
	ReportType      string `yaml:"report_type"`       // 上报类型
	ReportTimeout   int    `yaml:"report_timeout"`    // 上报超时时间(秒)
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Routines:        400,
		PingTimes:       4,
		TCPPort:         443,
		Httping:         false,
		HttpingCode:     0,
		HttpingCFColo:   "",
		URL:             "https://cf.xiu2.xyz/url",
		MaxDelay:        9999,
		MinDelay:        0,
		MaxLossRate:     1.0,
		TestCount:       10,
		DownloadTime:    10,
		MinSpeed:        0.0,
		DisableDownload: false,
		IPFile:          "ip.txt",
		IPText:          "",
		TestAll:         false,
		PrintNum:        10,
		Output:          "result.csv",
		Debug:           false,
		EnableSchedule:  true,
		CronExpression:  "0 0 0/3 * * ?", // 默认每3小时执行一次
		EnableReport:    true,
		ReportServerURL: "http://127.0.0.1:15920",
		ReportType:      "cloudflare",
		ReportTimeout:   10,
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := DefaultConfig()
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig 保存配置到文件
func SaveConfig(filename string, config *Config) error {
	// 添加注释说明
	content := `# CloudflareSpeedTest 配置文件
# 详细说明请参考: https://github.com/XIU2/CloudflareSpeedTest

# ===== 延迟测速相关 =====
# 延迟测速线程数 (默认 400，最多 1000)
routines: ` + fmt.Sprintf("%d", config.Routines) + `

# 延迟测速次数 (默认 4)
ping_times: ` + fmt.Sprintf("%d", config.PingTimes) + `

# 测速端口 (默认 443)
tcp_port: ` + fmt.Sprintf("%d", config.TCPPort) + `

# 是否使用 HTTPing 模式 (默认 false，使用 TCPing)
httping: ` + fmt.Sprintf("%v", config.Httping) + `

# HTTPing 有效状态码 (默认 0 表示 200/301/302，仅限一个)
httping_code: ` + fmt.Sprintf("%d", config.HttpingCode) + `

# 匹配指定地区 (IATA 机场代码，英文逗号分隔，如: HKG,KHH,NRT,LAX)
httping_cfcolo: "` + config.HttpingCFColo + `"

# 测速地址 (默认 https://cf.xiu2.xyz/url)
url: "` + config.URL + `"

# ===== 过滤条件 =====
# 平均延迟上限，单位ms (默认 9999)
max_delay: ` + fmt.Sprintf("%d", config.MaxDelay) + `

# 平均延迟下限，单位ms (默认 0)
min_delay: ` + fmt.Sprintf("%d", config.MinDelay) + `

# 丢包率上限，范围 0.00-1.00 (默认 1.00)
max_loss_rate: ` + fmt.Sprintf("%.2f", config.MaxLossRate) + `

# ===== 下载测速相关 =====
# 下载测速数量 (默认 10)
test_count: ` + fmt.Sprintf("%d", config.TestCount) + `

# 下载测速时间，单位秒 (默认 10)
download_time: ` + fmt.Sprintf("%d", config.DownloadTime) + `

# 下载速度下限，单位MB/s (默认 0.00)
min_speed: ` + fmt.Sprintf("%.2f", config.MinSpeed) + `

# 禁用下载测速 (默认 false)
disable_download: ` + fmt.Sprintf("%v", config.DisableDownload) + `

# ===== IP 相关 =====
# IP段数据文件 (默认 ip.txt)
ip_file: "` + config.IPFile + `"

# 指定IP段数据 (直接指定，英文逗号分隔，如: 1.1.1.1,2.2.2.2/24)
ip_text: "` + config.IPText + `"

# 测速全部IP (默认 false，每个 /24 段随机测速一个 IP)
test_all: ` + fmt.Sprintf("%v", config.TestAll) + `

# ===== 输出相关 =====
# 显示结果数量 (默认 10，设为 0 则不显示)
print_num: ` + fmt.Sprintf("%d", config.PrintNum) + `

# 输出结果文件 (默认 result.csv，设为空字符串则不输出)
output: "` + config.Output + `"

# ===== 其他 =====
# 调试模式 (默认 false)
debug: ` + fmt.Sprintf("%v", config.Debug) + `

# ===== 定时任务相关 =====
# 是否启用定时任务 (默认 false)
enable_schedule: ` + fmt.Sprintf("%v", config.EnableSchedule) + `

# Cron 表达式 (默认 "0 */1 * * *" 表示每1小时执行一次)
# 格式: 分 时 日 月 周
# 示例: "0 0 * * *" 每天0点, "0 */2 * * *" 每2小时, "0 0 * * 0" 每周日0点
cron_expression: "` + config.CronExpression + `"

# ===== 上报相关 =====
# 是否启用自动上报 (默认 false)
enable_report: ` + fmt.Sprintf("%v", config.EnableReport) + `

# 上报服务器地址 (默认 http://127.0.0.1:15920)
report_server_url: "` + config.ReportServerURL + `"

# 上报类型 (默认 "msn")
report_type: "` + config.ReportType + `"

# 上报超时时间，单位秒 (默认 10)
report_timeout: ` + fmt.Sprintf("%d", config.ReportTimeout) + `
`

	return os.WriteFile(filename, []byte(content), 0644)
}

// CreateDefaultConfigFile 创建默认配置文件
func CreateDefaultConfigFile(filename string) error {
	config := DefaultConfig()
	return SaveConfig(filename, config)
}
