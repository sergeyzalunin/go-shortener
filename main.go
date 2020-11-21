package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	h "github.com/sergeyzalunin/go-shortener/api"
	mr "github.com/sergeyzalunin/go-shortener/repository/mongodb"
	rr "github.com/sergeyzalunin/go-shortener/repository/redis"
	"github.com/sergeyzalunin/go-shortener/shortener"
	"go.uber.org/zap"
)

const (
	// SERVERPORT env variable should have port of service.
	SERVERPORT = "SERVERPORT"

	// URLDB env variable should have either 'redis' or 'mongo'.
	URLDB = "URLDB"

	// REDISURL env variable should have connection url to redis.
	REDISURL = "REDISURL"

	// REDISTIMEOUT env variable should have connection timeout to redis.
	REDISTIMEOUT = "REDISTIMEOUT"

	// MONGOURL env variable should have connection url to mongodb.
	MONGOURL = "MONGOURL"

	// MONGODB env variable should have mongodb database name.
	MONGODB = "MONGODB"

	// MONGOTIMEOUT env variable should have connection timeout to mongodb.
	MONGOTIMEOUT = "MONGOTIMEOUT"
)

// repo <- service -> serializer -> http

func main() {
	log, _ := zap.NewProduction()

	defer func() {
		_ = log.Sync()
	}()

	repo := chooseRepo(log)
	service := shortener.NewRedirectService(repo)
	handler := h.NewHandler(service, log)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	buffer := 3
	errs := make(chan error, buffer)

	go func() {
		port := httpPort()
		log.Info("Listening on port " + port)
		errs <- http.ListenAndServe(port, r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		val := <-c
		errs <- errors.New(val.String())
	}()

	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	port := os.Getenv(SERVERPORT)
	if port == "" {
		port = "8080"
	}
	return fmt.Sprintf(":%s", port)
}

func chooseRepo(log *zap.Logger) shortener.RedirectRepository {
	switch os.Getenv(URLDB) {
	case "redis":
		repo := getRedisRepository(log)
		return repo
	case "mongo":
		repo := getMongoRepository(log)
		return repo
	}
	return nil
}

func getRedisRepository(log *zap.Logger) shortener.RedirectRepository {
	redisURL := os.Getenv(REDISURL)
	redisTimeout, err := strconv.Atoi(os.Getenv(REDISTIMEOUT))
	if err != nil {
		log.Fatal(err.Error(), zap.Error(err))
	}
	repo, err := rr.NewRedisRepository(redisURL, redisTimeout)
	if err != nil {
		log.Fatal(err.Error(), zap.Error(err))
	}
	return repo
}

func getMongoRepository(log *zap.Logger) shortener.RedirectRepository {
	mongoURL := os.Getenv(MONGOURL)
	mongoDB := os.Getenv(MONGODB)
	mongoTimeout, err := strconv.Atoi(os.Getenv(MONGOTIMEOUT))
	if err != nil {
		log.Fatal(err.Error(), zap.Error(err))
	}

	repo, err := mr.NewMongoRepository(mongoURL, mongoDB, mongoTimeout)
	if err != nil {
		log.Fatal(err.Error(), zap.Error(err))
	}

	return repo
}
