package main

import (
	"log"
	"net/http"

	"3foodServer/routes"
)

func main() {
	http.HandleFunc("/", routes.HomeRoute)
	http.HandleFunc("/user", routes.UserRoute)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
