package main

import (
	"log"
	"net/http"

	"3foodServer/routes"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/food", routes.FoodRoute)
	mux.HandleFunc("/", routes.HomeRoute)
	mux.HandleFunc("/user", routes.UserRoute)
	mux.HandleFunc("/place", routes.PlaceRoute)
	mux.HandleFunc("/reviews", routes.ReviewsRoute)
	log.Println("Listening...")
	handler := cors.AllowAll().Handler(mux)
	http.ListenAndServe(":3000", handler)
}
