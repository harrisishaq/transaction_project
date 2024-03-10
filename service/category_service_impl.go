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

type categoryService struct {
	repoCategory repository.CategoryRepository
}

func NewCategoryService(repoCategory repository.CategoryRepository) CategoryService {
	return &categoryService{repoCategory}
}

func (svc *categoryService) CreateCategory(req *model.CreateCategoryRequest) error {
	timeNow := time.Now()

	var newData = &entity.Category{
		ID:       uuid.New(),
		Name:     req.Name,
		IsActive: true,
		Audit: &entity.Audit{
			CurrNo:    1,
			CreatedAt: &timeNow,
			CreatedBy: "SYSTEM",
		},
	}

	_, err := svc.repoCategory.Create(*newData)
	if err != nil {
		log.Println("Error while create new data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	}

	return nil
}

func (svc *categoryService) DeleteCategory(id string) error {
	oldData, err := svc.repoCategory.Get(id)
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

	err = svc.repoCategory.Delete(oldData)
	if err != nil {
		log.Println("Error while delete data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	}

	return nil
}

func (svc *categoryService) GetCategory(id string) (*model.DataCategoryResponse, error) {
	dataCategory, err := svc.repoCategory.Get(id)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return nil, model.NewError("500", "Internal server error.")
	} else if dataCategory == nil {
		return nil, model.NewError("404", "Data not found.")
	}

	return &model.DataCategoryResponse{
		ID:       dataCategory.ID.String(),
		Name:     dataCategory.Name,
		IsActive: dataCategory.IsActive,
	}, nil
}

func (svc *categoryService) ListCategory(req model.ListCategoryRequest) ([]model.DataCategoryResponse, int64, error) {
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	if req.Filter["is_active"] == nil {
		req.Filter["is_active"] = true
	}

	var offset = (req.Page - 1) * req.Limit
	dataCategory, total, err := svc.repoCategory.List(req.Limit, offset, req.Filter)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return make([]model.DataCategoryResponse, 0), 0, model.NewError("500", "Internal server error.")
	} else if len(dataCategory) == 0 {
		return make([]model.DataCategoryResponse, 0), 0, nil
	}

	var respData []model.DataCategoryResponse
	for _, data := range dataCategory {

		respData = append(respData, model.DataCategoryResponse{
			ID:       data.ID.String(),
			Name:     data.Name,
			IsActive: data.IsActive,
		})
	}

	return respData, total, nil
}

func (service *categoryService) saveLog(data *entity.Category) (err error) {
	dataLog := entity.CategoryLog{
		ID:       fmt.Sprintf("%s-%d", data.ID.String(), data.Audit.CurrNo),
		Name:     data.Name,
		IsActive: data.IsActive,
		Audit:    data.Audit,
	}

	err = service.repoCategory.CreateLog(dataLog)
	if err != nil {
		log.Printf("Error while creating log:%+v\n ", err)
		return model.NewError("500", "Internal server error.")
	}

	return
}

func (svc *categoryService) UpdateCategory(req *model.UpdateCategoryRequest) error {
	oldData, err := svc.repoCategory.Get(req.ID)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	} else if oldData == nil {
		return model.NewError("400", "Data not found.")
	}

	logReason := fmt.Sprintf("Perubahan data oleh %v", req.ID)
	oldData.Audit.LogReason = &logReason

	err = svc.saveLog(oldData)
	if err != nil {
		log.Printf("Error while creating log: %v\n", err)
		return err
	}

	timeNow := time.Now()

	var newData = &entity.Category{
		ID:       oldData.ID,
		Name:     req.Name,
		IsActive: req.IsActive,
		Audit: &entity.Audit{
			CurrNo:    oldData.Audit.CurrNo + 1,
			CreatedAt: oldData.Audit.CreatedAt,
			CreatedBy: oldData.Audit.CreatedBy,
			UpdatedAt: &timeNow,
			UpdatedBy: req.ID,
		},
	}

	err = svc.repoCategory.Update(newData)
	if err != nil {
		log.Println("Error while update data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	}

	return nil
}
