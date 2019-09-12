package service

import "github.com/spf13/viper"

const (
	EnvVarPrefix = "CART"
	// EnvVarPrefix is the env var prefix for environments specific to this
	// service.

	// EnvVarHttpPort is the key to an env var that specifies the port the cart
	// Service will listen on.
	EnvVarHttpPort = "HTTP_PORT"

	// EnvVarEnvironment is the key to an env var that specifies the service's
	// current environment.
	EnvVarEnvironment = "ENVIRONMENT"

	// EnvVarDbConnStr is the key to an env var that specifies the
	// database connection string.
	EnvVarDbConnStr = "DB_CONN_STR"

	// EnvVarDbMaxOpenConns is the key to an env var that specifies the
	// max number open connections in the db connection pool.
	EnvVarDbMaxOpenConns = "DB_MAX_OPEN_CONNS"

	// EnvVarDbMaxIdleConns is the key to an env var that specifies the
	// possible max number of idle connections in the db connection
	// pool.
	EnvVarDbMaxIdleConns = "DB_MAX_IDLE_CONNS"
)

const (
	prod = "prod"
	dev  = "dev"

	port    = ":8080"
	connStr = "admin:password@tcp(localhost:3306)/shoppingcart-db?tls=false&timeout=30s"
)

func ViperDefaults(v *viper.Viper) *viper.Viper {
	v.SetDefault(EnvVarEnvironment, dev)
	v.SetDefault(EnvVarHttpPort, port)
	v.SetDefault(EnvVarDbConnStr, connStr)
	v.SetDefault(EnvVarDbMaxOpenConns, 8)
	v.SetDefault(EnvVarDbMaxIdleConns, 0)
	return v
}
