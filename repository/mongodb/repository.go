package mongodb

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sergeyzalunin/go-shortener/shortener"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	RedirectTable string = "redirects"
)

// mongoRepository is an implementation of RedirectRepository interface.
type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

// newMongoClient gets the mongo client by provided URL.
// After timeout connecting cancels automatically.
func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, errors.Wrap(err, "repository.mongodb.newMongoClient.Connect")
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, errors.Wrap(err, "repository.mongodb.newMongoClient.Ping")
	}

	return client, nil
}

// NewMongoRepository is a constructor of mongoRepository.
func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (shortener.RedirectRepository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}

	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.mongodb.NewMongoRepository")
	}

	repo.client = client

	return repo, nil
}

// Find looks for shortened url in the mongo db by provided code.
func (m *mongoRepository) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	collection := m.client.Database(m.database).Collection(RedirectTable)
	filter := bson.M{"code": code}
	redirect := &shortener.Redirect{}

	err := collection.FindOne(ctx, filter).Decode(redirect) // ?? почему тут &
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.mongodb.Find")
		}

		return nil, errors.Wrap(err, "repository.mongodb.Find")
	}

	return redirect, nil
}

// Store saves the shortened url in the mongo db.
func (m *mongoRepository) Store(redirect *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	collection := m.client.Database(m.database).Collection(RedirectTable)
	newItem := bson.M{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}
	_, err := collection.InsertOne(ctx, newItem)
	//nolint:gofumt
	if err != nil {
		return errors.Wrap(err, "redirect.mongodb.Store")
	}

	return nil
}
