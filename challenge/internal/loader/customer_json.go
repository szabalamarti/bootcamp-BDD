package loader

import (
	"app/internal"
	"encoding/json"
	"os"
)

type CustomerLoaderJSON struct {
	rp       internal.RepositoryCustomer
	filepath string
}

func NewCustomerLoaderJSON(rp internal.RepositoryCustomer, filepath string) *CustomerLoaderJSON {
	return &CustomerLoaderJSON{rp, filepath}
}

type CustomerJSON struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Condition int    `json:"condition"`
}

// Load customers from JSON file.
func (l *CustomerLoaderJSON) Load() (c []internal.Customer, err error) {
	var cs []CustomerJSON
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

	// iterate over the slice and append the customers to the slice
	for _, v := range cs {
		c = append(c, internal.Customer{
			Id: v.Id,
			CustomerAttributes: internal.CustomerAttributes{
				FirstName: v.FirstName,
				LastName:  v.LastName,
				Condition: v.Condition,
			},
		})
	}

	return
}

// Migrate customers to the repository.
func (l *CustomerLoaderJSON) Migrate() (err error) {
	// load the customers
	c, err := l.Load()
	if err != nil {
		return err
	}

	// iterate over the customers and save them to the repository
	for _, v := range c {
		err = l.rp.Save(&v)
		if err != nil {
			return err
		}
	}

	return
}
