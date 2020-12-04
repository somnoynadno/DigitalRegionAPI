package utils

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

// 400
func HandleBadRequest(w http.ResponseWriter, err error) {
	log.Warn(err)
	w.WriteHeader(http.StatusBadRequest)
	Respond(w, Message(false, err.Error()))
}

// 401
func HandleUnauthorized(w http.ResponseWriter, err error) {
	log.Info(err)
	w.WriteHeader(http.StatusUnauthorized)
	Respond(w, Message(false, err.Error()))
}

// 403
func HandleForbidden(w http.ResponseWriter, err error) {
	log.Warn(err)
	w.WriteHeader(http.StatusForbidden)
	Respond(w, Message(false, err.Error()))
}

// 404
func HandleNotFound(w http.ResponseWriter, err error) {
	log.Warn(err)
	w.WriteHeader(http.StatusNotFound)
	Respond(w, Message(false, "not found"))
}

// 500
func HandleInternalError(w http.ResponseWriter, err error) {
	log.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
	Respond(w, Message(false, err.Error()))
}

