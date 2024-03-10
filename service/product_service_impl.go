package service

import (
	"fmt"
	"log"
	"test_project/entity"
	"test_project/model"
	"test_project/repository"
	"time"

	"github.com/google/uuid"
)

type productService struct {
	repoCategory repository.CategoryRepository
	repoProduct  repository.ProductRepository
}

func NewProductService(repoCategory repository.CategoryRepository, repoProduct repository.ProductRepository) ProductService {
	return &productService{
		repoCategory: repoCategory,
		repoProduct:  repoProduct,
	}
}

func (svc *productService) CreateProduct(req *model.CreateProductRequest) error {
	// Check if category_id is valid or not
	dataCategory, err := svc.repoCategory.Get(req.CategoryID)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	} else if dataCategory == nil {
		log.Println("Category not Exist")
		return model.NewError("400", "Category not found")
	} else if !dataCategory.IsActive {
		log.Println("Category not active")
		return model.NewError("400", "Category is not active")
	}

	timeNow := time.Now()

	var newData = &entity.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Qty:         req.Qty,
		Price:       req.Price,
		Description: req.Description,
		CategoryID:  uuid.MustParse(req.CategoryID),
		IsActive:    true,
		Audit: &entity.Audit{
			CurrNo:    1,
			CreatedAt: &timeNow,
			CreatedBy: "SYSTEM",
		},
	}

	_, err = svc.repoProduct.Create(*newData)
	if err != nil {
		log.Println("Error while create new data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	}

	return nil
}

func (svc *productService) DeleteProduct(id string) error {
	oldData, err := svc.repoProduct.Get(id)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	} else if oldData == nil {
		return model.NewError("404", "Data not found.")
	}

	logReason := fmt.Sprintf("Data dihapus oleh %v", id)
	oldData.Audit.LogReason = &logReason

	err = svc.saveLog(oldData)
	if err != nil {
		log.Printf("Error while creating log: %v\n", err)
		return err
	}

	err = svc.repoProduct.Delete(oldData)
	if err != nil {
		log.Println("Error while delete data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	}

	return nil
}

func (svc *productService) GetProduct(id string) (*model.DataProductResponse, error) {
	dataProduct, err := svc.repoProduct.Get(id)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return nil, model.NewError("500", "Internal server error.")
	} else if dataProduct == nil {
		return nil, model.NewError("404", "Data not found.")
	}

	return &model.DataProductResponse{
		ID:           dataProduct.ID.String(),
		Name:         dataProduct.Name,
		IsActive:     dataProduct.IsActive,
		CategoryID:   dataProduct.CategoryID.String(),
		CategoryName: dataProduct.Category.Name,
		Price:        dataProduct.Price,
		Qty:          dataProduct.Qty,
		Description:  dataProduct.Description,
	}, nil
}

func (svc *productService) ListProduct(req model.ListProductRequest) ([]model.DataProductResponse, int64, error) {
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	var offset = (req.Page - 1) * req.Limit
	dataProduct, total, err := svc.repoProduct.List(req.Limit, offset, req.Filter)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return make([]model.DataProductResponse, 0), 0, model.NewError("500", "Internal server error.")
	} else if len(dataProduct) == 0 {
		return make([]model.DataProductResponse, 0), 0, nil
	}

	var respData []model.DataProductResponse
	for _, data := range dataProduct {

		respData = append(respData, model.DataProductResponse{
			ID:           data.ID.String(),
			Name:         data.Name,
			IsActive:     data.IsActive,
			CategoryID:   data.CategoryID.String(),
			CategoryName: data.Category.Name,
			Price:        data.Price,
			Qty:          data.Qty,
			Description:  data.Description,
		})
	}

	return respData, total, nil
}

func (service *productService) saveLog(data *entity.Product) (err error) {
	dataLog := entity.ProductLog{
		ID:          fmt.Sprintf("%s-%d", data.ID.String(), data.Audit.CurrNo),
		Name:        data.Name,
		Qty:         data.Qty,
		Price:       data.Price,
		Description: data.Description,
		CategoryID:  data.CategoryID,
		Category:    data.Category,
		IsActive:    true,
		Audit:       data.Audit,
	}

	err = service.repoProduct.CreateLog(dataLog)
	if err != nil {
		log.Printf("Error while creating log:%+v\n ", err)
		return model.NewError("500", "Internal server error.")
	}

	return
}

func (svc *productService) UpdateProduct(req *model.UpdateProductRequest) error {
	oldData, err := svc.repoProduct.Get(req.ID)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	} else if oldData == nil {
		return model.NewError("400", "Data not found.")
	}

	// Check if category_id is different and neet to check valid or not
	if req.CategoryID != oldData.CategoryID.String() {
		dataCategory, err := svc.repoCategory.Get(req.CategoryID)
		if err != nil {
			log.Println("Error while get data, cause: ", err)
			return model.NewError("500", "Internal server error.")
		} else if dataCategory == nil {
			log.Println("Category not Exist")
			return model.NewError("400", "Category not found")
		} else if !dataCategory.IsActive {
			log.Println("Category not active")
			return model.NewError("400", "Category is not active")
		}
	}

	logReason := fmt.Sprintf("Perubahan data oleh %v", req.ID)
	oldData.Audit.LogReason = &logReason

	err = svc.saveLog(oldData)
	if err != nil {
		log.Printf("Error while creating log: %v\n", err)
		return err
	}

	timeNow := time.Now()

	var newData = &entity.Product{
		ID:          oldData.ID,
		Name:        req.Name,
		IsActive:    req.IsActive,
		Qty:         req.Qty,
		Price:       req.Price,
		Description: req.Description,
		CategoryID:  uuid.MustParse(req.CategoryID),
		Audit: &entity.Audit{
			CurrNo:    oldData.Audit.CurrNo + 1,
			CreatedAt: oldData.Audit.CreatedAt,
			CreatedBy: oldData.Audit.CreatedBy,
			UpdatedAt: &timeNow,
			UpdatedBy: "SYSTEM",
		},
	}

	err = svc.repoProduct.Update(newData)
	if err != nil {
		log.Println("Error while update data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	}

	return nil
}
