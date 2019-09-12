package cart

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	// Info is a typically used to wrap info logs.
	Info = "[cart INFO]"

	// Error is typically used to warp error logs and/or error types.
	Error = "[cart ERROR]"
)

// Service defines all service dependencies.
type Service struct {
	Viper *viper.Viper
	DB    *sql.DB
	Zap   *zap.Logger
}

// NewService initializes a new cart Service via option functions.
func NewService(v *viper.Viper, options ...ServiceOption) *Service {
	var svc = new(Service)

	for _, option := range options {
		option(svc)
	}
	return svc
}

// ServiceOption modifies a Service object. Typically used with Service
// initialization via NewService().
type ServiceOption func(*Service)

// WithDb returns a ServiceOption that initializes the Service.Db field.
func WithDB() ServiceOption {
	return func(svc *Service) {
		const driver = "mysql"
		db, err := sql.Open(driver, svc.Viper.GetString(EnvVarDbConnStr))
		if err != nil {
			panic(err)
		}
		for {
			log.Println("Attempting to connect to DB")
			err := db.Ping()
			if err == nil {
				log.Println("Connected to DB")
				svc.DB = db
				return
			}
			log.Printf("Failed to connect to DB\terr = %s\n", err)
			log.Println("sleeping...")
			time.Sleep(2 * time.Second)
		}
	}
}

// WithZap returns a ServiceOption that initializes the Service.Zap field.
func WithZap() ServiceOption {
	return func(svc *Service) {
		var (
			logger *zap.Logger
			err    error
		)
		switch env := svc.Viper.GetString(EnvVarEnvironment); env {
		case prod:
			logger, err = zap.NewProduction()
		case dev:
			logger, err = zap.NewDevelopment()
		default:
			panic("switch does not handle environment \"" + env + "\"")
		}
		if err != nil {
			panic(err)
		}
		svc.Zap = logger
	}
}

// ListenAndServe opens a set of HTTP endpoints as specified by svc.Routes() on
// the port specified in viper.
func (svc *Service) ListenAndServe() {
	var (
		srv             http.Server
		idleConnsClosed = make(chan struct{})
	)
	go func() {
		var sigint = make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			svc.Zap.Sugar().Error(errors.Wrap(err, "HTTP server shutdown"))
		}
		close(idleConnsClosed)
	}()

	if err := http.ListenAndServe(
		svc.Viper.GetString(EnvVarHttpPort),
		svc.Routes(chi.NewRouter()),
	); err != http.ErrServerClosed {
		svc.Zap.Sugar().Fatal(errors.Wrap(err, "HTTP server ListenAndServe"))
	}
	<-idleConnsClosed
}

// Close executes necessary cleanup and closing procedures for the Service.
func (svc *Service) Close() {
	svc.Zap.Sync()
}

// Error writes a status code and an optional message to the client. If an
// internal Server error has occurred, the error is logged.
func (svc Service) Error(w http.ResponseWriter, err error, code int, message ...string) {
	if code == http.StatusInternalServerError {
		svc.Zap.Sugar().Error(err)
	}
	http.Error(w, strings.Join(message, "\n"), code)
}
