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
	customerRepository struct {
		DB *gorm.DB
	}
)

func NewCustomerRepository(conn *gorm.DB) CustomerRepository {
	return &customerRepository{
		DB: conn,
	}
}

func (repo *customerRepository) Create(model entity.Customer) (string, error) {
	db := repo.DB.Create(&model)
	err := db.Error
	if err != nil {
		log.Println("error cause: ", err)
		return "", err
	}
	return model.ID.String(), nil
}

// Log
func (repo *customerRepository) CreateLog(model entity.CustomerLog) error {
	db := repo.DB.Create(&model)
	err := db.Error
	if err != nil {
		log.Println("error cause: ", err)
		return err
	}
	return nil
}

func (repo *customerRepository) Delete(model *entity.Customer) error {
	qResult := repo.DB.Clauses(clause.Returning{}).Delete(model)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if qResult.Error != nil {
		return qResult.Error
	}
	return nil
}

func (repo *customerRepository) Get(id string) (*entity.Customer, error) {
	var result entity.Customer
	qResult := repo.DB.First(&result, "id = ?", id)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if qResult.Error != nil {
		return nil, qResult.Error
	}
	return &result, nil
}

func (repo *customerRepository) GetByEmail(email string) (*entity.Customer, error) {
	var result entity.Customer
	qResult := repo.DB.First(&result, "email = ?", email)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if qResult.Error != nil {
		return nil, qResult.Error
	}
	return &result, nil
}

func (repo *customerRepository) GetByUsername(username string) (*entity.Customer, error) {
	var result entity.Customer
	qResult := repo.DB.First(&result, "username = ?", username)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if qResult.Error != nil {
		return nil, qResult.Error
	}
	return &result, nil
}

func (repo *customerRepository) GetByUsernameOrEmail(username, email string) (*entity.Customer, error) {
	var result entity.Customer
	qResult := repo.DB.First(&result, "username = ? OR email = ?", username, email)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if qResult.Error != nil {
		return nil, qResult.Error
	}
	return &result, nil
}

func (repo *customerRepository) List(limit, offset int, filters map[string]interface{}) ([]entity.Customer, int64, error) {
	var results []entity.Customer
	var tx = repo.DB.Model(&entity.Customer{})
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

			if field == "last_login_date" {
				var values = value.([]interface{})
				tx.Where("last_login_date between ? AND ?", values[0], values[1])
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

	err := tx.Find(&results).Error
	if err != nil {
		return make([]entity.Customer, 0), 0, err
	}

	total, err := repo.listCount(&entity.Customer{}, filters)
	if err != nil {
		return make([]entity.Customer, 0), 0, err
	}

	return results, total, nil
}

func (repo *customerRepository) listCount(model interface{}, filters map[string]interface{}, condition ...interface{}) (int64, error) {
	var results []entity.Customer
	var tx = repo.DB.Model(&model)

	if len(filters) > 0 {
		for field, value := range filters {
			// avoid sql injection in column name
			var checkField = strings.Split(field, " ")
			if len(checkField) > 1 {
				field = checkField[0]
			}

			if field == "last_login_date" {
				var values = value.([]interface{})
				tx.Where("last_login_date between ? AND ?", values[0], values[1])
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

func (repo *customerRepository) Update(model *entity.Customer) error {
	qResult := repo.DB.Select("*").Where("id = ?", model.ID).Updates(model)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if qResult.Error != nil {
		return qResult.Error
	}
	return nil
}
