package sqlite

import (
	"database/sql"
	"errors"

	"github.com/fsmiamoto/zcart/cart_service/internal/models"
	"github.com/fsmiamoto/zcart/cart_service/internal/repository"
)

var (
	ErrCartNotFound = errors.New("cart not found")
)

type sqlCartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) repository.CartRepository {
	return &sqlCartRepository{db}
}

func (c *sqlCartRepository) EmptyCart(cartId string) error {
	return c.emptyCart(cartId)
}

func (c *sqlCartRepository) RemoveProduct(cartId string, productId string) error {
	return c.removeProduct(cartId, productId)
}

func (c *sqlCartRepository) UpdateProductQuantity(cartId string, productId string, delta int) error {
	return c.updateQuantity(cartId, productId, delta)
}

func (c *sqlCartRepository) GetCartProduct(cartId string, productId string) (*models.CartProduct, error) {
	const query = `SELECT cart_id,product_id,quantity FROM cart_products WHERE cart_id = ? AND product_id = ?`

	row := c.db.QueryRow(query, cartId, productId)

	cartProduct := &models.CartProduct{}
	if err := row.Scan(&cartProduct.CartID, &cartProduct.ProductID, &cartProduct.Quantity); err != nil {
		return nil, err
	}

	return cartProduct, nil
}

func (c *sqlCartRepository) GetCart(cartId string) (*models.Cart, error) {
	const query = `
        SELECT
          cp.cart_id,
          cp.product_id,
          cp.quantity,
          p.name,
          p.price,
          p.id,
          p.description,
          p.image_url
        FROM
          cart_products cp
          JOIN products p ON cp.product_id = p.id
        WHERE
          cart_id = ?
        ORDER BY
          cp.updated_at;
`

	var cartProducts []*models.CartProduct

	rows, err := c.db.Query(query, cartId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		cp := &models.CartProduct{}
		if err := rows.Scan(
			&cp.CartID, &cp.ProductID, &cp.Quantity, &cp.Product.Name,
			&cp.Product.Price, &cp.Product.ID, &cp.Product.Description, &cp.Product.ImageURL,
		); err != nil {
			return nil, err
		}
		cartProducts = append(cartProducts, cp)
	}

	return &models.Cart{
		ID:       cartId,
		Products: cartProducts,
	}, nil
}

func (c *sqlCartRepository) emptyCart(cartId string) error {
	const query = `DELETE FROM cart_products WHERE cart_id = ?`
	_, err := c.db.Exec(query, cartId)
	return err
}

func (c *sqlCartRepository) removeProduct(cartId string, productId string) error {
	const query = `DELETE FROM cart_products WHERE cart_id = ? AND product_id = ?`
	_, err := c.db.Exec(query, cartId, productId)
	return err
}

func (c *sqlCartRepository) updateQuantity(cartId string, productId string, delta int) error {
	// Docs: https://sqlite.org/lang_upsert.html
	const query = `
        INSERT INTO
          cart_products(cart_id, product_id, quantity)
        VALUES
          (?, ?, ?) ON CONFLICT(cart_id, product_id) DO
        UPDATE
        SET
          quantity = quantity + excluded.quantity;

        DELETE FROM
            cart_products
        WHERE
            cart_id = ? AND product_id = ? AND quantity <= 0;
    `
	_, err := c.db.Exec(query, cartId, productId, delta, cartId, productId)

	return err
}
