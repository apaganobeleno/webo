package cmd

const mainTemplate = `package main
import (
  "log"
  webo "github.com/apaganobeleno/webo/server"
  "{{.Name}}/config"
)

func main(){
	s := webo.NewServer(config.DefineRoutes)
  s.Start("8080")
	defer log.Println("| {{.Name}} Closing")
}
`

var routesTemplate = `package config
import (
	"{{.Name}}/handlers"
	"github.com/apaganobeleno/webo/routing"
)

//defineRoutes gets called by the generated main.go, please define routes inside it.
func DefineRoutes(r *routing.Router){
  r.Get("/", handlers.Home)
}`

var handlersTemplate = `package handlers
import "net/http"
//Home renders a Hello message
func Home(rw http.ResponseWriter, req *http.Request){
  rw.Write([]byte("Hello From Webo!"))
}`

var handlerGenTemplate = `package handlers
import "net/http"
//{{.Name}}
func {{.Name}}(rw http.ResponseWriter, req *http.Request){
  rw.Write([]byte("{{.Name}} Handler"))
}
`

var middlewareGenTemplate = `package middlewares
import "net/http"
//{{.Name}}
func {{.Name}}(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc){
  next(rw, req)
}
`
