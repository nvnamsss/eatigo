package services

import (
	"context"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/nvnamsss/eatigo/dtos"
	"github.com/nvnamsss/eatigo/errors"
	"github.com/nvnamsss/eatigo/logger"
	"github.com/nvnamsss/eatigo/repositories"
)

type RestaurantService interface {
	Find(ctx context.Context, req *dtos.FindRestaurantsRequest) (*dtos.FindRestaurantsResponse, error)
}

type restaurantService struct {
	restaurantRepository repositories.RestaurantRepository
}

func (s *restaurantService) Find(ctx context.Context, req *dtos.FindRestaurantsRequest) (*dtos.FindRestaurantsResponse, error) {
	restaurants, err := s.restaurantRepository.Find(ctx, req)
	if err != nil {
		logger.Context(ctx).Errorf("find restaurants got error: %v", err)
		return nil, errors.New(errors.ErrInternalServer, err.Error())
	}

	var data []*dtos.FindRestaurantsData
	for _, v := range restaurants {
		d := dtos.FindRestaurantsData{}
		_ = copier.Copy(&d, v)
		data = append(data, &d)
	}

	return &dtos.FindRestaurantsResponse{
		Meta: dtos.Meta{
			Code:    http.StatusOK,
			Message: "OK",
			Cursor:  req.Cursor,
		},
		Data: data,
	}, nil
}

func NewRestaurantService(restaurantRepository repositories.RestaurantRepository) RestaurantService {
	return &restaurantService{
		restaurantRepository: restaurantRepository,
	}
}
