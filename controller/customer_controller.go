package controller

import (
	"net/http"
	"test_project/model"
	"test_project/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type customerController struct {
	service service.CustomerService
}

func NewCustomerController(service service.CustomerService) *customerController {
	return &customerController{service}
}

func (controller *customerController) CustomerRoutes(e *echo.Echo) {
	e.Use(middleware.CORS())

	// Customer EP
	var customerRoute = e.Group("/customers")
	customerRoute.Use(middleware.BodyDump(Logger))
	customerRoute.POST("/", controller.CreateCustomer)
	customerRoute.POST("/list", controller.ListCustomers)
	customerRoute.POST("/login", controller.Login)
	customerRoute.DELETE("/:id", controller.DeleteCustomer)
	customerRoute.GET("/:id", controller.GetCustomer)
	customerRoute.PUT("/:id", controller.UpdateCustomer)
}

func (ctrl *customerController) CreateCustomer(c echo.Context) error {
	request := new(model.CreateCustomerRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	} else if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", err.Error()))
	}

	err := ctrl.service.CreateCustomer(request)
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

	err := ctrl.service.DeleteCustomer(id)
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

	data, err := ctrl.service.GetCustomer(id)
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

	data, total, err := ctrl.service.ListCustomer(*request)
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

	token, err := ctrl.service.LoginCustomer(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).SetData(token))
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
	err := ctrl.service.UpdateCustomer(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}
