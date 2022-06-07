package repository

import (
	"strings"

	"github.com/fsmiamoto/zcart/cart_service/internal/entity"
)

type InMemoryProductRepo struct {
	products map[uint64]*entity.Product
}

var _ ProductRepository = (*InMemoryProductRepo)(nil)

func (repo *InMemoryProductRepo) GetByID(id uint64) (*entity.Product, error) {
	product, ok := repo.products[id]
	if !ok {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (repo *InMemoryProductRepo) SearchByName(name string) ([]*entity.Product, error) {
	var candidates []*entity.Product
	for _, product := range repo.products {
		if strings.Contains(product.Name, name) {
			candidates = append(candidates, product)
		}
	}
	return candidates, nil
}

type InMemoryCartRepo struct {
	carts map[uint64]*entity.Cart
}

var _ CartRepository = (*InMemoryCartRepo)(nil)

func (repo *InMemoryCartRepo) GetByID(id uint64) (*entity.Cart, error) {
	if _, ok := repo.carts[id]; !ok {
		repo.carts[id] = &entity.Cart{}
	}
	return repo.carts[id], nil
}

func (repo *InMemoryCartRepo) AddProduct(c *entity.Cart, cp *entity.CartProduct) error {
	return nil
}

func (repo *InMemoryCartRepo) RemoveProduct(c *entity.Cart, cp *entity.CartProduct) error {
	return nil
}
