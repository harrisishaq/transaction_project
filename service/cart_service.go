package service

import (
	"context"
	"test_project/model"
)

type (
	CartService interface {
		AddItemCart(ctx context.Context, req *model.AddItemCartRequest) error
		GetCart(ctx context.Context) (*model.GetCartResponse, error)
	}
)
