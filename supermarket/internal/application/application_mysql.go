package application

import (
	"app/internal/handler"
	"app/internal/repository"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewApplicationMySQL creates a new default application.
func NewApplicationMySQL(addr string, dbConfig mysql.Config) (a *ApplicationMySQL) {
	// default config
	defaultRouter := chi.NewRouter()
	defaultAddr := ":8080"
	if addr != "" {
		defaultAddr = addr
	}

	a = &ApplicationMySQL{
		rt:       defaultRouter,
		addr:     defaultAddr,
		dbConfig: dbConfig,
	}
	return
}

// ApplicationMySQL is the default application.
type ApplicationMySQL struct {
	// rt is the router.
	rt *chi.Mux
	// addr is the address to listen.
	addr string
	// dbConfig is the database config.
	dbConfig mysql.Config
	// db is the connection to the database.
	db *sql.DB
}

// TearDown tears down the application.
func (a *ApplicationMySQL) TearDown() (err error) {
	// close the database connection
	if err := a.db.Close(); err != nil {
		fmt.Printf("error closing database: %v", err)
	}
	return
}

// SetUp sets up the application.
func (a *ApplicationMySQL) SetUp() (err error) {
	// dependencies

	// - store
	db, err := sql.Open("mysql", a.dbConfig.FormatDSN())
	if err != nil {
		fmt.Printf("error connecting to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		fmt.Printf("error pinging database: %v", err)
	}
	a.db = db

	// - middlewares
	a.rt.Use(middleware.Logger)
	a.rt.Use(middleware.Recoverer)

	// - warehouse
	err = a.setUpWarehouse()
	if err != nil {
		return
	}
	// - product
	err = a.setUpProduct()
	if err != nil {
		return
	}

	return
}

// Run runs the application.
func (a *ApplicationMySQL) Run() (err error) {
	err = http.ListenAndServe(a.addr, a.rt)
	return
}

func (a *ApplicationMySQL) setUpWarehouse() (err error) {
	// dependencies
	// - repository
	rw := repository.NewRepositoryWarehouseMySQL(a.db)
	// - handler
	wh := handler.NewWarehouseHandler(rw)
	// routes
	a.rt.Route("/warehouses", func(r chi.Router) {
		// GET /warehouses
		r.Get("/", wh.GetAll())
		// GET /warehouses/{id}
		r.Get("/{id}", wh.GetByID())
		// POST /warehouses
		r.Post("/", wh.Create())
	})
	return
}

func (a *ApplicationMySQL) setUpProduct() (err error) {
	// - repository
	rp := repository.NewRepositoryProductMySQL(a.db)
	// - handler
	hd := handler.NewHandlerProduct(rp)

	// router
	// - endpoints
	a.rt.Route("/products", func(r chi.Router) {
		// GET /products
		r.Get("/", hd.GetAll())
		// GET /products/{id}
		r.Get("/{id}", hd.GetById())
		// POST /products
		r.Post("/", hd.Create())
		// PUT /products/{id}
		r.Put("/{id}", hd.UpdateOrCreate())
		// PATCH /products/{id}
		r.Patch("/{id}", hd.Update())
		// DELETE /products/{id}
		r.Delete("/{id}", hd.Delete())
	})
	return
}
