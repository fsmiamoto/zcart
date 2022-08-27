package sqlite_test

import (
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fsmiamoto/zcart/cart_service/internal/models"
	"github.com/fsmiamoto/zcart/cart_service/internal/repository"
	"github.com/fsmiamoto/zcart/cart_service/internal/repository/sqlite"
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
	return sqlite.NewCartRepository(db), db, mock
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

			cart, err := repo.GetCart(cartId)
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())

			assert.EqualValues(t, expectedCartProducts, cart.Products)
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
		t.Run("Success with product already added to cart", func(t *testing.T) {
			repo, _, mock := createCartSetup()

			cartId := "1"
			productId := "42"
			quantity := uint(25)
			quantityAlreadyAdded := uint(15)

			mock.ExpectQuery("SELECT cart_id,product_id,quantity FROM cart_products").WillReturnRows(
				sqlmock.NewRows([]string{"cart_id", "product_id", "quantity"}).AddRow(cartId, productId, quantityAlreadyAdded),
			)

			mock.ExpectExec("INSERT INTO cart_products").
				WithArgs(cartId, productId, quantity+quantityAlreadyAdded).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.AddProduct(cartId, productId, quantity)

			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("Success with product not added to cart", func(t *testing.T) {
			repo, _, mock := createCartSetup()

			cartId := "1"
			productId := "42"
			quantity := uint(25)

			mock.ExpectQuery("SELECT cart_id,product_id,quantity FROM cart_products").WillReturnRows(
				sqlmock.NewRows([]string{"cart_id", "product_id", "quantity"}),
			)

			mock.ExpectExec("INSERT INTO cart_products").
				WithArgs(cartId, productId, quantity).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.AddProduct(cartId, productId, quantity)

			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("Error when inserting", func(t *testing.T) {
			repo, _, mock := createCartSetup()

			cartId := "1"
			productId := "42"
			quantity := uint(25)
			quantityAlreadyAdded := uint(15)

			mock.ExpectQuery("SELECT cart_id,product_id,quantity FROM cart_products").WillReturnRows(
				sqlmock.NewRows([]string{"cart_id", "product_id", "quantity"}).AddRow(cartId, productId, quantityAlreadyAdded),
			)

			expectedError := errors.New("nope")
			mock.ExpectExec("INSERT INTO cart_products").
				WithArgs(cartId, productId, quantity+quantityAlreadyAdded).
				WillReturnError(expectedError)

			err := repo.AddProduct(cartId, productId, quantity)

			assert.ErrorIs(t, err, expectedError)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	})

	t.Run("RemoveProduct", func(t *testing.T) {
		t.Run("Success with product added to cart", func(t *testing.T) {
			repo, _, mock := createCartSetup()

			cartId := "1"
			productId := "42"
			quantity := uint(5)
			quantityAlreadyAdded := uint(15)

			mock.ExpectQuery("SELECT cart_id,product_id,quantity FROM cart_products").WillReturnRows(
				sqlmock.NewRows([]string{"cart_id", "product_id", "quantity"}).AddRow(cartId, productId, quantityAlreadyAdded),
			)

			mock.ExpectExec("INSERT INTO cart_products").
				WithArgs(cartId, productId, quantityAlreadyAdded-quantity).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.RemoveProduct(cartId, productId, quantity)
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("Success removing more quantity than what is added to the cart", func(t *testing.T) {
			repo, _, mock := createCartSetup()

			cartId := "1"
			productId := "42"
			quantity := uint(25)
			quantityAlreadyAdded := uint(15)

			mock.ExpectQuery("SELECT cart_id,product_id,quantity FROM cart_products").WillReturnRows(
				sqlmock.NewRows([]string{"cart_id", "product_id", "quantity"}).AddRow(cartId, productId, quantityAlreadyAdded),
			)

			mock.ExpectExec("DELETE FROM cart_products").
				WithArgs(cartId, productId).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.RemoveProduct(cartId, productId, quantity)
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("Success removing with product not added to cart", func(t *testing.T) {
			repo, _, mock := createCartSetup()

			cartId := "1"
			productId := "42"
			quantity := uint(25)

			mock.ExpectQuery("SELECT cart_id,product_id,quantity FROM cart_products").WillReturnRows(
				sqlmock.NewRows([]string{"cart_id", "product_id", "quantity"}),
			)

			mock.ExpectExec("DELETE FROM cart_products").
				WithArgs(cartId, productId).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.RemoveProduct(cartId, productId, quantity)
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("Error when deleting", func(t *testing.T) {
			repo, _, mock := createCartSetup()

			expectedError := errors.New("you have no power here")
			cartId := "1"
			productId := "42"
			quantity := uint(25)

			mock.ExpectQuery("SELECT cart_id,product_id,quantity FROM cart_products").WillReturnRows(
				sqlmock.NewRows([]string{"cart_id", "product_id", "quantity"}),
			)

			mock.ExpectExec("DELETE FROM cart_products").
				WithArgs(cartId, productId).
				WillReturnError(expectedError)

			err := repo.RemoveProduct(cartId, productId, quantity)
			assert.ErrorIs(t, err, expectedError)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("Error when updating", func(t *testing.T) {
			repo, _, mock := createCartSetup()

			expectedError := errors.New("you have no power here")
			cartId := "1"
			productId := "42"
			quantity := uint(25)
			quantityAlreadyAdded := uint(50)

			mock.ExpectQuery("SELECT cart_id,product_id,quantity FROM cart_products").WillReturnRows(
				sqlmock.NewRows([]string{"cart_id", "product_id", "quantity"}).
					AddRow(cartId, productId, quantityAlreadyAdded),
			)

			mock.ExpectExec("INSERT INTO cart_products").
				WithArgs(cartId, productId, quantityAlreadyAdded-quantity).
				WillReturnError(expectedError)

			err := repo.RemoveProduct(cartId, productId, quantity)
			assert.ErrorIs(t, err, expectedError)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	})
}
