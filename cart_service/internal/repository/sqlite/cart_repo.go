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

func (c *sqlCartRepository) AddProduct(cartId string, productId string, amount uint) error {
	var quantity uint

	quantity += amount

	cartProduct, err := c.GetCartProduct(cartId, productId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		// Product not added to cart, nothing to add
	} else {
		quantity += cartProduct.Quantity
	}

	return c.SetProductQuantity(cartId, productId, quantity)
}

func (c *sqlCartRepository) RemoveProduct(cartId string, productId string, amount uint) error {
	// Using signed int to detect underflow
	var quantity int

	quantity -= int(amount)

	cartProduct, err := c.GetCartProduct(cartId, productId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		// Product not added to cart, nothing to add
	} else {
		quantity += int(cartProduct.Quantity)
	}

	if quantity < 0 {
		quantity = 0
	}

	return c.SetProductQuantity(cartId, productId, uint(quantity))
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

func (c *sqlCartRepository) removeProduct(cartId string, productId string) error {
	const query = `DELETE FROM cart_products WHERE cart_id = ? AND product_id = ?`
	_, err := c.db.Exec(query, cartId, productId)
	return err
}

func (c *sqlCartRepository) SetProductQuantity(cartId string, productId string, quantity uint) error {
	// Docs: https://sqlite.org/lang_upsert.html
	const query = `
        INSERT INTO
          cart_products(cart_id, product_id, quantity)
        VALUES
          (?, ?, ?) ON CONFLICT(cart_id, product_id) DO
        UPDATE
        SET
          quantity = excluded.quantity;
    `

	if quantity == 0 {
		return c.removeProduct(cartId, productId)
	}

	_, err := c.db.Exec(query, cartId, productId, quantity)

	return err
}
