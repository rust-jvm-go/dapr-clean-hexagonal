package main

import (
	"context"
	"fmt"
	daprSdkClient "github.com/dapr/go-sdk/client"
	json "github.com/json-iterator/go"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	daprConfigStore      = "configstore"
	bunRouterDebugConfig = "BUNROUTER_DEBUG"
	mongoDBURIConfig     = "MONGODB_URI"
)

func main() {

	////////////////
	// Dapr client
	////////////////

	daprClient, err := daprSdkClient.NewClient()
	if err != nil {
		log.Panic(err)
	}

	ctx := context.Background()

	// Get config items from config store.
	var configurationItems = []string{bunRouterDebugConfig, mongoDBURIConfig}
	configItems, err := daprClient.GetConfigurationItems(ctx, daprConfigStore, configurationItems)
	if err != nil {
		fmt.Println("Could not get config item, err: ", err.Error())
		os.Exit(1)
	}
	debug, ok := configItems[bunRouterDebugConfig]
	if !ok {
		fmt.Printf("Could not get config item %s, err: %v\n", bunRouterDebugConfig, err.Error())
		os.Exit(1)
	}

	bunRouterDebug, err := strconv.ParseBool(debug.Value)
	if err != nil {
		fmt.Printf("Could not parse config item %s, err: %v\n", bunRouterDebugConfig, err.Error())
		os.Exit(1)
	}
	fmt.Printf("Configuration item \"%s\" = %v\n", bunRouterDebugConfig, bunRouterDebug)
	mongoDBURI, ok := configItems[mongoDBURIConfig]
	if !ok {
		fmt.Printf("Could not get config item %s, err: %v\n", bunRouterDebugConfig, err.Error())
		os.Exit(1)
	}
	fmt.Printf("Configuration item \"%s\" = %s\n", mongoDBURIConfig, mongoDBURI.Value)

	////////////////////
	// MongoDB client
	////////////////////

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBURI.Value))
	if err != nil {
		fmt.Printf("Could not connect to MongoDB, err: %v\n", err.Error())
		log.Panic(err)
	}
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			fmt.Printf("Could not disconnect from MongoDB, err: %v\n", err.Error())
			log.Panic(err)
		}
	}()

	// Database and collection read operations
	database := mongoClient.Database("local_db")
	fmt.Printf("database = %s\n", database.Name())
	collectionNames, err := database.ListCollectionNames(ctx, bson.D{}) // List all collections
	if err != nil {
		fmt.Printf("Could not list collection names, err: %v\n", err.Error())
	}
	fmt.Printf("collection size = %d\n", len(collectionNames))
	for _, collectionName := range collectionNames {
		fmt.Printf("collection = %s\n", collectionName)
	}

	collection := database.Collection("test_collection")
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Printf("Error finding documents, err: %v\n", err)
	}
	var documents []bson.M
	if err = cursor.All(ctx, &documents); err != nil {
		fmt.Printf("Error getting documents, err: %v\n", err)
	}
	jsonDocs, err := json.MarshalIndent(documents, "", "    ")
	if err != nil {
		fmt.Printf("Error marshaling documents, err: %v\n", err)
	}
	fmt.Printf("jsonDocs = %s\n", jsonDocs)

	///////////////
	// BunRouter
	///////////////

	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware(
			reqlog.WithVerbose(bunRouterDebug),
		)),
	).Verbose()

	router.GET("/", indexHandler)

	/*router.WithGroup("/api", func(g *bunrouter.VerboseGroup) {
		g.GET("/users/:id", debugHandler)
		g.GET("/users/current", debugHandler)
		g.GET("/users/*path", debugHandler)
	})*/

	fmt.Println("listening on http://localhost:9999")
	fmt.Println(http.ListenAndServe(":9999", router))
}

func indexHandler(w http.ResponseWriter, _ *http.Request, _ bunrouter.Params) {
	_ = indexTemplate().Execute(w, nil)
}

func indexTemplate() *template.Template {
	return template.Must(template.New("index").Parse(indexTmpl))
}

var indexTmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Welcome!</title>
</head>
<body>
  <h1>Hello from BunRouter!</h1>
  <ul>
    <li><a href="/api/users/123">/api/users/123</a></li>
    <li><a href="/api/users/current">/api/users/current</a></li>
    <li><a href="/api/users/foo/bar">/api/users/foo/bar</a></li>
  </ul>
</body>
</html>
`
