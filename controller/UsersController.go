package controller

import (
	"DigitalRegionAPI/db"
	"DigitalRegionAPI/models/entities"
	u "DigitalRegionAPI/utils"
	"encoding/json"
	"net/http"
)

var QueryUsers = func(w http.ResponseWriter, r *http.Request) {
	var data []entities.User

	err := db.GetDB().Find(&data).Error
	if err != nil {
		u.HandleInternalError(w, err)
		return
	}

	res, err := json.Marshal(data)
	if err != nil {
		u.HandleInternalError(w, err)
	} else {
		u.RespondJSON(w, res)
	}
}
