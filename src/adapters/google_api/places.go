package google_api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-kit/kit/endpoint"
	http_transport "github.com/go-kit/kit/transport/http"
	"github.com/nvnamsss/eatigo/configs"
	"github.com/nvnamsss/eatigo/logger"
)

const (
	dataPath = "./data.json"
)

type GooglePlace interface {
	FindRestaurants(ctx context.Context, req *FindRestaurantsRequest) (*FindRestaurantsResponse, error)
}

type googlePlace struct {
	findEndpoint endpoint.Endpoint
}

func NewGooglePlace() GooglePlace {
	return &googlePlace{
		findEndpoint: findEndPoint(context.TODO()),
	}
}

func (h *googlePlace) FindRestaurants(ctx context.Context, req *FindRestaurantsRequest) (*FindRestaurantsResponse, error) {
	if h.findEndpoint == nil {
		logger.Context(ctx).Errorf("find endpoint is nil")
		return nil, errors.New("invalid find endpoint config")
	}
	res, err := h.findEndpoint(ctx, req)
	if err != nil {
		logger.Context(ctx).Errorf("call to google place error %v with resp %#v", err, res)
		return nil, errors.New("invalid response")
	}
	r, ok := res.(*FindRestaurantsResponse)
	if !ok {
		logger.Context(ctx).Errorf("cast response invalid %#v", res)
		return nil, errors.New("invalid cast")
	}
	if r == nil {
		logger.Context(ctx).Infof("call value in group has empty response")
		return nil, errors.New("invalid response")
	}
	if r.Meta.Code != http.StatusOK {
		logger.Context(ctx).Infof("call value in group unsuccessfully")
	}

	return r, nil
}

func findEndPoint(ctx context.Context) endpoint.Endpoint {
	fullPath := configs.Config.GooglePlace.BaseURL + configs.Config.GooglePlace.Endpoint
	parserPath, err := url.Parse(fullPath)
	if err != nil {
		logger.Context(ctx).Errorf("parser full url for config %v has error : %v", configs.Config.GooglePlace)
		return nil
	}
	return http_transport.NewClient(
		http.MethodGet, parserPath, encodeFind, decodeFind,
		http_transport.SetClient(&http.Client{
			Timeout: 30 * time.Second,
		}),
		http_transport.ClientBefore(),
	).Endpoint()
}

func encodeFind(ctx context.Context, r *http.Request, req interface{}) error {
	request, ok := req.(*FindRestaurantsRequest)
	if !ok {
		return errors.New("missing request config")
	}

	query := r.URL.Query()

	query.Add("query", url.QueryEscape(fmt.Sprintf("restaurants in %v", request.Place)))
	query.Add("radius", strconv.Itoa(request.Radius))
	query.Add("key", configs.Config.GooglePlace.Key)
	query.Add("pagetoken", request.NextPageToken)
	query.Add("type", "restaurant")
	r.URL.RawQuery = query.Encode()

	return nil
}

func decodeFind(ctx context.Context, r *http.Response) (interface{}, error) {
	var (
		res   FindRestaurantsResponse
		adapt findRestaurantsDataAdapt
	)

	if err := json.NewDecoder(r.Body).Decode(&adapt); err != nil {
		logger.Context(ctx).Errorf("decode response has error : %v", err)
		return nil, err
	}

	res.Meta.Code = r.StatusCode
	res.Meta.Message = adapt.Status
	res.Meta.NextPageToken = adapt.NextPageToken
	res.Data = adapt.Results
	return &res, nil
}

type googlePlaceMock struct {
}

func (h *googlePlaceMock) FindRestaurants(ctx context.Context, req *FindRestaurantsRequest) (*FindRestaurantsResponse, error) {
	data, err := ioutil.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}

	var (
		res   FindRestaurantsResponse
		adapt findRestaurantsDataAdapt
	)

	_ = json.Unmarshal(data, &adapt)
	res.Meta.Code = http.StatusOK
	res.Meta.Message = adapt.Status
	res.Meta.NextPageToken = adapt.NextPageToken
	res.Data = adapt.Results

	return &res, nil
}

func NewGooglePlaceMock() GooglePlace {
	return &googlePlaceMock{}
}
