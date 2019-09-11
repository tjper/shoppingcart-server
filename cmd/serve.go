package cmd

import (
	"fmt"

	"github.com/tjper/shoppingcart-service/pkg/cart"

	"github.com/spf13/cobra"
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
			cart.WithViper(v),
			cart.WithDB(db),
		)
		cartService.ListenAndServe()

	},
}
