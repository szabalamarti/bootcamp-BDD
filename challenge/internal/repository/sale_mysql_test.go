package repository_test

import (
	"app/internal"
	"app/internal/repository"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func init() {
	cfg := mysql.Config{
		User:      "root",
		Passwd:    "",
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "fantasy_products_test",
		ParseTime: true,
	}

	txdb.Register("txdb_sale_repository", "mysql", cfg.FormatDSN())
}

func TestFindTopSold(t *testing.T) {
	t.Run("should return the top sold products", func(t *testing.T) {
		// ARRANGE
		db, err := sql.Open("txdb_sale_repository", "fantasy_products_test")
		require.NoError(t, err)
		defer db.Close()

		// reset database on close
		defer func() error {
			_, err = db.Exec("DELETE FROM sales")
			_, err = db.Exec("ALTER TABLE sales AUTO_INCREMENT = 1")
			_, err = db.Exec("DELETE FROM products")
			_, err = db.Exec("ALTER TABLE products AUTO_INCREMENT = 1")
			require.NoError(t, err)
			return err
		}()

		// repository
		rp := repository.NewSalesMySQL(db)

		// populate products
		_, err = db.Exec("INSERT INTO products (`description`, `price`) VALUES (?, ?)", "A", 100)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO products (`description`, `price`) VALUES (?, ?)", "B", 50)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO products (`description`, `price`) VALUES (?, ?)", "C", 50)
		require.NoError(t, err)

		// populate sales
		_, err = db.Exec("INSERT INTO sales (`quantity`, `product_id`) VALUES (?, ?)", 1, 1)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO sales (`quantity`, `product_id`) VALUES (?, ?)", 2, 2)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO sales (`quantity`, `product_id`) VALUES (?, ?)", 3, 3)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO sales (`quantity`, `product_id`) VALUES (?, ?)", 4, 3)
		require.NoError(t, err)

		// expected result
		expected := []internal.ProductSales{
			{
				ProductDescription: "C",
				Sales:              7,
			},
			{
				ProductDescription: "B",
				Sales:              2,
			},
		}

		// ACT
		s, err := rp.FindTopSold(2)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, expected, s)
	})
}
