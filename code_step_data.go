package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"log"
	// "os"
	// "time"
)

// var (
// 	ErrAlreadyExists = errors.New("code already exists")
// )
//db modle
// type Code_step_modle struct {
// 	Id          int
// 	Create_date string
// 	Name        string
// 	Description string
// 	Code_id     int
// 	Image_id    int
// }
// The Code step data structure, serializable in JSON, XML and text using the Stringer interface.
type Code_step struct {
	Id          int    `json:"id" xml:"id,attr"`
	Create_date string `json:"create_date" xml:"create_date"`
	Name        string `json:"name" xml:"name"`
	Description string `json:"description" xml:"description"`
	Code_id     int    `json:"code_id" xml:"code_id"`
	Image_id    int    `json:"image_id" xml:"image_id"`
}

type Code_detail struct {
	Id           int    `json:"id"`
	Code_content string `json:"code"`
	Post_content string `json:"post"`
	Time         int    `json:time`
}

func (a *Code_step) String() string {
	return fmt.Sprintf("%s - %s (%s)", a.Name, a.Description, a.Create_date)
}

type codeStepDB struct {
	m *gorp.DbMap
}

//db interface
type codeStepDB_inter interface {
	Get(id int) Code_step
	GetStepDetail(id int) Code_detail
	GetAll(id int) []Code_step
	Add(a *Code_step) (int, error)
	AddDetail(a *Code_detail) (int, error)
	Update(a *Code_step) error
	UpdateStepDetail(a *Code_detail) error
	Delete(id int)
}

//only one instance
var code_step_db codeStepDB_inter

func (db *codeStepDB) Get(id int) Code_step {
	var res Code_step
	cmd := fmt.Sprintf("select * from code_step_meta where id=%d", id)
	err := db.m.SelectOne(&res, cmd)
	checkErr(err, cmd+" failed")

	log.Println("code step query:", id)
	return res
}

func (db *codeStepDB) GetStepDetail(id int) Code_detail {
	var res Code_detail
	cmd := fmt.Sprintf("select * from code_step_detail where id=%d", id)
	err := db.m.SelectOne(&res, cmd)
	checkErr(err, cmd+" failed")

	log.Println("code step detail query:", id)
	return res
}

func (db *codeStepDB) GetAll(id int) []Code_step {
	var res []Code_step
	cmd := fmt.Sprintf("select * from code_step_meta where code_id=%d", id)

	_, err := db.m.Select(&res, cmd)
	checkErr(err, "error in get all")
	return res
}

func (db *codeStepDB) AddDetail(a *Code_detail) (int, error) {
	err := db.m.Insert(a)
	if err != nil {
		return 0, err
	}
	return a.Id, nil
}
func (db *codeStepDB) Add(a *Code_step) (int, error) {
	err := db.m.Insert(a)
	if err != nil {
		return 0, err
	}
	checkErr(err, "insert failed")
	_, err = db.m.Exec("insert into `code_step_detail` (`Id`,`Code_content`,`Post_content`,`Time`) values (?,'','',0)", a.Id)
	return a.Id, nil
	// trans, err := db.m.Begin()
	// if err != nil {
	// 	log.Fatal("open trans failed")
	// 	return 0, err
	// }
	// a.Create_date = time.Now().String()
	// trans.Insert(a)

	// detail := Code_detail{
	// 	Id: a.Id,
	// }
	// err = trans.Insert(&detail)
	// return a.Id, trans.Commit()
}

func (db *codeStepDB) Update(a *Code_step) error {
	flag := 1
	cmd := "update code_step_meta set"
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
	if a.Image_id != -1 {
		if flag == 0 {
			cmd += ","
		}
		cmd = fmt.Sprintf("%s image_id=%d", cmd, a.Image_id)
		flag = 0
	}
	if flag == 1 {
		return nil
	}
	cmd = fmt.Sprintf("%s where id=%d", cmd, a.Id)
	count, err := db.m.Exec(cmd)
	checkErr(err, "Update failed"+cmd)
	log.Println("Rows updated:", count)
	return nil
}
func (db *codeStepDB) UpdateStepDetail(a *Code_detail) error {
	flag := 1
	cmd := "update code_step_detail set"
	if a.Code_content != "" {
		log.Println("Code_content: " + a.Code_content)
		cmd += " Code_content='" + a.Code_content + "'"
		flag = 0
	}
	if a.Post_content != "" {
		log.Println("Post_content: " + a.Post_content)
		if flag == 0 {
			cmd += ","
		}
		cmd += " Post_content='" + a.Post_content + "'"
		flag = 0
	}
	if flag == 1 {
		return nil
	}
	cmd = fmt.Sprintf("%s where id=%d", cmd, a.Id)
	count, err := db.m.Exec(cmd)
	checkErr(err, "Update failed")
	log.Println("Rows updated:", count)
	return nil
}
func (db *codeStepDB) Delete(id int) {
	obj := Code_step{
		Id: id,
	}
	count, err := db.m.Delete(&obj)
	count, err = db.m.Delete(&Code_detail{Id: id})
	checkErr(err, "Delete failed")
	log.Println("Code Row deleted:", count)
}

func init_codestep(db *gorp.DbMap) {

	code_step_db = &codeStepDB{
		m: db,
	}

	db.AddTableWithName(Code_step{}, "code_step_meta").SetKeys(true, "Id")
	db.AddTableWithName(Code_detail{}, "code_step_detail").SetKeys(true, "Id")
}
