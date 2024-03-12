package service

import (
	"context"
	"fmt"
	"log"
	"test_project/entity"
	"test_project/model"
	"test_project/repository"

	"github.com/google/uuid"
)

type cartService struct {
	repoCart    repository.CartRepository
	repoProduct repository.ProductRepository
}

func NewCartService(repoCart repository.CartRepository, repoProduct repository.ProductRepository) CartService {
	return &cartService{repoCart, repoProduct}
}

func (svc *cartService) AddItemCart(ctx context.Context, req *model.AddItemCartRequest) error {
	// Get userContext
	var userCtx = model.GetUserContext(ctx)
	if userCtx == nil {
		log.Printf("userCtx nil")
		return model.NewError("401", "Invalid login session.")
	}

	cartData, err := svc.repoCart.Get(userCtx.UserID)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	} else if cartData == nil {
		cartData = &entity.Cart{
			ID:         uuid.New(),
			Total:      0,
			CustomerID: uuid.MustParse(userCtx.UserID),
			Products:   make([]entity.CartProduct, 0),
		}

		err := svc.repoCart.Create(*cartData)
		if err != nil {
			log.Println("Error while create data, cause: ", err)
			return model.NewError("500", "Internal server error.")
		}
	}

	dataProduct, err := svc.repoProduct.Get(req.ProductID)
	if err != nil {
		log.Println("Error while get data product, cause: ", err)
		return model.NewError("500", "Internal server error.")
	} else if dataProduct == nil || dataProduct.DeletedAt.Valid {
		log.Printf("Data Product %s not Found", req.ProductID)
		return model.NewError("404", "Product is not found.")
	} else if !dataProduct.IsActive {
		log.Printf("Product is %s not active", dataProduct.ID.String())
		return model.NewError("404", "Product is not active or archived.")
	} else if dataProduct.Qty < req.Qty {
		log.Printf("Product %s is less than required", dataProduct.ID.String())
		return model.NewError("400", "Insufficient product quantity.")
	}

	if len(cartData.Products) == 0 {
		var productCart []entity.CartProduct
		productCart = append(productCart, entity.CartProduct{
			ProductID: dataProduct.ID.String(),
			Qty:       req.Qty,
			Notes:     req.Notes,
		})

		cartData.Products = productCart
		cartData.Total = dataProduct.Price * req.Qty
	} else {
		var tempTotal int
		var isExist bool
		var newProductCart []entity.CartProduct
		for _, existingProduct := range cartData.Products {
			if dataProduct.ID.String() == existingProduct.ProductID {
				if existingProduct.Qty+req.Qty > dataProduct.Qty {
					log.Printf("Product %s is less than required", dataProduct.ID.String())
					return model.NewError("400", "Insufficient product quantity.")
				}

				newProductCart = append(newProductCart, entity.CartProduct{
					ProductID: existingProduct.ProductID,
					Qty:       existingProduct.Qty + req.Qty,
					Notes:     req.Notes,
				})

				tempTotal = tempTotal + (dataProduct.Price * (existingProduct.Qty + req.Qty))
				isExist = true
			}

			newProductCart = append(newProductCart, entity.CartProduct{
				ProductID: existingProduct.ProductID,
				Qty:       existingProduct.Qty,
				Notes:     req.Notes,
			})
			tempTotal = tempTotal + (dataProduct.Price * existingProduct.Qty)
		}

		if !isExist {
			newProductCart = append(newProductCart, entity.CartProduct{
				ProductID: dataProduct.ID.String(),
				Qty:       req.Qty,
				Notes:     req.Notes,
			})

			tempTotal = tempTotal + (dataProduct.Price * req.Qty)
		}

		cartData.Products = newProductCart
		cartData.Total = tempTotal
	}

	err = svc.repoCart.Update(cartData)
	if err != nil {
		log.Println("Error while update data cart, cause: ", err)
		return model.NewError("500", "Internal server error.")
	}

	return nil
}

func (svc *cartService) GetCart(ctx context.Context) (*model.GetCartResponse, error) {
	// Get userContext
	var userCtx = model.GetUserContext(ctx)
	if userCtx == nil {
		log.Printf("userCtx nil")
		return nil, model.NewError("401", "Invalid login session.")
	}

	cartData, err := svc.repoCart.Get(userCtx.UserID)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return nil, model.NewError("500", "Internal server error.")
	} else if cartData == nil {
		cartData = &entity.Cart{
			ID:         uuid.New(),
			Total:      0,
			CustomerID: uuid.MustParse(userCtx.UserID),
			Products:   make([]entity.CartProduct, 0),
		}

		err := svc.repoCart.Create(*cartData)
		if err != nil {
			log.Println("Error while create data, cause: ", err)
			return nil, model.NewError("500", "Internal server error.")
		}
	}

	if len(cartData.Products) > 0 {
		return &model.GetCartResponse{
			ID: cartData.ID.String(),
			Products: model.CartProduct{
				AvailableProduct:   make([]model.Product, 0),
				UnavailableProduct: make([]model.Product, 0),
			},
			Total:      0,
			CustomerID: cartData.CustomerID.String(),
		}, nil
	}

	var availableProduct []model.Product
	var unavailableProduct []model.Product
	var tempTotal int

	for _, product := range cartData.Products {
		var isAvailable = true
		var systemNotes string

		dataProduct, err := svc.repoProduct.Get(product.ProductID)
		if err != nil {
			log.Println("Error while get data product, cause: ", err)
			return nil, model.NewError("500", "Internal server error.")
		} else if dataProduct == nil || dataProduct.DeletedAt.Valid {
			systemNotes = fmt.Sprintf("Data Product %s not Found", product.ProductID)
			isAvailable = false
		} else if !dataProduct.IsActive {
			systemNotes = fmt.Sprintf("Product is %s not active", product.ProductID)
			isAvailable = false
		} else if dataProduct.Qty < product.Qty {
			systemNotes = "Insufficient product quantity."
			isAvailable = false
		}

		if isAvailable {
			availableProduct = append(availableProduct, model.Product{
				ID:         product.ProductID,
				Name:       dataProduct.Name,
				Category:   dataProduct.Category.Name,
				Qty:        product.Qty,
				Price:      dataProduct.Price,
				PriceTotal: product.Qty * dataProduct.Price,
				Notes:      product.Notes,
			})

			tempTotal = tempTotal + (dataProduct.Price * product.Qty)
		} else {
			unavailableProduct = append(unavailableProduct, model.Product{
				ID:          product.ProductID,
				Name:        dataProduct.Name,
				Category:    dataProduct.Category.Name,
				Qty:         product.Qty,
				Price:       dataProduct.Price,
				PriceTotal:  product.Qty * dataProduct.Price,
				Notes:       product.Notes,
				SystemNotes: systemNotes,
			})
		}
	}

	return &model.GetCartResponse{
		ID: cartData.ID.String(),
		Products: model.CartProduct{
			AvailableProduct:   availableProduct,
			UnavailableProduct: unavailableProduct,
		},
		Total:      tempTotal,
		CustomerID: cartData.CustomerID.String(),
	}, nil
}
