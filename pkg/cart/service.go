package cart

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

const (
	// EnvHttpPort is the key to an env var that specifies the port the cart
	// Service will listen on.
	EnvHttpPort = "http_port"
)

// Service defines all service dependencies.
type Service struct {
	Viper *viper.Viper
	DB    *sql.DB
}

// NewService initializes a new cart Service via option functions.
func NewService(port string, options ...ServiceOption) *Service {
	var svc = new(Service)

	for _, option := range options {
		option(svc)
	}
	return svc
}

// ServiceOption modifies a Service object. Typically used with Service
// initialization via NewService().
type ServiceOption func(*Service)

// WithViper returns a ServiceOption that sets the Service.Viper field to the
// specified viper instance.
func WithViper(v *viper.Viper) ServiceOption {
	return func(svc *Service) {
		svc.Viper = v
	}
}

// WithDb returns a ServiceOption that sets the Service.Db field to the
// specified Db instance.
func WithDb(db *sql.DB) ServiceOption {
	return func(svc *Service) {
		svc.DB = db
	}
}

func (svc *Service) ListenAndServe() {
	http.ListenAndServe(svc.Viper.GetString(EnvHttpPort), svc.Routes(chi.NewRouter()))
}
