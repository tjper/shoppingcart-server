package cmd

import (
	"fmt"

	"github.com/tjper/shoppingcart-server/pkg/cart"

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

		var cartService = cart.NewService(
			cart.ViperDefaults(viper.New()),
			cart.WithDB(),
			cart.WithZap(),
		)
		defer cartService.Close()

		cartService.ListenAndServe()
	},
}
