package main

import (
	"encoding/json"
	//	"fmt"
	//	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/martini"
	"net/http"
	"strconv"
	//	"time"
)

func updateUserInfo(w http.ResponseWriter, r *http.Request) {

}

func getUserInfo(w http.ResponseWriter, r *http.Request, parms martini.Params) {
	uid, err := strconv.ParseInt(parms["uid"], 10, 64)
	if err != nil {
		logger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var u UserInfo
	result, err := u.isExist(uid)
	if err != nil {
		logger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(u); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
