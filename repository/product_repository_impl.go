package repository

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
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

// Log
func (repo *productRepository) CreateLog(model entity.ProductLog) error {
	db := repo.DB.Create(&model)
	err := db.Error
	if err != nil {
		log.Println("error cause: ", err)
		return err
	}
	return nil
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
	qResult := repo.DB.Preload("Category").Unscoped().First(&result, "id = ?", id)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if qResult.Error != nil {
		return nil, qResult.Error
	}
	return &result, nil
}

func (repo *productRepository) List(limit, offset int, filters map[string]interface{}) ([]entity.Product, int64, error) {
	var results []entity.Product
	var tx = repo.DB.Model(&entity.Product{})
	tx.Order("id asc")

	if limit > 0 {
		tx.Limit(limit)
		tx.Offset(offset)
	}

	if len(filters) > 0 {
		for field, value := range filters {
			// avoid sql injection in column name
			var checkField = strings.Split(field, " ")
			if len(checkField) > 1 {
				field = checkField[0]
			}

			if field == "is_active" {
				tx.Where(fmt.Sprintf("%s = ?", field), value)
				continue
			}

			switch reflect.TypeOf(value).Kind() {
			case reflect.Float64, reflect.Int, reflect.Int64:
				tx.Where(fmt.Sprintf("%s = ?", field), value)
			case reflect.Slice:
				tx.Where(fmt.Sprintf("%s IN (?)", field), value)
			default:
				tx.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value.(string)+"%")
			}
		}
	}

	err := tx.Preload("Category").Unscoped().Find(&results).Error
	if err != nil {
		return make([]entity.Product, 0), 0, err
	}

	total, err := repo.listCount(&entity.Product{}, filters)
	if err != nil {
		return make([]entity.Product, 0), 0, err
	}

	return results, total, nil
}

func (repo *productRepository) listCount(model interface{}, filters map[string]interface{}, condition ...interface{}) (int64, error) {
	var results []entity.Product
	var tx = repo.DB.Model(&model)

	if len(filters) > 0 {
		for field, value := range filters {
			// avoid sql injection in column name
			var checkField = strings.Split(field, " ")
			if len(checkField) > 1 {
				field = checkField[0]
			}

			if field == "is_active" {
				tx.Where(fmt.Sprintf("%s = ?", field), value)
				continue
			}

			switch reflect.TypeOf(value).Kind() {
			case reflect.Float64, reflect.Int, reflect.Int64:
				tx.Where(fmt.Sprintf("%s = ?", field), value)
			case reflect.Slice:
				tx.Where(fmt.Sprintf("%s IN (?)", field), value)
			default:
				tx.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value.(string)+"%")
			}
		}
	}

	if len(condition) > 0 {
		tx.Where(condition[0], condition[1:]...)
	}

	err := tx.Find(&results).Error
	if err != nil {
		return 0, err
	}

	var total int64
	err = tx.Count(&total).Error
	if err != nil {
		return 0, err
	}

	return total, nil
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
