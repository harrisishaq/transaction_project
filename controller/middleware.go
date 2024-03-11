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

func (controller *userController) middlewareCheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get Token
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

		log.Println("Start Verify Token, token: ", tokenString)

		// Verify token signature
		token, err := jwt.Parse(tokenData, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})
		if err != nil || token == nil || !token.Valid {
			log.Println("Error Cause: ", err)
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
		var userID = c.Request().Header.Get("x-consumer-id")
		if userID != sub {
			log.Println("user does not match token owner")
			return c.JSON(http.StatusUnauthorized, model.NewJsonResponse(false).SetError("401", "Unauthorized"))
		}

		return next(c)
	}
}
