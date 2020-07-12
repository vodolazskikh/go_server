package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

// Food — так будем хранить JSON блюд
type Food struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Places []string `json:"places"`
	City   string   `json:"city"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// FoodRoute - роут для get/post еды
func FoodRoute(w http.ResponseWriter, r *http.Request) {

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

	foodCollection := client.Database("test").Collection("food")

	if r.Method == "GET" {
		getFood(foodCollection, w, r)
	}

	if r.Method == "POST" {
		postFood(foodCollection, w, r)
	}

	if r.Method == "PATCH" {
		patchFood(foodCollection, w, r)
	}

}

func getFood(collection *mongo.Collection, w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]
	cityKeys, cityOk := r.URL.Query()["city"]

	noID := !ok || len(keys[0]) < 1
	noCity := !cityOk || len(cityKeys[0]) < 1

	if noID && noCity {
		log.Println("Url Param 'id' and 'city' is missing")
		return
	}

	if !noID {
		filter := bson.D{primitive.E{Key: "id", Value: keys[0]}}
		var result Food
		err := collection.FindOne(context.TODO(), filter).Decode(&result)

		if err != nil {
			http.Error(w, "Нет такого блюда", http.StatusNotFound)
		}

		json, err := json.Marshal(result)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
	if !noCity {
		filter := bson.D{primitive.E{Key: "city", Value: cityKeys[0]}}
		options := options.Find()

		var results []*Food
		cur, err := collection.Find(context.TODO(), filter, options)

		for cur.Next(context.TODO()) {

			// create a value into which the single document can be decoded
			var elem Food
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
			http.Error(w, "В этом городе нет еды", http.StatusNotFound)
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

func patchFood(collection *mongo.Collection, w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]
	newPlaceKeys, placeOk := r.URL.Query()["place"]

	noID := !ok || len(keys[0]) < 1
	noPlace := !placeOk || len(newPlaceKeys[0]) < 1

	if noID || noPlace {
		log.Println("Url Param 'id' or 'place' is missing")
		return
	}

	filter := bson.D{primitive.E{Key: "id", Value: keys[0]}}

	var findedFood Food
	err := collection.FindOne(context.TODO(), filter).Decode(&findedFood)

	if err != nil {
		http.Error(w, "Нет такого блюда", http.StatusNotFound)
	}

	newPlaces := append(findedFood.Places[:], newPlaceKeys[0])

	fmt.Println(findedFood, newPlaces)

	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "places", Value: newPlaces}}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	json, err := json.Marshal(updateResult)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func postFood(collection *mongo.Collection, w http.ResponseWriter, r *http.Request) {
	titleSlugs, titleOk := r.URL.Query()["title"]
	placeSlugs, placeOk := r.URL.Query()["place"]
	citySlugs, cityOk := r.URL.Query()["city"]

	if !titleOk || len(titleSlugs[0]) < 1 {
		log.Println("Url Param 'title' is missing")
		return
	}

	if !placeOk || len(placeSlugs[0]) < 1 {
		log.Println("Url Param 'place' is missing")
		return
	}

	if !cityOk || len(citySlugs[0]) < 1 {
		log.Println("Url Param 'city' is missing")
		return
	}

	id := utils.GenerateUUID()
	var placeArr []string
	newAr := append(placeArr, placeSlugs[0])
	title := titleSlugs[0]
	city := citySlugs[0]
	newFood := Food{id, title, newAr, city}
	json, err := json.Marshal(newFood)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	insertResult, err := collection.InsertOne(context.TODO(), newFood)

	if err != nil {
		panic(err)
	}
	fmt.Println(insertResult.InsertedID)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
