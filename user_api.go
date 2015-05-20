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
	var user UserInfo
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var tmp UserInfo
	result, err := tmp.isExist(user.UserId)
	if err != nil {
		logger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if result {
		err = user.updateInfo()
	} else {
		err = user.insertInfo()
	}
	if err != nil {
		logger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getUserInfo(w http.ResponseWriter, r *http.Request, parms martini.Params) {
	uid, err := strconv.ParseInt(parms["uid"], 10, 64)
	if err != nil {
		logger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var u UserInfo
	_, err = u.isExist(uid)
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
