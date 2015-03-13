package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"log"
	// "os"
	"sync"
	"time"
)

var (
	ErrAlreadyExists = errors.New("code already exists")
)

type Code_modle struct {
	Id          int    `db:"id"`
	Create_date string `db:"create_date"`
	Name        string `db:"name"`
	Description string `db:"description"`
	User_id     int    `db:"user_id"`
	Star        int    `db:"star"`
}

// The Code data structure, serializable in JSON, XML and text using the Stringer interface.
type Code struct {
	XMLName     xml.Name `json:"-" xml:"album"`
	Id          int      `json:"id" xml:"id,attr"`
	Create_date string   `json:"create_date" xml:"create_date"`
	Name        string   `json:"name" xml:"name"`
	Description string   `json:"description" xml:"description"`
	User_id     int      `json:"user_id" xml:"user_id"`
	Star        int      `json:"star" xml:"star"`
}

func (a *Code) String() string {
	return fmt.Sprintf("%s - %s (%s)", a.Name, a.Description, a.Create_date)
}

// Thread-safe in-memory map of Code.
type codeDB struct {
	sync.RWMutex
	m   *gorp.DbMap
	seq int
}

// The DB interface defines methods to manipulate the code.
type codeDB_inter interface {
	Get(id int) Code
	GetAll() []Code
	Find(name string, description string, create_time string, userid int) []Code
	Add(a *Code) (int, error)
	Update(a *Code) error
	Delete(id int)
}

// The one and only database instance.
var code_db codeDB_inter

// GetAll returns all albums from the database.
func (db *codeDB) GetAll() []Code {
	db.RLock()
	defer db.RUnlock()
	var res []Code_modle
	var json_res []Code
	_, err := db.m.Select(&res, "select * from code")
	checkErr(err, "error in get all")
	for _, v := range res {
		json_res = append(json_res, convertModle2Json(v))
	}
	return json_res
}

// Find returns albums that match the search criteria.
func (db *codeDB) Find(name string, description string, create_date string, userid int) []Code {
	db.RLock()
	defer db.RUnlock()
	var res []Code
	var res_modle []Code_modle
	cmd := "select * from code "
	flag := 0
	if name != "" {
		if flag == 0 {
			cmd += "where "
			flag = 1
		}
		cmd += "name='" + name + "'"
	}
	if description != "" {
		if flag == 0 {
			cmd += "where "
			flag = 1
		} else {
			cmd += " and "
		}
		cmd += " description='" + description + "'"
	}
	if create_date != "" {
		if flag == 0 {
			cmd += "where "
			flag = 1
		} else {
			cmd += " and "
		}
		cmd += " create_date='" + create_date + "'"
	}
	if userid != -1 {
		if flag == 0 {
			cmd += "where "
			flag = 1
		} else {
			cmd += " and "
		}
		cmd = fmt.Sprintf("%s user_id=%d", cmd, userid)
	}
	_, err := db.m.Select(&res_modle, cmd)
	checkErr(err, "select condition failed")
	for _, v := range res_modle {
		res = append(res, convertModle2Json(v))
	}
	return res
}

// Get returns the album identified by the id, or nil.
func (db *codeDB) Get(id int) Code {
	db.RLock()
	defer db.RUnlock()
	var res Code_modle
	cmd := fmt.Sprintf("select * from code where id =%d", id)
	err := db.m.SelectOne(&res, cmd)
	checkErr(err, cmd+" failed")
	// obj, err := db.m.Get(Code_modle{}, id)
	// checkErr(err, "select faile")
	log.Println("code query:", id)
	// if obj == nil {
	// 	return Code{}
	// }
	return convertModle2Json(res)
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
	// Store
	//compute time
	a.Create_date = time.Now().String()
	obj := convertJson2Modle(*a)
	err := db.m.Insert(&obj)
	if checkErr(err, "Insert failed") == true {
		return 0, err
	}
	a.Id = obj.Id

	return a.Id, nil
}

// Update changes the album identified by the id. It returns an error if the
// updated album is a duplicate.
func (db *codeDB) Update(a *Code) error {
	db.Lock()
	defer db.Unlock()
	flag := 1
	cmd := "update code set "
	if a.Name != "" {
		log.Println("name: " + a.Name)
		cmd += " name='" + a.Name + "'"
		flag = 0
	}
	if a.Description != "" {
		log.Println("description: " + a.Description)
		if flag == 0 {
			cmd += ","
		}
		cmd += " description='" + a.Description + "'"
		flag = 0
	}
	if flag == 1 {
		return nil
	}
	cmd = fmt.Sprintf("%s where id=%d", cmd, a.Id)
	count, err := db.m.Exec(cmd)
	if checkErr(err, "Update failed") == true {
		return err
	}
	log.Println("Rows updated:", count)
	return nil
}

// Delete removes the album identified by the id from the database. It is a no-op
// if the id does not exist.
func (db *codeDB) Delete(id int) {
	db.Lock()
	defer db.Unlock()
	obj := Code_modle{
		Id: id,
	}
	count, err := db.m.Delete(obj)
	checkErr(err, "Delete failed")
	log.Println("Code Row deleted:", count)
}

// Checks if the album already exists in the database, based on the Band and Title
// fields.
func (db *codeDB) isUnique(a *Code) bool {
	var res []Code_modle
	var cmd string
	cmd = fmt.Sprintf("select * from code where description ='%s' and name='%s' and user_id=%d", a.Description, a.Name, a.User_id)
	_, err := db.m.Select(&res, cmd)
	checkErr(err, "check is unique failed")
	if len(res) == 0 {
		return true
	}
	return false
}

func init_code(db *gorp.DbMap) {
	// dbmap := initDb()
	// dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))

	// defer dbmap.Db.Close()
	code_db = &codeDB{
		m: db,
	}

	// p1 := Code_modle{
	// 	Name:        "zpl",
	// 	Description: "dsadas dsadasd",
	// 	User_id:     1,
	// }
	// p2 := Code_modle{
	// 	Name:        "zpl2",
	// 	Description: "dsadas dsadas2d",
	// 	User_id:     1,
	// }
	// err := dbmap.Insert(&p2)
	// checkErr(err, "Insert failed")
	// Fill the database
	// code_db.Add(&Code{Name: "zpl", Description: "Reign1 333", User_id: 1})

	// code_db.Add(&Code{Name: "zpl2", Description: "Reign2", User_id: 2})
	// code_db.Add(&Code{Name: "zpl3", Description: "Reign3", User_id: 1})
	db.AddTableWithName(Code_modle{}, "code").SetKeys(true, "Id")
}

func convertJson2Modle(code Code) Code_modle {
	return Code_modle{
		User_id:     code.User_id,
		Id:          code.Id,
		Name:        code.Name,
		Description: code.Description,
		Create_date: code.Create_date,
		Star:        0,
	}
}
func convertModle2Json(code Code_modle) Code {
	return Code{
		User_id:     code.User_id,
		Id:          code.Id,
		Name:        code.Name,
		Description: code.Description,
		Create_date: code.Create_date,
		Star:        code.Star,
	}
}
