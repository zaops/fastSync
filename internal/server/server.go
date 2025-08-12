package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fastdfs-migration-system/internal/config"
	"fastdfs-migration-system/internal/logger"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config *config.Config
	router *gin.Engine
	server *http.Server
}

func New(cfg *config.Config) *Server {
	// 初始化日志
	if err := logger.Init(&cfg.Logging); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	// 设置Gin模式
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 添加中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	server := &Server{
		config: cfg,
		router: router,
	}

	// 设置路由
	server.setupRoutes()

	return server
}

func (s *Server) Start() error {
	// 创建HTTP服务器
	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port),
		Handler: s.router,
	}

	logger.Infof("Starting server on %s", s.server.Addr)

	// 启动服务器
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
		return err
	}

	logger.Info("Server exited")
	return nil
}

func (s *Server) setupRoutes() {
	// 健康检查
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// API路由组
	api := s.router.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	// 静态文件服务 (后续会添加嵌入的前端资源)
	s.router.Static("/static", "./web/static")

	// 尝试加载HTML模板，如果失败则跳过（测试环境）
	if err := s.loadHTMLTemplates(); err == nil {
		// 默认路由
		s.router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "FastDFS Migration System",
			})
		})
	} else {
		// 测试环境的简单首页
		s.router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"title":  "FastDFS Migration System",
				"status": "running",
			})
		})
	}
}

func (s *Server) loadHTMLTemplates() error {
	defer func() {
		if r := recover(); r != nil {
			// 模板加载失败，忽略错误
		}
	}()

	s.router.LoadHTMLGlob("web/templates/*")
	return nil
}
