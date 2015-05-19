package main

import (
	"fmt"
	"gopkg.in/gorp.v1"
	//	"log"
	//	"strconv"
	//	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type UserInfo struct {
	Id      int64  `db:"id"`
	UserId  int64  `db:"user_id"`
	Avatar  string `db:"avatar"`
	Discrip string `db:"discription"`
}

func (u UserInfo) updateInfo() error {
	_, err := dbmap.Update(u)
	if err != nil {
		logger.Println("update userinfo error: ", err)
		return err
	}
	return nil
}

func (u UserInfo) insertInfo() error {
	err := dbmap.Insert(u)
	if err != nil {
		logger.Println("insert userinfo error: ", err)
		return err
	}
	return nil
}

func (u *UserInfo) isExist(uid int64) (bool, error) {
	cmd := fmt.Sprintf("select count(*) from user_info where user_id=%d", uid)
	count, err := dbmap.SelectInt(cmd)
	if err != nil {
		logger.Println("select count userinfo error: ", err)
		return false, err
	}
	if count <= 0 {
		return false, nil
	} else {
		cmd = fmt.Sprintf("select * from user_info where user_id =%d", uid)
		err := dbmap.SelectOne(u, cmd)
		if err != nil {
			logger.Println("select userinfo error: ", err)
			return true, err
		}
		return true, nil
	}
}

func init_userDb(db *gorp.DbMap) {
	db.AddTableWithName(UserInfo{}, "user_info").SetKeys(true, "Id")
}
