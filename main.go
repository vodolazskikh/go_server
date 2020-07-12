package main

import (
	"log"
	"net/http"

	"3foodServer/routes"
)

func main() {
	http.HandleFunc("/", routes.HomeRoute)
	http.HandleFunc("/user", routes.UserRoute)
	http.HandleFunc("/place", routes.PlaceRoute)
	http.HandleFunc("/reviews", routes.ReviewsRoute)
	http.HandleFunc("/food", routes.FoodRoute)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
