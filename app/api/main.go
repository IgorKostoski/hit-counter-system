package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"os"
	"time"
)

//triggering ci build

var db *sql.DB

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "hitcounter")
	appPort := getEnv("APP_PORT", "8080")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("Failed to connect to database, retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		err = db.Ping()
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}
	}

	log.Println("Successfully connected to database")

	if err := InitSchema(db); err != nil {
		log.Fatalf("Error initializing schema: %v", err)
	}

	router := gin.Default()

	promHandler := promhttp.Handler()
	router.GET("/metrics", func(c *gin.Context) {
		promHandler.ServeHTTP(c.Writer, c.Request)
	})

	api := router.Group("/api/v1")
	{
		api.POST("/hit", HitHandler)
		api.GET("/count/:key", CountHandler)
	}

	log.Printf("Starting server on port %s", appPort)
	if err := router.Run(":" + appPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
