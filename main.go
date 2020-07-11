package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Response — это JSON-представление сообщения
type Response struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Status string `json:"status"`
}

func main() {
	http.HandleFunc("/", postHandler)
	http.HandleFunc("/hi", getHi)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func getHi(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["name"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	name := keys[0]

	age := 18

	status := "Хороший человек"

	if name == "Артём" {
		age = 29
		status = "Пёс"

	}

	response := Response{
		Name:   name,
		Age:    age,
		Status: status,
	}
	kek, _ := json.Marshal(response)
	w.Write(kek)
}
