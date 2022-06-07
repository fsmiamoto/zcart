package repository

import "github.com/fsmiamoto/zcart/cart_service/internal/entity"

type CartRepository interface {
	GetByID(id uint64) (*entity.Cart, error)
	AddProduct(c *entity.Cart, cp *entity.CartProduct) error
	RemoveProduct(c *entity.Cart, cp *entity.CartProduct) error
}
