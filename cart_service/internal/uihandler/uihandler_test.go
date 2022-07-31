package uihandler

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	t.Run("testing", func(t *testing.T) {
		os.Remove("./zcart.db")
		db, err := sql.Open("sqlite3", "./zcart.db")

		assert.NoError(t, err)

		handler := New(db, log.Default())

		err = handler.applyMigration()
		assert.NoError(t, err)

		got, err := handler.getCart("2")

		assert.NoError(t, err)

		expect := []*CartProduct{
			{
				CartID:    "2",
				ProductID: "1",
				Quantity:  10,
				Product: Product{
					ID:    "1",
					Name:  "Coca Cola",
					Price: 5.99,
				},
			},
			{
				CartID:    "2",
				ProductID: "2",
				Quantity:  5,
				Product: Product{
					ID:    "2",
					Name:  "BomBril",
					Price: 1.99,
				},
			},
			{
				CartID:    "2",
				ProductID: "3",
				Quantity:  9,
				Product: Product{
					ID:    "3",
					Name:  "Leite Longa Vida 1L",
					Price: 4.99,
				},
			},
			{
				CartID:    "2",
				ProductID: "4",
				Quantity:  1,
				Product: Product{
					ID:    "4",
					Name:  "Biscoito Passatempo",
					Price: 8.99,
				},
			},
		}

		assert.ElementsMatch(t, got, expect)
	})
}
