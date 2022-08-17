package repository

import (
	"database/sql"
	"errors"

	"github.com/fsmiamoto/zcart/cart_service/internal/models"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductRepository interface {
	GetProduct(productId string) (models.Product, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db}
}

func (c *productRepository) GetProduct(productId string) (models.Product, error) {
	const query = `SELECT id, name, price, description, image_url FROM products WHERE id = ?`
	var product models.Product

	row := c.db.QueryRow(query, productId)

	if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.ImageURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product, ErrProductNotFound
		}
		return product, err
	}

	return product, nil
}
