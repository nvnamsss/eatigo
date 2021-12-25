package configs

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const (
	Prefix = ""
)

var Config AppConfig

type AppConfig struct {
	Host        string `default:"0.0.0.0" envconfig:"HOST"`
	Port        int    `default:"8080" envconfig:"PORT"`
	RunMode     string `default:"debug" envconfig:"RUN_MODE"`
	Env         string `default:"debug" envconfig:"ENV"`
	GooglePlace GooglePlace
	Redis       Redis
}

type GooglePlace struct {
	BaseURL  string `default:"https://maps.googleapis.com" envconfig:"GOOGLE_PLACE_BASE_URL"`
	Endpoint string `default:"/maps/api/place/textsearch/json" envconfig:"GOOGLE_PLACE_ENDPOINT"`
	Key      string `default:"" envconfig:"GOOGLE_PLACE_KEY"`
}

// AddressListener returns address listener of HTTP server.
func (c *AppConfig) AddressListener() string {
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}

type Redis struct {
	Host         string `default:"127.0.0.1" envconfig:"REDIS_HOST"`
	Port         int    `default:"6379" envconfig:"REDIS_PORT"`
	Password     string `default:"" envconfig:"REDIS_PASSWORD"`
	Database     int    `default:"0" envconfig:"REDIS_DB"`
	PoolSize     int    `default:"10" envconfig:"REDIS_POOL_SIZE"`
	MinIdleConns int    `default:"100" envconfig:"REDIS_MIN_IDLE_CONNS"`
}

// URL return redis connection URL.
func (c *Redis) URL() string {
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}

func New() (*AppConfig, error) {
	if err := envconfig.Process(Prefix, &Config); err != nil {
		return nil, err
	}
	return &Config, nil
}
