package main

import (
	"cf_opt/config"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cf_opt/task"
	"cf_opt/utils"

	"github.com/robfig/cron/v3"
)

var (
	version = "cfst-v2.3.4" // 版本号
)

func init() {
	var printVersion bool
	var configFile string
	var help = `
CloudflareSpeedTest ` + version + `
测试各个 CDN 或网站所有 IP 的延迟和速度，获取最快 IP (IPv4+IPv6)！
https://github.com/XIU2/CloudflareSpeedTest

参数：
    -c config.yaml
        指定配置文件；默认使用当前目录下的 config.yaml，不存在时自动创建；(默认 config.yaml)
    -n 200
        延迟测速线程；越多延迟测速越快，性能弱的设备 (如路由器) 请勿太高；(默认 200 最多 1000)
    -t 4
        延迟测速次数；单个 IP 延迟测速的次数；(默认 4 次)
    -dn 10
        下载测速数量；延迟测速并排序后，从最低延迟起下载测速的数量；(默认 10 个)
    -dt 10
        下载测速时间；单个 IP 下载测速最长时间，不能太短；(默认 10 秒)
    -tp 443
        指定测速端口；延迟测速/下载测速时使用的端口；(默认 443 端口)
    -url https://cf.xiu2.xyz/url
        指定测速地址；延迟测速(HTTPing)/下载测速时使用的地址，默认地址不保证可用性，建议自建；

    -httping
        切换测速模式；延迟测速模式改为 HTTP 协议，所用测试地址为 [-url] 参数；(默认 TCPing)
    -httping-code 200
        有效状态代码；HTTPing 延迟测速时网页返回的有效 HTTP 状态码，仅限一个；(默认 200 301 302)
    -cfcolo HKG,KHH,NRT,LAX,SEA,SJC,FRA,MAD
        匹配指定地区；IATA 机场地区码或国家/城市码，英文逗号分隔，仅 HTTPing 模式可用；(默认 所有地区)

    -tl 200
        平均延迟上限；只输出低于指定平均延迟的 IP，各上下限条件可搭配使用；(默认 9999 ms)
    -tll 40
        平均延迟下限；只输出高于指定平均延迟的 IP；(默认 0 ms)
    -tlr 0.2
        丢包几率上限；只输出低于/等于指定丢包率的 IP，范围 0.00~1.00，0 过滤掉任何丢包的 IP；(默认 1.00)
    -sl 5
        下载速度下限；只输出高于指定下载速度的 IP，凑够指定数量 [-dn] 才会停止测速；(默认 0.00 MB/s)

    -p 10
        显示结果数量；测速后直接显示指定数量的结果，为 0 时不显示结果直接退出；(默认 10 个)
    -f ip.txt
        IP段数据文件；如路径含有空格请加上引号；支持其他 CDN IP段；(默认 ip.txt)
    -ip 1.1.1.1,2.2.2.2/24,2606:4700::/32
        指定IP段数据；直接通过参数指定要测速的 IP 段数据，英文逗号分隔；(默认 空)
    -o result.csv
        写入结果文件；如路径含有空格请加上引号；值为空时不写入文件 [-o ""]；(默认 result.csv)

    -dd
        禁用下载测速；禁用后测速结果会按延迟排序 (默认按下载速度排序)；(默认 启用)
    -allip
        测速全部的IP；对 IP 段中的每个 IP (仅支持 IPv4) 进行测速；(默认 每个 /24 段随机测速一个 IP)

    -schedule
        启用定时任务；使程序以定时任务模式运行，按照 Cron 表达式定期执行测速；(默认 关闭)
    -cron "0 */6 * * *"
        Cron 表达式；指定定时任务的执行时间，支持标准 5 段格式或秒级 6 段格式；(默认 "0 */6 * * *")
        格式说明：
          标准格式(5段): 分 时 日 月 周
          秒级格式(6段): 秒 分 时 日 月 周
        示例：
          每6小时执行: 0 */6 * * *
          每天凌晨执行: 0 0 * * *
          每2小时执行: 0 */2 * * *
          每30分钟执行: */30 * * * *

    -debug
        调试输出模式；会在一些非预期情况下输出更多日志以便判断原因；(默认 关闭)

    -v
        打印程序版本 + 检查版本更新
    -h
        打印帮助说明

注意：命令行参数优先级高于配置文件，可以覆盖配置文件中的设置
`
	var minDelay, maxDelay, downloadTime int
	var maxLossRate float64
	var enableSchedule bool
	var cronExpression string

	// 先临时解析以获取配置文件路径
	var tempConfigFile string
	var tempVersion bool
	tempFlagSet := flag.NewFlagSet("temp", flag.ContinueOnError)
	tempFlagSet.StringVar(&tempConfigFile, "c", "config.yaml", "")
	tempFlagSet.BoolVar(&tempVersion, "v", false, "")
	tempFlagSet.Parse(os.Args[1:])

	// 如果是 -v，设置 printVersion 然后继续（稍后会退出）
	if tempVersion {
		printVersion = true
		configFile = tempConfigFile
	} else {
		configFile = tempConfigFile
	}

	// 加载配置文件（先设置默认值）
	globalConfig = loadConfigFile(configFile)
	config := globalConfig

	// 用配置文件的值初始化变量
	task.Routines = config.Routines
	task.PingTimes = config.PingTimes
	task.TestCount = config.TestCount
	downloadTime = config.DownloadTime
	task.TCPPort = config.TCPPort
	task.URL = config.URL
	task.Httping = config.Httping
	task.HttpingStatusCode = config.HttpingCode
	task.HttpingCFColo = config.HttpingCFColo
	maxDelay = config.MaxDelay
	minDelay = config.MinDelay
	maxLossRate = config.MaxLossRate
	task.MinSpeed = config.MinSpeed
	utils.PrintNum = config.PrintNum
	task.IPFile = config.IPFile
	task.IPText = config.IPText
	utils.Output = config.Output
	task.Disable = config.DisableDownload
	task.TestAll = config.TestAll
	utils.Debug = config.Debug
	enableSchedule = config.EnableSchedule
	cronExpression = config.CronExpression

	// 定义所有命令行参数
	flag.StringVar(&configFile, "c", configFile, "配置文件路径")
	flag.BoolVar(&printVersion, "v", printVersion, "打印程序版本")
	flag.IntVar(&task.Routines, "n", task.Routines, "延迟测速线程")
	flag.IntVar(&task.PingTimes, "t", task.PingTimes, "延迟测速次数")
	flag.IntVar(&task.TestCount, "dn", task.TestCount, "下载测速数量")
	flag.IntVar(&downloadTime, "dt", downloadTime, "下载测速时间")
	flag.IntVar(&task.TCPPort, "tp", task.TCPPort, "指定测速端口")
	flag.StringVar(&task.URL, "url", task.URL, "指定测速地址")

	flag.BoolVar(&task.Httping, "httping", task.Httping, "切换测速模式")
	flag.IntVar(&task.HttpingStatusCode, "httping-code", task.HttpingStatusCode, "有效状态代码")
	flag.StringVar(&task.HttpingCFColo, "cfcolo", task.HttpingCFColo, "匹配指定地区")

	flag.IntVar(&maxDelay, "tl", maxDelay, "平均延迟上限")
	flag.IntVar(&minDelay, "tll", minDelay, "平均延迟下限")
	flag.Float64Var(&maxLossRate, "tlr", maxLossRate, "丢包几率上限")
	flag.Float64Var(&task.MinSpeed, "sl", task.MinSpeed, "下载速度下限")

	flag.IntVar(&utils.PrintNum, "p", utils.PrintNum, "显示结果数量")
	flag.StringVar(&task.IPFile, "f", task.IPFile, "IP段数据文件")
	flag.StringVar(&task.IPText, "ip", task.IPText, "指定IP段数据")
	flag.StringVar(&utils.Output, "o", utils.Output, "输出结果文件")

	flag.BoolVar(&task.Disable, "dd", task.Disable, "禁用下载测速")
	flag.BoolVar(&task.TestAll, "allip", task.TestAll, "测速全部 IP")

	flag.BoolVar(&utils.Debug, "debug", utils.Debug, "调试输出模式")

	flag.BoolVar(&enableSchedule, "schedule", enableSchedule, "启用定时任务")
	flag.StringVar(&cronExpression, "cron", cronExpression, "Cron 表达式")

	flag.Usage = func() { fmt.Print(help) }
	// 解析命令行参数
	flag.Parse()

	if task.MinSpeed > 0 && time.Duration(maxDelay)*time.Millisecond == utils.InputMaxDelay {
		utils.Yellow.Println("[提示] 在使用 [-sl] 参数时，建议搭配 [-tl] 参数，以避免因凑不够 [-dn] 数量而一直测速...")
	}
	utils.InputMaxDelay = time.Duration(maxDelay) * time.Millisecond
	utils.InputMinDelay = time.Duration(minDelay) * time.Millisecond
	utils.InputMaxLossRate = float32(maxLossRate)
	task.Timeout = time.Duration(downloadTime) * time.Second
	task.HttpingCFColomap = task.MapColoMap()

	if printVersion {
		println(version)
		os.Exit(0)
	}

	// 如果启用定时任务，则启动 cron 调度器
	if enableSchedule {
		startScheduler(cronExpression)
	}
}

// loadConfigFile 加载配置文件，如果不存在则创建默认配置文件
func loadConfigFile(filename string) *config.Config {
	// 检查配置文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// 配置文件不存在，创建默认配置文件
		fmt.Printf("配置文件 %s 不存在，正在创建默认配置文件...\n", filename)
		if err := config.CreateDefaultConfigFile(filename); err != nil {
			fmt.Printf("创建配置文件失败: %v\n使用默认配置继续运行...\n", err)
			return config.DefaultConfig()
		}
		fmt.Printf("已创建默认配置文件: %s\n", filename)
	}

	// 加载配置文件
	yamlConfig, err := config.LoadConfig(filename)
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n使用默认配置继续运行...\n", err)
		return config.DefaultConfig()
	}

	return yamlConfig
}

// 全局配置变量
var globalConfig *config.Config

func run() {
	task.InitRandSeed() // 置随机数种子

	fmt.Printf("# XIU2/CloudflareSpeedTest %s \n\n", version)

	// 开始延迟测速 + 过滤延迟/丢包
	pingData := task.NewPing().Run().FilterDelay().FilterLossRate()
	// 开始下载测速
	speedData := task.TestDownloadSpeed(pingData)
	utils.ExportCsv(speedData) // 输出文件
	speedData.Print()

	// 如果启用了自动上报，则上报结果
	if globalConfig != nil && globalConfig.EnableReport && len(speedData) > 0 {
		fmt.Println("\n正在上报结果到服务器...")
		reportConfig := &task.ReportConfig{
			ServerURL: globalConfig.ReportServerURL,
			Type:      globalConfig.ReportType,
			Timeout:   globalConfig.ReportTimeout,
		}
		if err := task.Report(speedData, reportConfig); err != nil {
			fmt.Printf("上报失败: %v\n", err)
		}
	}
}

// startScheduler 启动定时任务调度器
func startScheduler(cronExpr string) {
	c := cron.New(cron.WithSeconds()) // 支持秒级精度的 cron 表达式

	// 添加定时任务
	_, err := c.AddFunc(cronExpr, func() {
		fmt.Printf("\n[%s] 定时任务开始执行...\n", time.Now().Format("2006-01-02 15:04:05"))
		run()
		fmt.Printf("[%s] 定时任务执行完成\n\n", time.Now().Format("2006-01-02 15:04:05"))
	})

	if err != nil {
		fmt.Printf("添加定时任务失败: %v\n", err)
		fmt.Printf("Cron 表达式格式错误，请检查配置\n")
		fmt.Printf("标准格式(5段): 分 时 日 月 周\n")
		fmt.Printf("秒级格式(6段): 秒 分 时 日 月 周\n")
		fmt.Printf("示例:\n")
		fmt.Printf("  每6小时: 0 */6 * * *\n")
		fmt.Printf("  每天0点: 0 0 * * *\n")
		fmt.Printf("  每2小时: 0 */2 * * *\n")
		fmt.Printf("  每30分钟: */30 * * * *\n")
		os.Exit(1)
	}

	// 启动调度器
	c.Start()

	fmt.Printf("定时任务已启动，Cron 表达式: %s\n", cronExpr)
	fmt.Printf("下次执行时间: %s\n", c.Entries()[0].Next.Format("2006-01-02 15:04:05"))
	fmt.Printf("按 Ctrl+C 停止服务\n\n")

	// 立即执行一次
	fmt.Printf("[%s] 首次执行任务...\n", time.Now().Format("2006-01-02 15:04:05"))
	run()
	fmt.Printf("[%s] 首次执行完成\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n收到停止信号，正在关闭服务...")
	c.Stop()
	fmt.Println("服务已停止")
}

func main() {
	run()
}
