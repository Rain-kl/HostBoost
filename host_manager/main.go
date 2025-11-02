package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"hostMgr/internal/config"
	"hostMgr/internal/host"
	"hostMgr/internal/server"
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

	// 初始化 repository
	repo, err := host.NewFileRepository(cfg.Data.HostFile)
	if err != nil {
		log.Fatalf("init repository: %v", err)
	}

	svc := host.NewService(repo)
	handler := server.NewHandler(svc)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), buildCorsMiddleware(cfg))
	handler.RegisterRoutes(router)

	port := cfg.Server.NormalizePort()
	log.Printf("Host service listening on %s using data file %s", port, cfg.Data.HostFile)
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
