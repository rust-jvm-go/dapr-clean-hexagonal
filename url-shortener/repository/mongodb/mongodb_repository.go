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

type mongoRepository struct {
	ctx      context.Context
	client   *mongo.Client
	database *mongo.Database
	timeout  time.Duration
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
	defer func() {
		fmt.Println("Disconnecting from MongoDB")
		if err := mongoClient.Disconnect(ctx); err != nil {
			fmt.Printf("Could not disconnect from MongoDB, err: %v\n", err.Error())
		}
	}()

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	database := mongoClient.Database(initConfig.MongoDatabase)
	r := &mongoRepository{
		ctx:      ctx,
		client:   mongoClient,
		database: database,
		timeout:  timeout,
	}

	// Some initial database test operations
	err = initDBTests(ctx, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *mongoRepository) Find(code string) (*domain.Redirect, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	redirect := &domain.Redirect{}

	collection := r.database.Collection("redirects")

	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(domain.ErrRedirectNotFound, "repository.Redirect.Find")
		}
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	return redirect, nil
}

func (r *mongoRepository) Store(redirect *domain.Redirect) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	collection := r.database.Collection("redirects")
	result, err := collection.InsertOne(
		ctx,
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

func (r *mongoRepository) Info() string {
	return fmt.Sprintf("MongoDB database = %s", r.database.Name())
}

func initDBTests(ctx context.Context, r *mongoRepository) error {
	collectionNames, err := r.database.ListCollectionNames(ctx, bson.D{}) // List all collections
	if err != nil {
		fmt.Printf("Could not list collection names, err: %v\n", err.Error())
		return err
	}
	fmt.Printf("collection size = %d\n", len(collectionNames))
	for _, collectionName := range collectionNames {
		fmt.Printf("collection = %s\n", collectionName)
	}

	collection := r.database.Collection("test_collection")
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Printf("Error finding documents, err: %v\n", err)
		return err
	}

	var documents []bson.M
	if err = cursor.All(ctx, &documents); err != nil {
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
