package product

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/yosephbuitrago/garagesale/business/auth"
	"github.com/yosephbuitrago/garagesale/foundation/database"
)

var (

	// ErrNotFound is use when a specific Producto is requested but does not exist.
	ErrNotFound = errors.New("not found")

	ErrInvalidID = errors.New("ID is not its proper form")

	ErrForbidden = errors.New("attempted action is not allowed")
)

// Product manage the set of API's for product access
type Product struct {
	log *log.Logger
	db  *sqlx.DB
}

// New constructs a Product for api access
func New(log *log.Logger, db *sqlx.DB) Product {
	return Product{
		log: log,
		db:  db,
	}
}

// Add a product to the database. It return the created Product with
// fields like ID and DateCreated populated.
func (p Product) Create(ctx context.Context, traceID string, claims auth.Claims, np NewProduct, now time.Time) (Info, error) {
	prd := Info{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		UserID:      claims.Subject,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `INSERT INTO products
			(product_id, user_id, name, cost, quantity, date_created, date_updated)
			($1, $2, $3, $4, $5, $6, $7)`

	p.log.Printf("%s : %s : query : %s", traceID, "product.Create",
		database.Log(q, prd.ID, prd.UserID, prd.Name, prd.Cost, prd.Quantity, prd.DateCreated, prd.DateUpdated),
	)

	if _, err := p.db.ExecContext(ctx, q, prd.ID, prd.UserID, prd.Name, prd.Cost, prd.Quantity, prd.DateCreated, prd.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting product")
	}

	return prd, nil
}

func (p Product) Query(ctx context.Context, TraceID string) ([]Info, error) {
	const q = `SELECT
					p.*,
					COALESCE(SUM(s.quantity), 0) AS sold,
					COALESCE(SUN(s.paid), 0) AS revenue
			FROM product AS p
			LEFT JOIN sales AS s ON p.Product_id = s.product_id
			GROUP BY  p.product_id`

	p.log.Printf("%s : %s : query : %s", TraceID, "product.Query", database.Log(q))

	products := []Info{}

	if err := p.db.SelectContext(ctx, &products, q); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}
