package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	cognito "auth-api-cognito/internal/auth"
	"auth-api-cognito/internal/auth/jwt"
	validatorUtil "auth-api-cognito/internal/utils/validator"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Server struct {
	log          *zap.Logger
	gin          *gin.Engine
	publicRoute  *gin.RouterGroup
	privateRoute *gin.RouterGroup
	cognito      *cognito.Cognito
	validator    *validatorUtil.Validator
	server       *http.Server
	jwtVerify    *jwt.Auth
	host         string
	port         int
}

type Options struct {
	Log               *zap.Logger
	AwsConfig         aws.Config
	CognitoClientId   string
	CognitoRegion     string
	CognitoUserPoolID string
	Port              int
	Host              string
}

func New(opts Options) *Server {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	g := gin.Default()

	c := cognito.New(cognito.Options{
		AWSConfig: opts.AwsConfig,
		ClientId:  opts.CognitoClientId,
	})

	v := validatorUtil.New(validatorUtil.Options{
		Validate: validator.New(),
	})

	j := jwt.NewAuth(&jwt.Config{
		CognitoRegion:     opts.CognitoRegion,
		CognitoUserPoolID: opts.CognitoUserPoolID,
		Log:               opts.Log,
	})

	return &Server{
		log:       opts.Log,
		gin:       g,
		cognito:   c,
		validator: v,
		host:      opts.Host,
		port:      opts.Port,
		jwtVerify: j,
	}
}

func (s *Server) Start() error {
	s.log.Info("Starting server", zap.Int("port", s.port))
	s.SetupCors()
	s.SetupMiddlewareAndRouteGroup()
	s.SetupRoutes()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		s.log.Info("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {
			s.log.Error("Server Shutdown:", zap.Error(err))
		}
		s.log.Info("Server exiting")
	}()

	s.server = &http.Server{
		Addr:    s.host + ":" + strconv.Itoa(s.port),
		Handler: s.gin,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	if s.server != nil {
		s.log.Info("Stopping server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}
