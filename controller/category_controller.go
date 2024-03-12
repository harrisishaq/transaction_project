package controller

import (
	"net/http"
	"test_project/model"
	"test_project/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type categoryController struct {
	Service     service.CategoryService
	UserService service.UserService
}

func NewCategoryController(service service.CategoryService, userService service.UserService) *categoryController {
	return &categoryController{
		Service:     service,
		UserService: userService,
	}
}

func (controller *categoryController) CategoryRoutes(e *echo.Echo) {
	e.Use(middleware.CORS())

	// Category EP
	var categoryRoute = e.Group("/categories")
	categoryRoute.POST("/", controller.CreateCategory, controller.middlewareCheckAuthAdmin)
	categoryRoute.POST("/list", controller.ListCategory)
	categoryRoute.DELETE("/:id", controller.DeleteCategory, controller.middlewareCheckAuthAdmin)
	categoryRoute.GET("/:id", controller.GetCategory, controller.middlewareCheckAuthAdmin)
	categoryRoute.PUT("/:id", controller.UpdateCategory, controller.middlewareCheckAuthAdmin)
}

func (ctrl *categoryController) CreateCategory(c echo.Context) error {
	request := new(model.CreateCategoryRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	} else if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", err.Error()))
	}

	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	err := ctrl.Service.CreateCategory(ctx, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}

func (ctrl *categoryController) DeleteCategory(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	err := ctrl.Service.DeleteCategory(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}

func (ctrl *categoryController) GetCategory(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	data, err := ctrl.Service.GetCategory(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).SetData(data))
}

func (ctrl *categoryController) ListCategory(c echo.Context) error {
	request := new(model.ListCategoryRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	if request.Filter == nil {
		request.Filter = make(map[string]interface{})
	}

	data, total, err := ctrl.Service.ListCategory(*request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true).List(data, total))
}

func (ctrl *categoryController) UpdateCategory(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	}

	request := new(model.UpdateCategoryRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", "Bad Request"))
	} else if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewJsonResponse(false).SetError("400", err.Error()))
	}

	request.ID = id
	var ctx = model.SetUserContext(c.Request().Context(), c.Get("userCtx"))
	err := ctrl.Service.UpdateCategory(ctx, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.(*model.JsonResponse))
	}

	return c.JSON(http.StatusOK, model.NewJsonResponse(true))
}
