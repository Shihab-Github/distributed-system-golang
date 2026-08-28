package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Video struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Genre    string `json:"genre"`
	Duration int    `json:"duration_minutes"`
}

var port string

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

var ctx = context.Background()
var dbPool *pgxpool.Pool

func fetchCatalogFromDB() ([]Video, error) {
	rows, err := dbPool.Query(ctx, "SELECT * FROM videos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []Video

	for rows.Next() {
		var v Video
		err := rows.Scan(&v.ID, &v.Title, &v.Genre, &v.Duration)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}

	return videos, nil

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

	data, err := fetchCatalogFromDB()
	if err != nil {
		http.Error(w, "Failed to fetch catalog", http.StatusInternalServerError)
		log.Printf("DB error: %v", err)
		return
	}
	jsonData, _ := json.Marshal(data)

	rdb.Set(ctx, "catalog", jsonData, 10*time.Second)

	w.Write(jsonData)
}

func main() {
	var err error
	dbPool, err = pgxpool.New(ctx, "postgres://strteamhub:streamhub@localhost:5432/streamhub")
	if err != nil {
		log.Fatalf("Unabnle to connect to database: %v", err)
	}

	defer dbPool.Close()

	http.HandleFunc("/catalog", catalogHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Streamhub server is running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
