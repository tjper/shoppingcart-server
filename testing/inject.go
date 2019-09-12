package testing

import (
	"github.com/tjper/shoppingcart-server/service"

	"github.com/spf13/viper"
)

type inject struct {
	Svc *service.Service
}

func newInject() *inject {
	var svc = service.New(
		service.ViperDefaults(viper.New()),
		service.WithDB(),
		service.WithZap(),
	)
	service.WithRouters(
		svc.CartRoutes,
		svc.ItemRoutes,
	)(svc)

	return &inject{
		Svc: svc,
	}
}

func (i inject) Close() {
	i.Svc.Close()
}
