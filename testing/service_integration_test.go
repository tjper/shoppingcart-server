// +build integration

package testing

import (
	"context"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
	"github.com/tjper/shoppingcart-server/service/cart"
	testutil "github.com/tjper/testing"
)

var golden = flag.Bool("golden", false, "overwrite the existing golden files")

func TestGetItems(t *testing.T) {
	t.Parallel()
	i := newInject(t)
	defer i.Close(t)

	t.Run("GET items", func(t *testing.T) {
		var ts = httptest.NewServer(i.Svc.GetItemsHandler())
		defer ts.Close()

		resp, err := http.Get(ts.URL)
		require.Nil(t, err)
		defer resp.Body.Close()

		actual, err := ioutil.ReadAll(resp.Body)
		require.Nil(t, err)

		if *golden {
			testutil.GoldenUpdate(t, actual)
		}
		expected := testutil.GoldenGet(t)

		require.Equal(t, expected, actual)
	})
}

func TestPostCartItem(t *testing.T) {
	t.Parallel()
	var i = newInject(t)
	defer i.Close(t)

	var ts = httptest.NewServer(i.Svc.PostCartItemHandler())
	defer ts.Close()

	tests := []struct {
		Name        string
		RequestBody io.Reader
	}{
		{
			Name:        "Baseline",
			RequestBody: strings.NewReader(`{"itemId": 1, "userId": 1, "count": 1}`),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			resp, err := http.Post(ts.URL, "application/json", test.RequestBody)
			require.Nil(t, err)
			defer resp.Body.Close()

			actual, err := ioutil.ReadAll(resp.Body)
			require.Nil(t, err)

			if *golden {
				testutil.GoldenUpdate(t, actual)
			}
			expected := testutil.GoldenGet(t)

			require.Equal(t, expected, actual)
		})
	}
}

func TestGetCart(t *testing.T) {
	t.Parallel()
	var i = newInject(t)
	defer i.Close(t)

	tests := []struct {
		Name   string
		Rels   []cart.UserCartItemRel
		UserId string
	}{
		{
			Name: "Baseline",
			Rels: []cart.UserCartItemRel{
				cart.UserCartItemRel{ItemId: 1, UserId: 1, Count: 1},
			},
			UserId: "1",
		},
		{
			Name: "Two items",
			Rels: []cart.UserCartItemRel{
				cart.UserCartItemRel{ItemId: 1, UserId: 1, Count: 1},
				cart.UserCartItemRel{ItemId: 2, UserId: 1, Count: 1},
			},
			UserId: "1",
		},
		{
			Name: "Three items, counts > 1",
			Rels: []cart.UserCartItemRel{
				cart.UserCartItemRel{ItemId: 1, UserId: 1, Count: 2},
				cart.UserCartItemRel{ItemId: 2, UserId: 1, Count: 1},
				cart.UserCartItemRel{ItemId: 3, UserId: 1, Count: 1},
			},
			UserId: "1",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			for _, rel := range test.Rels {
				_, err := cart.CreateUserCartItemRel(
					context.Background(),
					i.Svc.DB,
					rel)
				require.Nil(t, err)
			}

			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/cart/"+test.UserId, nil)
			)

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("userId", test.UserId)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

			i.Svc.GetCartHandler()(w, r)

			actual, err := ioutil.ReadAll(w.Result().Body)
			require.Nil(t, err)

			if *golden {
				testutil.GoldenUpdate(t, actual)
			}
			expected := testutil.GoldenGet(t)

			require.Equal(t, expected, actual)
		})
	}
}

func TestPutCartItem(t *testing.T) {
	t.Parallel()
	var i = newInject(t)
	defer i.Close(t)

	tests := []struct {
		Name           string
		Rel            cart.UserCartItemRel
		PutRequestBody io.Reader
	}{
		{
			Name:           "Baseline",
			Rel:            cart.UserCartItemRel{ItemId: 1, UserId: 1, Count: 1},
			PutRequestBody: strings.NewReader(`{"itemId": 1, "userId": 1, "count": 6}`),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			id, err := cart.CreateUserCartItemRel(
				context.Background(),
				i.Svc.DB,
				test.Rel)
			require.Nil(t, err)

			var (
				idStr = strconv.Itoa(id)
				w     = httptest.NewRecorder()
				r     = httptest.NewRequest(http.MethodPut, "/cart/item/"+idStr, test.PutRequestBody)
			)

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", idStr)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

			i.Svc.PutCartItemHandler()(w, r)

			actual, err := ioutil.ReadAll(w.Result().Body)
			require.Nil(t, err)

			if *golden {
				testutil.GoldenUpdate(t, actual)
			}
			expected := testutil.GoldenGet(t)

			require.Equal(t, expected, actual)
		})
	}
}

func TestDeleteItem(t *testing.T) {
	t.Parallel()
	var i = newInject(t)
	defer i.Close(t)

	tests := []struct {
		Name string
		Rel  cart.UserCartItemRel
	}{
		{
			Name: "Baseline",
			Rel:  cart.UserCartItemRel{ItemId: 1, UserId: 1, Count: 1},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			id, err := cart.CreateUserCartItemRel(context.Background(), i.Svc.DB, test.Rel)
			require.Nil(t, err)

			var (
				idStr = strconv.Itoa(id)
				w     = httptest.NewRecorder()
				r     = httptest.NewRequest(http.MethodDelete, "/cart/item/"+idStr, nil)
			)

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", idStr)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

			i.Svc.DeleteCartItemHandler()(w, r)
			require.Equal(t, http.StatusOK, w.Result().StatusCode)
		})
	}
}
