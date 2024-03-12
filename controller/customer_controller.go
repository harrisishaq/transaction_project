package controller

import (
	"net/http"
	"test_project/model"
	"test_project/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type customerController struct {
	service     service.CustomerService
	UserService service.UserService
}

func NewCustomerController(service service.CustomerService, userService service.UserService) *customerController {
	return &customerController{service, userService}
}

func (controller *customerController) CustomerRoutes(e *echo.Echo) {
	e.Use(middleware.CORS())

	// Customer EP
	var customerRoute = e.Group("/customers")
	customerRoute.Use(middleware.BodyDump(Logger))
	customerRoute.POST("/register", controller.CreateCustomer)
	customerRoute.POST("/", controller.CreateCustomer, controller.middlewareCheckAuthAdmin)
	customerRoute.POST("/list", controller.ListCustomers, controller.middlewareCheckAuthAdmin)
	customerRoute.POST("/login", controller.Login)
	customerRoute.POST("/logout/:id", controller.Logout, controller.middlewareCheckAuthCust)
	customerRoute.DELETE("/:id", controller.DeleteCustomer, controller.middlewareCheckAuthAdmin)
	customerRoute.GET("/:id", controller.GetCustomer, controller.middlewareCheckAuth)
	customerRoute.PUT("/:id", controller.UpdateCustomer, controller.middlewareCheckAuth)
}

func (ctrl *customerController) CreateCustomer(c echo.Context) error {
	request := new(model.CreateCustomerRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	} else if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", err.Error()))
	}

	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	err := ctrl.service.CreateCustomer(ctx, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}

func (ctrl *customerController) DeleteCustomer(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	err := ctrl.service.DeleteCustomer(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}

func (ctrl *customerController) GetCustomer(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	data, err := ctrl.service.GetCustomer(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).SetData(data))
}

func (ctrl *customerController) ListCustomers(c echo.Context) error {
	request := new(model.ListCustomerRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	if request.Filter == nil {
		request.Filter = make(map[string]interface{})
	}

	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	data, total, err := ctrl.service.ListCustomer(ctx, *request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).List(data, total))
}

func (ctrl *customerController) Login(c echo.Context) error {
	request := new(model.LoginCustomerRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	} else if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", err.Error()))
	}

	if (request.Email == "" && request.Username == "") || (request.Email != "" && request.Username != "") {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	token, err := ctrl.service.LoginCustomer(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).SetData(token))
}

func (ctrl *customerController) Logout(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	err := ctrl.service.LogoutCustomer(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}

func (ctrl *customerController) UpdateCustomer(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	request := new(model.UpdateCustomerRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	} else if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", err.Error()))
	}

	request.ID = id
	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	err := ctrl.service.UpdateCustomer(ctx, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}
