package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/AnikBarua007/http_server_go/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: could not load .env file: %v", err)
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database connection: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	apicfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
	}
	mux := http.NewServeMux()
	//mux.Handle("/assets/", http.FileServer(http.Dir(".")))
	fshandler := apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fshandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apicfg.handlermatrics)
	mux.HandleFunc("POST /admin/reset", apicfg.handlerReset)
	//mux.HandleFunc("POST /api/validate_chirp", handlerValidate)
	mux.HandleFunc("POST /api/users", apicfg.handleruser)
	mux.HandleFunc("POST /api/chirps", apicfg.handlerChirp)
	mux.HandleFunc("GET /api/chirps", apicfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apicfg.handlerGetIDchirps)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("MISS: method=%s path=%s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	})
	log.Fatal(server.ListenAndServe())
}
