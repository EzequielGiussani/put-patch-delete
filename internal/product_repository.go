package internal

import "errors"

var (
	ErrProductCodeAlreadyExists = errors.New("product code already exists")
	ErrProductNotFound          = errors.New("product not found")
)

type ProductRepository interface {
	Save(product *Product) error
	GetById(id int) (*Product, error)
}
