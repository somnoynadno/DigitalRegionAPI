package controller

import (
	"DigitalRegionAPI/db"
	"DigitalRegionAPI/models/auxiliary"
	u "DigitalRegionAPI/utils"
	"encoding/json"
	"errors"
	"net/http"
)

type School struct {
	Name  string
	Score float64
}

type Subject struct {
	Name string
}

type Student struct {
	Name  string
	Score int
}

type Grade struct {
	Average float64
}

var GetStats = func(w http.ResponseWriter, r *http.Request) {
	exam := r.FormValue("exam")
	period := r.FormValue("period")

	if exam == "" || period == "" {
		u.HandleBadRequest(w, errors.New("not enough parameters"))
		return
	}

	db := db.GetDB()

	var totalWorks int
	var numberOfTopWorks int
	var passedWorks int
	var averageGrade Grade
	var topSchool School
	var worstSchool School
	var topSubject Subject
	var worstSubject Subject
	var topStudent Student

	db.Table("data").Where("exam = ?", exam).Where("period = ?", period).Count(&totalWorks)
	db.Table("data").Where("exam = ?", exam).Where("period = ?", period).
		Where("is_passed = 1").Count(&passedWorks)
	db.Table("data").Where("exam = ?", exam).Where("period = ?", period).
		Select("avg(grade) as average").First(&averageGrade)
	db.Table("data").Where("exam = ?", exam).Where("period = ?", period).
		Where("score = 100").Count(&numberOfTopWorks)
	db.Table("data").Select("school as name, avg(score) as score").
		Where("exam = ?", exam).Where("period = ?", period).
		Order("avg(score) desc").
		Group("school").First(&topSchool)
	db.Table("data").Select("school as name, avg(score) as score").
		Where("exam = ?", exam).Where("period = ?", period).
		Order("avg(score) asc").
		Group("school").First(&worstSchool)
	db.Table("data").Select("subject as name").
		Where("exam = ?", exam).Where("period = ?", period).
		Order("avg(score) desc").
		Group("subject").First(&topSubject)
	db.Table("data").Select("subject as name").
		Where("exam = ?", exam).Where("period = ?", period).
		Order("avg(score) asc").
		Group("subject").First(&worstSubject)
	db.Table("data").Select("student as name, sum(score) as score").
		Where("exam = ?", exam).Where("period = ?", period).
		Order("sum(score) desc").
		Group("student").First(&topStudent)

	passedPercentage := 0.0
	if totalWorks > 0 {
		passedPercentage = float64(passedWorks * 100 / totalWorks)
	}

	data := auxiliary.Stats{
		TotalWorks:  totalWorks,
		TopSchool:   topSchool.Name,
		WorstSchool: worstSchool.Name,
		NumberOfTopWorks: numberOfTopWorks,
		TopSubject: topSubject.Name,
		WorstSubject: worstSubject.Name,
		AverageGrade: averageGrade.Average,
		TopStudent: topStudent.Name,
		TopTotalScore: topStudent.Score,
		PassedPercentage: passedPercentage,
		TopScore: int(topSchool.Score),
		WorstScore: int(worstSchool.Score),
	}

	res, err := json.Marshal(data)
	if err != nil {
		u.HandleInternalError(w, err)
	} else {
		u.RespondJSON(w, res)
	}
}

