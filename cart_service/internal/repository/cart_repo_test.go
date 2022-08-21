package repository_test

import (
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fsmiamoto/zcart/cart_service/internal/models"
	"github.com/fsmiamoto/zcart/cart_service/internal/repository"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func createCartSetup() (repository.CartRepository, *sql.DB, sqlmock.Sqlmock) {
	db, mock := NewMock()
	return repository.NewCartRepository(db), db, mock
}

func TestCartRepo(t *testing.T) {
	t.Run("GetCart", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo, _, mock := createCartSetup()

			cartId := "2"
			rows := sqlmock.NewRows([]string{
				"cp.cart_id", "cp.product_id", "cp.quantity", "p.name",
				"p.price", "p.id", "p.description", "p.image_url",
			})

			expectedCartProducts := []*models.CartProduct{
				{
					ProductID: "1",
					Quantity:  3,
					Product: models.Product{
						ID: "1", Name: "Calzone", Price: 5.99, Description: optional("PedRão"),
					},
				},
				{
					ProductID: "2",
					Quantity:  1,
					Product: models.Product{
						ID: "2", Name: "Pão de Batata", Price: 2.99, ImageURL: optional("https://example.com"),
					},
				},
			}

			for _, cp := range expectedCartProducts {
				rows.AddRow(
					cp.CartID, cp.ProductID, cp.Quantity, cp.Product.Name,
					cp.Product.Price, cp.Product.ID, cp.Product.Description,
					cp.Product.ImageURL,
				)
			}

			mock.ExpectQuery(`SELECT .* FROM cart_products cp JOIN products p`).
				WithArgs(cartId).
				WillReturnRows(rows)

			cartProducts, err := repo.GetCart(cartId)
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())

			assert.EqualValues(t, expectedCartProducts, cartProducts)
		})

		t.Run("Error", func(t *testing.T) {
			repo, _, mock := createCartSetup()

			cartId := "2"

			expectedError := errors.New("ooops")
			mock.ExpectQuery(`SELECT .* FROM cart_products cp JOIN products p`).
				WithArgs(cartId).
				WillReturnError(expectedError)

			_, err := repo.GetCart(cartId)
			assert.ErrorIs(t, err, expectedError)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	})

	t.Run("AddProduct", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo, _, mock := createCartSetup()

			cartId := "1"
			productId := "42"
			quantity := uint(25)

			mock.ExpectExec("INSERT INTO cart_products").
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.AddProduct(cartId, productId, quantity)

			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("Error", func(t *testing.T) {
			repo, _, mock := createCartSetup()

			cartId := "1"
			productId := "42"
			quantity := uint(25)

			expectedError := errors.New("kaboom")

			mock.ExpectExec("INSERT INTO cart_products").
				WillReturnError(expectedError)

			err := repo.AddProduct(cartId, productId, quantity)

			assert.ErrorIs(t, err, expectedError)
			assert.NoError(t, mock.ExpectationsWereMet())

		})
	})
}
