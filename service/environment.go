package service

import "github.com/spf13/viper"

const (
	// EnvVarHttpPort is the key to an env var that specifies the port the cart
	// Service will listen on.
	EnvVarHttpPort = "http_port"

	// EnvVarDbConnStr is the key to an env var that specifies the
	// database connection string.
	EnvVarDbConnStr = "db_conn_str"

	// EnvVarEnvironment is the key to an env var that specifies the service's
	// current environment.
	EnvVarEnvironment = "environment"
)

const (
	prod = "prod"
	dev  = "dev"

	port    = ":8080"
	connStr = "admin:password@tcp(localhost:3306)/shoppingcart-db?tls=false"
)

func ViperDefaults(v *viper.Viper) *viper.Viper {
	v.SetDefault(EnvVarEnvironment, dev)
	v.SetDefault(EnvVarHttpPort, port)
	v.SetDefault(EnvVarDbConnStr, connStr)
	return v
}
