// backend/cmd/api/main.go
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fynn404/gin-demo/backend/internal/common/config"
	"github.com/fynn404/gin-demo/backend/internal/handler"
	"github.com/fynn404/gin-demo/backend/internal/middleware"
	"github.com/fynn404/gin-demo/backend/internal/repository"
	v1 "github.com/fynn404/gin-demo/backend/internal/route/v1"
	"github.com/gin-gonic/gin"
)

// setupConfig 初始化配置
func setupConfig() error {
	if err := config.Init("/Users/fu.xie/personal/gin-demo/backend/config/config.ini"); err != nil {
		return err
	}
	return nil
}

// setupDatabase 初始化数据库连接
func setupDatabase() {
	config.InitMySQL()
	config.InitRedis()
}

// setupGin 配置 Gin 框架
func setupGin() *gin.Engine {
	// 设置 gin 模式
	if config.GlobalConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 添加全局中间件
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandlingMiddleware())

	return r
}

// setupHandlers 初始化处理器和路由
func setupHandlers(r *gin.Engine) error {
	db := config.GetDB()
	if db == nil {
		return errors.New("database connection not initialized")
	}

	// 初始化 repositories
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)

	// 创建 handler 配置
	handlerConfig := &handler.HandlerConfig{
		UserRepo: userRepo,
		TodoRepo: todoRepo,
	}

	// 创建处理器
	h := handler.NewHandler(handlerConfig)

	// 设置 API 路由
	v1.SetupRoutes(r, h)

	return nil
}

// gracefulShutdown 处理服务优雅关闭
func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// 关闭数据库连接
	config.CloseDB()
	config.CloseRedis()

	log.Println("Server exiting")
}

func main() {
	// 1. 初始化配置
	if err := setupConfig(); err != nil {
		log.Fatalf("Failed to setup config: %v", err)
	}
	log.Println("Config initialized successfully")

	// 2. 初始化数据库连接
	setupDatabase()
	log.Println("Database connections established successfully")

	// 3. 设置 Gin 框架
	r := setupGin()
	log.Println("Gin framework initialized successfully")

	// 4. 设置处理器和路由
	if err := setupHandlers(r); err != nil {
		log.Fatalf("Failed to setup handlers: %v", err)
	}
	log.Println("Handlers and routes setup successfully")

	// 5. 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    config.GlobalConfig.Server.Port,
		Handler: r,
	}

	// 6. 启动服务器
	go func() {
		log.Printf("Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 7. 处理优雅关闭
	gracefulShutdown(srv)
}
