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

// // get code steps
func GetCodeSteps(r *http.Request, enc Encoder, db codeStepDB_inter, parms martini.Params) (int, string) {
	// Otherwise, return all Codes
	id, err := strconv.Atoi(parms["codeid"])
	if err != nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the Code with id %s does not exist", parms["codeid"]))))
	}
	return http.StatusOK, Must(enc.Encode(toIfaceStep(db.GetAll(id))))
}

// //get code step
func GetCodeStep(r *http.Request, enc Encoder, db codeStepDB_inter, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["stepid"])
	if err != nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the Code step with id %s does not exist", parms["stepid"]))))
	}
	al := db.Get(id)
	if al.Meta.Id == 0 {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the Code step with id %s does not exist", parms["stepid"]))))
	}
	return http.StatusOK, Must(enc.Encode(al))
}
func GetCodeStepCmd(r *http.Request, enc Encoder, db codeStepDB_inter, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["stepid"])
	if err != nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the Code step with id %s cmd does not exist", parms["stepid"]))))
	}
	al := db.GetCodeCmds(id)
	if len(al) == 0 {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the Code step with id %s cmd does not exist", parms["stepid"]))))
	}
	return http.StatusOK, Must(enc.Encode(al))
}

// // Addcode creates the posted code step.
func AddCodeStep(w http.ResponseWriter, r *http.Request, enc Encoder, parms martini.Params, db codeStepDB_inter) (int, string) {
	codeid := parms["codeid"]
	// userid := parms["userid"]
	al := getPostCodeStep(r, codeid)
	id, err := db.Add(al)
	switch err {
	case ErrAlreadyExists:
		// Duplicate
		return http.StatusConflict, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the code step '%s' from '%s' already exists", al.Name, al.Description))))
	case nil:
		// TODO : Location is expected to be an absolute URI, as per the RFC2616
		res := db.Get(id)
		if res.Meta.Id == 0 {
			return http.StatusNotFound, Must(enc.Encode(
				NewError(ErrCodeNotExist, fmt.Sprintf("create failed"))))
		}
		return http.StatusCreated, Must(enc.Encode(res))
	default:
		panic(err)
	}
}

func UpdateCodeStepCmd(r *http.Request, enc Encoder, db codeStepDB_inter, parms martini.Params) (int, string) {
	stepid, err := strconv.Atoi(parms["stepid"])
	if err != nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the Code step with id %s cmd does not exist", parms["stepid"]))))
	}
	al, err := getPutCodeStepCmd(r, parms)
	if err != nil {
		// Invalid id, 404
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("bad request"))))
	}
	err = db.UpdateCodeCmd(stepid, al)
	res := db.Get(stepid)
	switch err {
	case nil:
		return http.StatusOK, Must(enc.Encode(res))
	default:
		panic(err)
	}
}

// UpdateCode changes the specified code.
func UpdateCodeStep(r *http.Request, enc Encoder, db codeStepDB_inter, parms martini.Params) (int, string) {
	al, err := getPutCodeSetp(r, parms)
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
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the code step '%s' from '%s' already exists", a.Meta.Name, a.Meta.Description))))
	case nil:
		if a.Meta.Id == 0 {
			return http.StatusNotFound, Must(enc.Encode(
				NewError(ErrCodeNotExist, fmt.Sprintf("not found obj %d", al.Id))))
		}
		return http.StatusOK, Must(enc.Encode(a))
	default:
		panic(err)
	}
}

func DeleteCodeStep(enc Encoder, db codeStepDB_inter, parms martini.Params) (int, string) {
	userid, err := strconv.Atoi(parms["userid"])
	stepid, err := strconv.Atoi(parms["stepid"])
	al := db.Get(stepid)
	if err != nil || al.Meta.Name == "" {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the code step with id %d does not exist,user %d", stepid, userid))))
	}
	db.Delete(stepid)
	return http.StatusOK, fmt.Sprintf("delete stepid=%d is ok", stepid)
}

// //get code step
func GetCodeStepDetail(r *http.Request, enc Encoder, db codeStepDB_inter, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["stepid"])
	if err != nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the Code step detail with id %s does not exist", parms["stepid"]))))
	}
	al := db.GetStepDetail(id)
	if al.Meta.Id == 0 {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the Code step detail with id %s does not exist", parms["stepid"]))))
	}
	return http.StatusOK, Must(enc.Encode(al))
}

// UpdateCode changes the specified code.
func UpdateCodeStepDetail(r *http.Request, enc Encoder, db codeStepDB_inter, parms martini.Params) (int, string) {
	al, err := getPutCodeSetpDetail(r, parms)
	if err != nil {
		// Invalid id, 404
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the code with id %s does not exist", parms["codeid"]))))
	}
	err = db.UpdateStepDetail(al)
	a := db.GetStepDetail(al.Id)
	switch err {
	case ErrAlreadyExists:
		return http.StatusConflict, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the code step detail '%d' already exists", a.Meta.Id))))
	case nil:
		return http.StatusOK, Must(enc.Encode(a))
	default:
		panic(err)
	}
}

// Parse the request body, load into an Code structure.
func getPostCodeStep(r *http.Request, code_id string) *Code_step {
	decoder := json.NewDecoder(r.Body)
	var t Code_step
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t)
	create_date := time.Now().Local().Format("2006-01-02 15:04:05 +0800")
	codeId, err := strconv.Atoi(code_id)

	if t.Status == 0 {
		t.Status = 1
	}
	if err != nil {
		return nil
	}
	t.Code_id = codeId
	t.Create_date = create_date
	return &t
}

// Like getPostCode, but additionnally, parse and store the `id` query string.
func getPutCodeSetp(r *http.Request, parms martini.Params) (*Code_step, error) {
	al := getPostCodeStep(r, parms["codeid"])
	id, err := strconv.Atoi(parms["stepid"])
	if err != nil {
		return nil, err
	}
	al.Id = id
	return al, nil
}

func getPutCodeSetpDetail(r *http.Request, parms martini.Params) (*Code_detail, error) {
	id, err := strconv.Atoi(parms["stepid"])
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(r.Body)
	var t Code_detail
	err = decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t)
	t.Id = id
	return &t, nil
}

func getPutCodeStepCmd(r *http.Request, parms martini.Params) ([]Code_step_cmd, error) {
	id, err := strconv.Atoi(parms["stepid"])
	if err != nil {
		return nil, err
	}
	log.Println("stepid" + parms["stepid"])
	decoder := json.NewDecoder(r.Body)
	var t []Code_step_cmd
	err = decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t)
	for i, _ := range t {
		t[i].Stepid = id
		t[i].Id = 0
	}
	log.Println(t)
	return t, nil
}
func toIfaceStep(v []Code_step) []interface{} {
	if len(v) == 0 {
		return nil
	}
	ifs := make([]interface{}, len(v))
	for i, v := range v {
		ifs[i] = v
	}
	return ifs
}
