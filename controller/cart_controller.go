package controller

import (
	"net/http"
	"test_project/model"
	"test_project/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type cartController struct {
	Service         service.CartService
	CustomerService service.CustomerService
}

func NewCartController(service service.CartService, customerService service.CustomerService) *cartController {
	return &cartController{
		Service:         service,
		CustomerService: customerService,
	}
}

func (controller *cartController) CartRoutes(e *echo.Echo) {
	e.Use(middleware.CORS())

	// User EP
	var userRoute = e.Group("/carts")
	userRoute.Use(middleware.BodyDump(Logger))
	userRoute.GET("/", controller.Get, controller.middlewareCheckAuthCust)
	userRoute.POST("/add", controller.AddItem, controller.middlewareCheckAuthCust)
}

func (ctrl *cartController) AddItem(c echo.Context) error {
	request := new(model.AddItemCartRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	} else if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", err.Error()))
	}

	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	err := ctrl.Service.AddItemCart(ctx, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}

func (ctrl *cartController) Get(c echo.Context) error {
	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	data, err := ctrl.Service.GetCart(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).SetData(data))
}
