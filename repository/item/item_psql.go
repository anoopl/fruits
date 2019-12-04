package itemRepository

import (
	"database/sql"
	"items/models"
)

type ItemRepository struct{}

func (i ItemRepository) GetItems(db *sql.DB, item models.Item, items []models.Item) ([]models.Item, error) {
	rows, err := db.Query("select * from items")

	if err != nil {
		return []models.Item{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&item.ID, &item.Title, &item.Owner, &item.Year)
		items = append(items, item)

	}
	if err != nil {
		return []models.Item{}, err

	}

	return items, nil

}

func (i ItemRepository) GetItem(db *sql.DB, item models.Item, id int) (models.Item, error) {
	rows := db.QueryRow("select * from items where id=$1", id)
	err := rows.Scan(&item.ID, &item.Title, &item.Owner, &item.Year)

	return item, err
}

func (i ItemRepository) AddItem(db *sql.DB, item models.Item) (int, error) {
	err := db.QueryRow("insert into items (title, owner, year) values($1, $2, $3) RETURNING id;",
		item.Title, item.Owner, item.Year).Scan(&item.ID)

	if err != nil {
		return 0, err
	}

	return item.ID, nil
}

func (b ItemRepository) UpdateItem(db *sql.DB, item models.Item) (int64, error) {
	result, err := db.Exec("update items set title=$1, author=$2, year=$3 where id=$4 RETURNING id",
		&item.Title, &item.Owner, &item.Year, &item.ID)

	if err != nil {
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsUpdated, nil
}

func (b ItemRepository) RemoveItem(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("delete from items where id = $1", id)

	if err != nil {
		return 0, err
	}

	rowsDeleted, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsDeleted, nil
}
