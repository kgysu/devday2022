package database

import (
	"context"
	"github.com/jackc/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d Directory) AddProduct(ctx context.Context, req AddProductParams) (*Product, error) {
	product, err := d.querier.AddProduct(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unexpected error adding product: %s", err.Error())
	}
	return &product, nil
}

func (d Directory) GetProducts(ctx context.Context) ([]Product, error) {
	products, err := d.querier.GetProducts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unexpected error get products: %s", err.Error())
	}
	return products, nil
}

func (d Directory) GetProduct(ctx context.Context, id pgtype.UUID) (*Product, error) {
	product, err := d.querier.GetProduct(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unexpected error get product: %s", err.Error())
	}
	return &product, nil
}

func (d Directory) DeleteProduct(ctx context.Context, id pgtype.UUID) (*Product, error) {
	product, err := d.querier.DeleteProduct(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unexpected error deleting product: %s", err.Error())
	}
	return &product, nil
}
