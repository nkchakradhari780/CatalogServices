package postgres

import (
	"database/sql"
	"fmt"
)

func (p *Postgres) AddToWishList(user_id int, product_id int) (int, error) {
	stmt, err := p.Db.Prepare("INSERT INTO wishList (product_id, user_id) VALUES ($1, $2) ON CONFLICT (user_id, product_id) DO NOTHING RETURNING wish_list_id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var wishListId int

	err = stmt.QueryRow(product_id, user_id).Scan(&wishListId)

	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("product already added to wish list")
	}

	if err != nil {
		return 0, err
	}

	return int(wishListId), nil
}