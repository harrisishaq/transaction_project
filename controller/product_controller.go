package controller

import (
	"log"
	"net/http"
	"test_project/model"
	"test_project/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type productController struct {
	service     service.ProductService
	userService service.UserService
}

func NewProductController(service service.ProductService, userService service.UserService) *productController {
	return &productController{
		service:     service,
		userService: userService,
	}
}

func (controller *productController) ProductRoutes(e *echo.Echo) {
	e.Use(middleware.CORS())

	// Product EP
	var productRoute = e.Group("/products")
	productRoute.Use(middleware.BodyDump(Logger))
	productRoute.POST("/", controller.CreateProduct, controller.middlewareCheckAuthAdmin)
	productRoute.POST("/list", controller.ListProduct)
	productRoute.DELETE("/:id", controller.DeleteProduct, controller.middlewareCheckAuthAdmin)
	productRoute.GET("/:id", controller.GetProduct)
	productRoute.PUT("/:id", controller.UpdateProduct, controller.middlewareCheckAuthAdmin)
}

func (ctrl *productController) CreateProduct(c echo.Context) error {
	request := new(model.CreateProductRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	} else if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", err.Error()))
	}

	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	err := ctrl.service.CreateProduct(ctx, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}

func (ctrl *productController) DeleteProduct(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	err := ctrl.service.DeleteProduct(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}

func (ctrl *productController) GetProduct(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	data, err := ctrl.service.GetProduct(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).SetData(data))
}

func (ctrl *productController) ListProduct(c echo.Context) error {
	request := new(model.ListProductRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	if request.Filter == nil {
		request.Filter = make(map[string]interface{})
	}

	data, total, err := ctrl.service.ListProduct(*request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).List(data, total))
}

func (ctrl *productController) UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	request := new(model.UpdateProductRequest)
	if err := c.Bind(&request); err != nil {
		log.Println("Error:  ", err)
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	} else if err := c.Validate(request); err != nil {
		log.Println("Error:  ", err)
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", err.Error()))
	}

	request.ID = id
	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	err := ctrl.service.UpdateProduct(ctx, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}
