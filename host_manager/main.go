package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"awesomeProject/internal/host"
	"awesomeProject/internal/server"
)

var (
	portFlag     = flag.String("port", "", "HTTP port for the server (overrides HOSTBOOST_PORT)")
	dataFileFlag = flag.String("data-file", "", "path to the simulated hosts file (overrides HOSTBOOST_DATA_FILE)")
)

func main() {
	flag.Parse()

	port := resolvePort(*portFlag)
	hostFile := resolveHostFile(*dataFileFlag)

	repo, err := host.NewFileRepository(hostFile)
	if err != nil {
		log.Fatalf("init repository: %v", err)
	}

	svc := host.NewService(repo)
	handler := server.NewHandler(svc)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), buildCorsMiddleware())
	handler.RegisterRoutes(router)

	log.Printf("Host service listening on %s using data file %s", port, hostFile)
	if err := router.Run(port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func resolvePort(flagValue string) string {
	defaultPort := ":15920"

	envPort := os.Getenv("HOSTBOOST_PORT")
	if envPort != "" {
		if envPort[0] != ':' {
			return fmt.Sprintf(":%s", envPort)
		}
		return envPort
	}

	if flagValue != "" {
		if flagValue[0] != ':' {
			return fmt.Sprintf(":%s", flagValue)
		}
		return flagValue
	}

	return defaultPort
}

func resolveHostFile(flagValue string) string {
	defaultFile := "data/hosts.json"

	if file := os.Getenv("HOSTBOOST_DATA_FILE"); file != "" {
		return file
	}

	if flagValue != "" {
		return flagValue
	}

	return defaultFile
}

func buildCorsMiddleware() gin.HandlerFunc {
	cfg := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(cfg)
}
