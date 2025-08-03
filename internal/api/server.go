package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bxrne/goforge/internal/api/handlers"
	"github.com/bxrne/goforge/internal/api/middleware"
	"github.com/bxrne/goforge/internal/api/routes"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
	port   string
}

type Config struct {
	Port        string
	Environment string
}

func NewServer(config Config) *Server {
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Prometheus())

	server := &Server{
		router: router,
		port:   config.Port,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", handlers.HealthCheck)

	// Metrics endpoint for Prometheus
	s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger documentation
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := s.router.Group("/api/v1")
	routes.SetupRoutes(api)
}

func (s *Server) Start() error {
	log.Printf("Starting server on port %s", s.port)
	return s.router.Run(":" + s.port)
}

func (s *Server) StartWithContext(ctx context.Context) error {
	srv := &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", s.port)

	<-ctx.Done()

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited")
	return nil
}
