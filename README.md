[![Circle CI](https://circleci.com/gh/apaganobeleno/webo.svg?style=svg)](https://circleci.com/gh/apaganobeleno/webo)

### WEBO

This is a proof-of-concept for a very small framework for web development with **golang**.

  - Routing
    - Route Definition      (GET/POST/PUT/DELETE)
    - Route Blocks                                
    - Route RegExp                                
    - Route Parameters                            
    - Block/Route Middlewares (Before)            
    - Handle 404 Gracefuly                        

  - Server
    - Server construct to make main.go simpler    
    - Serve static files                          
    - Port from $PORT                      
    - Handle Panic

  - CLI
    - Init command
    - Generate handler file
    - List Routes               [TODO]
    - Generate middleware file

#### Principles:

  - Things should be as simple as possible
  - Prefer code-gen instead of reflection
  - Stay as closer to the GO way as possible
  - Return errors instead of panic
  - Test as much as possible


#### Folder structure and files:

To make things simpler we propose the following folder structure.

  - app_folder
    * handlers/         // handlers used by the routes
    * middlewares/      // middlewares used by the routes
    * config            // package for configuration files
      - routes.go       // route definition file
    * main.go           // where the app starts
    * vendor/           // vendor libraries folder (GO15VENDOREXPERIMENT)


##### Handlers
Handlers simply need to meet the `http.HandlerFunc` type which is something like:

```go
import "net/http"

func MyHanlderFunc(rw http.RequestWriter, req *http.Request){
  //DO SOMETHING
}
```

And are stored inside your app `handlers` folder, but could be anywhere.

##### Middlewares
##### Routes.go
You define your app routing inside the `routes.go` file, basic routes can be defined like:

```go
package main
import (
  "app/handlers"
  "github.com/apaganobeleno/webo"
)

func DefineRoutes(wr *webo.Router){
  wr.Get("/home", handlers.Home);

  wr.Group("/users", func( rg *webo.RouteGroup){
    rg.Post("/login", handlers.Login);
    rg.Put("/", handlers.UserUpdate);
  })
}
```
But you can also define more complex routes with `route params` and `route regexp`:

```go
package main
import (
  "app/handlers"
  "github.com/apaganobeleno/webo"
)

func defineRoutes(wr *webo.Router){
  wr.Get("/home/", handlers.Home);

  wr.Group("/tweets", func( rg *webo.RouteGroup){
    rg.Get("/{:id}", handlers.ShowTweet);
    rg.Get("/{id:[0-9]+}/details", handlers.TweetDetails); //This will only accept integer id's
  })
}
```
And you can also define middlewares to be executed __before__ and __after__ your routes:

```go
package main
import (
  "app/middlewares"
  "app/handlers"
  "github.com/apaganobeleno/webo"
)

func defineRoutes(wr webo.Router){
  wr.Get("/secure", handlers.SecureRoute).Before(middlewares.SecurityVerification);
  wr.Post("/logged_route", handlers.LoggedRoute).After(middlewares.LogRequestDetails);
}
```

Note, After handlers happen in random order, and are not chained, that means these will app execute after the route handler.

#### Static files

If you need to serve static files from your app, you can do it by adding:

```go
s.AddStatic("/static_dir", "/path_you_want")
```
And Webo will point that path to the directory specified.

__Note:__ Ensure your static paths does not colapse with a route path, other wise, route has presedence over the static directory.

#### CLI
##### Init
To create a Webo project simply run `webo init myproject` inside the parent folder of your project.

##### Generate middleware file
##### Generate handler file
To generate a handler you can use the CLI, it provides a command that helps with this task.

```bash
$ webo g handler MyHandler
```

#### Upcomming Features

This is a list of the features i would like to add into webo.

  - Database connection
  - Routes documentation
  - Vendor fetch
  - Vendor update
  - Vendor remove
  - After callbacks
  - Testing framework
  - Database Connection support
    runs application tests injecting common things like app instance with the DB.
