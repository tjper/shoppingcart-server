package testing

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/tjper/shoppingcart-server/service"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type inject struct {
	Svc     *service.Service
	Closers []func(*testing.T)
}

func newInject(t *testing.T) *inject {
	var i = new(inject)

	dbport, remove := newDb(t)
	i.Closers = append(i.Closers, remove)

	i.Svc = newService(dbport)

	return i
}

func (i inject) Close(t *testing.T) {
	i.Svc.Close()
	for _, closer := range i.Closers {
		closer(t)
	}
}

// newService returns an initialized service.
func newService(dbport string) *service.Service {
	var connStr = fmt.Sprintf(
		"admin:password@tcp(localhost:%s)/shoppingcart-db?tls=false&timeout=30s",
		dbport)

	var v = service.ViperDefaults(viper.New())
	v.Set(service.EnvVarDbConnStr, connStr)

	var svc = service.New(v, service.WithDB(), service.WithZap())
	service.WithRouters(
		svc.CartRoutes,
		svc.ItemRoutes,
	)(svc)
	return svc
}

// newDb creates and starts a shoppingcart-db docker container and returns a
// cleanup function to be called to stop and remove said container as the 2nd
// return value. The first return value is the host's port for the db instance.
func newDb(t *testing.T) (string, func(t *testing.T)) {
	var ctx = context.Background()

	client, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion("1.38"),
	)
	assert.Nil(t, err)

	const image = "tjperr/shoppingcart-db:latest"
	_, err = client.ImagePull(ctx, image, types.ImagePullOptions{})
	assert.Nil(t, err)

	resp, err := client.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			Tty:   true,
			Env: []string{
				"MYSQL_ALLOW_EMPTY_PASSWORD=true",
			},
		},
		&container.HostConfig{
			PublishAllPorts: true,
		},
		nil,
		"")
	assert.Nil(t, err)

	err = client.ContainerStart(
		ctx,
		resp.ID,
		types.ContainerStartOptions{})
	assert.Nil(t, err)

	i, err := client.ContainerInspect(ctx, resp.ID)
	assert.Nil(t, err)

	const containerPort = "3306/tcp"
	var hostPort = i.NetworkSettings.NetworkSettingsBase.Ports[containerPort][0].HostPort
	waitForPort(t, hostPort)

	return hostPort,
		func(t *testing.T) {
			var ctx = context.Background()
			err := client.ContainerStop(ctx, resp.ID, nil)
			assert.Nil(t, err)

			err = client.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
			assert.Nil(t, err)
		}
}

func waitForPort(t *testing.T, port string) {
	for {
		conn, err := net.Dial("tcp", "localhost:"+port)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		conn.Close()
		return
	}
}
