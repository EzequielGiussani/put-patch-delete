package repository

import "app/internal"

type ProductMap struct {
	db     map[int]internal.Product
	lastId int
}

func NewProductMap(db map[int]internal.Product, startingId int) *ProductMap {
	return &ProductMap{
		db:     make(map[int]internal.Product),
		lastId: startingId,
	}
}

func (pm *ProductMap) Save(product *internal.Product) error {

	for _, prod := range pm.db {
		if prod.CodeValue == product.CodeValue {
			return internal.ErrProductCodeAlreadyExists
		}
	}

	pm.lastId++

	product.ID = pm.lastId

	pm.db[product.ID] = *product

	return nil
}

func (pm *ProductMap) GetById(id int) (internal.Product, error) {
	product, ok := pm.db[id]

	if !ok {
		return internal.Product{}, internal.ErrProductNotFound
	}

	return product, nil
}

func (pm *ProductMap) Update(product *internal.Product) error {
	_, ok := pm.db[product.ID]

	if !ok {
		return internal.ErrProductNotFound
	}

	for _, prod := range pm.db {
		if prod.CodeValue == product.CodeValue && prod.ID != product.ID {
			return internal.ErrProductCodeAlreadyExists
		}
	}

	pm.db[product.ID] = *product

	return nil
}

func (pm *ProductMap) Delete(id int) error {
	_, ok := pm.db[id]

	if !ok {
		return internal.ErrProductNotFound
	}

	delete(pm.db, id)

	return nil
}
