package main

import (
	//	"fmt"
	"gopkg.in/gorp.v1"
	//	"log"
	//	"strconv"
	//	"strings"
	"time"

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

func NewMessage(replyTo int64, author int64, content string, level int8) (Message, error) {
	m := Message{
		ReplyTo: replyTo,
		Author:  author,
		Content: content,
		Date:    time.Now().Format("2006-01-02"),
		Status:  1,
		Level:   level,
	}
	if level == 3 {
		sendMail(replyTo, author, content)
	}
	err := dbmap.Insert(&m)
	return m, err
}

func (m Message) readMessage() error {
	m.Status = 2
	_, err := dbmap.Update(m)
	return err
}

func queryMessage(replyTo int64) ([]Message, error) {
	var m []Message
	_, err := dbmap.Select(&m, "select * from message where replyto = ? and status = 1", replyTo)
	return m, err
}

func sendMail(replyTo int64, author int64, content string) {

}
