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

// Place — так будем хранить JSON заведения в БД
type Place struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Rating int    `json:"rating"`
	City   string `json:"city"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// PlaceRoute - роуты для get/post заведения
func PlaceRoute(w http.ResponseWriter, r *http.Request) {

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

	placesCollection := client.Database("test").Collection("places")

	if r.Method == "GET" {
		getPlace(placesCollection, w, r)
	}

	if r.Method == "POST" {
		postPlace(placesCollection, w, r)
	}

}

func getPlace(collection *mongo.Collection, w http.ResponseWriter, r *http.Request) {
	keysID, ok := r.URL.Query()["id"]
	keysCity, cityOk := r.URL.Query()["city"]

	noID := !ok || len(keysID[0]) < 1
	noCity := !cityOk || len(keysCity[0]) < 1

	if noID && noCity {
		log.Println("Url Param 'id' or 'city' is missing")
		return
	}
	// Если есть айдишник и нет города - вернем место byId
	if !noID && noCity {
		filter := bson.D{primitive.E{Key: "id", Value: keysID[0]}}

		var result Place
		err := collection.FindOne(context.TODO(), filter).Decode(&result)

		fmt.Println(result)

		if err != nil {
			http.Error(w, "Нет такого заведения", http.StatusNotFound)
		}

		json, err := json.Marshal(result)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}

	// Если есть город и нет айдишника - вернем все места для города
	if !noCity {
		filter := bson.D{primitive.E{Key: "city", Value: keysCity[0]}}
		options := options.Find()

		var results []*Place
		cur, err := collection.Find(context.TODO(), filter, options)

		for cur.Next(context.TODO()) {

			// create a value into which the single document can be decoded
			var elem Place
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}

			results = append(results, &elem)
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		// Close the cursor once finished
		cur.Close(context.TODO())

		if err != nil {
			http.Error(w, "Нет таких заведений", http.StatusNotFound)
		}

		json, err := json.Marshal(results)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}

func postPlace(collection *mongo.Collection, w http.ResponseWriter, r *http.Request) {
	nameSlugs, nameOk := r.URL.Query()["name"]
	ratingSlug, ratingOk := r.URL.Query()["rating"]
	citySlugs, cityOk := r.URL.Query()["city"]

	if !nameOk || len(nameSlugs[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	if !ratingOk || len(ratingSlug[0]) < 1 {
		log.Println("Url Param 'rating' is missing")
		return
	}

	if !cityOk || len(citySlugs[0]) < 1 {
		log.Println("Url Param 'city' is missing")
		return
	}

	id := utils.GenerateUUID()
	name := nameSlugs[0]
	rating, _ := strconv.Atoi(ratingSlug[0])
	city := citySlugs[0]
	newPlace := Place{id, name, rating, city}
	json, err := json.Marshal(newPlace)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	insertResult, err := collection.InsertOne(context.TODO(), newPlace)

	if err != nil {
		panic(err)
	}
	fmt.Println(insertResult.InsertedID)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
