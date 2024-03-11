package controller

import (
	"net/http"
	"test_project/model"
	"test_project/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type userController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *userController {
	return &userController{service}
}

func (controller *userController) UserRoutes(e *echo.Echo) {
	e.Use(middleware.CORS())

	// User EP
	var userRoute = e.Group("/users")
	userRoute.Use(middleware.BodyDump(Logger))
	userRoute.POST("/", controller.CreateUser)
	userRoute.POST("/list", controller.ListUsers)
	userRoute.POST("/login", controller.Login)
	userRoute.DELETE("/:id", controller.DeleteUser)
	userRoute.GET("/:id", controller.GetUser, controller.middlewareCheckAuth)
	userRoute.PUT("/:id", controller.UpdateUser)
}

func (ctrl *userController) CreateUser(c echo.Context) error {
	request := new(model.CreateUserRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	} else if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", err.Error()))
	}

	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	err := ctrl.service.CreateUser(ctx, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}

func (ctrl *userController) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	err := ctrl.service.DeleteUser(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}

func (ctrl *userController) GetUser(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	data, err := ctrl.service.GetUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).SetData(data))
}

func (ctrl *userController) ListUsers(c echo.Context) error {
	request := new(model.ListUserRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	if request.Filter == nil {
		request.Filter = make(map[string]interface{})
	}

	data, total, err := ctrl.service.ListUser(*request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).List(data, total))
}

func (ctrl *userController) Login(c echo.Context) error {
	request := new(model.LoginUserRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	} else if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", err.Error()))
	}

	token, err := ctrl.service.LoginUser(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).SetData(token))
}

func (ctrl *userController) UpdateUser(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	request := new(model.UpdateUserRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	} else if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", err.Error()))
	}

	request.ID = id
	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	err := ctrl.service.UpdateUser(ctx, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}
