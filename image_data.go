package main

import (
	"gopkg.in/gorp.v1"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type CRImage struct {
	// db tag lets you specify the column name if it differs from the struct field.
	// remember to keep the first letter of the fields in the struct uppercase
	// because All the fields in a struct are exported or hidden simply based on the first letter
	// if it is uppercase, the field is exported. Otherwise, it is not, then the sql operation will return error.
	ImageId   int64  `db:"Image_id"`
	UserId    int64  `db:"User_id"`
	ImageName string `db:"Image_name"`
	Tag       int    `db:"Tag"`
	Star      int    `db:"Star"`
	Fork      int    `db:"Fork"`
	Comm      int    `db:"Comment"`
	Status    int8   `db:"Status`
	Descrip   string `db:"Description"`
	Date      string `db:"Date"`
}

type CRComments struct {
	CommentId int64  `db:"comment_id"`
	ImageID   int64  `db:"image_id"`
	Author    string `db:"author"`
	Reply     string `db:"replyto"`
	Content   string `db:"content"`
}

type CRStar struct {
	StarId  int64 `db:"star_id"`
	ImageId int64 `db:"image_id"`
	UserId  int64 `db:"user_id"`
}

type CRFork struct {
	ForkId  int64 `db:"fork_id"`
	ImageId int64 `db:"image_id"`
	UserId  int64 `db:"user_id"`
}

type SqlOperation interface {
	Add() error
	QuerybyUser(uid int64) []CRImage
	QueryVerify(name string) bool
	Querylog(imageid int64)
	DeleteImg()
	UpdateStatus(status int8) error
	UpdateImage() error
	UpdateStar(uid int64) error
	UpdateFork(uid int64, uname string) error
}

//return a new CRImage struct by the input data
func newImage(uid int64, imgname string, tag int, des string) CRImage {
	return CRImage{
		UserId:    uid,
		ImageName: imgname,
		Tag:       tag,
		Star:      0,
		Fork:      0,
		Comm:      0,
		Status:    0,
		Descrip:   des,
		Date:      time.Now().Format("2006-01-02"),
	}
}

//list all the images
func QueryImage() []CRImage {
	var image []CRImage
	_, err := dbmap.Select(&image, "select * from cr_image")
	checkErr(err, "Select failed")
	return image
}

//fuzzy search of image list by image name
func QuerybyName(name string) []CRImage {
	var image []CRImage
	pattern := string("%" + name + "%")
	_, err := dbmap.Select(&image, "select * from cr_image where Image_name like ?", pattern)
	checkErr(err, "Select failed")
	return image
}

//insert a new record into cr_image table
func (c CRImage) Add() error {
	count, err := dbmap.SelectInt("select count(1) from cr_image where Image_name = ?", c.ImageName)
	if count > 0 || err != nil {
		return err
	}
	err = dbmap.Insert(&c)
	return err
}

//Query the image list by userid, return an array of CRImage struct
func (c CRImage) QuerybyUser(uid int64) []CRImage {
	var image []CRImage
	_, err := dbmap.Select(&image, "select * from cr_image where User_id = ?", uid)
	checkErr(err, "Select list failed")
	return image
}

//Query the log of an image by its id
func (c *CRImage) Querylog(imageid int64) *CRImage {
	obj, err := dbmap.Get(CRImage{}, imageid)
	if err != nil {
		log.Fatalln("Select log failed", err)
	}
	c = obj.(*CRImage)
	return obj.(*CRImage)
}

//Verify whether the name of image is existed
func QueryVerify(name string) bool {
	count, err := dbmap.SelectInt("select count(1) from cr_image where Image_name = ?", name)
	if err != nil {
		log.Fatalln("Verify failed", err)
		return false
	}
	if count < 1 {
		return true
	}
	return false
}

//Delete an image by its id, if it is forked from another image, delete the fork record too
func (c CRImage) DeleteImg() {
	_, err := dbmap.Delete(&c)
	if err != nil {
		log.Println("Delete failed", err)
		return
	}
	cf := new(CRFork)
	err = dbmap.SelectOne(&cf, "select fork_id from cr_fork where user_id = ? and image_id = ?", c.UserId, c.ImageId)
	if err != nil {
		return
	}
	_, err = dbmap.Delete(&cf)
	if err != nil {
		log.Println("Delete failed", err)
		return
	}
}

//set the status of image
func (c CRImage) UpdateStatus(status int8) error {
	log.Println(c.ImageName)
	_, err := dbmap.Exec("update cr_image set Status = ? WHERE Image_name = ? ", status, c.ImageName)
	return err
}

//Update the details of an image
func (c CRImage) UpdateImage() error {
	_, err := dbmap.Update(&c)
	return err
}

//update the star list of an image
func (c CRImage) UpdateStar(uid int64) error {
	//	if _, err := dbmap.Update(&c); err != nil {
	//		log.Println("Update image log failed", err)
	//	}
	var cs CRStar
	star := true
	trans, err := dbmap.Begin()
	if err != nil {
		return err
	}
	//查询是否存在star记录
	count, err := trans.SelectInt("select count(1) from cr_star where user_id = ? and image_id = ?", uid, c.ImageId)
	if err != nil {
		return err
	}
	//如果存在则查出那条记录
	if count > 0 {
		star = false
		err = trans.SelectOne(&cs, "select * from cr_star where user_id = ? and image_id = ?", uid, c.ImageId)
		if err != nil {
			trans.Rollback()
			return err
		}
	}
	//如果不存在记录，则插入一条，并使star数加一
	if star {
		//		err = dbmap.Insert(&cs)
		cs = CRStar{UserId: uid, ImageId: c.ImageId}
		if err := trans.Insert(&cs); err != nil {
			log.Println("Star failed", err)
			trans.Rollback()
			return err
		}
		_, err := trans.Exec("update cr_image set Star = Star + 1 WHERE Image_id = ? ", c.ImageId)
		if err != nil {
			log.Println("Star failed", err)
			trans.Rollback()
			return err
		}
	} else {
		//如果存在则将该记录删除，并使star数减一
		if _, err := trans.Delete(&cs); err != nil {
			trans.Rollback()
			return err
		}
		_, err := trans.Exec("update cr_image set Star = Star - 1 WHERE Image_id = ? ", c.ImageId)
		if err != nil {
			trans.Rollback()
			return err
		}
	}
	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		return err
	}
	return nil
}

//insert a fork record of an image
func (c CRImage) UpdateFork(uid int64, uname string) error {
	var cf CRFork
	//事务开始
	trans, err := dbmap.Begin()
	if err != nil {
		return err
	}
	//获得新镜像名称
	oldName := strings.Split(c.ImageName, "-")
	newName := uname + "-" + oldName[1]
	//检查是否已存在同名镜像
	count, err := trans.SelectInt("select count(1) from cr_image where User_id = ? and Image_name = ?", uid, newName)
	if err != nil || count > 0 {
		return err
	}
	//检查是否存在该fork记录
	count, err = trans.SelectInt("select count(1) from cr_fork where user_id = ? and image_id = ?", uid, c.ImageId)
	if err != nil {
		return err
	}
	//存在则退出
	if count > 0 {
		trans.Rollback()
		return err
	}
	//不存在，先插入一条cr_fork表记录
	cf = CRFork{UserId: uid, ImageId: c.ImageId}
	if err := trans.Insert(&cf); err != nil {
		trans.Rollback()
		return err
	}
	//镜像fork数量加一
	_, err = trans.Exec("update cr_image set Fork = Fork + 1 WHERE Image_id = ? ", c.ImageId)
	if err != nil {
		trans.Rollback()
		return err
	}
	//调用docker API，tag新的镜像
	oldImageName := c.ImageName + ":" + strconv.Itoa(c.Tag)
	ni := newImage(uid, newName, 1, c.Descrip)
	if err = ni.dockerFork(oldImageName); err != nil {
		trans.Rollback()
		return err
	}
	//插入新镜像记录
	//	oldName := strings.Split(c.ImageName, "-")
	//	newName := uname + "-" + oldName[1]
	//	ni := newImage(uid, newName, 1, c.Descrip)
	if err = trans.Insert(&ni); err != nil {
		trans.Rollback()
		return err
	}
	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		return err
	}
	return nil
}

//query whether there is a star log, if is, return the starid, else return 0
func (c CRStar) QueryStar() int64 {
	var cs CRStar
	//c.UserId here is the current user's id, not the image owner's id
	err := dbmap.SelectOne(&cs, "select star_id from cr_star where user_id = ? and image_id = ?", c.UserId, c.ImageId)
	//	count, err := dbmap.SelectInt("select count(1) from cr_star where user_id = ? and image_id = ?", cs.UserId, cs.ImageId)
	if err != nil {
		log.Println("Query starlog failed", err)
		return 0
	}
	return cs.StarId
}

//not consider the situation that the user is owner of image, but it is controller by the front end, and the function is only for query
func (c CRFork) QueryFork() bool {
	//c.UserId here is the current user's id, not the image owner's id
	count, err := dbmap.SelectInt("select count(1) from cr_fork where user_id = ? and image_id = ?", c.UserId, c.ImageId)
	if err != nil {
		log.Println("Query starlog failed", err)
		return true
	}
	if count > 0 {
		return true
	}
	return false
}

func init_imangeDb(db *gorp.DbMap) {
	db.AddTableWithName(CRImage{}, "cr_image").SetKeys(true, "ImageId")
	db.AddTableWithName(CRComments{}, "cr_comment").SetKeys(true, "CommentId")
	db.AddTableWithName(CRStar{}, "cr_star").SetKeys(true, "StarId")
	db.AddTableWithName(CRFork{}, "cr_fork").SetKeys(true, "ForkId")

}

/*
func checkErr(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
	}
}
*/
