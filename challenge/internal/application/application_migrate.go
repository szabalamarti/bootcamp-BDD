package application

import (
	"app/internal/loader"
	"app/internal/repository"
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

// ConfigApplicationMigrate is the configuration for NewApplicationDefault.
type ConfigApplicationMigrate struct {
	// Db is the database configuration.
	Db *mysql.Config
	// Addr is the server address.
	Addr string
}

// NewApplicationMigrate creates a new ApplicationMigrate.
func NewApplicationMigrate(config *ConfigApplicationMigrate) *ApplicationMigrate {
	// default values
	defaultCfg := &ConfigApplicationMigrate{
		Db:   nil,
		Addr: ":8080",
	}
	if config != nil {
		if config.Db != nil {
			defaultCfg.Db = config.Db
		}
		if config.Addr != "" {
			defaultCfg.Addr = config.Addr
		}
	}

	return &ApplicationMigrate{
		cfgDb:   defaultCfg.Db,
		cfgAddr: defaultCfg.Addr,
	}
}

// ApplicationMigrate is an implementation of the Application interface.
type ApplicationMigrate struct {
	// cfgDb is the database configuration.
	cfgDb *mysql.Config
	// cfgAddr is the server address.
	cfgAddr string
	// db is the database connection.
	db *sql.DB
}

// SetUp sets up the application.
func (a *ApplicationMigrate) SetUp() (err error) {
	// dependencies
	// - db: init
	a.db, err = sql.Open("mysql", a.cfgDb.FormatDSN())
	if err != nil {
		return
	}
	// - db: ping
	err = a.db.Ping()
	if err != nil {
		return
	}
	return
}

// Run runs the application.
func (a *ApplicationMigrate) Run() (err error) {
	// customer
	customerRepository := repository.NewCustomersMySQL(a.db)
	customerLoader := loader.NewCustomerLoaderJSON(customerRepository, "docs/db/json/customers.json")

	// product
	productRepository := repository.NewProductsMySQL(a.db)
	productLoader := loader.NewProductLoaderJSON(productRepository, "docs/db/json/products.json")

	// invoice
	invoiceRepository := repository.NewInvoicesMySQL(a.db)
	invoiceLoader := loader.NewInvoiceLoaderJSON(invoiceRepository, "docs/db/json/invoices.json")

	// sale
	saleRepository := repository.NewSalesMySQL(a.db)
	saleLoader := loader.NewSaleLoaderJSON(saleRepository, "docs/db/json/sales.json")

	// migrate
	err = customerLoader.Migrate()
	if err != nil {
		return
	}

	err = productLoader.Migrate()
	if err != nil {
		return
	}

	err = invoiceLoader.Migrate()
	if err != nil {
		return
	}

	err = saleLoader.Migrate()
	if err != nil {
		return
	}

	return
}
