package item

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type Queryer interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
}

// Item is an item that may be purchased.
type Item struct {
	Id          int64
	Name        string
	Description string
	Price       float64
}

// Items retrieves all items from the db.
func Items(ctx context.Context, db Queryer) ([]Item, error) {
	var sql = `
    SELECT 
      id,
      name,
      description,
      price
    FROM item
  `
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to Items/QueryContext\tsql=%s", sql)
	}
	defer rows.Close()

	var (
		items = make([]Item, 0)
		item  Item
	)
	for rows.Next() {
		if err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
			&item.Price,
		); err != nil {
			return nil, errors.Wrapf(err, "failed to Items/Scan\tsql=%s", sql)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to Items/Err\tsql=%s", sql)
	}
	return items, nil
}
