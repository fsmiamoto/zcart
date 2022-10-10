package repository

import (
	"github.com/fsmiamoto/zcart/cart_service/internal/models"
)

type CartRepository interface {
	GetCart(cartId string) (*models.Cart, error)
	GetCartProduct(cartId string, productId string) (*models.CartProduct, error)
	UpdateProductQuantity(cartId string, productId string, delta int) error
	RemoveProduct(cartId string, productId string) error
	EmptyCart(cartId string) error
}

type ProductRepository interface {
	GetProduct(productId string) (models.Product, error)
}
