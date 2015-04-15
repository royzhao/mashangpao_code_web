package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
)

type Run_code struct {
	Code string
	Meta Code_step_meta
	Cmds []Code_step_cmd
}
type Run_res struct {
	Res    string `json:"res"`
	Status int    `json:"status"`
	Run_id string `json:"run_id"`
}

func GetRunResult(w http.ResponseWriter, r *http.Request, enc Encoder, parms martini.Params, db codeStepDB_inter) (int, string) {
	runid := parms["runid"]

	//mock rpc
	log.Println("query cache use ", runid)
	status, code_res := GetValue(runid)
	res := Run_res{
		Res:    code_res,
		Status: status,
		Run_id: runid,
	}
	return http.StatusOK, Must(enc.Encode(res))
}
func RunCodeStep(w http.ResponseWriter, r *http.Request, enc Encoder, parms martini.Params, db codeStepDB_inter) (int, string) {
	imageid := parms["imageid"]
	decoder := json.NewDecoder(r.Body)
	var t Run_code
	err := decoder.Decode(&t)
	if err != nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("bad request"))))
	}
	log.Println(t)
	log.Println("imageid", imageid)
	cmdstr := ""
	for _, v := range t.Cmds {
		cmdstr += v.Cmd + v.Args
	}
	cmdstr += t.Code
	cmdstr += imageid
	//compute md5 as id
	id := GetMd5String(cmdstr)
	status, code_run_res := GetValue(id)
	res := Run_res{
		Run_id: id,
		Res:    code_run_res,
		Status: status,
	}
	// if al.Id == 0 {
	// 	return http.StatusNotFound, Must(enc.Encode(
	// 		NewError(ErrCodeNotExist, fmt.Sprintf("the Code step detail with id %s does not exist", parms["stepid"]))))
	// }
	return http.StatusOK, Must(enc.Encode(res))
}
