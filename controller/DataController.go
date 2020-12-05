package controller

import (
	"DigitalRegionAPI/db"
	"DigitalRegionAPI/models/auxiliary"
	"DigitalRegionAPI/models/entities"
	u "DigitalRegionAPI/utils"
	"encoding/csv"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var ImportDataCSV = func(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("error reading body: %v", err)
		u.HandleInternalError(w, err)
		return
	}

	var data auxiliary.CSV
	err = json.Unmarshal(body, &data)
	if err != nil {
		u.HandleBadRequest(w, err)
		return
	}

	log.Debug("CSV length: ", len(data.Data))

	res := csv.NewReader(strings.NewReader(string(data.Data)))
	records, err := res.ReadAll()
	if err != nil {
		u.HandleInternalError(w, err)
	}

	legend := make(map[string]int)
	for i, v := range records[0] {
		legend[v] = i
	}

	tx := db.GetDB().Begin()
	defer tx.Rollback()

	schools := make(map[string]uint)
	for i, v := range records {
		if i > 0 {
			d := entities.Data{}

			d.School = v[legend["Школа"]]
			d.Student = v[legend["Ученик"]]
			d.Exam = v[legend["Экзамен"]]
			d.Subject = v[legend["Предмет"]]

			value, err := strconv.ParseUint(v[legend["Период"]], 10, 64)
			if err != nil {
				u.HandleBadRequest(w, errors.New("bad period"))
				return
			}
			d.Period = uint(value)

			value, err = strconv.ParseUint(v[legend["Баллы"]], 10, 64)
			if err != nil {
				u.HandleBadRequest(w, errors.New("bad score"))
				return
			}
			d.Score = uint(value)

			value, err = strconv.ParseUint(v[legend["Оценка"]], 10, 64)
			if err != nil {
				u.HandleBadRequest(w, errors.New("bad grade"))
				return
			}
			d.Grade = uint(value)

			value, err = strconv.ParseUint(v[legend["Признак успешной сдачи экзамена"]], 10, 64)
			if err != nil {
				u.HandleBadRequest(w, errors.New("bad parameter"))
				return
			}
			d.IsPassed = uint(value)

			if _, ok := schools[d.School]; ok {
				schools[d.School]++
			} else {
				schools[d.School] = 0
			}

			err = tx.Create(&d).Error
			if err != nil {
				u.HandleInternalError(w, err)
				return
			}
		}
	}

	var users []entities.User
	err = tx.Find(&users).Error
	if err != nil {
		u.HandleInternalError(w, err)
		return
	}

	for _, v := range users {
		v.DataUploads += schools[v.School]
		tx.Save(v)
	}

	tx.Commit()
	u.Respond(w, u.Message(true, "OK"))
}

var QueryData = func(w http.ResponseWriter, r *http.Request) {
	var data []entities.Data

	err := db.GetDB().Find(&data).Error
	if err != nil {
		u.HandleInternalError(w, err)
	}

	res, err := json.Marshal(data)
	if err != nil {
		u.HandleInternalError(w, err)
	} else {
		u.RespondJSON(w, res)
	}
}