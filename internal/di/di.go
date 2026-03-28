package di

import (
	"log"
	"net/http"
	"tictactoe/internal/handler"
	"tictactoe/internal/service"
	"tictactoe/internal/storage"
)

func Init() *handler.GameHandler {
	repo := storage.NewGameRepository()
	svc := &service.GameService{}
	h := &handler.GameHandler{
		Service: svc,
		Repo:    repo,
	}
	return h
}

func StartServer(handler *handler.GameHandler) {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./gui"))

	mux.Handle("/", fs)
	mux.HandleFunc("POST /game", handler.CreateGame)
	mux.HandleFunc("POST /game/{id}", handler.MakeMove)
	mux.HandleFunc("GET /game/{id}", handler.GetGame)

	log.Println("Сервер запущен на :80")
	log.Fatal(http.ListenAndServe(":80", enableCORS(mux)))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
