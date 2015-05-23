package main

import (
	"encoding/json"
	//	"fmt"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/dylanzjy/coderun-request-client"
	"net/http"
	"net/url"
	"strconv"
)

var (
	ssoEndpoint  = "http://sso.learn4me.com"
	ssoClient, _ = client.NewSSOClient(ssoEndpoint)
)

type appInfo struct {
	App_id  int    `json:"app_id" yaml:"app_id"`
	App_key string `json:"app_key" yaml:"app_key"`
	Token   string `json:"token" yaml:"token"`
}

func isLogin(w http.ResponseWriter, r *http.Request) {
	var info appInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		logger.Warnf("error decoding login info: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//	fmt.Println(info)
	formInfo := url.Values{"app_id": {conf.App_id}, "app_key": {conf.App_key}, "token": {info.Token}}
	//	fmt.Println(formInfo)
	userData, err := ssoClient.IsLogin(formInfo)
	if err != nil {
		logger.Warnf("error querying login status: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//	if true {
	//		http.Redirect(w, r, "http://www.baidu.com", http.StatusFound)
	//		//		http.RedirectHandler("http://www.baidu.com", http.StatusMovedPermanently)
	//		return
	//	}
	if err := json.NewEncoder(w).Encode(userData); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetUserInfoByID(r *http.Request, enc Encoder, parms martini.Params) (int, string) {
	id, err := strconv.ParseInt(parms["userid"], 10, 64)
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the user with id %s does not exist", parms["userid"]))))
	}
	user, err := GetUserTotalInfoByID(id)
	fmt.Println(err)
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the user with id %s does not exist", parms["userid"]))))
	}
	return http.StatusOK, Must(enc.Encode(user))
}

//func logout(w http.ResponseWriter, r *http.Request) {

//}

//func main() {
//	c, err := client.NewSSOClient(ssoEndpoint)
//	cook := url.Values{"app_id": {"1"}, "app_key": {"Ei1F4LeTIUmJeFdO1MfbdkGQpZMeQ0CUX3aQD4kMOMVsRz7IAbjeBpurD6LTvNoI"}, "token": {"07016283de1ee8b2f55db4af920edd75"}}
//	data, err := c.IsLogin("POST", "/html/baigoSSO/mypage/user_identification.php", cook)
//	fmt.Println(data)
//	fmt.Println(err)
//	fmt.Println("Hello World!")

//}
