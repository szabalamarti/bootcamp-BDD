package internal

type SaleLoader interface {
	Load() (c []Sale, err error)
	Migrate() (err error)
}
