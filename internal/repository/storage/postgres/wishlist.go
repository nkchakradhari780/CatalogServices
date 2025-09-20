package postgres

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/nkchakradhari780/catalogServices/internal/modules"
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

func (p *Postgres) FetchWishListItems(user_id int) ([]modules.WishList, []modules.Product, error) {

	rows, err := p.Db.Query(`SELECT wi.wish_list_id, wi.product_id, wi.user_id, wi.added_at,
					p.product_id, p.name, p.price, p.stock, p.category_id, p.quantity, p.brand, p.images
				FROM wishList wi
				JOIN products p ON wi.product_id = p.product_id
				WHERE wi.user_id = $1
	`, user_id)

	if err != nil {
		return nil, nil, fmt.Errorf("error fetching wishlist items: %w", err)
	}

	defer rows.Close()

	var wishlistItems []modules.WishList
	var products []modules.Product

	for rows.Next() {
		var wi modules.WishList
		var p modules.Product

		err := rows.Scan(&wi.WishListId, &wi.ProductId, &wi.UserId, &wi.AddedAt, &p.ProductId, &p.Name, &p.Price, &p.Quantity, &p.Stock, &p.Brand, &p.CategoryID, pq.Array(&p.Images))

		
		if err != nil {
			return nil, nil, fmt.Errorf("error scanning row: %w", err)
		}

		wishlistItems = append(wishlistItems, wi)
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil,nil, fmt.Errorf("row iteration err: %w", err)
	}

	return wishlistItems, products, nil
}
