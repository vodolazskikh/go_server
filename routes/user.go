package routes

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

// UserRoute - роут для get/post пользователя
func UserRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getUser(w, r)
	}

	if r.Method == "POST" {
		postUser(w, r)
	}
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

func postUser(w http.ResponseWriter, r *http.Request) {

	namesSlugs, nameOk := r.URL.Query()["name"]
	ageSlugs, ageOk := r.URL.Query()["age"]
	statusSlugs, statusOk := r.URL.Query()["status"]
	fmt.Println(namesSlugs)
	if !nameOk || len(namesSlugs[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	if !ageOk || len(ageSlugs[0]) < 1 {
		log.Println("Url Param 'age' is missing")
		return
	}

	if !statusOk || len(statusSlugs[0]) < 1 {
		log.Println("Url Param 'status' is missing")
		return
	}

	name := namesSlugs[0]
	age, _ := strconv.Atoi(ageSlugs[0])
	status := statusSlugs[0]
	newUser := Human{name, age, status}
	json, err := json.Marshal(newUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
