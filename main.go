package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Video struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Genre    string `json:"genre"`
	Duration int    `json:"duration_minutes"`
}

var catalog = []Video{
	{
		ID:       1,
		Title:    "Inception",
		Genre:    "Sci-Fi",
		Duration: 148,
	},
	{
		ID:       2,
		Title:    "The Godfather",
		Genre:    "Crime",
		Duration: 175,
	},
	{
		ID:       3,
		Title:    "The Dark Knight",
		Genre:    "Action",
		Duration: 152,
	},
}

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

var ctx = context.Background()

func fetchCatalogFromDB() []Video {
	time.Sleep(200 * time.Millisecond)
	return catalog
}

func catalogHandler(w http.ResponseWriter, r *http.Request) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Handling request from %s", port)
	w.Header().Set("Content-Type", "application/json")

	cached, err := rdb.Get(ctx, "catalog").Result()
	if err == nil {
		w.Write([]byte(cached))
		return
	}

	data := fetchCatalogFromDB()
	jsonData, _ := json.Marshal(data)

	rdb.Set(ctx, "catalog", jsonData, 10*time.Second)

	w.Write(jsonData)
}

func main() {
	http.HandleFunc("/catalog", catalogHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Streamhub server is running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
