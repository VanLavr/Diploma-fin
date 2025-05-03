package rest

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/VanLavr/Diploma-fin/utils/config"
)

type Server struct {
	gin            *gin.Engine
	cfg            *config.Config
	studentHandler *StudentHandler
	teacherHandler *TeacherHandler
	examHandler    *ExamHandler
	groupHandler   *GroupHandler
	authHandler    *AuthHandler
	fileHandler    *FileHandler
}

func NewServer(
	cfg *config.Config,
	studentHandler *StudentHandler,
	teacherHandler *TeacherHandler,
	authHandler *AuthHandler,
	examHandler *ExamHandler,
	groupHandler *GroupHandler,
	fileHandler *FileHandler,
) *Server {
	return &Server{
		cfg:            cfg,
		studentHandler: studentHandler,
		teacherHandler: teacherHandler,
		examHandler:    examHandler,
		groupHandler:   groupHandler,
		authHandler:    authHandler,
		fileHandler:    fileHandler,
		gin:            gin.Default(),
	}
}

func (s *Server) Start(c context.Context) error {
	s.gin.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"ResponseType",
			"accept",
			"origin",
			"Cache-Control",
			"X-Request-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	s.setV1Routes(s.gin.Group("/api"))

	server := &http.Server{
		Addr:           ":" + s.cfg.Port,
		Handler:        s.gin,
		ReadTimeout:    time.Second * 20,
		WriteTimeout:   time.Second * 20,
		MaxHeaderBytes: 1 << 20,
	}

	ctx, stop := signal.NotifyContext(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	ctxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxt); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Server) setV1Routes(group *gin.RouterGroup) {
	var v1 *gin.RouterGroup
	var v1Auth *gin.RouterGroup

	group.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Health-check": "ok",
		})
	})

	if !s.cfg.WithJWTAuth {
		v1 = group.Group("/v1")
	} else {
		v1Auth = group.Group("/auth")
		v1 = group.Group("/v1", s.authHandler.auth.ValidateAccessToken())
	}

	s.authHandler.RegisterRoutes(v1Auth)
	s.studentHandler.RegisterRoutes(v1)
	s.examHandler.RegisterRoutes((v1))
	s.groupHandler.RegisterRoutes((v1))
	s.teacherHandler.RegisterRoutes(v1)
	s.fileHandler.RegisterRoutes(v1)
}
