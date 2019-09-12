package cmd

import (
	"fmt"

	"github.com/tjper/shoppingcart-server/service"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve provides access to the shoppingcart resource via http",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running shoppingcart-server...")

		var v = viper.New()
		v.AutomaticEnv()
		v.SetEnvPrefix(service.EnvVarPrefix)

		var svc = service.New(
			service.ViperDefaults(v),
			service.WithDB(),
			service.WithZap(),
		)
		service.WithRouters(
			svc.CartRoutes,
			svc.ItemRoutes,
		)(svc)

		defer svc.Close()

		svc.ListenAndServe()
	},
}
