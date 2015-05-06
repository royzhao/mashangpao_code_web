package main

import (
	//	"fmt"
	"gopkg.in/gorp.v1"
	//	"log"
	//	"strconv"
	//	"strings"
	//	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Message struct {
	Id      int64  `db:"id"`
	ReplyTo int64  `db:"replyto"`
	Author  int64  `db:"author"`
	Content string `db:"content"`
	Date    string `db:"date"`
	Status  int8   `db:"status"`
	Level   int8   `db:"level"`
}

func init_messageDb(db *gorp.DbMap) {
	db.AddTableWithName(Message{}, "message").SetKeys(true, "Id")
}
