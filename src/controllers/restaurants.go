package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nvnamsss/eatigo/dtos"
	"github.com/nvnamsss/eatigo/errors"
	"github.com/nvnamsss/eatigo/logger"
	"github.com/nvnamsss/eatigo/services"
)

type RestaurantController struct {
	Base
	restaurantService services.RestaurantService
}

// @Summary Find restaurants around a specific place
// @Description Find restaurants around a specific place
// @Tags Restaurants
// @Accept json
// @Produce json
// @Param values	query dtos.FindRestaurantsRequest	true "query"
// @Success 200 {object} dtos.FindRestaurantsResponse
// @Failure 401 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /restaurants/ [get]
func (h *RestaurantController) Find(c *gin.Context) {
	var (
		req dtos.FindRestaurantsRequest
		res *dtos.FindRestaurantsResponse
		err error
	)

	if err = c.ShouldBindQuery(&req); err != nil {
		logger.Context(c.Request.Context()).Errorf("validation error: %v", err)
		h.HandleError(c, errors.New(errors.ErrInvalidRequest, err.Error()))
		return
	}

	if res, err = h.restaurantService.Find(c.Request.Context(), &req); err != nil {
		h.HandleError(c, err)
		return
	}

	h.JSON(c, res)
}

func NewRestaurantController(restaurantService services.RestaurantService) *RestaurantController {
	return &RestaurantController{
		restaurantService: restaurantService,
	}
}
