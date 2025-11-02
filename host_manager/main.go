package main

import (
	"flag"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"hostMgr/internal/config"
	"hostMgr/internal/host"
	"hostMgr/internal/server"
)

var (
	configFlag = flag.String("config", "config.yaml", "path to the configuration file")
)

func main() {
	flag.Parse()

	// 加载配置文件
	cfg, err := config.Load(*configFlag)
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
	log.Printf("Config file: %s", *configFlag)
	if err := router.Run(port); err != nil {
		log.Fatalf("server error: %v", err)
	}
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
