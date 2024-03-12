package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"test_project/config"
	"test_project/entity"
	"test_project/model"
	"test_project/repository"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type customerService struct {
	repoCustomer repository.CustomerRepository
}

func NewCustomerService(repoCustomer repository.CustomerRepository) CustomerService {
	return &customerService{repoCustomer}
}

func (svc *customerService) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (svc *customerService) CreateCustomer(ctx context.Context, req *model.CreateCustomerRequest) error {
	// Get userContext
	var userCtx = model.GetUserContext(ctx)
	if userCtx == nil {
		log.Printf("userCtx nil")
		return model.NewError("401", "Invalid login session.")
	}

	dataExist, err := svc.repoCustomer.GetByUsernameOrEmail(req.Username, req.Email)
	if err != nil {
		log.Println("Error while existing data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	} else if dataExist != nil {
		if dataExist.Email == req.Email {
			return model.NewError("400", "Email already exist.")
		} else if dataExist.Username == req.Username {
			return model.NewError("400", "Username already exist.")
		}
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error encrypting password:", err)
		return model.NewError("500", "Internal server error.")
	}

	timeNow := time.Now()

	var createdBy string
	if userCtx.UserID != "" {
		createdBy = fmt.Sprintf("%s|%s", userCtx.UserID, userCtx.Username)
	} else {
		createdBy = "SYSTEM"
	}

	var newData = &entity.Customer{
		ID:          uuid.New(),
		Name:        req.Name,
		Email:       req.Email,
		Username:    req.Username,
		Password:    fmt.Sprintf("%x", hashPassword),
		PhoneNumber: req.PhoneNumber,
		Audit: &entity.Audit{
			CurrNo:    1,
			CreatedAt: &timeNow,
			CreatedBy: createdBy,
		},
	}

	_, err = svc.repoCustomer.Create(*newData)
	if err != nil {
		log.Println("Error while create new data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	}

	return nil
}

func (svc *customerService) DeleteCustomer(ctx context.Context, id string) error {
	// Get userContext
	var userCtx = model.GetUserContext(ctx)
	if userCtx == nil {
		log.Printf("userCtx nil")
		return model.NewError("401", "Invalid login session.")
	} else if !userCtx.IsAdmin {
		return model.NewError("401", "Forbidden.")
	}

	oldData, err := svc.repoCustomer.Get(id)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	} else if oldData == nil {
		return model.NewError("404", "Data not found.")
	}

	logReason := fmt.Sprintf("Data dihapus oleh %v", userCtx.UserID)
	oldData.Audit.LogReason = &logReason

	err = svc.saveLog(oldData)
	if err != nil {
		log.Printf("Error while creating log: %v\n", err)
		return err
	}

	err = svc.repoCustomer.Delete(oldData)
	if err != nil {
		log.Println("Error while delete data user, cause: ", err)
		return model.NewError("500", "Internal server error.")
	}

	return nil
}

func (svc *customerService) GenerateTokenAndSession(dataCustomer entity.Customer) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = dataCustomer.ID.String()
	claims["exp"] = time.Now().Add(time.Hour * 6).Unix() // Expires in 24 hours
	signedToken, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		log.Println("Error when generate JWT Token, cause: ", err)
		return "", model.NewError("500", "Internal server error.")
	}

	return signedToken, nil
}

func (svc *customerService) GetCustomer(ctx context.Context, id string) (*model.DataCustomerResponse, error) {
	// Get userContext
	var userCtx = model.GetUserContext(ctx)
	if userCtx == nil {
		log.Printf("userCtx nil")
		return nil, model.NewError("401", "Invalid login session.")
	} else if !userCtx.IsAdmin && userCtx.UserID != id {
		return nil, model.NewError("401", "Forbidden.")
	}

	dataCustomer, err := svc.repoCustomer.Get(id)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return nil, model.NewError("500", "Internal server error.")
	} else if dataCustomer == nil {
		return nil, model.NewError("404", "Data not found.")
	}

	return &model.DataCustomerResponse{
		ID:            dataCustomer.ID.String(),
		Name:          dataCustomer.Name,
		Email:         dataCustomer.Email,
		Username:      dataCustomer.Username,
		PhoneNumber:   dataCustomer.PhoneNumber,
		LastLoginDate: dataCustomer.LastLoginDate,
	}, nil
}

// Specific for middleware auth
func (svc *customerService) GetCustomerByID(id string) (*model.DataCustomerResponse, error) {
	dataCustomer, err := svc.repoCustomer.Get(id)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return nil, model.NewError("500", "Internal server error.")
	} else if dataCustomer == nil {
		return nil, model.NewError("404", "Data not found.")
	}

	return &model.DataCustomerResponse{
		ID:            dataCustomer.ID.String(),
		Name:          dataCustomer.Name,
		Email:         dataCustomer.Email,
		Username:      dataCustomer.Username,
		PhoneNumber:   dataCustomer.PhoneNumber,
		LastLoginDate: dataCustomer.LastLoginDate,
	}, nil
}

func (svc *customerService) ListCustomer(ctx context.Context, req model.ListCustomerRequest) ([]model.DataCustomerResponse, int64, error) {
	// Get userContext
	var userCtx = model.GetUserContext(ctx)
	if userCtx == nil {
		log.Printf("userCtx nil")
		return make([]model.DataCustomerResponse, 0), 0, model.NewError("401", "Invalid login session.")
	} else if !userCtx.IsAdmin {
		return make([]model.DataCustomerResponse, 0), 0, model.NewError("401", "Forbidden.")
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	var offset = (req.Page - 1) * req.Limit
	dataCustomers, total, err := svc.repoCustomer.List(req.Limit, offset, req.Filter)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return make([]model.DataCustomerResponse, 0), 0, model.NewError("500", "Internal server error.")
	} else if len(dataCustomers) == 0 {
		return make([]model.DataCustomerResponse, 0), 0, nil
	}

	var respData []model.DataCustomerResponse
	for _, data := range dataCustomers {

		respData = append(respData, model.DataCustomerResponse{
			ID:            data.ID.String(),
			Name:          data.Name,
			Email:         data.Email,
			Username:      data.Username,
			PhoneNumber:   data.PhoneNumber,
			LastLoginDate: data.LastLoginDate,
		})
	}

	return respData, total, nil
}

func (svc *customerService) LoginCustomer(req *model.LoginCustomerRequest) (string, error) {
	dataCust, err := svc.repoCustomer.GetByUsernameOrEmail(req.Username, req.Email)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return "", model.NewError("500", "Internal server error.")
	} else if dataCust == nil {
		return "", model.NewError("404", "Wrong Email/Username.")
	}

	decodedBytes, _ := hex.DecodeString(dataCust.Password)
	passMatch := svc.CheckPassword(string(decodedBytes), req.Password)
	if !passMatch {
		return "", model.NewError("401", "Wrong Password")
	}

	token, err := svc.GenerateTokenAndSession(*dataCust)
	if err != nil {
		return "", err
	}

	timeNow := time.Now()
	dataCust.LastLoginDate = &timeNow

	// Split token
	splitToken := strings.Split(token, ".")
	dataCust.Session = splitToken[2]

	err = svc.repoCustomer.Update(dataCust)
	if err != nil {
		log.Println("Error while update data user, cause: ", err)
		return "", model.NewError("500", "Internal server error.")
	}

	return token, nil
}

func (service *customerService) saveLog(data *entity.Customer) (err error) {
	dataLog := entity.CustomerLog{
		ID:              fmt.Sprintf("%s-%d", data.ID.String(), data.Audit.CurrNo),
		Name:            data.Name,
		Email:           data.Email,
		Username:        data.Username,
		Password:        data.Password,
		PhoneNumber:     data.PhoneNumber,
		ShippingAddress: data.ShippingAddress,
		LastLoginDate:   data.LastLoginDate,
		Session:         data.Session,
		Audit:           data.Audit,
	}

	err = service.repoCustomer.CreateLog(dataLog)
	if err != nil {
		log.Printf("Error while creating log:%+v\n ", err)
		return model.NewError("500", "Internal server error.")
	}

	return
}

func (svc *customerService) UpdateCustomer(ctx context.Context, req *model.UpdateCustomerRequest) error {
	// Get userContext
	var userCtx = model.GetUserContext(ctx)
	if userCtx == nil {
		log.Printf("userCtx nil")
		return model.NewError("401", "Invalid login session.")
	} else if !userCtx.IsAdmin && userCtx.UserID != req.ID {
		return model.NewError("401", "Forbidden.")
	}

	oldData, err := svc.repoCustomer.Get(req.ID)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	} else if oldData == nil {
		return model.NewError("400", "Data not found.")
	}

	if req.Email != oldData.Email {
		emailExist, err := svc.repoCustomer.GetByEmail(req.Email)
		if err != nil {
			log.Println("Error while check email, cause: ", err)
			return model.NewError("500", "Internal server error.")
		} else if emailExist != nil && fmt.Sprintf("%d", emailExist.ID) != req.ID {
			return model.NewError("400", "Email already exist.")
		}
	}

	logReason := fmt.Sprintf("Perubahan data oleh %v", userCtx.UserID)
	oldData.Audit.LogReason = &logReason

	err = svc.saveLog(oldData)
	if err != nil {
		log.Printf("Error while creating log: %v\n", err)
		return err
	}

	timeNow := time.Now()

	var newData = &entity.Customer{
		ID:              oldData.ID,
		Name:            req.Name,
		Email:           req.Email,
		Username:        oldData.Username,
		Password:        oldData.Password,
		PhoneNumber:     oldData.PhoneNumber,
		ShippingAddress: oldData.ShippingAddress,
		LastLoginDate:   oldData.LastLoginDate,
		Session:         oldData.Session,
		Audit: &entity.Audit{
			CurrNo:    oldData.Audit.CurrNo + 1,
			CreatedAt: oldData.Audit.CreatedAt,
			CreatedBy: oldData.Audit.CreatedBy,
			UpdatedAt: &timeNow,
			UpdatedBy: fmt.Sprintf("%s|%s", userCtx.UserID, userCtx.Username),
		},
	}

	err = svc.repoCustomer.Update(newData)
	if err != nil {
		log.Println("Error while update data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	}

	return nil
}

func (svc *customerService) UpdateSesion(req *model.UpdateSessionCustomerRequest) error {
	oldData, err := svc.repoCustomer.Get(req.ID)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	} else if oldData == nil {
		return model.NewError("400", "Data not found.")
	}

	oldData.Session = req.Session

	err = svc.repoCustomer.Update(oldData)
	if err != nil {
		log.Println("Error while update data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	}

	return nil
}
