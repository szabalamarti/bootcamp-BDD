package internal

type CustomerLoader interface {
	Load() (c []Customer, err error)
	Migrate() (err error)
}
