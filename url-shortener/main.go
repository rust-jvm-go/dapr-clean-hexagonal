package main

import (
	"context"
	"fmt"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"html/template"
	"net/http"
	"os"
	"strings"
	"url-shortener/repository/mongodb"
)

func main() {

	ctx := context.Background()

	//////////////////////
	// Setup environment
	//////////////////////
	env, err := InitializeSetup(ctx)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(env.Info())

	repository := *(*env.RedirectRepository).(*mongodb.MongoRepository)
	defer func() {
		fmt.Println("Disconnecting from MongoDB...")
		if err := repository.Client.Disconnect(repository.Ctx); err != nil {
			fmt.Printf("Could not disconnect from MongoDB, err: %v\n", err.Error())
		}
	}()

	// Test repository client
	_, err = repository.Find("Test 1")
	if err != nil {
		if strings.Contains(err.Error(), "cancelled") {
			fmt.Printf("Could not find \"Test 1\", err: %v\n", err.Error())
			os.Exit(1)
		}
	}

	///////////////
	// BunRouter
	///////////////
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware(
			reqlog.WithVerbose(env.InitConfig.BunRouterDebug),
		)),
	).Verbose()

	router.GET("/", indexHandler)

	/*router.WithGroup("/api", func(g *bunrouter.VerboseGroup) {
		g.GET("/users/:id", debugHandler)
		g.GET("/users/current", debugHandler)
		g.GET("/users/*path", debugHandler)
	})*/

	port := fmt.Sprintf(":%s", env.InitConfig.BunRouterPort)
	fmt.Printf("listening on HTTP port %s\n", port)
	fmt.Println(http.ListenAndServe(port, router))
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
