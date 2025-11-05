package main

import (
	"flag"
	"fmt"
	"hostMgr/config"
	"hostMgr/internal/extSvc"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"hostMgr/internal/host"
	"hostMgr/internal/opt"
	"hostMgr/internal/server"
	"hostMgr/internal/tool"
)

const (
	version = "1.0.0"
)

var (
	configFlag   = flag.String("config", "config.yaml", "path to the configuration file")
	configShort  = flag.String("c", "config.yaml", "path to the configuration file (shorthand)")
	helpFlag     = flag.Bool("help", false, "show help message")
	helpShort    = flag.Bool("h", false, "show help message (shorthand)")
	versionFlag  = flag.Bool("version", false, "show version information")
	versionShort = flag.Bool("v", false, "show version information (shorthand)")
)

func main() {
	flag.Parse()

	// 处理 help 参数
	if *helpFlag || *helpShort {
		printHelp()
		os.Exit(0)
	}

	// 处理 version 参数
	if *versionFlag || *versionShort {
		fmt.Printf("Host Manager version %s\n", version)
		os.Exit(0)
	}

	// 优先使用短参数 -c
	configPath := *configFlag
	if *configShort != "config.yaml" {
		configPath = *configShort
	}

	// 加载配置文件
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// 初始化 host repository
	repo, err := host.NewFileRepository(cfg.Data.HostFile)
	if err != nil {
		log.Fatalf("init repository: %v", err)
	}

	// 初始化 host service
	hostSvc := host.NewService(repo)
	extSvc.HostService = hostSvc

	// 初始化 opt repository
	optRepo, err := opt.NewRepository(cfg.Data.OptFile)
	if err != nil {
		log.Fatalf("init opt repository: %v", err)
	}

	// 初始化 opt service
	optSvc := opt.NewService(optRepo)
	extSvc.OptService = optSvc

	// 初始化 tool service
	toolSvc := tool.NewToolService()

	handler := server.NewHandler(hostSvc, optSvc, toolSvc)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), buildCorsMiddleware(cfg))
	handler.RegisterRoutes(router)

	port := cfg.Server.NormalizePort()
	log.Printf("Host service listening on %s using data file %s", port, cfg.Data.HostFile)
	log.Printf("Opt service using data file %s", cfg.Data.OptFile)
	log.Printf("Tool service initialized with DNS resolver and IP geolocation service")
	log.Printf("Config file: %s", configPath)
	if err := router.Run(port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func printHelp() {
	fmt.Println("Host Manager - A host configuration management service")
	fmt.Printf("Version: %s\n\n", version)
	fmt.Println("Usage:")
	fmt.Println("  host_manager [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -c, --config <file>    Path to the configuration file (default: config.yaml)")
	fmt.Println("  -h, --help             Show this help message")
	fmt.Println("  -v, --version          Show version information")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  host_manager")
	fmt.Println("  host_manager -c /path/to/config.yaml")
	fmt.Println("  host_manager --config /path/to/config.yaml")
	fmt.Println("  host_manager --version")
}

func buildCorsMiddleware(cfg *config.Config) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		ExposeHeaders:    cfg.CORS.ExposeHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           cfg.CORS.GetMaxAge(),
	}

	return cors.New(corsConfig)
}
