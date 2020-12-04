package entities

import "DigitalRegionAPI/models/auxiliary"

type Data struct {
	auxiliary.BaseModelCompact
	School   string
	Student  string
	Exam     string
	Period   uint
	Subject  string
	Score    uint
	Grade    uint
	IsPassed uint
}
