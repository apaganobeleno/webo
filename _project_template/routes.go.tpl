package main
import "{.package_root}/handlers"

//defineRoutes gets called by the generated main.go, please define routes inside it.
func defineRoutes(r *Router){
  r.Get("/", handlers.Home)
}
