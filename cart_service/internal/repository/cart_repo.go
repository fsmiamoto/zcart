package repository

import (
	"database/sql"

	"github.com/fsmiamoto/zcart/cart_service/internal/models"
)

type CartRepository interface {
	GetCart(cartId string) ([]*models.CartProduct, error)
	AddProduct(cartId string, productId string, quantity uint) error
}

type sqlCartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &sqlCartRepository{db}
}

func (c *sqlCartRepository) AddProduct(cartId string, productId string, quantity uint) error {
	const query = `INSERT INTO cart_products (cart_id,product_id,quantity) VALUES (?,?,?)`

	_, err := c.db.Exec(query, cartId, productId, quantity)

	return err
}

func (c *sqlCartRepository) GetCart(cartId string) ([]*models.CartProduct, error) {
	const query = `
    SELECT 
        cp.cart_id, cp.product_id, cp.quantity, 
        p.name, p.price, p.id, p.description, p.image_url
        FROM cart_products cp 
        JOIN products p ON cp.product_id = p.id 
        WHERE cart_id = ?
        ORDER BY cp.updated_at;
    `

	var result []*models.CartProduct

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
		result = append(result, cp)
	}

	return result, nil
}
