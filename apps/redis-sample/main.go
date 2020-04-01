package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/bitrise-devops-microservice-common/common/entry"
	redis "github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const (
	appName = "redis-sample"
	version = "v0.0.1"

	redisCounterKey = "counter"
)

var (
	gRedisClient *redis.Client
	gConfig      Config
)

// Config ...
type Config struct {
	Port          string `envconfig:"PORT" required:"true" default:"8182"`
	DeployVersion string `envconfig:"DEPLOY_VERSION"`
	RedisURL      string `envconfig:"REDIS_URL" required:"true"`
}

func createRedisClient() (*redis.Client, error) {
	opt, err := redis.ParseURL(gConfig.RedisURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	client := redis.NewClient(opt)

	pong, err := client.Ping().Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println("Redis Ping:", pong)

	return client, err
}

func handleIncrement(w http.ResponseWriter, r *http.Request) error {
	result, err := gRedisClient.Incr(redisCounterKey).Result()
	if err != nil {
		return errors.WithStack(err)
	}
	log.Println("value:", result)

	return errors.WithStack(httpresponse.RespondWithSuccess(w, map[string]string{
		"value": fmt.Sprintf("%d", result),
	}))
}

func handleDecrement(w http.ResponseWriter, r *http.Request) error {
	result, err := gRedisClient.Decr(redisCounterKey).Result()
	if err != nil {
		return errors.WithStack(err)
	}
	log.Println("value:", result)

	return errors.WithStack(httpresponse.RespondWithSuccess(w, map[string]string{
		"value": fmt.Sprintf("%d", result),
	}))
}

func handleRoot(w http.ResponseWriter, r *http.Request) error {
	return errors.WithStack(httpresponse.RespondWithSuccess(w, map[string]string{
		"message":        "Welcome!",
		"app":            appName,
		"version":        version,
		"deploy.version": gConfig.DeployVersion,
	}))
}

func handleNotFound(w http.ResponseWriter, r *http.Request) error {
	return errors.WithStack(httpresponse.RespondWithSuccess(w, map[string]string{
		"message":        fmt.Sprintf("Not found [%s] %s", r.Method, r.RequestURI),
		"app":            appName,
		"version":        version,
		"deploy.version": gConfig.DeployVersion,
	}))
}

func mainE() error {
	// configs
	if err := envconfig.Process("", &gConfig); err != nil {
		return errors.Wrap(err, "configuration error")
	}

	{
		redisClient, err := createRedisClient()
		if err != nil {
			return errors.WithStack(err)
		}
		gRedisClient = redisClient
	}

	// Setup routing
	r := mux.NewRouter().StrictSlash(true)
	middlewareProvider := NewMiddlewareProvider(appName, version)
	//
	r.Handle("/+", middlewareProvider.CommonMiddleware().Then(
		httpresponse.InternalErrHandlerFuncAdapter(handleIncrement))).Methods("GET")
	r.Handle("/-", middlewareProvider.CommonMiddleware().Then(
		httpresponse.InternalErrHandlerFuncAdapter(handleDecrement))).Methods("GET")
	r.Handle("/", middlewareProvider.CommonMiddleware().Then(
		httpresponse.InternalErrHandlerFuncAdapter(handleRoot))).Methods("GET")
	//
	r.NotFoundHandler = middlewareProvider.CommonMiddleware().Then(
		httpresponse.InternalErrHandlerFuncAdapter(handleNotFound))
	http.Handle("/", r)

	// Start the server
	log.Printf("Starting (on port: %s) ...", gConfig.Port)
	return errors.WithStack(http.ListenAndServe(":"+gConfig.Port, nil))
}

func main() {
	entry.HandleMainE(mainE)
}
