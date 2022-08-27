package repository

import (
	"github.com/fsmiamoto/zcart/cart_service/internal/models"
)

type CartRepository interface {
	GetCart(cartId string) (*models.Cart, error)
	GetCartProduct(cartId string, productId string) (*models.CartProduct, error)
	AddProduct(cartId string, productId string, amount uint) error
	RemoveProduct(cartId string, productId string, amount uint) error
	SetProductQuantity(cartId string, productId string, quantity uint) error
}

type ProductRepository interface {
	GetProduct(productId string) (models.Product, error)
}
