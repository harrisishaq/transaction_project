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

func (repo *categoryRepository) List(limit, offset int, filters map[string]interface{}) ([]entity.Category, int64, error) {
	var results []entity.Category
	var tx = repo.DB.Model(&entity.Category{})
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

	err := tx.Find(&results).Error
	if err != nil {
		return make([]entity.Category, 0), 0, err
	}

	total, err := repo.listCount(&entity.Category{}, filters)
	if err != nil {
		return make([]entity.Category, 0), 0, err
	}

	return results, total, nil
}

func (repo *categoryRepository) listCount(model interface{}, filters map[string]interface{}, condition ...interface{}) (int64, error) {
	var results []entity.Category
	var tx = repo.DB.Model(&model)

	if len(filters) > 0 {
		for field, value := range filters {
			// avoid sql injection in column name
			var checkField = strings.Split(field, " ")
			if len(checkField) > 1 {
				field = checkField[0]
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

func (repo *categoryRepository) Update(model *entity.Category) error {
	qResult := repo.DB.Select("*").Where("id = ?", model.ID).Updates(model)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if qResult.Error != nil {
		return qResult.Error
	}
	return nil
}
