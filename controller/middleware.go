package controller

import (
	"log"
	"net/http"
	"strings"
	"test_project/config"
	"test_project/model"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func Logger(c echo.Context, reqBody, resBody []byte) {
	log.Printf("[LOGGER]: Path: %s\n Request: %s\nResponse: %s \n\n", c.Request().RequestURI, string(reqBody), string(resBody))
}

// Middleware User
func (controller *userController) middlewareCheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get Token
		var userID = c.Request().Header.Get("x-consumer-id")
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		tokenSplit := strings.Split(tokenString, " ")
		if len(tokenSplit) < 2 {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		// Expecting Format "Bearer {token}"
		if len(tokenSplit) != 2 || tokenSplit[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		tokenData := tokenSplit[1]

		dataUser, err := controller.service.GetUser(userID)
		if err != nil {
			log.Println("Error Cause: ", err)
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
		} else if dataUser == nil {
			log.Println("user not found")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		log.Println("Start Verify Token, token: ", tokenString)

		// Verify token signature
		token, err := jwt.Parse(tokenData, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})
		if err != nil || token == nil || !token.Valid {
			log.Println("Error Cause: ", err)
			log.Println("Token not valid:", !token.Valid)

			if err.Error() == "Token is expired" && dataUser.Session != "" {
				err = controller.service.UpdateSesionUser(&model.UpdateSessionUserRequest{
					ID:      dataUser.ID,
					Session: "",
				})
				if err != nil {
					log.Println("Error Cause: ", err)
					return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
				}
			}

			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Invalid Token"))
		}

		// Extract user identity information from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Invalid Token Claims"))
		}

		log.Println(claims)

		sub, ok := claims["sub"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
		}

		// Verify user's identity with your application's authentication system
		if userID != sub {
			log.Println("user does not match token owner")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		tokenArr := strings.Split(tokenData, ".")
		if dataUser.Session != tokenArr[2] {
			log.Println("user session is ended/logout")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		// Set UserCtx
		var userCtx = model.UserContext{
			UserID:  sub,
			Name:    dataUser.Name,
			Token:   tokenData,
			Email:   dataUser.Email,
			IsAdmin: true,
		}

		if dataUser.Email == config.AppConfig.DefaultEmail {
			userCtx.Username = "SUPER ADMIN"
		} else {
			userCtx.Username = dataUser.Email
		}

		c.Set("userCtx", userCtx)

		return next(c)
	}
}

// Middleware Category
func (controller *categoryController) middlewareCheckAuthAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get Token
		var userID = c.Request().Header.Get("x-consumer-id")
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		tokenSplit := strings.Split(tokenString, " ")
		if len(tokenSplit) < 2 {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		// Expecting Format "Bearer {token}"
		if len(tokenSplit) != 2 || tokenSplit[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		tokenData := tokenSplit[1]

		dataUser, err := controller.UserService.GetUser(userID)
		if err != nil {
			log.Println("Error Cause: ", err)
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
		} else if dataUser == nil {
			log.Println("user not found")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		log.Println("Start Verify Token, token: ", tokenString)

		// Verify token signature
		token, err := jwt.Parse(tokenData, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})
		if err != nil || token == nil || !token.Valid {
			log.Println("Error Cause: ", err)
			log.Println("Token not valid:", !token.Valid)

			if err.Error() == "Token is expired" && dataUser.Session != "" {
				err = controller.UserService.UpdateSesionUser(&model.UpdateSessionUserRequest{
					ID:      dataUser.ID,
					Session: "",
				})
				if err != nil {
					log.Println("Error Cause: ", err)
					return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
				}
			}

			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Invalid Token"))
		}

		// Extract user identity information from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Invalid Token Claims"))
		}

		log.Println(claims)

		sub, ok := claims["sub"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
		}

		// Verify user's identity with your application's authentication system
		if userID != sub {
			log.Println("user does not match token owner")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		tokenArr := strings.Split(tokenData, ".")
		if dataUser.Session != tokenArr[2] {
			log.Println("user session is ended/logout")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		// Set UserCtx
		var userCtx = model.UserContext{
			UserID:  sub,
			Name:    dataUser.Name,
			Token:   tokenData,
			Email:   dataUser.Email,
			IsAdmin: true,
		}

		if dataUser.Email == config.AppConfig.DefaultEmail {
			userCtx.Username = "SUPER ADMIN"
		} else {
			userCtx.Username = dataUser.Email
		}

		c.Set("userCtx", userCtx)

		return next(c)
	}
}

// Middleware Product
func (controller *productController) middlewareCheckAuthAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get Token
		var userID = c.Request().Header.Get("x-consumer-id")
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		tokenSplit := strings.Split(tokenString, " ")
		if len(tokenSplit) < 2 {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		// Expecting Format "Bearer {token}"
		if len(tokenSplit) != 2 || tokenSplit[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		tokenData := tokenSplit[1]

		dataUser, err := controller.userService.GetUser(userID)
		if err != nil {
			log.Println("Error Cause: ", err)
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
		} else if dataUser == nil {
			log.Println("user not found")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		log.Println("Start Verify Token, token: ", tokenString)

		// Verify token signature
		token, err := jwt.Parse(tokenData, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})
		if err != nil || token == nil || !token.Valid {
			log.Println("Error Cause: ", err)
			log.Println("Token not valid:", !token.Valid)

			if err.Error() == "Token is expired" && dataUser.Session != "" {
				err = controller.userService.UpdateSesionUser(&model.UpdateSessionUserRequest{
					ID:      dataUser.ID,
					Session: "",
				})
				if err != nil {
					log.Println("Error Cause: ", err)
					return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
				}
			}

			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Invalid Token"))
		}

		// Extract user identity information from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Invalid Token Claims"))
		}

		log.Println(claims)

		sub, ok := claims["sub"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
		}

		// Verify user's identity with your application's authentication system
		if userID != sub {
			log.Println("user does not match token owner")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		tokenArr := strings.Split(tokenData, ".")
		if dataUser.Session != tokenArr[2] {
			log.Println("user session is ended/logout")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		// Set UserCtx
		var userCtx = model.UserContext{
			UserID:  sub,
			Name:    dataUser.Name,
			Token:   tokenData,
			Email:   dataUser.Email,
			IsAdmin: true,
		}

		if dataUser.Email == config.AppConfig.DefaultEmail {
			userCtx.Username = "SUPER ADMIN"
		} else {
			userCtx.Username = dataUser.Email
		}

		c.Set("userCtx", userCtx)

		return next(c)
	}
}

// Middleware Customer
func (controller *customerController) middlewareCheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get Token
		var userID = c.Request().Header.Get("x-consumer-id")
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		tokenSplit := strings.Split(tokenString, " ")
		if len(tokenSplit) < 2 {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		// Expecting Format "Bearer {token}"
		if len(tokenSplit) != 2 || tokenSplit[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		tokenData := tokenSplit[1]

		var userCtx = model.UserContext{}

		dataUser, err := controller.UserService.GetUser(userID)
		if err != nil {
			log.Println("Error Cause: ", err)
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
		} else if dataUser == nil {
			dataCust, err := controller.service.GetCustomerByID(userID)
			if err != nil {
				log.Println("Error Cause: ", err)
				return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
			} else if dataCust != nil {
				userCtx.Name = dataCust.Name
				userCtx.Username = dataCust.Username
				userCtx.Email = dataCust.Email
				userCtx.Name = dataCust.Name
				userCtx.IsAdmin = false
				userCtx.Session = dataCust.Session
			}

			log.Println("user not found")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		} else {
			userCtx.Name = dataUser.Name
			userCtx.Email = dataUser.Email
			userCtx.Name = dataUser.Name
			userCtx.IsAdmin = false
			userCtx.Session = dataUser.Session

			if dataUser.Email == config.AppConfig.DefaultEmail {
				userCtx.Username = "SUPER ADMIN"
			} else {
				userCtx.Username = dataUser.Email
			}
		}

		// Set UserCtx

		log.Println("Start Verify Token, token: ", tokenString)

		// Verify token signature
		token, err := jwt.Parse(tokenData, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})
		if err != nil || token == nil || !token.Valid {
			log.Println("Error Cause: ", err)
			log.Println("Token not valid:", !token.Valid)

			if err.Error() == "Token is expired" && userCtx.Session != "" {
				if userCtx.IsAdmin {
					err = controller.UserService.UpdateSesionUser(&model.UpdateSessionUserRequest{
						ID:      userID,
						Session: "",
					})
					if err != nil {
						log.Println("Error Cause: ", err)
						return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
					}
				} else {
					err = controller.service.UpdateSesion(&model.UpdateSessionCustomerRequest{
						ID:      userID,
						Session: "",
					})
					if err != nil {
						log.Println("Error Cause: ", err)
						return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
					}
				}

			}

			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Invalid Token"))
		}

		// Extract user identity information from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Invalid Token Claims"))
		}

		log.Println(claims)

		sub, ok := claims["sub"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
		}

		// Verify user's identity with your application's authentication system
		if userID != sub {
			log.Println("user does not match token owner")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		tokenArr := strings.Split(tokenData, ".")
		if userCtx.Session != tokenArr[2] {
			log.Println("user session is ended/logout")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		userCtx.UserID = sub
		userCtx.Token = tokenData

		c.Set("userCtx", userCtx)

		return next(c)
	}
}

func (controller *customerController) middlewareCheckAuthAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get Token
		var userID = c.Request().Header.Get("x-consumer-id")
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		tokenSplit := strings.Split(tokenString, " ")
		if len(tokenSplit) < 2 {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		// Expecting Format "Bearer {token}"
		if len(tokenSplit) != 2 || tokenSplit[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		tokenData := tokenSplit[1]

		dataUser, err := controller.UserService.GetUser(userID)
		if err != nil {
			log.Println("Error Cause: ", err)
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
		} else if dataUser == nil {
			log.Println("user not found")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		log.Println("Start Verify Token, token: ", tokenString)

		// Verify token signature
		token, err := jwt.Parse(tokenData, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})
		if err != nil || token == nil || !token.Valid {
			log.Println("Error Cause: ", err)
			log.Println("Token not valid:", !token.Valid)

			if err.Error() == "Token is expired" && dataUser.Session != "" {
				err = controller.UserService.UpdateSesionUser(&model.UpdateSessionUserRequest{
					ID:      dataUser.ID,
					Session: "",
				})
				if err != nil {
					log.Println("Error Cause: ", err)
					return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
				}
			}

			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Invalid Token"))
		}

		// Extract user identity information from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Invalid Token Claims"))
		}

		log.Println(claims)

		sub, ok := claims["sub"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
		}

		// Verify user's identity with your application's authentication system
		if userID != sub {
			log.Println("user does not match token owner")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		tokenArr := strings.Split(tokenData, ".")
		if dataUser.Session != tokenArr[2] {
			log.Println("user session is ended/logout")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		// Set UserCtx
		var userCtx = model.UserContext{
			UserID:  sub,
			Name:    dataUser.Name,
			Token:   tokenData,
			Email:   dataUser.Email,
			IsAdmin: true,
		}

		if dataUser.Email == config.AppConfig.DefaultEmail {
			userCtx.Username = "SUPER ADMIN"
		} else {
			userCtx.Username = dataUser.Email
		}

		c.Set("userCtx", userCtx)

		return next(c)
	}
}

func (controller *customerController) middlewareCheckAuthCust(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get Token
		var userID = c.Request().Header.Get("x-consumer-id")
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		tokenSplit := strings.Split(tokenString, " ")
		if len(tokenSplit) < 2 {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		// Expecting Format "Bearer {token}"
		if len(tokenSplit) != 2 || tokenSplit[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		tokenData := tokenSplit[1]

		dataCust, err := controller.service.GetCustomerByID(userID)
		if err != nil {
			log.Println("Error Cause: ", err)
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
		} else if dataCust == nil {
			log.Println("customer not found")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		log.Println("Start Verify Token, token: ", tokenString)
		log.Println("Data Customer: ", dataCust)

		// Verify token signature
		token, err := jwt.Parse(tokenData, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})
		if err != nil || token == nil || !token.Valid {
			log.Println("Error Cause: ", err)
			log.Println("Token not valid:", !token.Valid)

			if err.Error() == "Token is expired" && dataCust.Session != "" {
				err = controller.service.UpdateSesion(&model.UpdateSessionCustomerRequest{
					ID:      dataCust.ID,
					Session: "",
				})
				if err != nil {
					log.Println("Error Cause: ", err)
					return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
				}
			}

			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Invalid Token"))
		}

		// Extract user identity information from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Invalid Token Claims"))
		}

		log.Println(claims)

		sub, ok := claims["sub"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
		}

		// Verify user's identity with your application's authentication system
		if userID != sub {
			log.Println("user does not match token owner")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		tokenArr := strings.Split(tokenData, ".")
		log.Println("Data Session: ", dataCust.Session)
		log.Println("Token Data: ", tokenArr[2])

		if dataCust.Session != tokenArr[2] {
			log.Println("user session is ended/logout")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		// Set UserCtx
		var userCtx = model.UserContext{
			UserID:   sub,
			Name:     dataCust.Name,
			Username: dataCust.Username,
			Token:    tokenData,
			Email:    dataCust.Email,
			IsAdmin:  false,
		}

		c.Set("userCtx", userCtx)

		return next(c)
	}
}

// Middleware Cart
func (controller *cartController) middlewareCheckAuthCust(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get Token
		var userID = c.Request().Header.Get("x-consumer-id")
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		tokenSplit := strings.Split(tokenString, " ")
		if len(tokenSplit) < 2 {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		// Expecting Format "Bearer {token}"
		if len(tokenSplit) != 2 || tokenSplit[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Not Authorized"))
		}

		tokenData := tokenSplit[1]

		dataCust, err := controller.CustomerService.GetCustomerByID(userID)
		if err != nil {
			log.Println("Error Cause: ", err)
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
		} else if dataCust == nil {
			log.Println("customer not found")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		log.Println("Start Verify Token, token: ", tokenString)
		log.Println("Data Customer: ", dataCust)

		// Verify token signature
		token, err := jwt.Parse(tokenData, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})
		if err != nil || token == nil || !token.Valid {
			log.Println("Error Cause: ", err)
			log.Println("Token not valid:", !token.Valid)

			if err.Error() == "Token is expired" && dataCust.Session != "" {
				err = controller.CustomerService.UpdateSesion(&model.UpdateSessionCustomerRequest{
					ID:      dataCust.ID,
					Session: "",
				})
				if err != nil {
					log.Println("Error Cause: ", err)
					return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
				}
			}

			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Invalid Token"))
		}

		// Extract user identity information from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Invalid Token Claims"))
		}

		log.Println(claims)

		sub, ok := claims["sub"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("500", "Internal server error"))
		}

		// Verify user's identity with your application's authentication system
		if userID != sub {
			log.Println("user does not match token owner")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		tokenArr := strings.Split(tokenData, ".")
		log.Println("Data Session: ", dataCust.Session)
		log.Println("Token Data: ", tokenArr[2])

		if dataCust.Session != tokenArr[2] {
			log.Println("user session is ended/logout")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		// Set UserCtx
		var userCtx = model.UserContext{
			UserID:   sub,
			Name:     dataCust.Name,
			Username: dataCust.Username,
			Token:    tokenData,
			Email:    dataCust.Email,
			IsAdmin:  false,
		}

		c.Set("userCtx", userCtx)

		return next(c)
	}
}
