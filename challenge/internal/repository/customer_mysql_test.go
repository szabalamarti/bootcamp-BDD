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

	txdb.Register("txdb_customer_repository", "mysql", cfg.FormatDSN())
}

func TestFindTotalByCondition(t *testing.T) {
	t.Run("should return the total of customers by condition", func(t *testing.T) {
		// ARRANGE
		db, err := sql.Open("txdb_customer_repository", "fantasy_products_test")
		require.NoError(t, err)
		defer db.Close()

		// reset database on close
		defer func() error {
			_, err = db.Exec("DELETE FROM customers")
			_, err = db.Exec("ALTER TABLE customers AUTO_INCREMENT = 1")
			_, err = db.Exec("DELETE FROM invoices")
			_, err = db.Exec("ALTER TABLE invoices AUTO_INCREMENT = 1")
			require.NoError(t, err)
			return err
		}()

		// repository
		rp := repository.NewCustomersMySQL(db)

		// populate customers
		_, err = db.Exec("INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)", "customer", "1", 1)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)", "customer", "2", 0)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)", "customer", "3", 1)

		// populate invoices
		_, err = db.Exec("INSERT INTO invoices (`customer_id`, `datetime`, `total`) VALUES (?, ?, ?)", 1, "2021-01-01 00:00:00", 100)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO invoices (`customer_id`, `datetime`, `total`) VALUES (?, ?, ?)", 2, "2021-01-01 00:00:00", 50)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO invoices (`customer_id`, `datetime`, `total`) VALUES (?, ?, ?)", 3, "2021-01-01 00:00:00", 50)
		require.NoError(t, err)

		// expected output
		expected := []internal.TotalByCondition{
			{
				Condition: 1,
				Total:     150,
			},
			{
				Condition: 0,
				Total:     50,
			},
		}

		// ACT
		result, err := rp.FindTotalByCondition()

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})
}

func TestFindTopActiveCustomers(t *testing.T) {
	t.Run("should return the top active customers", func(t *testing.T) {
		// ARRANGE
		db, err := sql.Open("txdb_customer_repository", "fantasy_products_test")
		require.NoError(t, err)
		defer db.Close()

		// reset database on close
		defer func() error {
			_, err = db.Exec("DELETE FROM customers")
			_, err = db.Exec("ALTER TABLE customers AUTO_INCREMENT = 1")
			_, err = db.Exec("DELETE FROM invoices")
			_, err = db.Exec("ALTER TABLE invoices AUTO_INCREMENT = 1")
			require.NoError(t, err)
			return err
		}()

		// repository
		rp := repository.NewCustomersMySQL(db)

		// populate customers
		_, err = db.Exec("INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)", "customer", "1", 1)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)", "customer", "2", 0)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)", "customer", "3", 1)
		require.NoError(t, err)

		// populate invoices
		_, err = db.Exec("INSERT INTO invoices (`customer_id`, `datetime`, `total`) VALUES (?, ?, ?)", 1, "2021-01-01 00:00:00", 100)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO invoices (`customer_id`, `datetime`, `total`) VALUES (?, ?, ?)", 2, "2021-01-01 00:00:00", 50)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO invoices (`customer_id`, `datetime`, `total`) VALUES (?, ?, ?)", 3, "2021-01-01 00:00:00", 10)
		require.NoError(t, err)

		// expected output
		expected := []internal.CustomerAmount{
			{
				FirstName: "customer",
				LastName:  "1",
				Amount:    100,
			},
			{
				FirstName: "customer",
				LastName:  "2",
				Amount:    50,
			},
		}

		// ACT
		result, err := rp.FindTopActive(2)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})
	t.Run("should return all entries if limit is greater than the number of customers", func(t *testing.T) {
		// ARRANGE
		db, err := sql.Open("txdb_customer_repository", "fantasy_products_test")
		require.NoError(t, err)
		defer db.Close()

		// reset database on close
		defer func() error {
			_, err = db.Exec("DELETE FROM customers")
			_, err = db.Exec("ALTER TABLE customers AUTO_INCREMENT = 1")
			_, err = db.Exec("DELETE FROM invoices")
			_, err = db.Exec("ALTER TABLE invoices AUTO_INCREMENT = 1")
			require.NoError(t, err)
			return err
		}()

		// repository
		rp := repository.NewCustomersMySQL(db)

		// populate customers
		_, err = db.Exec("INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)", "customer", "1", 1)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)", "customer", "2", 0)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)", "customer", "3", 1)
		require.NoError(t, err)

		// populate invoices
		_, err = db.Exec("INSERT INTO invoices (`customer_id`, `datetime`, `total`) VALUES (?, ?, ?)", 1, "2021-01-01 00:00:00", 100)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO invoices (`customer_id`, `datetime`, `total`) VALUES (?, ?, ?)", 2, "2021-01-01 00:00:00", 50)
		require.NoError(t, err)
		_, err = db.Exec("INSERT INTO invoices (`customer_id`, `datetime`, `total`) VALUES (?, ?, ?)", 3, "2021-01-01 00:00:00", 10)
		require.NoError(t, err)

		// expected output
		expected := []internal.CustomerAmount{
			{
				FirstName: "customer",
				LastName:  "1",
				Amount:    100,
			},
			{
				FirstName: "customer",
				LastName:  "2",
				Amount:    50,
			},
			{
				FirstName: "customer",
				LastName:  "3",
				Amount:    10,
			},
		}

		// ACT
		result, err := rp.FindTopActive(4)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})
}
