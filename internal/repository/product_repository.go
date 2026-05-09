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
	ProductExists(ctx context.Context, id int) (bool, error)
	ProductOwnedByUser(ctx context.Context, productID int, userID int) (bool, error)
	CreateProductImage(ctx context.Context, productID int, imageURL string) (int, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

// get all products
func (r *productRepository) GetAllProducts(ctx context.Context, query model.ProductListQuery) ([]model.Product, error) {
	baseQuery := `SELECT id, seller_id, name, description, price, stock, category, is_verified, created_at FROM products`
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
		err := rows.Scan(&p.ID, &p.SellerID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.Category, &p.IsVerified, &p.CreatedAt)
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
	query := `SELECT id, seller_id, name, description, price, stock, category, is_verified, created_at FROM products WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var p model.Product
	err := row.Scan(&p.ID, &p.SellerID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.Category, &p.IsVerified, &p.CreatedAt)
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
		INSERT INTO products (seller_id, name, description, price, stock, category, is_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		product.SellerID, product.Name, product.Description, product.Price, product.Stock, product.Category, product.IsVerified,
	).Scan(&product.ID, &product.CreatedAt)

	return err
}

func (r *productRepository) ProductExists(ctx context.Context, id int) (bool, error) {
	query := `SELECT 1 FROM products WHERE id = $1`
	var exists int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// check if product is owned by user
func (r *productRepository) ProductOwnedByUser(ctx context.Context, productID int, userID int) (bool, error) {
	query := `
		SELECT 1
		FROM products p
		JOIN sellers s ON s.id = p.seller_id
		WHERE p.id = $1 AND s.user_id = $2
	`
	var exists int
	err := r.db.QueryRowContext(ctx, query, productID, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *productRepository) CreateProductImage(ctx context.Context, productID int, imageURL string) (int, error) {
	query := `
		INSERT INTO product_images (product_id, image_url)
		VALUES ($1, $2)
		RETURNING id
	`

	var imageID int
	err := r.db.QueryRowContext(ctx, query, productID, imageURL).Scan(&imageID)
	if err != nil {
		return 0, err
	}

	return imageID, nil
}
