package repository

import (
	"errors"
	"log"
	"test_project/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	categoryRepository struct {
		DB *gorm.DB
	}
)

func NewCategoryRepository(conn *gorm.DB) CategoryRepository {
	return &categoryRepository{
		DB: conn,
	}
}

func (repo *categoryRepository) Create(model entity.Category) (string, error) {
	db := repo.DB.Create(&model)
	err := db.Error
	if err != nil {
		log.Println("error cause: ", err)
		return "", err
	}
	return model.ID.String(), nil
}

// Log
func (repo *categoryRepository) CreateLog(model entity.CategoryLog) error {
	db := repo.DB.Create(&model)
	err := db.Error
	if err != nil {
		log.Println("error cause: ", err)
		return err
	}
	return nil
}

func (repo *categoryRepository) Delete(model *entity.Category) error {
	qResult := repo.DB.Clauses(clause.Returning{}).Delete(model)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if qResult.Error != nil {
		return qResult.Error
	}
	return nil
}

func (repo *categoryRepository) Get(id string) (*entity.Category, error) {
	var result entity.Category
	qResult := repo.DB.First(&result, "id = ?", id)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if qResult.Error != nil {
		return nil, qResult.Error
	}
	return &result, nil
}

func (repo *categoryRepository) List(limit, offset int) ([]entity.Category, int64, error) {
	var results []entity.Category
	var tx = repo.DB.Model(&entity.Category{})
	tx.Order("id asc")

	if limit > 0 {
		tx.Limit(limit)
		tx.Offset(offset)
	}

	err := tx.Find(&results).Error
	if err != nil {
		return make([]entity.Category, 0), 0, err
	}

	var total int64
	err = tx.Count(&total).Error
	if err != nil {
		return make([]entity.Category, 0), 0, err
	}

	return results, total, nil
}

func (repo *categoryRepository) Update(model *entity.Category) error {
	qResult := repo.DB.Select("*").Where("id = ?", model.ID).Updates(model)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if qResult.Error != nil {
		return qResult.Error
	}
	return nil
}
