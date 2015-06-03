package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	qiniuconf "github.com/qiniu/api/conf"
	"github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
	"net/http"
)

func initqiniu() {
	qiniuconf.ACCESS_KEY = "c2mKwQp3ijoVCmatldqlhatX0hjAr00vhG1S-R5p"
	qiniuconf.SECRET_KEY = "Ih14QSNZ0d3BkNksUgMyto5n8HiGcSdVfg9SYt96"
}
func uptoken(bucketName string) string {
	putPolicy := rs.PutPolicy{
		Scope: bucketName,
	}
	return putPolicy.Token(nil)
}
func upload2qiniuHandler(r *http.Request) (key string, err error) {
	file, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	var ret io.PutRet
	var extra = &io.PutExtra{}
	token := uptoken("learn4me")

	// ret       变量用于存取返回的信息，详情见 io.PutRet
	// uptoken   为业务服务器端生成的上传口令
	// r         为io.Reader类型，用于从其读取数据
	// extra     为上传文件的额外信息,可为空， 详情见 io.PutExtra, 可选
	err = io.PutWithoutKey(nil, &ret, token, file, extra)

	if err != nil {
		//上传产生错误
		return "", err
	}
	// log.Print(ret.Hash, ret.Key)
	// fmt.Fprintf(w, "File uploaded successfully : ")
	// fmt.Fprintf(w, filename)
	return ret.Key, nil
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func DelKeyValue(key string) {
	conn := pool.Get()
	defer conn.Close()
	conn.Do("DEL", key)
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

//func convert Code_issue_comment to Code_issue_comment_info
func convertCode2Codeinfo(src []Code_issue_comment) []Code_issue_comment_info {
	lens := len(src)
	var u UserInfo
	des := make([]Code_issue_comment_info, lens)
	for i := 0; i < lens; i++ {
		des[i].Id = src[i].Id
		des[i].Issue_id = src[i].Issue_id
		des[i].Status = src[i].Status
		des[i].Reply_to = src[i].Reply_to
		des[i].Create_date = src[i].Create_date
		des[i].Content = src[i].Content
		des[i].Author, _ = u.getInfoFilter(int64(src[i].Author))
	}
	return des
}

func convertCode2CodeOne(src Code_issue) Code_issue_info {
	var des Code_issue_info
	var u UserInfo
	des.Author, _ = u.getInfoFilter(int64(src.Author))
	des.Code_id = src.Code_id
	des.Content = src.Content
	des.Create_date = src.Create_date
	des.Id = src.Id
	des.Status = src.Status
	des.Title = src.Title
	return des
}

func convertImage2Imageinfo(src []Image_issue_comment) []Image_issue_comment_info {
	lens := len(src)
	var u UserInfo
	des := make([]Image_issue_comment_info, lens)
	for i := 0; i < lens; i++ {
		des[i].Id = src[i].Id
		des[i].Issue_id = src[i].Issue_id
		des[i].Status = src[i].Status
		des[i].Reply_to = src[i].Reply_to
		des[i].Create_date = src[i].Create_date
		des[i].Content = src[i].Content
		des[i].Author, _ = u.getInfoFilter(src[i].Author)
	}
	return des
}

func convertImage2ImageOne(src Image_issue) Image_issue_info {
	var des Image_issue_info
	var u UserInfo
	des.Author, _ = u.getInfoFilter(src.Author)
	des.Image_id = src.Image_id
	des.Content = src.Content
	des.Create_date = src.Create_date
	des.Id = src.Id
	des.Status = src.Status
	des.Title = src.Title
	return des
}
