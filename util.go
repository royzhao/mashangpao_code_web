package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func SetValue(key string, value interface{}) error {
	conn := pool.Get()
	defer conn.Close()
	buf, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, buf)
	return err
}

//return =1 not found,retrun =5 ok
func GetValue(key string) (int, string) {
	//	code_run_res_byte, _ := redis_client.Get(key)
	conn := pool.Get()
	defer conn.Close()
	code_run_res_byte, _ := conn.Do("GET", key)
	status := 1
	tmp, _ := code_run_res_byte.([]byte)
	//	code_run_res := string(code_run_res_byte)
	code_run_res := string(tmp)
	if code_run_res == "" {
		code_run_res = "commit successful, now running"
		//pull request
	} else {
		status = 5
	}
	return status, code_run_res
}
