package service

import "github.com/spf13/viper"

const (
	// EnvVarHttpPort is the key to an env var that specifies the port the cart
	// Service will listen on.
	EnvVarHttpPort = "http_port"

	// EnvVarEnvironment is the key to an env var that specifies the service's
	// current environment.
	EnvVarEnvironment = "environment"

	// EnvVarDbConnStr is the key to an env var that specifies the
	// database connection string.
	EnvVarDbConnStr = "db_conn_str"

	// EnvVarDbMaxOpenConns is the key to an env var that specifies the
	// max number open connections in the db connection pool.
	EnvVarDbMaxOpenConns = "db_max_open_conns"

	// EnvVarDbMaxIdleConns is the key to an env var that specifies the
	// possible max number of idle connections in the db connection
	// pool.
	EnvVarDbMaxIdleConns = "db_max_idle_conns"
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
