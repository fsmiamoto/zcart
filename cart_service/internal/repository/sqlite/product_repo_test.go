package sqlite_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fsmiamoto/zcart/cart_service/internal/models"
	"github.com/fsmiamoto/zcart/cart_service/internal/repository"
	"github.com/fsmiamoto/zcart/cart_service/internal/repository/sqlite"
	"github.com/stretchr/testify/assert"
)

func optional(s string) *string {
	return &s
}

func createProductSetup() (repository.ProductRepository, *sql.DB, sqlmock.Sqlmock) {
	db, mock := NewMock()
	return sqlite.NewProductRepository(db), db, mock
}

func TestProductRepo(t *testing.T) {
	t.Run("GetProduct", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo, _, mock := createProductSetup()

			productId := "1"

			expectedProduct := models.Product{
				ID:          "1",
				Name:        "Pureisteixo 5",
				Price:       8999.99,
				Description: optional("asdf"),
				ImageURL:    optional("https://someurl.com/pureisteixo5"),
			}

			rows := sqlmock.NewRows([]string{"id", "name", "price", "description", "image_url"}).
				AddRow(expectedProduct.ID, expectedProduct.Name, expectedProduct.Price, expectedProduct.Description, expectedProduct.ImageURL)

			mock.ExpectQuery(`SELECT .* FROM products WHERE id = ?`).WithArgs(productId).WillReturnRows(rows)

			product, err := repo.GetProduct(productId)
			assert.NoError(t, err)

			assert.EqualValues(t, product, expectedProduct)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("Error", func(t *testing.T) {
			repo, _, mock := createProductSetup()

			productId := "2"

			expectedError := errors.New("not found")
			mock.ExpectQuery(`SELECT .* FROM products WHERE id = ?`).WithArgs(productId).WillReturnError(expectedError)

			_, err := repo.GetProduct(productId)
			assert.ErrorIs(t, err, expectedError)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	})
}
