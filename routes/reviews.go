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

// Review — так будем хранить JSON отзывы в БД
type Review struct {
	ID      string `json:"id"`
	Text    string `json:"text"`
	Rating  int    `json:"rating"`
	UserID  string `json:"userId"`
	PlaceID string `json:"placeId"`
	FoodID  string `json:"foodId"`
}

// Empty — пустой результат
type Empty struct {
	IsEmpty bool `json:"isEmpty"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// ReviewsRoute - роуты для get/post отзывов
func ReviewsRoute(w http.ResponseWriter, r *http.Request) {

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

	reviewsCollection := client.Database("test").Collection("reviews")

	if r.Method == "GET" {
		getReview(reviewsCollection, w, r)
	}

	if r.Method == "POST" {
		postReview(reviewsCollection, w, r)
	}

}

func getReview(collection *mongo.Collection, w http.ResponseWriter, r *http.Request) {
	keysID, ok := r.URL.Query()["id"]
	keysUser, userOk := r.URL.Query()["user"]
	keysPlace, placeOk := r.URL.Query()["place"]
	keysFood, foodOk := r.URL.Query()["food"]

	noID := !ok || len(keysID[0]) < 1
	noUser := !userOk || len(keysUser[0]) < 1
	noPlace := !placeOk || len(keysPlace[0]) < 1
	noFood := !foodOk || len(keysFood[0]) < 1

	if noID && noUser && noPlace && noFood {
		log.Println("Url Param 'id' and 'user' and 'place' is missing")
		return
	}
	// Если есть айдишник и нет города - вернем отзывы byUserId
	if !noID && noUser && noPlace && noFood {
		filter := bson.D{primitive.E{Key: "id", Value: keysID[0]}}

		var result Review
		err := collection.FindOne(context.TODO(), filter).Decode(&result)

		fmt.Println(result)

		if err != nil {
			http.Error(w, "Нет такого отзыва", http.StatusNotFound)
		}

		json, err := json.Marshal(result)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}

	// Если есть юзер - вернем все отзывы юзера
	if !noUser {
		filter := bson.D{primitive.E{Key: "userid", Value: keysUser[0]}}
		options := options.Find()

		var results []*Review
		cur, err := collection.Find(context.TODO(), filter, options)

		for cur.Next(context.TODO()) {

			// create a value into which the single document can be decoded
			var elem Review
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
			http.Error(w, "У юзера нет отзывов", http.StatusNotFound)
		}

		json, err := json.Marshal(results)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}

	// Если есть место - вернем все отзывы о месте
	if !noPlace {
		filter := bson.D{primitive.E{Key: "placeid", Value: keysPlace[0]}}
		options := options.Find()

		var results []*Review
		cur, err := collection.Find(context.TODO(), filter, options)

		for cur.Next(context.TODO()) {

			// create a value into which the single document can be decoded
			var elem Review
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
			http.Error(w, "У места нет отзывов", http.StatusNotFound)
		}

		json, err := json.Marshal(results)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}

	// Если есть айдишник блюда - вернем все отзывы о блюде
	if !noFood {
		filter := bson.D{primitive.E{Key: "foodid", Value: keysFood[0]}}
		options := options.Find()

		var results []*Review
		cur, err := collection.Find(context.TODO(), filter, options)

		for cur.Next(context.TODO()) {

			// create a value into which the single document can be decoded
			var elem Review
			error := cur.Decode(&elem)
			if error != nil {
				log.Fatal(error)
			}

			results = append(results, &elem)
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}
		// Close the cursor once finished
		cur.Close(context.TODO())

		if err != nil {
			http.Error(w, "У блюда нет отзывов", http.StatusNotFound)
		}

		if len(results) == 0 {
			k := make([]*Review, 0)
			empty, _ := json.Marshal(k)
			w.Header().Set("Content-Type", "application/json")
			w.Write(empty)
			return
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

func postReview(collection *mongo.Collection, w http.ResponseWriter, r *http.Request) {
	textSlugs, textOk := r.URL.Query()["text"]
	ratingSlug, ratingOk := r.URL.Query()["rating"]
	userIDSlugs, userIDOk := r.URL.Query()["user"]
	placeIDSlugs, placeIDOk := r.URL.Query()["place"]
	foodIDSlugs, foodIDOk := r.URL.Query()["food"]
	fmt.Println(placeIDSlugs)
	if !textOk || len(textSlugs[0]) < 1 {
		log.Println("Url Param 'text' is missing")
		return
	}

	if !ratingOk || len(ratingSlug[0]) < 1 {
		log.Println("Url Param 'rating' is missing")
		return
	}

	if !userIDOk || len(userIDSlugs[0]) < 1 {
		log.Println("Url Param 'user' is missing")
		return
	}

	if !placeIDOk || len(placeIDSlugs[0]) < 1 {
		log.Println("Url Param 'place' is missing")
		return
	}

	if !foodIDOk || len(foodIDSlugs[0]) < 1 {
		log.Println("Url Param 'food' is missing")
		return
	}

	id := utils.GenerateUUID()
	text := textSlugs[0]
	rating, _ := strconv.Atoi(ratingSlug[0])
	user := userIDSlugs[0]
	place := placeIDSlugs[0]
	food := foodIDSlugs[0]
	newReview := Review{id, text, rating, user, place, food}
	json, err := json.Marshal(newReview)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	insertResult, err := collection.InsertOne(context.TODO(), newReview)

	if err != nil {
		panic(err)
	}
	fmt.Println(insertResult.InsertedID)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
