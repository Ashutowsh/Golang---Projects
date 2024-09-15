package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

type ApiResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"message"`
}

func endPointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	apiResponse := ApiResponse{
		Success: true,
		Msg:     "Hello User! How can I help you?",
	}

	err := json.NewEncoder(w).Encode(&apiResponse)
	if err != nil {
		panic(err)
	}
}

func rateLimiter(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 4)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			apiResponse := ApiResponse{
				Success: false,
				Msg:     "The API capacity is at capacity. Try again later.",
			}

			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&apiResponse)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	r.Handle("/", rateLimiter(http.HandlerFunc(endPointHandler)))
	log.Fatal(http.ListenAndServe(":4000", r))
}
