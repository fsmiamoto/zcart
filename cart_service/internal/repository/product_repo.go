package repository

import (
	"errors"

	"github.com/fsmiamoto/zcart/cart_service/internal/entity"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductRepository interface {
	GetByID(id uint64) (*entity.Product, error)
	SearchByName(name string) ([]*entity.Product, error)
}
