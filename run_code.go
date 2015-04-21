package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/dylanzjy/coderun-request-client"
	"log"
	"net/http"
)

type Run_code struct {
	Code Code_detail
	Meta Code_step
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
	imageid := parms["imagename"]
	decoder := json.NewDecoder(r.Body)
	var t Run_code
	err := decoder.Decode(&t)
	if err != nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("bad request"))))
	}
	log.Println(t)

	log.Println("imageid", imageid)
	//get image name
	cmdstr := imageid
	cmds := make([]client.Cmd_type, 4)
	for i, v := range t.Cmds {
		cmdstr += v.Cmd + v.Args
		cmds[i].Cmd = v.Cmd
		cmds[i].Args = v.Args
	}
	cmdstr += t.Code.Code_content
	cmdstr += imageid
	//compute md5 as id
	id := GetMd5String(cmdstr)
	data := &client.RunData{
		Id:      id,
		Workdir: t.Meta.Work_dir,
		Code: client.Code_type{
			Filename: t.Meta.Code_name,
			Content:  t.Code.Code_content,
		},
		Cmds: cmds,
	}
	status, code_run_res := GetValue(id)
	res := Run_res{
		Run_id: id,
		Res:    code_run_res,
		Status: status,
	}
	if dc == nil {
		dc, err = client.NewDockerClient(docker_end_point)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
	runout, err := dc.DockerRun(*data, imageid)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(runout)
	}

	// if al.Id == 0 {
	// 	return http.StatusNotFound, Must(enc.Encode(
	// 		NewError(ErrCodeNotExist, fmt.Sprintf("the Code step detail with id %s does not exist", parms["stepid"]))))
	// }
	return http.StatusOK, Must(enc.Encode(res))
}
