package main

import (
	"auth-api-cognito/cmd/server"
	aConfig "auth-api-cognito/config"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/smithy-go/logging"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var release string

func main() {
	os.Exit(start())
}

func start() int {
	appConfig, err := aConfig.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return 1
	}

	log, err := createLogger(appConfig.AppEnv)
	if err != nil {
		fmt.Println("Error setting up the logger:", err)
		return 1
	}
	log = log.With(zap.String("release", release))
	defer func() {
		_ = log.Sync()
	}()

	awsConfig, err := config.LoadDefaultConfig(context.Background(),
		config.WithLogger(createAWSLogAdapter(log)),
	)
	if err != nil {
		log.Info("Error creating AWS config", zap.Error(err))
		return 1
	}

	s := server.New(server.Options{
		Log:               log,
		AwsConfig:         awsConfig,
		Host:              appConfig.Host,
		Port:              appConfig.Port,
		CognitoClientId:   appConfig.CognitoClientId,
		CognitoRegion:     appConfig.CognitoRegion,
		CognitoUserPoolID: appConfig.CognitoUserPoolID,
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
			log.Info("Error starting server", zap.Error(err))
			return err
		}
		return nil
	})

	<-ctx.Done()

	eg.Go(func() error {
		if err := s.Stop(); err != nil {
			log.Info("Error stopping server", zap.Error(err))
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return 1
	}
	return 0
}

func createLogger(env string) (*zap.Logger, error) {
	switch env {
	case "production":
		return zap.NewProduction()
	case "development":
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}
}

func createAWSLogAdapter(log *zap.Logger) logging.LoggerFunc {
	return func(classification logging.Classification, format string, v ...interface{}) {
		switch classification {
		case logging.Debug:
			log.Sugar().Debugf(format, v...)
		case logging.Warn:
			log.Sugar().Warnf(format, v...)
		}
	}
}
