package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

// Human — это JSON-представление человека
type Human struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Status string `json:"status"`
}

func main() {
	http.HandleFunc("/", getHelloPage)
	http.HandleFunc("/user", getUser)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func getHelloPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func getNameByAgeMock(age int) string {
	name := "Сашка"
	if age < 70 {
		name = "Юрий"
	}
	if age < 50 {
		name = "Оксана"
	}
	if age < 40 {
		name = "Афанасий"
	}
	if age < 20 {
		name = "Глебка"
	}
	return name
}

func getUser(w http.ResponseWriter, r *http.Request) {
	oleg := Human{"Олег", 30, "хороший чувак"}
	sashka := Human{"Сашка", 24, "хороший чувак"}
	petr := Human{"Петр", 12, "поросенок"}
	humans := make(map[string]Human)
	humans["0"] = oleg
	humans["1"] = sashka
	humans["2"] = petr

	for i := 3; i < 1000; i++ {
		age := rand.Intn(100)
		ind := strconv.Itoa(i)
		name := getNameByAgeMock(age)

		humans[ind] = Human{name, age, "статус"}
	}

	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}

	finedUser, isFined := humans[keys[0]]

	response := finedUser
	if !isFined {

		http.Error(w, "Нет такого юзера", http.StatusNotFound)

	}
	json, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
