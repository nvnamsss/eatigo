package repositories

import (
	"context"
	"net/http"
	"time"

	"github.com/nvnamsss/eatigo/adapters/cache"
	"github.com/nvnamsss/eatigo/adapters/google_api"
	"github.com/nvnamsss/eatigo/dtos"
	"github.com/nvnamsss/eatigo/logger"
	"github.com/nvnamsss/eatigo/models"
	"github.com/nvnamsss/eatigo/utils"
)

type RestaurantRepository interface {
	Find(ctx context.Context, req *dtos.FindRestaurantsRequest) ([]*models.Restaurant, error)
}

type restaurantRepository struct {
	placeAdapter google_api.GooglePlace
	cacheAdapter cache.CacheAdapter
}

func (r *restaurantRepository) Find(ctx context.Context, req *dtos.FindRestaurantsRequest) ([]*models.Restaurant, error) {
	var (
		rs  []*models.Restaurant
		res *google_api.FindRestaurantsResponse = &google_api.FindRestaurantsResponse{}
		err error
	)

	if err = r.cacheAdapter.Get(ctx, utils.GetRestaurantKey(req.Place, req.Cursor), res); err != nil {
		if res, err = r.placeAdapter.FindRestaurants(ctx, &google_api.FindRestaurantsRequest{
			Place:         req.Place,
			Radius:        2000,
			NextPageToken: req.Cursor,
		}); err != nil {
			return nil, err
		}

		if res.Meta.Code != http.StatusOK {
			return nil, err
		}

		if err = r.cacheAdapter.Set(ctx, utils.GetRestaurantKey(req.Place, req.Cursor), res, time.Hour*72); err != nil {
			logger.Context(ctx).Errorf("set cache error: %v", err)
		}
	}

	for _, v := range res.Data {
		m := models.Restaurant{
			Name:    v.Name,
			Address: v.Address,
		}
		rs = append(rs, &m)
	}
	req.Cursor = res.Meta.NextPageToken
	return rs, nil
}

func NewRestaurantRepository(placeAdapter google_api.GooglePlace, cacheAdapter cache.CacheAdapter) RestaurantRepository {
	return &restaurantRepository{
		placeAdapter: placeAdapter,
		cacheAdapter: cacheAdapter,
	}
}
