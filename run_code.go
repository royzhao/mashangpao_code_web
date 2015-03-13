package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
)

type Run_code struct {
	Code    string
	Imageid int
}
type Run_res struct {
	Res    string
	Run_id string
}

func RunCodeStep(w http.ResponseWriter, r *http.Request, enc Encoder, parms martini.Params, db codeStepDB_inter) (int, string) {
	decoder := json.NewDecoder(r.Body)
	var t Run_code
	err := decoder.Decode(&t)
	if err != nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("bad request"))))
	}
	log.Println(t)
	//compute md5 as id
	id := "dsdasdadasd"
	res := Run_res{
		Run_id: id,
	}
	// if al.Id == 0 {
	// 	return http.StatusNotFound, Must(enc.Encode(
	// 		NewError(ErrCodeNotExist, fmt.Sprintf("the Code step detail with id %s does not exist", parms["stepid"]))))
	// }
	return http.StatusOK, Must(enc.Encode(res))
}
