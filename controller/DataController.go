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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

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
		return
	}

	legend := make(map[string]int)
	for i, v := range records[0] {
		legend[v] = i
	}

	if len(records) <= 1 {
		u.HandleBadRequest(w, errors.New("bad rec"))
		return
	}

	if len(legend) != 8 {
		u.HandleBadRequest(w, errors.New("bad legends"))
		return
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

			boolVal := stringInSlice(d.Exam, []string{"ОГЭ", "ЕГЭ"})
			if boolVal == false {
				u.HandleBadRequest(w, errors.New("bad exam"))
				return
			}

			boolVal = stringInSlice(d.Subject, []string{"Математика", "Русский язык", "Физика", "Химия", "Биология", "Литература", "География", "История", "Информатика", "Обществознание", "Информатика", "Ангийский язык", "Французский язык", "Немецкий язык", "Испанский язык", "Китайский язык"})
			if boolVal == false {
				u.HandleBadRequest(w, errors.New("bad subject"))
				return
			}

			value, err := strconv.ParseUint(v[legend["Период"]], 10, 64)
			if err != nil {
				u.HandleBadRequest(w, errors.New("bad period"))
				return
			}
			d.Period = uint(value)

			if d.Period < 2001 || d.Period >= 2020 {
				u.HandleBadRequest(w, errors.New("bad period"))
				return
			}

			value, err = strconv.ParseUint(v[legend["Баллы"]], 10, 64)
			if err != nil {
				u.HandleBadRequest(w, errors.New("bad score"))
				return
			}
			d.Score = uint(value)

			if d.Score > 100 || d.Score < 0 {
				u.HandleBadRequest(w, errors.New("bad score"))
				return
			}

			value, err = strconv.ParseUint(v[legend["Оценка"]], 10, 64)
			if err != nil {
				u.HandleBadRequest(w, errors.New("bad grade"))
				return
			}
			d.Grade = uint(value)

			if d.Grade > 5 || d.Grade < 2 {
				u.HandleBadRequest(w, errors.New("bad grade"))
				return
			}

			value, err = strconv.ParseUint(v[legend["Признак успешной сдачи экзамена"]], 10, 64)
			if err != nil {
				u.HandleBadRequest(w, errors.New("bad parameter"))
				return
			}
			d.IsPassed = uint(value)

			if d.IsPassed != 1 && d.IsPassed != 0 {
				u.HandleBadRequest(w, errors.New("bad passed"))
				return
			}

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