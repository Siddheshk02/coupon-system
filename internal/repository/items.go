package repository

import (
	"context"
	"database/sql"
)

type Item struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

type ItemRepository struct {
	DB *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{DB: db}
}

func (i *ItemRepository) CreateItem(ctx context.Context, req Item) error {
	query := `INSERT INTO items (name, category, price) 
              VALUES ($1, $2, $3)`
	_, err := i.DB.ExecContext(ctx, query, req.Name, req.Category, req.Price)
	return err
}

func (i *ItemRepository) GetItems(ctx context.Context, id int, category string) ([]Item, error) {
	rows, err := i.DB.QueryContext(ctx, `SELECT id, name, category, price FROM items WHERE ($1 = 0 OR id = $1) and ($2 = '' OR category = $2) ORDER BY id`, id, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.Price); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
