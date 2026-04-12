package services

import (
	"context"
	"shop_service/config"
	"shop_service/models"
)

type ProductService interface {
	Create(name string, price float64) error
	GetAll() ([]models.Product, error)
	GetByID(id int) (models.Product, error)
	Update(id int, name string, price float64) error
	Delete(id int) error
}

type productService struct{}

func NewProductService() ProductService {
	return &productService{}
}

// CREATE
func (s *productService) Create(name string, price float64) error {
	ctx := context.Background()

	_, err := config.DB.Exec(ctx,
		`INSERT INTO products (name, price) VALUES ($1,$2)`,
		name, price,
	)

	return err
}

// GET ALL
func (s *productService) GetAll() ([]models.Product, error) {
	ctx := context.Background()

	rows, err := config.DB.Query(ctx,
		`SELECT id, name, price FROM products`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// GET BY ID
func (s *productService) GetByID(id int) (models.Product, error) {
	ctx := context.Background()

	var p models.Product

	err := config.DB.QueryRow(ctx,
		`SELECT id, name, price FROM products WHERE id=$1`,
		id,
	).Scan(&p.ID, &p.Name, &p.Price)

	return p, err
}

// UPDATE
func (s *productService) Update(id int, name string, price float64) error {
	ctx := context.Background()

	_, err := config.DB.Exec(ctx,
		`UPDATE products SET name=$1, price=$2 WHERE id=$3`,
		name, price, id,
	)

	return err
}

// DELETE
func (s *productService) Delete(id int) error {
	ctx := context.Background()

	_, err := config.DB.Exec(ctx,
		`DELETE FROM products WHERE id=$1`,
		id,
	)

	return err
}
