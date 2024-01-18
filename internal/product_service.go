package internal

import "errors"

var (
	ErrFieldRequired        = errors.New("field is required")
	ErrFieldFormat          = errors.New("field has an invalid format")
	ErrProductAlreadyExists = errors.New("product already exists")

	ErrProductID = errors.New("product id provided is invalid")
)

type ProductService interface {
	Save(product *Product) error
	GetById(id int) (Product, error)
	Update(Product *Product) error
	Delete(id int) error
}
