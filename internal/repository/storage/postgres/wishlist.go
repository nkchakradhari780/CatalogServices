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

func (p *Postgres) RemoveFromWishList(user_id int, product_id int) error {

	var wishListId int

	err := p.Db.QueryRow(`SELECT wish_list_id FROM wishList WHERE user_id = $1 AND product_id = $2`, user_id, product_id).Scan(&wishListId)

	if err == sql.ErrNoRows {
		return fmt.Errorf("item not found on wishlist")
	} else if err != nil {
		return fmt.Errorf("error fetching item from wishlist: %w", err)
	}

	stmt, err := p.Db.Prepare(`DELETE FROM wishList WHERE user_id = $1 AND product_id = $2`)

	if err != nil {
		return fmt.Errorf("error removing product from wishlist")
	}

	_, err = stmt.Exec(user_id, product_id)

	if err != nil {
		return fmt.Errorf("error removing product from wishlist")
	}

	return nil

}
