package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"3foodServer/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// User — так будем хранить JSON юзера в БД
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	City  string `json:"city"`
	Email string `json:"email"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// UserRoute - роут для get/post пользователя
func UserRoute(w http.ResponseWriter, r *http.Request) {

	username, _ := os.LookupEnv("DB_USER")

	password, _ := os.LookupEnv("DB_PASSWORD")

	clusterAddress, _ := os.LookupEnv("DB_ADDRESS")

	uriStrings := []string{"mongodb+srv://", username, ":", password, "@", clusterAddress, "/test?w=majority"}
	uri := strings.Join(uriStrings, "")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	usersCollection := client.Database("test").Collection("users")

	if r.Method == "GET" {
		getUser(usersCollection, w, r)
	}

	if r.Method == "POST" {
		postUser(usersCollection, w, r)
	}

}

func getUser(collection *mongo.Collection, w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}

	filter := bson.D{primitive.E{Key: "id", Value: keys[0]}}

	var result User
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	fmt.Println(result)

	if err != nil {
		http.Error(w, "Нет такого юзера", http.StatusNotFound)
	}

	json, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func postUser(collection *mongo.Collection, w http.ResponseWriter, r *http.Request) {
	namesSlugs, nameOk := r.URL.Query()["name"]
	ageSlugs, ageOk := r.URL.Query()["age"]
	citySlugs, cityOk := r.URL.Query()["city"]
	emailSlugs, emailOk := r.URL.Query()["email"]

	if !nameOk || len(namesSlugs[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	if !ageOk || len(ageSlugs[0]) < 1 {
		log.Println("Url Param 'age' is missing")
		return
	}

	if !cityOk || len(citySlugs[0]) < 1 {
		log.Println("Url Param 'city' is missing")
		return
	}

	if !emailOk || len(emailSlugs[0]) < 1 {
		log.Println("Url Param 'email' is missing")
		return
	}

	id := utils.GenerateUUID()
	name := namesSlugs[0]
	age, _ := strconv.Atoi(ageSlugs[0])
	city := citySlugs[0]
	email := emailSlugs[0]
	newUser := User{id, name, age, city, email}
	json, err := json.Marshal(newUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	insertResult, err := collection.InsertOne(context.TODO(), newUser)

	if err != nil {
		panic(err)
	}
	fmt.Println(insertResult.InsertedID)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
