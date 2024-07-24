package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/yizhong187/EduMind-backend/handlers"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/internal/database"
	"github.com/yizhong187/EduMind-backend/middlewares"
	"github.com/yizhong187/EduMind-backend/routers"

	_ "github.com/lib/pq"
)

func main() {

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	fmt.Println(err)
	// 	log.Fatal("Error loading .env file")
	// }

	// retrieving the environment variables, if not set a fatal error will be logged and programme will be terminated.
	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "8080" // Default port if not set
		log.Printf("No PORT environment variable detected. Defaulting to %s\n", portString)
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY is not found in the environment")
	}

	// establishes a connection with the database. note that connection is lazily established and most errors will only be
	// thrown during a query and not during the the opening of the connection.
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Unable to connect to database:", err)
	}

	dbQueries := database.New(db)

	// used to configure API handlers by encapsulating various dependencies they might need.
	// in this case, the database connection.
	apiCfg := config.ApiConfig{
		DB:        dbQueries,
		DBConn:    db,
		SecretKey: secretKey,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Use(middleware.Logger)
	v1Router.Use(func(next http.Handler) http.Handler {
		return middlewares.MiddlewareApiConfig(next, &apiCfg)
	})
	v1Router.Get("/healthz", handlers.HandlerReadiness)
	v1Router.Get("/error", handlers.HandlerError)

	v1Router.Get("/subjects", handlers.HandlerGetAllSubjects)

	v1Router.Mount("/students", routers.StudentRouter(&apiCfg))
	v1Router.Mount("/tutors", routers.TutorRouter(&apiCfg))
	v1Router.Mount("/chats", routers.ChatRouter(&apiCfg))
	v1Router.Mount("/", routers.UtilRouter(&apiCfg))

	router.Mount("/v1", v1Router)

	// configuring http server with router and port
	// Addr is the TCP address to listen on (listening for HTTP requests at the port)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}
