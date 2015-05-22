package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/gorp.v1"
	//	"log"
	//	"strconv"
	//	"strings"
	"github.com/dylanzjy/coderun-request-client"

	_ "github.com/go-sql-driver/mysql"
)

type UserInfo struct {
	Id      int64  `db:"id"`
	UserId  int64  `db:"user_id"`
	Avatar  string `db:"avatar"`
	Discrip string `db:"discription"`
}

type UserTotalData struct {
	SSOmeta client.UserInfo `json:"sso_meta"`
	Info    UserInfo        `json:"info"`
}

//get user info from cache by user id
func GetUserinfoByCache(id int64) (*UserTotalData, error) {
	key := fmt.Sprintf("user_%d", id)
	var user UserTotalData
	var err error
	status, data := GetValue(key)
	if status == 5 {
		//ok
		err = json.Unmarshal([]byte(data), &user)
		if err == nil {
			return &user, nil
		}
	}
	return nil, err
}

//get user info from sso by user id
func GetUserTotalInfoByID(id int64) (*client.UserInfo, error) {
	user_key := fmt.Sprintf("act_get=get&user_by=user_id&user_id=%d", id)
	user, err := ssoClient.GetUserInfo(conf.App_id, conf.App_key, user_key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &user, nil
}
func (u UserInfo) updateInfo() error {
	_, err := dbmap.Update(&u)
	if err != nil {
		logger.Println("update userinfo error: ", err)
		return err
	}
	return nil
}

func (u UserInfo) insertInfo() error {
	err := dbmap.Insert(&u)
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

func (u *UserInfo) getInfo(uid int64) (*UserTotalData, error) {
	//get info from redis
	user, _ := GetUserinfoByCache(uid)
	if user != nil {
		return user, nil
	}
	userchan := make(chan *client.UserInfo, 1)
	//query it
	go func() {
		total, err := GetUserTotalInfoByID(uid)
		if err != nil {
			userchan <- nil
		}
		userchan <- total
	}()
	res, _ := u.isExist(uid)
	if res == false {
		return nil, NewError(1, "no such user info in web ")
	}
	var data UserTotalData
	data.Info = *u
	ssodata := <-userchan
	if ssodata == nil {
		return nil, NewError(1, "no such user info in sso")
	}
	data.SSOmeta = *ssodata
	return &data, nil

}
func init_userDb(db *gorp.DbMap) {
	db.AddTableWithName(UserInfo{}, "user_info").SetKeys(true, "Id")
}
