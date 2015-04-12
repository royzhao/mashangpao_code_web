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
	Code_name   string `json:"code_name" xml:"code_name"`
	Status      int    `json:"status" xml:"status"`
}

type Code_detail struct {
	Id           int    `json:"id"`
	Code_content string `json:"code"`
	Post_content string `json:"post"`
	Time         int    `json:time`
}
type Code_step_cmd struct {
	Id         int    `json:id`
	Seq        int    `json:seq`
	Cmd        string `json:cmd`
	Args       string `json:args`
	Is_replace int    `json:is_replace`
	Stepid     int    `json:stepid`
}

func (a *Code_step) String() string {
	return fmt.Sprintf("%s - %s (%s) image=%d code=%d status=%d", a.Name, a.Description, a.Create_date, a.Image_id, a.Code_id, a.Status)
}

type codeStepDB struct {
	m *gorp.DbMap
}

//db interface
type codeStepDB_inter interface {
	Get(id int) Code_step
	GetStepDetail(id int) Code_detail
	Find(image_id int, code_id int, name string, status int) []Code_step
	GetAll(id int) []Code_step
	Add(a *Code_step) (int, error)
	AddDetail(a *Code_detail) (int, error)
	Update(a *Code_step) error
	UpdateStepDetail(a *Code_detail) error
	UpdateCodeCmd(stepid int, a []Code_step_cmd) error
	GetCodeCmds(stepid int) []Code_step_cmd
	GetCodeCmdBySeq(stepid int, seqid int) *Code_step_cmd
	DeleteCodeCmd(stepid int, seqid int) error
	Delete(id int)
}

//only one instance
var code_step_db codeStepDB_inter

func (db *codeStepDB) Find(image_id int, code_id int, name string, status int) []Code_step {
	var res []Code_step
	cmd := fmt.Sprintf("select * from code_step_meta where status=%d ", status)

	if image_id != -1 {
		cmd = fmt.Sprintf(cmd+" and image_id=%d", image_id)
	}
	if code_id != -1 {
		cmd = fmt.Sprintf(cmd+" and code_id=%d", code_id)
	}

	if name != "" {
		cmd += " and name='" + name + "'"
	}
	_, err := db.m.Select(&res, cmd)
	checkErr(err, "error in find")
	return res
}
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
	log.Println(a.String())
	flag := 1
	cmd := "update code_step_meta set"
	if a.Name != "" {
		log.Println("name: " + a.Name)
		cmd += " name='" + a.Name + "'"
		flag = 0
	}
	if a.Code_name != "" {
		log.Println("code_name: " + a.Code_name)
		cmd += " code_name='" + a.Code_name + "'"
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
	if a.Image_id != 0 {
		if flag == 0 {
			cmd += ","
		}
		cmd = fmt.Sprintf("%s image_id=%d", cmd, a.Image_id)
		flag = 0
	}
	if a.Status != 0 {
		if flag == 0 {
			cmd += ","
		}
		cmd = fmt.Sprintf("%s status=%d", cmd, a.Status)
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
func (db *codeStepDB) GetCodeCmds(stepid int) []Code_step_cmd {
	var res []Code_step_cmd
	cmd := fmt.Sprintf("select * from code_step_cmd where stepid=%d", stepid)
	_, err := db.m.Select(&res, cmd)
	checkErr(err, "error in get all cmd")
	return res
}
func (db *codeStepDB) GetCodeCmdBySeq(stepid int, seqid int) *Code_step_cmd {
	var res Code_step_cmd
	cmd := fmt.Sprintf("select * from code_step_cmd where stepid=%d and seq=%d", stepid, seqid)
	_, err := db.m.Select(&res, cmd)
	checkErr(err, "error in get one cmd")
	return &res
}
func (db *codeStepDB) UpdateCodeCmd(stepid int, a []Code_step_cmd) error {
	for _, v := range a {
		tmp := db.GetCodeCmdBySeq(stepid, v.Seq)
		if tmp.Stepid != 0 && tmp.Seq != 0 {
			cmd := "update code_step_cmd set"
			flag := 1
			if v.Cmd != "" {
				cmd += " cmd=" + v.Cmd
				flag = 0
			}
			if v.Args != "" {
				if flag == 0 {
					cmd += ","
				}
				cmd += " args=" + v.Args
				flag = 0
			}
			if v.Is_replace != 0 {
				if flag == 0 {
					cmd += ","
				}
				cmd += fmt.Sprint("is_replace=%d", v.Is_replace)
				flag = 0
			}
			if flag == 0 {
				cmd = fmt.Sprintf("%s where stepid=%d and seq=%d", cmd, stepid, v.Seq)
				_, err := db.m.Exec(cmd)
				checkErr(err, "Update failed")
			}
		} else {
			cmd := fmt.Sprintf("insert into code_step_cmd (seq,cmd,args,is_replace,stepid) values(%d,'%s','%s',%d,%d)", v.Seq, v.Cmd, v.Args, v.Is_replace, v.Stepid)

			_, err := db.m.Exec(cmd)
			checkErr(err, "insert failed")
		}
	}
	return nil
}

func (db *codeStepDB) DeleteCodeCmd(stepid int, seqid int) error {
	tmp := db.GetCodeCmdBySeq(stepid, seqid)
	if tmp == nil {
		return nil
	}
	count, err := db.m.Delete(tmp)
	checkErr(err, "Delete failed")
	log.Println("Code Row deleted:", count)
	return err
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
	db.AddTableWithName(Code_step_cmd{}, "code_step_cmd").SetKeys(true, "Id")
}
