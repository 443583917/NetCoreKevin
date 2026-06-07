package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/kevin-ai/go-kevin/internal/application/command"
	"github.com/kevin-ai/go-kevin/internal/application/query"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/eventbus"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/persistence"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/persistence/gorm/repository"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/persistence/migration"
	"github.com/kevin-ai/go-kevin/internal/interfaces/http/handler"
	"github.com/kevin-ai/go-kevin/internal/interfaces/http/middleware"
	"github.com/kevin-ai/go-kevin/internal/interfaces/http/router"
	"github.com/kevin-ai/go-kevin/pkg/config"
	"github.com/kevin-ai/go-kevin/pkg/logger"
)

func main() {
	// 1. Load config
	configPath := "configs/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}
	if err := config.Init(configPath); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg := config.Get()

	// 2. Initialize logger
	logger.Init()
	logger.Info("应用启动中...")

	// 3. Initialize database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := persistence.InitMySQL(dsn)
	if err != nil {
		logger.Fatal("数据库连接失败", zap.Error(err))
	}

	// 4. Run migrations
	migrator := migration.NewMigrator(db)
	if err := migrator.Run(); err != nil {
		logger.Fatal("数据库迁移失败", zap.Error(err))
	}

	// 5. Initialize event bus
	bus := eventbus.NewInMemoryEventBus()

	// 6. Initialize repositories
	userRepo := repository.NewUserRepositoryImpl(db)

	// 7. Initialize command and query handlers
	createUserHandler := command.NewCreateUserHandler(userRepo, bus)
	getUserHandler := query.NewGetUserHandler(userRepo)
	listUsersHandler := query.NewListUsersHandler(userRepo)
	updateUserHandler := command.NewUpdateUserHandler(userRepo)
	deleteUserHandler := command.NewDeleteUserHandler(userRepo)

	// 8. Initialize HTTP handlers
	authHandler := handler.NewAuthHandler(userRepo, cfg.JWT.Secret, cfg.JWT.ExpireHour)
	userHandler := handler.NewUserHandler(createUserHandler, getUserHandler, listUsersHandler, updateUserHandler, deleteUserHandler)

	// 9. Initialize HTTP server
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.CORS())

	// 10. Register routes
	router.RegisterRoutes(r, userHandler, authHandler)

	// 11. Start server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("服务器启动在端口 %d", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 12. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("服务器关闭失败", zap.Error(err))
	}

	logger.Info("服务器已关闭")
}
