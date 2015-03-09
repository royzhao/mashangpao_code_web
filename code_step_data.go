package code_web

import (
	"encoding/xml"
	"errors"
	"fmt"
	"sync"
)

// The Album data structure, serializable in JSON, XML and text using the Stringer interface.
type Code struct {
	XMLName xml.Name `json:"-" xml:"album"`
	Id      int      `json:"id" xml:"id,attr"`
	Create_date   string   `json:"create_date" xml:"create_date"`
	Name    string      `json:"name" xml:"name"`
	Description	string `json:"description" xml:"description"`
	User_id	int	`json:"user_id" xml:"user_id"`
}
