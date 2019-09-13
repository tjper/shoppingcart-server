package cart

import (
	"context"
	"database/sql"

	"github.com/tjper/shoppingcart-server/service/item"

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

// CartItems retrieves the specified userId's cart items from the db. On
// success, the cart items and a nil error is returned. On failure, a non nil
// error is returned.
func CartItems(ctx context.Context, db Queryer, userId int) ([]CartItem, error) {
	var sql = `
  SELECT 
    cart.id,
    item.id,
    item.name,
    item.price,
    cart.count
  FROM 
    cart
  JOIN
    item
    ON item.id = cart.item_id
  WHERE cart.user_id = ?
  `

	rows, err := db.QueryContext(ctx, sql, userId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to CartItems/QueryContext\tsql=%s\tuserId=%v", sql, userId)
	}
	defer rows.Close()

	var (
		cartItems = make([]CartItem, 0)
		cartItem  CartItem
	)
	for rows.Next() {
		if err := rows.Scan(
			&cartItem.Id,
			&cartItem.Item.Id,
			&cartItem.Item.Name,
			&cartItem.Item.Price,
			&cartItem.Count,
		); err != nil {
			return nil, errors.Wrapf(err, "failed to CartItems/Scan\tsql=%s\tuserId=%v", sql, userId)
		}
		cartItems = append(cartItems, cartItem)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to CartItems/Err\tsql=%s\tuserId=%v", sql, userId)
	}
	return cartItems, nil
}

// CartItem is a item that belong's to a cart.
type CartItem struct {
	Id    int       `json:"id"`
	Count int       `json:"count"`
	Item  item.Item `json:"item"`
}

// FindCartItem retrieves the CartItem with the id passed from the db.
func FindCartItem(ctx context.Context, db QueryRower, id int) (*CartItem, error) {
	var sql = `
  SELECT 
    cart.id,
    item.id,
    item.name,
    item.price,
    cart.count
  FROM 
    cart
  JOIN
    item
    ON item.id = cart.item_id
  WHERE cart.id = ?
  `

	var cartItem CartItem
	if err := db.QueryRowContext(ctx, sql, id).Scan(
		&cartItem.Id,
		&cartItem.Item.Id,
		&cartItem.Item.Name,
		&cartItem.Item.Price,
		&cartItem.Count,
	); err != nil {
		return nil, errors.Wrapf(err, "failed to FindCartItem\tsql=%s\tid=%v", sql, id)
	}
	return &cartItem, nil
}

// UserCartItemRel is represents a many-to-many relationship between the item and
// user resource.
type UserCartItemRel struct {
	Id     int
	ItemId int
	UserId int
	Count  int
}

// CreateUserCartItemRel adds a cart item to a cart in the db based on the Create.
func CreateUserCartItemRel(ctx context.Context, db Execer, rel UserCartItemRel) (int, error) {
	var sql = `
  INSERT INTO cart (item_id, user_id, count)
  VALUES (?, ?, ?)
  `
	var args = []interface{}{rel.ItemId, rel.UserId, rel.Count}
	res, err := db.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to Insert/ExecContext\tsql=%s\targs=%v", sql, args)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Wrapf(err, "failed to Insert/LastInsertId\tres=%v", res)
	}
	return int(id), nil
}

// DeleteCartItem deletes the cart item associated with the id passed in the
// db.
func DeleteCartItem(ctx context.Context, db ExecQueryer, id int) error {
	if _, err := FindCartItem(ctx, db, id); err != nil {
		return errors.Wrap(err, "failed to DeleteCartItem")
	}

	var sql = `
  DELETE FROM cart
  WHERE id = ?
  `
	if _, err := db.ExecContext(ctx, sql, id); err != nil {
		return errors.Wrapf(err, "failed to DeleteCartItem/ExecContext\tsql=%s\tid=%v", sql, id)
	}
	return nil
}

// Update updates the cart item in the db based on the CartItem object.
// WARNING: All fields of CartItem must be non-empty.
func UpdateUserCartItemRel(ctx context.Context, db ExecQueryer, id int, rel UserCartItemRel) error {
	if _, err := FindCartItem(ctx, db, id); err != nil {
		return errors.Wrap(err, "failed to Update")
	}

	var sql = `
  UPDATE cart
  SET item_id = ?,
      user_id = ?,
      count = ?
  WHERE id = ?
  `
	var args = []interface{}{rel.ItemId, rel.UserId, rel.Count, id}
	if _, err := db.ExecContext(ctx, sql, args...); err != nil {
		return errors.Wrapf(err, "failed to Update/ExecContext\tsql=%s\targs=%v", sql, args)
	}
	return nil
}

// UserCartItemRelExists checks to see if a cart item exists for the userId and
// itemId pair, and if it does the cart item id is returned. If a cart item
// does not exist a non-valid (0) id it returned.
func UserCartItemRelExists(ctx context.Context, db QueryRower, userId, itemId int) (int, error) {
	var SQL = `
    SELECT cart.id
    FROM cart
    WHERE cart.item_id = ?
          AND cart.user_id = ?
  `
	var args = []interface{}{itemId, userId}
	var id int
	err := db.QueryRowContext(ctx, SQL, args...).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, errors.Wrapf(err, "failed to UserCartItemRelExists\tSQL=%s\targs=%v", SQL, args)
	}
	return id, nil
}
