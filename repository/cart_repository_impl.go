package repository

import (
	"errors"
	"log"
	"test_project/entity"

	"gorm.io/gorm"
)

type (
	cartRepository struct {
		DB *gorm.DB
	}
)

func NewCartRepository(conn *gorm.DB) CartRepository {
	return &cartRepository{
		DB: conn,
	}
}

func (repo *cartRepository) Create(model entity.Cart) error {
	db := repo.DB.Create(&model)
	err := db.Error
	if err != nil {
		log.Println("error cause: ", err)
		return err
	}

	return nil
}

func (repo *cartRepository) Get(custID string) (*entity.Cart, error) {
	var result entity.Cart
	qResult := repo.DB.First(&result, "customer_id = ?", custID)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if qResult.Error != nil {
		return nil, qResult.Error
	}
	return &result, nil
}

func (repo *cartRepository) Update(model *entity.Cart) error {
	qResult := repo.DB.Select("*").Where("id = ?", model.ID).Updates(model)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if qResult.Error != nil {
		return qResult.Error
	}
	return nil
}
