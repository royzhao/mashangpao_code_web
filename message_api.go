package main

import (
	//	"fmt"
	//	"log"
	"github.com/codegangsta/martini"
	"strconv"
	//	"strings"
	"encoding/json"
	"net/http"
)

func queryNotice(w http.ResponseWriter, r *http.Request, parms martini.Params) {
	replyto, err := strconv.ParseInt(parms["id"], 10, 64)
	message, err := queryMessage(replyto)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(message); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type newMessageData struct {
	ReplyTo int64
	Author  int64
	Content string
	Level   int8
}

func addMessage(w http.ResponseWriter, r *http.Request) {
	var n newMessageData
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		logger.Warnf("error decoding params: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err := NewMessage(n.ReplyTo, n.Author, n.Content, n.Level)
	if err != nil {
		logger.Warnf("error sending message: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}
func readMessageAPI(w http.ResponseWriter, r *http.Request, parms martini.Params) {
	id, err := strconv.ParseInt(parms["id"], 10, 64)
	if err != nil {
		logger.Warnf("error converting id: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := readMessage(id); err != nil {
		logger.Warnf("error updating message status: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}
