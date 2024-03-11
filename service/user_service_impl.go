package service

import (
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

type userService struct {
	repoUser repository.UserRepository
}

func NewUserService(repoUser repository.UserRepository) UserService {
	return &userService{repoUser}
}

func (svc *userService) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (svc *userService) CreateUser(req *model.CreateUserRequest) error {
	dataExist, err := svc.repoUser.GetByEmail(req.Email)
	if err != nil {
		log.Println("Error while check email, cause: ", err)
		return model.NewError("500", "Internal server error.")
	} else if dataExist != nil {
		return model.NewError("400", "Email already exist.")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error encrypting password:", err)
		return model.NewError("500", "Internal server error.")
	}

	timeNow := time.Now()

	var newData = &entity.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: fmt.Sprintf("%x", hashPassword),
		Audit: &entity.Audit{
			CurrNo:    1,
			CreatedAt: &timeNow,
			CreatedBy: "SYSTEM",
		},
	}

	_, err = svc.repoUser.Create(*newData)
	if err != nil {
		log.Println("Error while create new data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	}

	return nil
}

func (svc *userService) DeleteUser(id string) error {
	oldData, err := svc.repoUser.Get(id)
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

	err = svc.repoUser.Delete(oldData)
	if err != nil {
		log.Println("Error while delete data user, cause: ", err)
		return model.NewError("500", "Internal server error.")
	}

	return nil
}

func (svc *userService) GenerateTokenAndSession(dataUser entity.User) (string, error) {
	// claims := model.Claims{
	// 	UserID: dataUser.ID.String(),
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: jwt.TimeFunc().Add(time.Hour * 24).Unix(),
	// 	},
	// }
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// signedToken, err := token.SignedString([]byte(config.AppConfig.Auth0Secret))
	// if err != nil {
	// 	log.Println("Error failed to generate JWT token, cause: ", err)
	// 	return "", model.NewError("500", "Internal server error.")
	// }

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = dataUser.ID.String()
	claims["exp"] = time.Now().Add(time.Hour * 6).Unix() // Expires in 24 hours
	signedToken, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		log.Println("Error when generate JWT Token, cause: ", err)
		return "", model.NewError("500", "Internal server error.")
	}

	return signedToken, nil
}

func (svc *userService) GetUser(id string) (*model.DataUserResponse, error) {
	dataUser, err := svc.repoUser.Get(id)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return nil, model.NewError("500", "Internal server error.")
	} else if dataUser == nil {
		return nil, model.NewError("404", "Data not found.")
	}

	return &model.DataUserResponse{
		ID:            dataUser.ID.String(),
		Name:          dataUser.Name,
		Email:         dataUser.Email,
		LastLoginDate: dataUser.LastLoginDate,
	}, nil
}

func (svc *userService) ListUser(req model.ListUserRequest) ([]model.DataUserResponse, int64, error) {
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	var offset = (req.Page - 1) * req.Limit
	dataUsers, total, err := svc.repoUser.List(req.Limit, offset, req.Filter)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return make([]model.DataUserResponse, 0), 0, model.NewError("500", "Internal server error.")
	} else if len(dataUsers) == 0 {
		return make([]model.DataUserResponse, 0), 0, nil
	}

	var respData []model.DataUserResponse
	for _, data := range dataUsers {

		respData = append(respData, model.DataUserResponse{
			ID:            data.ID.String(),
			Name:          data.Name,
			Email:         data.Email,
			LastLoginDate: data.LastLoginDate,
		})
	}

	return respData, total, nil
}

func (svc *userService) LoginUser(req *model.LoginUserRequest) (string, error) {
	dataUser, err := svc.repoUser.GetByEmail(req.Email)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return "", model.NewError("500", "Internal server error.")
	} else if dataUser == nil {
		return "", model.NewError("404", "Data not found.")
	}

	decodedBytes, _ := hex.DecodeString(dataUser.Password)
	passMatch := svc.CheckPassword(string(decodedBytes), req.Password)
	if !passMatch {
		return "", model.NewError("401", "Wrong Password")
	}

	token, err := svc.GenerateTokenAndSession(*dataUser)
	if err != nil {
		return "", err
	}

	timeNow := time.Now()
	dataUser.LastLoginDate = &timeNow

	// Split token
	splitToken := strings.Split(token, ".")
	dataUser.Session = splitToken[2]

	err = svc.repoUser.Update(dataUser)
	if err != nil {
		log.Println("Error while update data user, cause: ", err)
		return "", model.NewError("500", "Internal server error.")
	}

	return token, nil
}

func (service *userService) saveLog(data *entity.User) (err error) {
	dataLog := entity.UserLog{
		ID:            fmt.Sprintf("%s-%d", data.ID.String(), data.Audit.CurrNo),
		Name:          data.Name,
		Email:         data.Email,
		Password:      data.Password,
		LastLoginDate: data.LastLoginDate,
		Session:       data.Session,
		Audit:         data.Audit,
	}

	err = service.repoUser.CreateLog(dataLog)
	if err != nil {
		log.Printf("Error while creating log:%+v\n ", err)
		return model.NewError("500", "Internal server error.")
	}

	return
}

func (svc *userService) UpdateUser(req *model.UpdateUserRequest) error {
	oldData, err := svc.repoUser.Get(req.ID)
	if err != nil {
		log.Println("Error while get data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	} else if oldData == nil {
		return model.NewError("400", "Data not found.")
	}

	emailExist, err := svc.repoUser.GetByEmail(req.Email)
	if err != nil {
		log.Println("Error while check email, cause: ", err)
		return model.NewError("500", "Internal server error.")
	} else if emailExist != nil && fmt.Sprintf("%d", emailExist.ID) != req.ID {
		return model.NewError("400", "Email already exist.")
	}

	logReason := fmt.Sprintf("Perubahan data oleh %v", req.ID)
	oldData.Audit.LogReason = &logReason

	err = svc.saveLog(oldData)
	if err != nil {
		log.Printf("Error while creating log: %v\n", err)
		return err
	}

	timeNow := time.Now()

	var newData = &entity.User{
		ID:            oldData.ID,
		Name:          req.Name,
		Email:         req.Email,
		Password:      oldData.Password,
		LastLoginDate: oldData.LastLoginDate,
		Session:       oldData.Session,
		Audit: &entity.Audit{
			CurrNo:    oldData.Audit.CurrNo + 1,
			CreatedAt: oldData.Audit.CreatedAt,
			CreatedBy: oldData.Audit.CreatedBy,
			UpdatedAt: &timeNow,
			UpdatedBy: req.ID,
		},
	}

	err = svc.repoUser.Update(newData)
	if err != nil {
		log.Println("Error while update data, cause: ", err)
		return model.NewError("500", "Internal server error.")
	}

	return nil
}
