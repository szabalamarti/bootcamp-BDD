package repository_test

import (
	"database/sql"
	"supermarket/internal"
	"supermarket/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func init() {
	cfg := mysql.Config{
		User:      "root",
		Passwd:    "test_password",
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "supermarket_test",
		ParseTime: true,
	}

	txdb.Register("txdb", "mysql", cfg.FormatDSN())
}

func TestRepositoryWarehouseSave(t *testing.T) {
	t.Run("should save a warehouse", func(t *testing.T) {
		// ARRANGE
		db, err := sql.Open("txdb", "supermarket_test")
		require.NoError(t, err)
		defer db.Close()

		// repository
		rw := repository.NewRepositoryWarehouseMySQL(db)

		// Warehouse
		warehouse := internal.Warehouse{
			Name:      "warehouse 1",
			Address:   "address 1",
			Telephone: "telephone 1",
			Capacity:  1,
		}

		// ACT
		err = rw.Save(&warehouse)

		// reset auto increment
		_, incrErr := db.Exec("ALTER TABLE warehouses AUTO_INCREMENT = 1")
		require.NoError(t, incrErr)

		// ASSERT
		require.NoError(t, err)
		require.NotEqual(t, 0, warehouse.Id)
		require.Equal(t, "warehouse 1", warehouse.Name)

	})
}

func TestRepositoryWarehouseFindById(t *testing.T) {
	t.Run("should get a warehouse by id", func(t *testing.T) {
		// ARRANGE
		db, err := sql.Open("txdb", "supermarket_test")
		require.NoError(t, err)
		defer db.Close()

		// setup db
		err = func() error {
			_, err = db.Exec("INSERT INTO warehouses (id, name, address, telephone, capacity) VALUES (1, 'warehouse 1', 'address 1', 'telephone 1', 1)")
			return err
		}()
		require.NoError(t, err)

		// repository
		rw := repository.NewRepositoryWarehouseMySQL(db)

		// ACT
		warehouse, err := rw.FindById(1)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, "warehouse 1", warehouse.Name)
		require.Equal(t, "address 1", warehouse.Address)
		require.Equal(t, "telephone 1", warehouse.Telephone)
		require.Equal(t, 1, warehouse.Capacity)
		require.Equal(t, 1, warehouse.Id)
	})
	t.Run("should return error if warehouse not found", func(t *testing.T) {
		// ARRANGE
		db, err := sql.Open("txdb", "supermarket_test")
		require.NoError(t, err)
		defer db.Close()

		// repository
		rw := repository.NewRepositoryWarehouseMySQL(db)

		// ACT
		warehouse, err := rw.FindById(1)

		// ASSERT
		require.Error(t, err)
		require.Equal(t, internal.ErrRepositoryWarehouseNotFound, err)
		require.Equal(t, internal.Warehouse{}, warehouse)
	})
}
