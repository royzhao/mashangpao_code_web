package main

import (
	"encoding/xml"
	"fmt"
)

// var (
// 	ErrAlreadyExists = errors.New("code already exists")
// )

// The Code step data structure, serializable in JSON, XML and text using the Stringer interface.
type Code_step struct {
	XMLName     xml.Name `json:"-" xml:"code_step"`
	Id          int      `json:"id" xml:"id,attr"`
	Create_date string   `json:"create_date" xml:"create_date"`
	Name        string   `json:"name" xml:"name"`
	Description string   `json:"description" xml:"description"`
	Code_id     int      `json:"code_id" xml:"code_id"`
	Image_id    int      `json:"image_id xml:"image_id"`
}

func (a *Code_step) String() string {
	return fmt.Sprintf("%s - %s (%s)", a.Name, a.Description, a.Create_date)
}
