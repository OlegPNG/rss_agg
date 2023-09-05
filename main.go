package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/QuestPNG/rss_agg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	dotenv "github.com/joho/godotenv"

        _ "github.com/lib/pq"
)

type apiConfig struct {
    DB *database.Queries
}

func main() {

    dotenv.Load()
    portString := os.Getenv("PORT")
    if portString == "" {
        log.Fatal("PORT is not found in the environment")
    }

    dbUrl := os.Getenv("DB_URL")
    if dbUrl == "" {
        log.Fatal("DB_URL is not found in the environment") 
    }

    db, err := sql.Open("postgres", dbUrl)
    if err != nil {
        log.Fatal("Can't connect to database")
    }

    apiCfg := apiConfig {
        DB: database.New(db),
    }
    router := chi.NewRouter()

    router.Use(cors.Handler(cors.Options{
        AllowedOrigins:     []string{"https://*", "http://*"},
        AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:     []string{"*"},
        ExposedHeaders:     []string{"Link"},
        AllowCredentials:   false,
        MaxAge:             300,
    }))

    v1Router := chi.NewRouter()

    v1Router.Get("/ready", handlerReadiness)
    v1Router.Get("/err", handlerErr)
    v1Router.Post("/user", apiCfg.handlerCreateUser)

    router.Mount("/v1", v1Router)

    srv := &http.Server {
        Handler:    router,
        Addr:       ":" + portString,
    }

    log.Printf("Server starting on port %v", portString)
    err = srv.ListenAndServe()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("PORT: %s\n", portString)
}
