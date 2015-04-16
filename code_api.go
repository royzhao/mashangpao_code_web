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

// GetCodes returns the list of codes (possibly filtered).
func GetCodes(r *http.Request, enc Encoder, db codeDB_inter) string {
	// Get the query string arguments, if any
	qs := r.URL.Query()
	name, description, create_date := qs.Get("name"), qs.Get("description"), qs.Get("create_date")

	if name != "" || description != "" || create_date != "" {
		// At least one filter, use Find()
		return Must(enc.Encode(toIface(db.Find(name, description, create_date, -1))...))
	}
	// Otherwise, return all Codes
	return Must(enc.Encode(toIface(db.GetAll())...))
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
	name, description, create_date := qs.Get("name"), qs.Get("description"), qs.Get("create_date")

	if name != "" || description != "" || create_date != "" {
		// At least one filter, use Find()
		return http.StatusOK, Must(enc.Encode(toIface(db.Find(name, description, create_date, id))...))
	}
	// Otherwise, return all Codes
	return http.StatusOK, Must(enc.Encode(toIface(db.Find("", "", "", id))...))
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
func getPostCode(r *http.Request, user_id string) *Code {
	decoder := json.NewDecoder(r.Body)
	var t Code
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t)
	// name, description := r.FormValue("name"), r.FormValue("description")
	// id := user_id
	// create_date := time.Now().Local().Format("2006-01-02 15:04:05 +0800")
	// userid, err := strconv.Atoi(id)
	// if err != nil {
	// 	return nil
	// }
	// return &Code{
	// 	Name:        name,
	// 	Description: description,
	// 	Create_date: create_date,
	// 	User_id:     userid,
	// }
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
