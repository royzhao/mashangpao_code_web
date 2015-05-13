package main

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetValue(key string) (int, string) {
	//	code_run_res_byte, _ := redis_client.Get(key)
	code_run_res_byte, _ := redis_client.Do("GET", key)
	status := 1
	//	code_run_res := string(code_run_res_byte)
	code_run_res := code_run_res_byte.(string)
	if code_run_res == "" {
		code_run_res = "commit successful, now running"
		//pull request
	} else {
		status = 5
	}
	return status, code_run_res
}
