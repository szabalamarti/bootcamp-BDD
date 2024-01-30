package internal

type ProductLoader interface {
	Load() (c []Product, err error)
	Migrate() (err error)
}
