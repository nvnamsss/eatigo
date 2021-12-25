package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nvnamsss/eatigo/adapters/cache"
	"github.com/nvnamsss/eatigo/adapters/google_api"
	_ "github.com/nvnamsss/eatigo/cmd/docs"
	"github.com/nvnamsss/eatigo/configs"
	"github.com/nvnamsss/eatigo/controllers"
	"github.com/nvnamsss/eatigo/logger"
	"github.com/nvnamsss/eatigo/repositories"
	"github.com/nvnamsss/eatigo/services"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.elastic.co/apm/module/apmgin"
	"go.elastic.co/apm/module/apmgoredisv8"
)

// @title Eatigo
// @version 1.0
// @description Eatigo API documentation
// @termsOfService http://swagger.io/terms/

// @contact.name Nam Nguyen
// @contact.email nvnam.c@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /eatigo/v1
func main() {
	redisOpt := &redis.Options{
		Addr:     configs.Config.Redis.URL(),
		Password: configs.Config.Redis.Password,
		DB:       configs.Config.Redis.Database,
	}

	redisClient := redis.NewClient(redisOpt)
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatalf(err, "Creating connection to redis: %v", err)
	}
	redisClient.AddHook(apmgoredisv8.NewHook())

	var r = gin.Default()
	corsConfig := cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "TimezoneOffset"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})

	var (
		placeAdapter = google_api.NewGooglePlace()
		cacheAdaper  = cache.NewCacheAdapter(redisClient)
	)

	var (
		restaurantRepository = repositories.NewRestaurantRepository(placeAdapter, cacheAdaper)
	)

	var (
		restaurantService = services.NewRestaurantService(restaurantRepository)
	)

	var (
		restaurantController = controllers.NewRestaurantController(restaurantService)
	)

	r.Use(
		corsConfig,
		apmgin.Middleware(r),
	)

	v1 := r.Group("/eatigo/v1")
	{
		restaurant := v1.Group("/restaurants")
		{
			restaurant.GET("/", restaurantController.Find)
		}
	}

	if configs.Config.RunMode == gin.DebugMode && configs.Config.Env != "PRODUCTION" {
		r.GET("/eatigo/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	server := &http.Server{
		Addr:    configs.Config.AddressListener(),
		Handler: r,
	}

	go func() {
		logger.Infof("Starting Server on %v", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf(err, "Opening HTTP server: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logger.Infof("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Shutdown error: %v", err)
	}

	os.Exit(0)
}

func init() {
	if _, err := configs.New(); err != nil {
		os.Exit(1)
	}
}
