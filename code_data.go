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

type CodeStar struct {
	StarId int `db:"star_id"`
	CodeId int `db:"code_id"`
	UserId int `db:"user_id"`
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

type HotCode struct {
	List  []Code `json:"list"`
	Total int64  `json:"total"`
	Page  int    `json:"page"`
	Num   int    `json:"num"`
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
	GetAll(page int, num int) HotCode
	Find(key string, page int, num int, userid int) HotCode
	Add(a *Code) (int, error)
	Update(a *Code) error
	UpdateStar(user int, codeid int) (*Code, error)
	Delete(id int)
}

// The one and only database instance.
var code_db codeDB_inter

//add star
func (db *codeDB) UpdateStar(userid int, codeid int) (*Code, error) {
	var cs CodeStar
	star := true
	trans, err := dbmap.Begin()
	if err != nil {
		return nil, err
	}
	//查询是否存在star记录
	count, err := trans.SelectInt("select count(1) from code_star where user_id = ? and code_id = ?", userid, codeid)
	if err != nil {
		trans.Rollback()
		return nil, err
	}
	//如果存在则查出那条记录
	if count > 0 {
		star = false
		err = trans.SelectOne(&cs, "select * from code_star where user_id = ? and code_id = ?", userid, codeid)
		if err != nil {
			trans.Rollback()
			return nil, err
		}
	}
	//如果不存在记录，则插入一条，并使star数加一
	if star {
		//		err = dbmap.Insert(&cs)
		cs = CodeStar{UserId: userid, CodeId: codeid}
		if err := trans.Insert(&cs); err != nil {
			log.Println("Star failed", err)
			trans.Rollback()
			return nil, err
		}
		_, err := trans.Exec("update code set Star = Star + 1 WHERE id = ? ", codeid)
		if err != nil {
			log.Println("Star failed", err)
			trans.Rollback()
			return nil, err
		}
	} else {
		//如果存在则将该记录删除，并使star数减一
		if _, err := trans.Delete(&cs); err != nil {
			trans.Rollback()
			return nil, err
		}
		_, err := trans.Exec("update code set Star = Star - 1 WHERE id = ? ", codeid)
		if err != nil {
			trans.Rollback()
			return nil, err
		}
	}
	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		return nil, err
	}
	res := db.Get(codeid)
	return &res, nil
}

// GetAll returns all albums from the database.
func (db *codeDB) GetAll(page int, num int) HotCode {
	var res []Code_modle
	var json_res HotCode
	json_res.Num = num
	json_res.Page = page
	var total int64
	total, err := db.m.SelectInt("select count(*) from code")
	checkErr(err, "error in get all")
	json_res.Total = total
	cmd := fmt.Sprintf("select * from code order by star DESC limit %d,%d", (page-1)*num, num)
	_, err = db.m.Select(&res, cmd)
	checkErr(err, "error in get all")
	for _, v := range res {
		json_res.List = append(json_res.List, convertModle2Json(v))
	}
	return json_res
}

// Find returns albums that match the search criteria.
func (db *codeDB) Find(key string, page int, num int, userid int) HotCode {
	var res HotCode
	var err error
	res.Num = num
	res.Page = page
	cmd := "select * from code where "
	if key == "" {
		if userid == -1 {
			return db.GetAll(page, num)
		} else {
			cmd = fmt.Sprintf("%s user_id=%d", cmd, userid)
		}
	} else {
		cmd += "name like '%" + key + "%' or description like '%" + key + "%'"
		if userid != -1 {
			cmd = fmt.Sprintf("%s and userid=%d", cmd, userid)
		}
	}
	var total int64
	if userid == -1 {
		total, err = db.m.SelectInt("select count(*) from code")
		checkErr(err, "error in get all")
	} else {
		tem := fmt.Sprintf("select count(*) from code where user_id=%d", userid)
		total, err = db.m.SelectInt(tem)
		checkErr(err, "error in get all")
	}
	res.Total = total
	var res_modle []Code_modle
	cmd = fmt.Sprintf("%s order by star DESC limit  %d,%d", cmd, (page-1)*num, num)
	_, err = db.m.Select(&res_modle, cmd)
	checkErr(err, "select condition failed")
	for _, v := range res_modle {
		res.List = append(res.List, convertModle2Json(v))
	}
	return res
}

// Get returns the album identified by the id, or nil.
func (db *codeDB) Get(id int) Code {
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
	db.AddTableWithName(CodeStar{}, "code_star").SetKeys(true, "StarId")
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
