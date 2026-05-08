package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/morning-glorys/drav-backend/internal/model"
)

var ErrProductNotFound = errors.New("product not found")

type ProductRepository interface {
	GetAllProducts(ctx context.Context, query model.ProductListQuery) ([]model.Product, error)
	GetProductByID(ctx context.Context, id int) (*model.Product, error)
	CreateProduct(ctx context.Context, product *model.Product) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

// get all products
func (r *productRepository) GetAllProducts(ctx context.Context, query model.ProductListQuery) ([]model.Product, error) {
	baseQuery := `SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products`
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	argPos := 1

	if query.Search != "" {
		conditions = append(conditions, "name ILIKE $"+strconv.Itoa(argPos))
		args = append(args, "%"+query.Search+"%")
		argPos++
	}

	if query.MinPrice != nil {
		conditions = append(conditions, "price >= $"+strconv.Itoa(argPos))
		args = append(args, *query.MinPrice)
		argPos++
	}

	if query.MaxPrice != nil {
		conditions = append(conditions, "price <= $"+strconv.Itoa(argPos))
		args = append(args, *query.MaxPrice)
		argPos++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY id DESC LIMIT $" + strconv.Itoa(argPos) + " OFFSET $" + strconv.Itoa(argPos+1)
	args = append(args, query.Limit, (query.Page-1)*query.Limit)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// get product by id
func (r *productRepository) GetProductByID(ctx context.Context, id int) (*model.Product, error) {
	query := `SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var p model.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return &p, nil
}

// create product
func (r *productRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO products (name, description, price, stock, image_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		product.Name, product.Description, product.Price, product.Stock, product.ImageURL,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	return err
}
