package service

import (
	"app/internal"
	"errors"
	"fmt"
	"time"
)

type ProductDefault struct {
	rp internal.ProductRepository
}

func NewProductDefault(rp internal.ProductRepository) *ProductDefault {
	return &ProductDefault{
		rp: rp,
	}
}

func (pd *ProductDefault) Save(product *internal.Product) error {
	if err := pd.validateProduct(product); err != nil {
		return err
	}

	err := pd.rp.Save(product)

	if err != nil {
		switch err {
		case internal.ErrProductCodeAlreadyExists:
			err = fmt.Errorf("%w: code_value", internal.ErrProductCodeAlreadyExists)
		}
	}

	return err

}

func (pd *ProductDefault) validateProduct(p *internal.Product) error {
	switch {
	case p.Name == "":
		return fmt.Errorf("%w: name", internal.ErrFieldRequired)
	case p.Quantity == 0:
		return fmt.Errorf("%w: quantity", internal.ErrFieldRequired)
	case p.CodeValue == "":
		return fmt.Errorf("%w: code_value", internal.ErrFieldRequired)
	case p.Expiration == "":
		return fmt.Errorf("%w: expiration", internal.ErrFieldRequired)
	case p.Price == 0:
		return fmt.Errorf("%w: price", internal.ErrFieldRequired)
	}

	_, err := time.Parse("02/02/2006", p.Expiration)

	if err != nil {
		return fmt.Errorf("%w: expiration", internal.ErrFieldFormat)
	}

	return nil
}

func (pd *ProductDefault) GetById(id int) (internal.Product, error) {
	prod, err := pd.rp.GetById(id)

	if err != nil {
		switch {
		case errors.Is(err, internal.ErrProductNotFound):
			err = fmt.Errorf("%w: id", internal.ErrProductID)
		}
	}

	return prod, err
}

func (pd *ProductDefault) Update(product *internal.Product) error {

	if err := pd.validateProduct(product); err != nil {
		return err
	}

	err := pd.rp.Update(product)

	if err != nil {
		switch err {
		case internal.ErrProductNotFound:
			err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		}
	}

	return err
}

func (pd *ProductDefault) Delete(id int) error {
	err := pd.rp.Delete(id)

	if err != nil {
		switch err {
		case internal.ErrProductNotFound:
			err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		}
	}

	return err
}
