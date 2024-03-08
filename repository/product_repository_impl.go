package repository

import (
	"errors"
	"log"
	"test_project/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	productRepository struct {
		DB *gorm.DB
	}
)

func NewProductRepository(conn *gorm.DB) ProductRepository {
	return &productRepository{
		DB: conn,
	}
}

func (repo *productRepository) Create(model entity.Product) (string, error) {
	db := repo.DB.Create(&model)
	err := db.Error
	if err != nil {
		log.Println("error cause: ", err)
		return "", err
	}
	return model.ID.String(), nil
}

func (repo *productRepository) Delete(model *entity.Product) error {
	qResult := repo.DB.Clauses(clause.Returning{}).Delete(model)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if qResult.Error != nil {
		return qResult.Error
	}
	return nil
}

func (repo *productRepository) Get(id string) (*entity.Product, error) {
	var result entity.Product
	qResult := repo.DB.Preload("Category").First(&result, "id = ?", id)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if qResult.Error != nil {
		return nil, qResult.Error
	}
	return &result, nil
}

func (repo *productRepository) List(limit, offset int) ([]entity.Product, int64, error) {
	var results []entity.Product
	var tx = repo.DB.Model(&entity.Product{})
	tx.Order("id asc")

	if limit > 0 {
		tx.Limit(limit)
		tx.Offset(offset)
	}

	err := tx.Preload("Category").Find(&results).Error
	if err != nil {
		return make([]entity.Product, 0), 0, err
	}

	var total int64
	err = tx.Count(&total).Error
	if err != nil {
		return make([]entity.Product, 0), 0, err
	}

	return results, total, nil
}

func (repo *productRepository) Update(model *entity.Product) error {
	qResult := repo.DB.Select("*").Where("id = ?", model.ID).Updates(model)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if qResult.Error != nil {
		return qResult.Error
	}
	return nil
}
