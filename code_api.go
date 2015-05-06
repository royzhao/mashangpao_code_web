package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetIssues(r *http.Request, enc Encoder, db codeDB_inter, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["codeid"])
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the Code with id %s does not exist", parms["userid"]))))
	}
	// Get the query string arguments, if any
	qs := r.URL.Query()
	key := qs.Get("key")
	page := qs.Get("page")
	num := qs.Get("num")
	_num, err := strconv.Atoi(num)
	if err != nil {
		_num = 5
	}
	if _num == 0 {
		_num = 5
	}
	_page, err := strconv.Atoi(page)
	if err != nil {
		_page = 1
	}
	if _page <= 0 {
		_page = 1
	}
	// Otherwise, return all Codes
	return http.StatusOK, Must(enc.Encode(db.FindIssues(key, _page, _num, id)))
}
func AddIssue(r *http.Request, enc Encoder, db codeDB_inter, parms martini.Params) (int, string) {
	codeid, err := strconv.Atoi(parms["codeid"])
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the code with id %s does not exist", parms["codeid"]))))
	}
	al, err := getPostCodeIssue(r)
	if err != nil {
		return http.StatusBadRequest, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the issue create failed no codeid %s", parms["codeid"]))))
	}
	al.Code_id = codeid
	id, err := db.AddOneIssue(al)
	if err != nil {
		return http.StatusBadRequest, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the issue create failed"))))
	}
	al.Id = id
	return http.StatusCreated, Must(enc.Encode(al))
}
func DeleteIssue(r *http.Request, enc Encoder, db codeDB_inter, parms martini.Params) (int, string) {
	issue, err := strconv.Atoi(parms["issueid"])
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the issue with id %s does not exist", parms["issueid"]))))
	}
	err = db.DeleteIssueById(issue)
	if err != nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the issue with id %s does not exist", parms["issueid"]))))
	}
	return http.StatusOK, fmt.Sprintf("delete issue=%d is ok", issue)
}
func UpdateIssue(r *http.Request, enc Encoder, db codeDB_inter, parms martini.Params) (int, string) {
	codeid, err := strconv.Atoi(parms["codeid"])
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the code with id %s does not exist", parms["codeid"]))))
	}
	issue, err := strconv.Atoi(parms["issueid"])
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the issue with id %s does not exist", parms["issueid"]))))
	}
	al, err := getPostCodeIssue(r)
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("Update failed"))))
	}
	al.Code_id = codeid
	al.Id = issue
	err = db.UpdateIssueById(al)
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("Update failed"))))
	}
	return http.StatusOK, Must(enc.Encode(al))
}
func GetIssueComments(r *http.Request, enc Encoder, db codeDB_inter, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["issueid"])
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the issue with id %s does not exist", parms["userid"]))))
	}
	// Get the query string arguments, if any
	qs := r.URL.Query()
	key := qs.Get("key")
	page := qs.Get("page")
	num := qs.Get("num")
	_num, err := strconv.Atoi(num)
	if err != nil {
		_num = 5
	}
	if _num == 0 {
		_num = 5
	}
	_page, err := strconv.Atoi(page)
	if err != nil {
		_page = 1
	}
	if _page <= 0 {
		_page = 1
	}
	// Otherwise, return all Codes
	return http.StatusOK, Must(enc.Encode(db.FindIssueComment(key, _page, _num, id)))
}
func AddIssueComment(r *http.Request, enc Encoder, db codeDB_inter, parms martini.Params) (int, string) {
	issueid, err := strconv.Atoi(parms["issueid"])
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("Add failed"))))
	}
	al, err := getPostCodeIssueComment(r)
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("Add failed"))))
	}
	al.Issue_id = issueid
	id, err := db.AddOneIssueComment(al)
	if err != nil {
		return http.StatusBadRequest, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the issue create failed"))))
	}
	al.Id = id
	return http.StatusCreated, Must(enc.Encode(al))
}
func DeleteIssueComment(r *http.Request, enc Encoder, db codeDB_inter, parms martini.Params) (int, string) {
	commentid, err := strconv.Atoi(parms["commentid"])
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the commentid with id %s does not exist", parms["issue"]))))
	}
	err = db.DeleteIssueComment(commentid)
	if err != nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the commentid with id %s does not exist", parms["issue"]))))
	}
	return http.StatusOK, fmt.Sprintf("delete commentid=%d is ok", commentid)
}
func UpdateIssueComment(r *http.Request, enc Encoder, db codeDB_inter, parms martini.Params) (int, string) {
	issueid, err := strconv.Atoi(parms["issueid"])
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("Add failed"))))
	}
	commentid, err := strconv.Atoi(parms["commentid"])
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the commentid with id %s does not exist", parms["issue"]))))
	}
	al, err := getPostCodeIssueComment(r)
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("Add failed"))))
	}
	al.Issue_id = issueid
	al.Id = commentid
	err = db.UpdateIssueComment(al)
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("Update failed"))))
	}
	return http.StatusOK, Must(enc.Encode(al))
}

// GetCodes returns the list of codes (possibly filtered).
func GetCodes(r *http.Request, enc Encoder, db codeDB_inter) (int, string) {
	// Get the query string arguments, if any
	qs := r.URL.Query()
	key := qs.Get("key")
	page := qs.Get("page")
	num := qs.Get("num")
	_num, err := strconv.Atoi(num)
	if err != nil {
		_num = 5
	}
	if _num == 0 {
		_num = 5
	}
	_page, err := strconv.Atoi(page)
	if err != nil {
		_page = 1
	}
	if _page <= 0 {
		_page = 1
	}
	// Otherwise, return all Codes
	return http.StatusOK, Must(enc.Encode(db.Find(key, _page, _num, -1)))
}

// GetCodes returns the list of codes (possibly filtered).
func GetCodesByUser(r *http.Request, enc Encoder, db codeDB_inter, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["userid"])
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the Code with id %s does not exist", parms["userid"]))))
	}
	// Get the query string arguments, if any
	qs := r.URL.Query()
	key, page, num := qs.Get("key"), qs.Get("page"), qs.Get("num")
	_num, err := strconv.Atoi(num)
	if err != nil {
		_num = 5
	}
	if _num == 0 {
		_num = 5
	}
	_page, err := strconv.Atoi(page)
	if err != nil {
		_page = 1
	}
	if _page <= 0 {
		_page = 1
	}
	// return http.StatusOK, Must(enc.Encode(toIface(db.Find(key, _page, _num, id))...))
	return http.StatusOK, Must(enc.Encode(db.Find(key, _page, _num, id)))
}

// GetCode returns the requested Code.
func GetCode(enc Encoder, db codeDB_inter, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["codeid"])
	al := db.Get(id)
	if err != nil || al.Id == 0 {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the Code with id %s does not exist", parms["codeid"]))))
	}

	return http.StatusOK, Must(enc.Encode(al))
}

// Addcode creates the posted code.
func AddCode(w http.ResponseWriter, r *http.Request, enc Encoder, parms martini.Params, db codeDB_inter) (int, string) {
	userid := parms["userid"]
	al := getPostCode(r, userid)
	id, err := db.Add(al)
	switch err {
	case ErrAlreadyExists:
		// Duplicate
		return http.StatusConflict, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the code '%s' from '%s' already exists", al.Name, al.Description))))
	case nil:
		// TODO : Location is expected to be an absolute URI, as per the RFC2616
		w.Header().Set("Location", fmt.Sprintf("/code/%s/%d", userid, id))
		return http.StatusCreated, Must(enc.Encode(al))
	default:
		panic(err)
	}
}

func UpdateCodeStar(r *http.Request, enc Encoder, db codeDB_inter, parms martini.Params) (int, string) {
	log.Println("update code star")
	userid, err := strconv.Atoi(parms["userid"])
	if err != nil {
		// Invalid id, 404
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("invalid user %s", parms["userid"]))))
	}
	codeid, err := strconv.Atoi(parms["codeid"])
	if err != nil {
		// Invalid id, 404
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("invalid codeid %s", parms["codeid"]))))
	}
	ret, err := db.UpdateStar(userid, codeid)
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("star failed"))))
	}
	return http.StatusOK, Must(enc.Encode(ret))
}

// UpdateCode changes the specified code.
func UpdateCode(r *http.Request, enc Encoder, db codeDB_inter, parms martini.Params) (int, string) {
	al, err := getPutCode(r, parms)
	if err != nil {
		// Invalid id, 404
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the code with id %s does not exist", parms["codeid"]))))
	}
	err = db.Update(al)
	a := db.Get(al.Id)
	switch err {
	case ErrAlreadyExists:
		return http.StatusConflict, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the code '%s' from '%s' already exists", a.Name, a.Description))))
	case nil:
		return http.StatusOK, Must(enc.Encode(a))
	default:
		panic(err)
	}
}

func DeleteCode(enc Encoder, db codeDB_inter, parms martini.Params) (int, string) {
	userid, err := strconv.Atoi(parms["userid"])
	id, err := strconv.Atoi(parms["codeid"])
	al := db.Get(id)
	if err != nil || al.Name == "" {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the code with id %d does not exist,user %d", id, userid))))
	}
	db.Delete(id)
	return http.StatusOK, fmt.Sprintf("delete code=%d is ok", id)
}

// Parse the request body, load into an Code structure.
func getPostCodeIssue(r *http.Request) (*Code_issue, error) {
	decoder := json.NewDecoder(r.Body)
	var t Code_issue
	err := decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	t.Create_date = time.Now().String()
	return &t, nil
}
func getPostCodeIssueComment(r *http.Request) (*Code_issue_comment, error) {
	decoder := json.NewDecoder(r.Body)
	var t Code_issue_comment
	err := decoder.Decode(&t)
	if err != nil {
		return nil, err
	}
	t.Create_date = time.Now().String()
	return &t, nil
}
func getPostCode(r *http.Request, user_id string) *Code {
	decoder := json.NewDecoder(r.Body)
	var t Code
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t)
	t.User_id, err = strconv.Atoi(user_id)
	if err != nil {
		panic(err)
	}
	t.Create_date = time.Now().Local().Format("2006-01-02 15:04:05 +0800")
	return &t
}

// Like getPostCode, but additionnally, parse and store the `id` query string.
func getPutCode(r *http.Request, parms martini.Params) (*Code, error) {
	al := getPostCode(r, parms["userid"])
	id, err := strconv.Atoi(parms["codeid"])
	if err != nil {
		return nil, err
	}
	al.Id = id
	return al, nil
}
func toIface(v []Code) []interface{} {
	if len(v) == 0 {
		return nil
	}
	ifs := make([]interface{}, len(v))
	for i, v := range v {
		ifs[i] = v
	}
	return ifs
}
