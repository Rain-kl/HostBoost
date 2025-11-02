package main

import (
	"cf_opt/config"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"cf_opt/task"
	"cf_opt/utils"
)

var (
	version, versionNew string
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

	// 添加配置文件参数
	flag.StringVar(&configFile, "c", "config.yaml", "配置文件路径")
	flag.BoolVar(&printVersion, "v", false, "打印程序版本")

	// 先解析 -c 和 -v 参数
	flag.Parse()

	// 加载配置文件（先设置默认值）
	config := loadConfigFile(configFile)

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

	// 重新定义命令行参数（用于覆盖配置文件）
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

	flag.Usage = func() { fmt.Print(help) }
	// 再次解析以应用命令行参数覆盖
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
		fmt.Println("检查版本更新中...")
		checkUpdate()
		if versionNew != "" {
			utils.Yellow.Printf("*** 发现新版本 [%s]！请前往 [https://github.com/XIU2/CloudflareSpeedTest] 更新！ ***", versionNew)
		} else {
			utils.Green.Println("当前为最新版本 [" + version + "]！")
		}
		os.Exit(0)
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

func main() {
	task.InitRandSeed() // 置随机数种子

	fmt.Printf("# XIU2/CloudflareSpeedTest %s \n\n", version)

	// 开始延迟测速 + 过滤延迟/丢包
	pingData := task.NewPing().Run().FilterDelay().FilterLossRate()
	print("ok")
	// 开始下载测速
	speedData := task.TestDownloadSpeed(pingData)
	print("ok")
	utils.ExportCsv(speedData) // 输出文件
	speedData.Print()          // 打印结果
	endPrint()                 // 根据情况选择退出方式（针对 Windows）
}

// 根据情况选择退出方式（针对 Windows）
func endPrint() {
	if utils.NoPrintResult() { // 如果不需要打印测速结果，则直接退出
		return
	}
	if runtime.GOOS == "windows" { // 如果是 Windows 系统，则需要按下 回车键 或 Ctrl+C 退出（避免通过双击运行时，测速完毕后直接关闭）
		fmt.Printf("按下 回车键 或 Ctrl+C 退出。")
		fmt.Scanln()
	}
}

// 检查更新
func checkUpdate() {
	timeout := 10 * time.Second
	client := http.Client{Timeout: timeout}
	res, err := client.Get("https://api.xiu2.xyz/ver/cloudflarespeedtest.txt")
	if err != nil {
		return
	}
	// 读取资源数据 body: []byte
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	// 关闭资源流
	defer res.Body.Close()
	if string(body) != version {
		versionNew = string(body)
	}
}
