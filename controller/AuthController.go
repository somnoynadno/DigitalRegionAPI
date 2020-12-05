package controller

import (
	"DigitalRegionAPI/db"
	"DigitalRegionAPI/models/auxiliary"
	"DigitalRegionAPI/models/entities"
	u "DigitalRegionAPI/utils"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func login(username, password string) map[string]interface{} {
	account := &entities.User{}
	err := db.GetDB().Table("users").Where("username = ?", username).First(account).Error

	if err != nil {
		log.Warn(err)
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "User not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	if password != account.Password { // Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	// Worked! Logged In
	db.GetDB().Model(&account).Update("UpdatedAt", time.Now())

	// Create JWT token
	tk := &auxiliary.Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte("hackathon"))

	resp := u.Message(true, "Logged In")
	resp["token"] = tokenString
	resp["user_id"] = account.ID

	return resp
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	account := &auxiliary.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.HandleBadRequest(w, err)
		return
	}

	resp := login(account.Username, account.Password)

	if resp["token"] == nil {
		u.HandleBadRequest(w, errors.New("wrong credentials"))
		return
	}

	u.Respond(w, resp)
}
