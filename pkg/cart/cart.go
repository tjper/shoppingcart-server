package cart

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type Execer interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

type QueryRower interface {
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Queryer interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
}

type ExecQueryer interface {
	Execer
	QueryRower
}

// Cart is a user's shopping cart. Typically populated by a set of CartItem
// objects.
type Cart struct {
	Items []CartItem
}

// Get retrieves the specified userId's cart from the Queryer. On success, a
// nil error is returned. On failure, a non nil error is returned.
func (c *Cart) Get(ctx context.Context, db Queryer, userId int64) error {
	var sql = `
  SELECT id, item_id, user_id, count
  FROM cart
  WHERE user_id = $1
  `

	rows, err := db.QueryContext(ctx, sql, userId)
	if err != nil {
		return errors.Wrapf(err, "failed to Cart.Get/Query\tsql=%s\tuserId=%v", sql, userId)
	}
	defer rows.Close()

	var (
		cartItems = make([]CartItem, 0)
		cartItem  = new(CartItem)
	)
	for rows.Next() {
		if err := rows.Scan(
			cartItem.Id,
			cartItem.ItemId,
			cartItem.UserId,
			cartItem.Count,
		); err != nil {
			return errors.Wrapf(err, "failed to Cart.Get/Scan\tsql=%s\tuserId=%v", sql, userId)
		}
		cartItems = append(cartItems, *cartItem)
	}
	if err := rows.Err(); err != nil {
		return errors.Wrapf(err, "failed to Cart.Get/Err\tsql=%s\tuserId=%v", sql, userId)
	}
	c.Items = cartItems

	return nil
}

// CartItem is a item that belong's to a cart.
type CartItem struct {
	Id     int64 `json:"id"`
	ItemId int64 `json:"itemId"`
	UserId int64 `json:"userId"`
	Count  int64 `json:"count"`
}

// FindCartItem retrieves the CartItem with the id passed from the db.
func FindCartItem(ctx context.Context, db QueryRower, id int64) (*CartItem, error) {
	var sql = `
  SELECT id, item_id, user_id, count
  FROM cart
  WHERE id = $1 
  `

	var cartItem = new(CartItem)
	if err := db.QueryRowContext(ctx, sql, id).Scan(
		cartItem.Id,
		cartItem.ItemId,
		cartItem.UserId,
		cartItem.Count,
	); err != nil {
		return nil, errors.Wrapf(err, "failed to FindCartItem\tsql=%s\tid=%v", sql, id)
	}
	return cartItem, nil
}

// Insert adds a cart item to a cart in the db.
func (i *CartItem) Insert(ctx context.Context, db Execer) error {
	var sql = `
  INSERT INTO cart (item_id, user_id, count)
  VALUES ($1, $2, $3)
  `
	var args = []interface{}{i.ItemId, i.UserId, i.Count}
	res, err := db.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrapf(err, "failed to Insert/Exec\tsql=%s\targs=%v", sql, args)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return errors.Wrapf(err, "failed to Insert/LastInsertId\tres=%v", res)
	}
	i.Id = id
	return nil
}

// DeleteCartItem deletes the cart item associated with the id passed in the
// db.
func DeleteCartItem(ctx context.Context, db ExecQueryer, id int64) error {
	if _, err := FindCartItem(ctx, db, id); err != nil {
		return errors.Wrap(err, "failed to DeleteCartItem")
	}

	var sql = `
  DELETE FROM cart
  WHERE id = $1
  `
	if _, err := db.ExecContext(ctx, sql, id); err != nil {
		return errors.Wrapf(err, "failed to DeleteCartItem/Exec\tsql=%s\tid=%v", sql, id)
	}
	return nil
}

// Update updates the cart item in the db based on the CartItem object.
// WARNING: All fields of CartItem must be non-empty.
func (i CartItem) Update(ctx context.Context, db ExecQueryer) error {
	if _, err := FindCartItem(ctx, db, i.Id); err != nil {
		return errors.Wrap(err, "failed to Update")
	}

	var sql = `
  UPDATE cart
  SET item_id = $1 
      user_id = $2
      count = $3
  WHERE id = $4
  `
	var args = []interface{}{i.ItemId, i.UserId, i.Count, i.Id}
	if _, err := db.ExecContext(ctx, sql, args...); err != nil {
		return errors.Wrapf(err, "failed to Update/Exec\tsql=%s\ti=%v", sql, i)
	}
	return nil
}
