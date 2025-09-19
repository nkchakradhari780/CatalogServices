package postgres

import (
	"database/sql"
	"fmt"
)

func (p *Postgres) AddToCart(user_id int, product_id int, quantity int, discount int) (int, error) {
	var cartID int
	err := p.Db.QueryRow(`
		SELECT cart_id FROM cartTable WHERE user_id = $1 AND status = 'active' LIMIT 1
	`, user_id).Scan(&cartID)

	if err == sql.ErrNoRows {
		err = p.Db.QueryRow(`
			INSERT INTO cartTable (user_id, status) VALUES ($1, 'active') RETURNING cart_id
		`, user_id).Scan(&cartID)
		if err != nil {
			return 0, fmt.Errorf("failed to create cart: %w", err)
		}
	} else if err != nil {
		return 0, fmt.Errorf("failed to fetch cart: %w", err)
	}

	var price float64
	var stock int
	err = p.Db.QueryRow(`
		SELECT price, stock FROM products WHERE product_id = $1
	`, product_id).Scan(&price, &stock)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch product details: %w", err)
	}

	if quantity > stock {
		return 0, fmt.Errorf("cannot add %d items, only %d available in stock", quantity, stock)
	}

	var existingQty int
	var cartItemID int
	err = p.Db.QueryRow(`
		SELECT cart_item_id, quantity FROM cartItems WHERE cart_id = $1 AND product_id = $2
	`, cartID, product_id).Scan(&cartItemID, &existingQty)

	if err == nil && existingQty+quantity > stock {
		return 0, fmt.Errorf("cannot add %d more, only %d left in stock", quantity, stock-existingQty)
	}

	subtotal := (price * float64(quantity)) - float64(discount)
	if subtotal < 0 {
		subtotal = 0
	}

	switch err {
	case sql.ErrNoRows:
		err = p.Db.QueryRow(`
			INSERT INTO cartItems (cart_id, product_id, quantity, price_at_time, discount, subtotal)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING cart_item_id
		`, cartID, product_id, quantity, price, discount, subtotal).Scan(&cartItemID)
		if err != nil {
			return 0, fmt.Errorf("failed to add item: %w", err)
		}
	case nil:
		_, err = p.Db.Exec(`
			UPDATE cartItems
			SET quantity = quantity + $1,
			    discount = $2,
			    subtotal = subtotal + $3
			WHERE cart_item_id = $4
		`, quantity, discount, subtotal, cartItemID)
		if err != nil {
			return 0, fmt.Errorf("failed to update cart item: %w", err)
		}
	default:
		return 0, fmt.Errorf("failed to check cart item: %w", err)
	}

	return cartItemID, nil
}

func (p *Postgres) RemoveFromCart(user_id int, product_id int) error {

	var cartId int

	err := p.Db.QueryRow(`SELECT cart_id FROM cartTable WHERE user_id = $1`, user_id).Scan(&cartId)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no cart found for the user")
		}
		return fmt.Errorf("error removing from the cart")
	}

	stmt, err := p.Db.Prepare(`DELETE FROM cartItems WHERE cart_id = $1 AND product_id = $2`)

	if err != nil {
		return fmt.Errorf("error removing from the cart")
	}

	_, err = stmt.Exec(cartId, product_id)

	if err != nil {
		return fmt.Errorf("error removing from the cart")
	}

	return nil

}
