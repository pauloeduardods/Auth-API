package server

import (
	"auth-api-cognito/internal/utils"
	"auth-api-cognito/pkg/jwtToken"
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Server struct {
	log             *zap.Logger
	gin             *gin.Engine
	cognitoClient   *cognito.Client
	utils           *utils.Utils
	server          *http.Server
	jwtToken        *jwtToken.JwtToken
	host            string
	port            int
	cognitoClientId string
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
		Credentials: opts.AwsConfig.Credentials,
		Region:      opts.AwsConfig.Region,
	})

	u := utils.NewUtils(validator.New())

	j := jwtToken.NewAuth(&jwtToken.Config{
		CognitoRegion:     opts.CognitoRegion,
		CognitoUserPoolID: opts.CognitoUserPoolID,
		Log:               opts.Log,
	})

	return &Server{
		log:             opts.Log,
		gin:             g,
		cognitoClient:   c,
		host:            opts.Host,
		port:            opts.Port,
		jwtToken:        j,
		utils:           u,
		cognitoClientId: opts.CognitoClientId,
	}
}

func (s *Server) Start() error {
	s.log.Info("Starting server", zap.Int("port", s.port))
	s.SetupCors()
	s.SetupMiddlewares()
	s.SetupApi()

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
