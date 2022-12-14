// Code generated by sqlc. DO NOT EDIT.

package database

import (
	"context"

	"github.com/jackc/pgtype"
)

type Querier interface {
	AddProduct(ctx context.Context, arg AddProductParams) (Product, error)
	DeleteProduct(ctx context.Context, id pgtype.UUID) (Product, error)
	GetProduct(ctx context.Context, id pgtype.UUID) (Product, error)
	GetProducts(ctx context.Context) ([]Product, error)
}

var _ Querier = (*Queries)(nil)
