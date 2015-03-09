package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrAlreadyExists = errors.New("code already exists")
)

// The Code data structure, serializable in JSON, XML and text using the Stringer interface.
type Code struct {
	XMLName xml.Name `json:"-" xml:"album"`
	Id      int      `json:"id" xml:"id,attr"`
	Create_date   string   `json:"create_date" xml:"create_date"`
	Name    string      `json:"name" xml:"name"`
	Description	string `json:"description" xml:"description"`
	User_id	int	`json:"user_id" xml:"user_id"`
}

func (a *Code) String() string {
	return fmt.Sprintf("%s - %s (%s)", a.Name, a.Description, a.Create_date)
}

// Thread-safe in-memory map of albums.
type codeDB struct {
	sync.RWMutex
	m   map[int]*Code
	seq int
}


// The DB interface defines methods to manipulate the code.
type codeDB_inter interface {
	Get(id int) *Code
	GetAll() []*Code
	Find(name string,description string, create_time string,userid int) []*Code
	Add(a *Code) (int, error)
	Update(a *Code) error
	Delete(id int)
}

// The one and only database instance.
var code_db codeDB_inter

// GetAll returns all albums from the database.
func (db *codeDB) GetAll() []*Code {
	db.RLock()
	defer db.RUnlock()
	if len(db.m) == 0 {
		return nil
	}
	ar := make([]*Code, len(db.m))
	i := 0
	for _, v := range db.m {
		ar[i] = v
		i++
	}
	return ar
}

// Find returns albums that match the search criteria.
func (db *codeDB) Find(name string, description string ,create_date string,userid int) []*Code {
	db.RLock()
	defer db.RUnlock()
	var res []*Code
	for _, v := range db.m {
		if v.Name == name || name == "" {
			if v.Description == description|| description == "" {
				if v.Create_date == create_date || create_date == "" {
					if v.User_id == userid || userid==-1 {
						res = append(res, v)
					}
				}
			}
		}
	}
	return res
}

// Get returns the album identified by the id, or nil.
func (db *codeDB) Get(id int) *Code {
	db.RLock()
	defer db.RUnlock()
	return db.m[id]
}

// Add creates a new album and returns its id, or an error.
func (db *codeDB) Add(a *Code) (int, error) {
	db.Lock()
	defer db.Unlock()
	// Return an error if band-title already exists
	if !db.isUnique(a) {
		return 0, ErrAlreadyExists
	}
	// Get the unique ID
	db.seq++
	a.Id = db.seq
	// Store
	db.m[a.Id] = a
	return a.Id, nil
}

// Update changes the album identified by the id. It returns an error if the
// updated album is a duplicate.
func (db *codeDB) Update(a *Code) error {
	db.Lock()
	defer db.Unlock()
	if !db.isUnique(a) {
		return ErrAlreadyExists
	}
	db.m[a.Id] = a
	return nil
}

// Delete removes the album identified by the id from the database. It is a no-op
// if the id does not exist.
func (db *codeDB) Delete(id int) {
	db.Lock()
	defer db.Unlock()
	delete(db.m, id)
}

// Checks if the album already exists in the database, based on the Band and Title
// fields.
func (db *codeDB) isUnique(a *Code) bool {
	for _, v := range db.m {
		if v.Name == a.Name  && v.Description == a.Description&& v.Id != a.Id&&v.User_id != a.User_id {
			return false
		}
	}
	return true
}

func init() {
	code_db = &codeDB{
		m: make(map[int]*Code),
	}
	// Fill the database
	code_db.Add(&Code{Id: 1, Name: "zpl", Description: "Reign In Blood", Create_date: "1986-10-22",User_id:1})
	
	code_db.Add(&Code{Id: 2, Name: "zpl2", Description: "Reign In Blood2", Create_date: "1986-10-22",User_id:2})
	code_db.Add(&Code{Id: 3, Name: "zpl3", Description: "Reign In Blood3", Create_date: "1986-10-22",User_id:1})
}