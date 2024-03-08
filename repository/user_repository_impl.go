package repository

import (
	"errors"
	"log"
	"test_project/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	userRepository struct {
		DB *gorm.DB
	}
)

func NewUserRepository(conn *gorm.DB) UserRepository {
	return &userRepository{
		DB: conn,
	}
}

func (repo *userRepository) Create(model entity.User) (string, error) {
	db := repo.DB.Create(&model)
	err := db.Error
	if err != nil {
		log.Println("error cause: ", err)
		return "", err
	}
	return model.ID.String(), nil
}

// Log
func (repo *userRepository) CreateLog(model entity.UserLog) error {
	db := repo.DB.Create(&model)
	err := db.Error
	if err != nil {
		log.Println("error cause: ", err)
		return err
	}
	return nil
}

func (repo *userRepository) Delete(model *entity.User) error {
	qResult := repo.DB.Clauses(clause.Returning{}).Delete(model)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if qResult.Error != nil {
		return qResult.Error
	}
	return nil
}

func (repo *userRepository) Get(id string) (*entity.User, error) {
	var result entity.User
	qResult := repo.DB.First(&result, "id = ?", id)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if qResult.Error != nil {
		return nil, qResult.Error
	}
	return &result, nil
}

func (repo *userRepository) GetByEmail(email string) (*entity.User, error) {
	var result entity.User
	qResult := repo.DB.First(&result, "email = ?", email)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if qResult.Error != nil {
		return nil, qResult.Error
	}
	return &result, nil
}

func (repo *userRepository) List(limit, offset int) ([]entity.User, int64, error) {
	var results []entity.User
	var tx = repo.DB.Model(&entity.User{})
	tx.Order("id asc")

	if limit > 0 {
		tx.Limit(limit)
		tx.Offset(offset)
	}

	err := tx.Find(&results).Error
	if err != nil {
		return make([]entity.User, 0), 0, err
	}

	var total int64
	err = tx.Count(&total).Error
	if err != nil {
		return make([]entity.User, 0), 0, err
	}

	return results, total, nil
}

func (repo *userRepository) Update(model *entity.User) error {
	qResult := repo.DB.Select("*").Where("id = ?", model.ID).Updates(model)
	if errors.Is(qResult.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if qResult.Error != nil {
		return qResult.Error
	}
	return nil
}
