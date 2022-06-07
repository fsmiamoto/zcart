package usecases

import (
	"github.com/fsmiamoto/zcart/cart_service/internal/entity"
	"github.com/fsmiamoto/zcart/cart_service/internal/repository"
)

type Cart interface {
	ListProducts(c *entity.Cart) ([]*entity.CartProduct, error)
	GetByID(id uint64) (*entity.Cart, error)
	AddProduct(c *entity.Cart, cp *entity.CartProduct) error
	RemoveProduct(c *entity.Cart, cp *entity.CartProduct) error
}

type cartUsecase struct {
	cartRepository repository.CartRepository
}

var _ Cart = (*cartUsecase)(nil)

func NewCart(repo repository.CartRepository) Cart {
	return &cartUsecase{
		repo,
	}
}

func (u *cartUsecase) ListProducts(c *entity.Cart) ([]*entity.CartProduct, error) {
	return c.Products, nil
}

func (u *cartUsecase) GetByID(id uint64) (*entity.Cart, error) {
	return u.cartRepository.GetByID(id)
}

func (u *cartUsecase) AddProduct(c *entity.Cart, cp *entity.CartProduct) error {
	for _, cartProd := range c.Products {
		if cartProd.Product.ID == cp.Product.ID {
			cartProd.Quantity += cp.Quantity
			return u.cartRepository.AddProduct(c, cartProd)
		}
	}
	c.Products = append(c.Products, cp)
	return u.cartRepository.AddProduct(c, cp)
}

func (u *cartUsecase) RemoveProduct(c *entity.Cart, cp *entity.CartProduct) error {
	for i := range c.Products {
		if c.Products[i].Product.ID == cp.Product.ID {
			// Same product, decrease quantity
			c.Products[i].Quantity -= cp.Quantity
			if c.Products[i].Quantity == 0 {
				// Remove element
				c.Products = append(c.Products[:i], c.Products[i+1:]...)
			}
		}
	}
	return u.cartRepository.RemoveProduct(c, cp)
}
