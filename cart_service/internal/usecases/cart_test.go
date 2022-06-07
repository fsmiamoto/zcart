package usecases_test

import (
	"testing"

	"github.com/fsmiamoto/zcart/cart_service/internal/entity"
	"github.com/fsmiamoto/zcart/cart_service/internal/repository"
	"github.com/fsmiamoto/zcart/cart_service/internal/usecases"
	"github.com/stretchr/testify/assert"
)

func newUsecase() usecases.Cart {
	return usecases.NewCart(&repository.InMemoryCartRepo{})
}

func TestCartUsecases(t *testing.T) {
	t.Run("add new product", func(t *testing.T) {
		cartUsecase := newUsecase()

		cart := &entity.Cart{
			ID: 1,
			Products: []*entity.CartProduct{
				{Product: &entity.Product{ID: 5}, Quantity: 1},
			},
		}

		product := &entity.CartProduct{
			Product:  &entity.Product{ID: 1},
			Quantity: 1,
		}

		err := cartUsecase.AddProduct(cart, product)

		assert.NoError(t, err)
		assert.Equal(t, len(cart.Products), 2)
	})

	t.Run("add product that is already on cart", func(t *testing.T) {
		cartUsecase := newUsecase()

		cart := &entity.Cart{
			ID: 1,
			Products: []*entity.CartProduct{
				{Product: &entity.Product{ID: 5}, Quantity: 1},
			},
		}

		product := &entity.CartProduct{
			Product:  &entity.Product{ID: 5},
			Quantity: 10,
		}

		err := cartUsecase.AddProduct(cart, product)

		assert.NoError(t, err)
		assert.Equal(t, len(cart.Products), 1)
		assert.Equal(t, cart.Products[0].Quantity, uint64(11))
	})

	t.Run("remove product", func(t *testing.T) {
		cartUsecase := newUsecase()

		cart := &entity.Cart{
			ID: 1,
			Products: []*entity.CartProduct{
				{Product: &entity.Product{ID: 5}, Quantity: 3},
			},
		}

		product := &entity.CartProduct{
			Product:  &entity.Product{ID: 5},
			Quantity: 3,
		}

		err := cartUsecase.RemoveProduct(cart, product)

		assert.NoError(t, err)
		assert.Equal(t, len(cart.Products), 0)
	})

	t.Run("remove product partially", func(t *testing.T) {
		cartUsecase := newUsecase()

		cart := &entity.Cart{
			ID: 1,
			Products: []*entity.CartProduct{
				{Product: &entity.Product{ID: 5}, Quantity: 3},
			},
		}

		product := &entity.CartProduct{
			Product:  &entity.Product{ID: 5},
			Quantity: 2,
		}

		err := cartUsecase.RemoveProduct(cart, product)

		assert.NoError(t, err)
		assert.Equal(t, len(cart.Products), 1)
		assert.Equal(t, uint64(cart.Products[0].Quantity), uint64(1))
	})
}
