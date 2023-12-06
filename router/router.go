package router

import (
	"github.com/gorilla/mux"
	"github.com/ife-oluwa/go-postres/middleware"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/stock/{id}", middleware.GetStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/stock/", middleware.GetAllStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/newstock", middleware.CreateStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/stock/{id}", middleware.UpdateStock).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/v1/deletestock/{id}", middleware.DeleteStock).Methods("DELETE", "OPTIONS")

	return router
}
