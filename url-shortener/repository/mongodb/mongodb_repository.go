package mongodb

import (
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	json "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
	"url-shortener/domain"
	"url-shortener/xsetup"
)

var redirects = "redirects"

type MongoRepository struct {
	Ctx      context.Context
	Client   *mongo.Client
	Database *mongo.Database
	Timeout  time.Duration
}

func NewMongoRepository(initConfig xsetup.InitConfig) (domain.IRedirectRepository, error) {
	timeout := time.Duration(initConfig.MongoDBTimeout) * time.Second
	ctx, cancel := context.WithTimeout(initConfig.Ctx, timeout)
	defer cancel()

	mongoDBURI := fmt.Sprintf("%s&authSource=%s", initConfig.MongoDBURI, initConfig.MongoDatabase)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBURI))
	if err != nil {
		fmt.Printf("Could not connect to MongoDB, err: %v\n", err.Error())
		return nil, err
	}

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	database := mongoClient.Database(initConfig.MongoDatabase)
	r := &MongoRepository{
		Ctx:      initConfig.Ctx,
		Client:   mongoClient,
		Database: database,
		Timeout:  timeout,
	}

	// Some initial Database test operations
	err = initDBTests(r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *MongoRepository) Find(code string) (*domain.Redirect, error) {
	redirect := &domain.Redirect{}

	collection := r.Database.Collection(redirects)

	filter := bson.M{"code": code}
	err := collection.FindOne(r.Ctx, filter).Decode(&redirect)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.Wrap(domain.ErrRedirectNotFound, "repository.Redirect.Find")
		}
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	return redirect, nil
}

func (r *MongoRepository) Store(redirect *domain.Redirect) error {
	collection := r.Database.Collection(redirects)
	result, err := collection.InsertOne(
		r.Ctx,
		bson.M{
			"code":       redirect.Code,
			"url":        redirect.URL,
			"created_at": redirect.CreatedAt,
		},
	)
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}

	fmt.Printf("inserted ID = %v\n", result.InsertedID)

	return nil
}

func initDBTests(r *MongoRepository) error {
	collectionNames, err := r.Database.ListCollectionNames(r.Ctx, bson.D{}) // List all collections
	if err != nil {
		fmt.Printf("Could not list collection names, err: %v\n", err.Error())
		return err
	}
	fmt.Printf("collection size = %d\n", len(collectionNames))
	for _, collectionName := range collectionNames {
		fmt.Printf("collection = %s\n", collectionName)
	}

	collection := r.Database.Collection("test_collection")
	cursor, err := collection.Find(r.Ctx, bson.D{})
	if err != nil {
		fmt.Printf("Error finding documents, err: %v\n", err)
		return err
	}

	var documents []bson.M
	if err = cursor.All(r.Ctx, &documents); err != nil {
		fmt.Printf("Error getting documents, err: %v\n", err)
		return err
	}
	_, err = json.MarshalIndent(documents, "", "    ")
	if err != nil {
		fmt.Printf("Error marshaling documents, err: %v\n", err)
		return err
	}
	// fmt.Printf("jsonDocs = %s\n", jsonDocs)

	return nil
}
