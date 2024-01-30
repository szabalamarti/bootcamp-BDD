package loader

import (
	"app/internal"
	"encoding/json"
	"os"
)

type ProductLoaderJSON struct {
	rp       internal.RepositoryProduct
	filepath string
}

func NewProductLoaderJSON(rp internal.RepositoryProduct, filepath string) *ProductLoaderJSON {
	return &ProductLoaderJSON{rp, filepath}
}

type ProductJSON struct {
	Id          int     `json:"id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// Load products from JSON file.
func (l *ProductLoaderJSON) Load() (c []internal.Product, err error) {
	var cs []ProductJSON
	// open the file
	f, err := os.Open(l.filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// the file has a JSON array, so we need to decode it into a slice
	err = json.NewDecoder(f).Decode(&cs)
	if err != nil {
		return nil, err
	}

	// iterate over the slice and append the products to the slice
	for _, v := range cs {
		c = append(c, internal.Product{
			Id: v.Id,
			ProductAttributes: internal.ProductAttributes{
				Description: v.Description,
				Price:       v.Price,
			},
		})
	}

	return
}

// Migrate products to the repository.
func (l *ProductLoaderJSON) Migrate() (err error) {
	// load the products
	c, err := l.Load()
	if err != nil {
		return err
	}

	// migrate the products to the repository
	for _, v := range c {
		err = l.rp.Save(&v)
		if err != nil {
			return err
		}
	}

	return
}
