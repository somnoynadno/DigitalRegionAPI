package main

import (
	"DigitalRegionAPI/controller"
	"DigitalRegionAPI/db"
	"DigitalRegionAPI/middleware"
	u "DigitalRegionAPI/utils"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/auth/login", nil).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/data/query", controller.QueryData).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/data/send_csv", controller.ImportDataCSV).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/ping", u.HandlePing).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	router.Use(middleware.CORS)
	router.Use(middleware.LogPath)
	router.Use(middleware.LogBody)

	// check connection
	con := db.GetDB()
	errors := con.GetErrors()
	if errors != nil && len(errors) > 0 {
		panic(errors[0])
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "7020" // localhost
	}

	log.Info("Listening on: ", port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		panic(err)
	}
}
