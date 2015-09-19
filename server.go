package main

import (
	"flag"
	"github.com/codegangsta/martini"
	"github.com/dylanzjy/coderun-request-client"
	"github.com/fsouza/go-dockerclient"
	//	"github.com/hoisie/redis"
	"gopkg.in/gorp.v1"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	//	"runtime"
	// "strconv"
	"strings"
	"time"
	//	"github.com/codegangsta/martini-contrib/auth"

	"github.com/garyburd/redigo/redis"
	//	"github.com/youtube/vitess/go/pools"
	//	"golang.org/x/net/context"
)

type ResourceConn struct {
	redis.Conn
}

func (r ResourceConn) Close() {
	r.Conn.Close()
}

var (
	addr  = flag.String("p", ":9000", "Address and port to serve dockerui")
	dbmap *gorp.DbMap

	// 只有一个martini实例
	m *martini.Martini
	//	redis_client redis.Client
	//	redis_pool     *pools.ResourcePool
	//	redis_resource pools.Resource
	//	redis_client   ResourceConn
	pool          *redis.Pool
	redisServer   string
	redisPassword string

	//docker proxy

	dc   *client.DockerClient
	lb   *client.LBClient
	conf Configuration

	//hot image list timer
	timer = time.NewTicker(12 * time.Hour)
)

func init() {
	initqiniu()
	var err error
	conf, err = ReadConfigure("conf.json")
	if err != nil {
		panic(err)
	}
	endpoint = conf.Endpoint
	dockerclient, _ = docker.NewClient(endpoint)
	ssoClient, _ = client.NewSSOClient(conf.SSO)
	browserEndpoint = conf.BrowserEndpoint
	dockerhub = conf.Dockerhub
	redis_addr := conf.Redis_addr
	lb, err = client.NewLBClient(conf.LB_Addr)
	if err != nil {
		panic(err)
	}
	dc = nil
	log.Println(dc)
	redisServer = redis_addr
	redisPassword = ""
	// redis_addr := os.Getenv("REDIS_ADDR")
	// log.Println("redis addr is:" + redis_addr)
	// if redis_addr == "" {
	// 	redis_addr = "redis.peilong.me:6379"
	// }
	//	redis_client.Addr = redis_addr

	//	redis_pool = pools.NewResourcePool(func() (pools.Resource, error) {
	//		c, err := redis.Dial("tcp", redis_addr)
	//		return ResourceConn{c}, err
	//	}, 1, 5, time.Minute)
	//	ctx := context.TODO()
	//	redis_resource, err = redis_pool.Get(ctx)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	redis_client = redis_resource.(ResourceConn)

	// dbmap = nil
	// //init database
	dbmap = initDb(conf.DB_addr)
	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
	log.Println("init db is successful!")
	pool = newPool(redisServer, redisPassword)

	// initrundb()
	//init code module
	//init_code()

	m = martini.New()
	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	//static
	m.Use(martini.Static("public"))
	//m.Use(auth.Basic(AuthToken, ""))

	m.Use(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			var Origin = req.Header.Get("Origin")
			if strings.Contains(Origin, "learn4me.com") {
				res.Header().Set("Access-Control-Allow-Origin", Origin)
				res.Header().Set("Access-Control-Allow-Methods", req.Header.Get("Access-Control-Request-Method"))
				res.Header().Set("Access-Control-Max-Age", "60")
				res.Header().Set("Access-Control-Allow-Headers", req.Header.Get("Access-Control-Request-Headers"))
				res.Header().Set("Access-Control-Allow-Credentials", "true")
				if req.Method == "OPTIONS" {
					res.WriteHeader(http.StatusOK)
					return
				}
			}
			is_auth := strings.Split(req.URL.Path, "coderun")
			log.Println(is_auth)
			if len(is_auth) <= 1 {
				token := req.Header.Get("x-session-token")
				log.Println(token)
				if token == "" {
					res.WriteHeader(http.StatusUnauthorized)
					return
				}
				log.Println(conf.App_key)
				formInfo := url.Values{"app_id": {conf.App_id}, "app_key": {conf.App_key}, "token": {token}}
				userData, err := ssoClient.IsLogin(formInfo)
				if err != nil {
					log.Println(err)
					res.WriteHeader(http.StatusUnauthorized)
					return
				}
				log.Println(userData)
				if userData.Is_login == "false" {
					res.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

		}
	})

	m.Use(MapEncoder)
	m.Use(AddCorsHandler)
	r := martini.NewRouter()
	//得到所有的代码
	r.Get(`/api/code`, GetCodes)
	//得到某一个用户的所有代码
	r.Get(`/api/user/code/:userid`, GetCodesByUser)
	//get user info by id
	r.Get(`/api/user/:userid/info`, GetUserAvatarByID)
	//一个用户增加一个代码
	r.Post(`/api/user/code/:userid`, AddCode)
	//查询一个指定的代码
	r.Get(`/api/code/:codeid`, GetCode)
	//修改指定的代码
	r.Put(`/api/code/:userid/:codeid`, UpdateCode)
	//删除指定的代码
	r.Delete(`/api/code/:userid/:codeid`, DeleteCode)

	r.Put(`/api/code/star/:userid/:codeid`, UpdateCodeStar)
	//code issue
	r.Get(`/api/code/:codeid/issues`, GetIssues)
	r.Post(`/api/code/:userid/:codeid/issue`, AddIssue)
	r.Put(`/api/code/:userid/:codeid/issue/:issueid`, UpdateIssue)
	r.Delete(`/api/code/:userid/:codeid/issue/:issueid`, DeleteIssue)
	//code issue comment
	r.Get(`/api/issue/:issueid/comments`, GetIssueComments)
	r.Post(`/api/issue/:userid/:issueid/comment`, AddIssueComment)
	r.Put(`/api/issue/:userid/:issueid/comment/:commentid`, UpdateIssueComment)
	r.Delete(`/api/issue/:userid/:issueid/comment/:commentid`, DeleteIssueComment)
	//得到代码的具体步骤
	//得到全部的代码步骤元数据
	r.Get(`/api/code/:codeid/step`, GetCodeSteps)
	//得到某一个具体步骤
	r.Get(`/api/code/:codeid/step/:stepid`, GetCodeStep)
	//增加一个代码的具体步骤
	r.Post(`/api/code/:userid/:codeid/step`, AddCodeStep)
	//修改置顶的代码步骤的元数据
	r.Put(`/api/code/:userid/:codeid/step/:stepid`, UpdateCodeStep)
	//删除
	r.Delete(`/api/code/:userid/:codeid/step/:stepid`, DeleteCodeStep)

	//具体内容操作
	r.Get(`/api/code/:codeid/step/:stepid/detail`, GetCodeStepDetail)
	//修改具体内容
	r.Put(`/api/code/:userid/:codeid/step/:stepid/detail`, UpdateCodeStepDetail)

	//add command
	r.Put(`/api/code/:userid/:codeid/step/:stepid/cmd`, UpdateCodeStepCmd)
	r.Get(`/api/code/:codeid/step/:stepid/cmd`, GetCodeStepCmd)

	//code run
	r.Put(`/api/coderun/:imagename`, RunCodeStep)
	r.Get(`/api/coderun/:runid`, GetRunResult)
	r.Get(`/api/prepare/:imagename`, PrePareImage)

	//image api
	r.Get("/api/images/:id/name", getImageName)
	r.Get("/api/images", listImages)
	r.Post("/api/image/:name/search", searchImage)
	r.Get("/api/images/:id/list", listMyImages)
	r.Get("/api/images/:id/log", imageLogs)
	r.Get("/api/images/:name/verify", imageVerify)
	//r.Delete("/dockerapi/images/{id}/delete", deleteImage)
	//add userid
	r.Post("/api/image/create/:userid", createImage)
	r.Post("/api/image/commit/:userid", commitImage)
	r.Post("/api/image/push/:userid", pushImage)
	r.Post("/api/image/edit/:userid", editImage)
	//
	r.Post("/api/image/star", starImage)
	r.Post("/api/image/fork", forkImage)
	r.Get("/api/star/:uid/:id", queryStarid)
	r.Get("/api/fork/:uid/:id", queryFork)
	r.Post("/api/sso/islogin", isLogin)
	r.Get("/api/message/query/:id", queryNotice)
	r.Get("/api/message/read/:id", readMessageAPI)
	r.Post("/api/message/add", addMessage)
	//	r.Post("/api/sso/logout", logout)

	//code issue
	r.Get(`/api/image/:imageid/issues`, GetImageIssues)
	r.Post(`/api/image/:userid/:imageid/issue`, AddImageIssue)
	r.Put(`/api/image/:userid/:imageid/issue/:issueid`, UpdateImageIssue)
	r.Delete(`/api/image/:userid/:imageid/issue/:issueid`, DeleteImageIssue)
	//code issue comment
	r.Get(`/api/image/issue/:issueid/comments`, GetImageIssueComments)
	r.Post(`/api/image/issue/:userid/:issueid/comment`, AddImageIssueComment)
	r.Put(`/api/image/issue/:userid/:issueid/comment/:commentid`, UpdateImageIssueComment)
	r.Delete(`/api/image/issue/:userid/:issueid/comment/:commentid`, DeleteImageIssueComment)

	r.Post(`/api/user/info/update/:uid`, updateUserInfo)
	r.Get(`/api/user/info/get/:uid`, getUserInfo)

	r.Get(`/api/user/upload/:userid`, checkPic)
	r.Post(`/api/user/upload/:userid`, uploadPic)

	// Inject database
	m.MapTo(code_db, (*codeDB_inter)(nil))
	m.MapTo(code_step_db, (*codeStepDB_inter)(nil))
	// Add the router action
	m.Action(r.Handle)

}

// slash).
var rxExt = regexp.MustCompile(`(\.(?:xml|text|json))\/?$`)

//add cors
func AddCorsHandler(c martini.Context, w http.ResponseWriter, r *http.Request) {
	var Origin = r.Header.Get("Origin")
	if strings.Contains(Origin, "learn4me.com") {
		w.Header().Set("Access-Control-Allow-Origin", Origin)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

	}

}

// MapEncoder intercepts the request's URL, detects the requested format,
// and injects the correct encoder dependency for this request. It rewrites
// the URL to remove the format extension, so that routes can be defined
// without it.
func MapEncoder(c martini.Context, w http.ResponseWriter, r *http.Request) {
	// Get the format extension
	matches := rxExt.FindStringSubmatch(r.URL.Path)
	ft := ".json"
	if len(matches) > 1 {
		// Rewrite the URL without the format extension
		l := len(r.URL.Path) - len(matches[1])
		if strings.HasSuffix(r.URL.Path, "/") {
			l--
		}
		r.URL.Path = r.URL.Path[:l]
		ft = matches[1]
	}
	// Inject the requested encoder
	switch ft {
	case ".xml":
		c.MapTo(xmlEncoder{}, (*Encoder)(nil))
		w.Header().Set("Content-Type", "application/xml")
	case ".text":
		c.MapTo(textEncoder{}, (*Encoder)(nil))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	case ".html":
		c.MapTo(textEncoder{}, (*Encoder)(nil))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	default:
		c.MapTo(jsonEncoder{}, (*Encoder)(nil))
		w.Header().Set("Content-Type", "application/json")
	}
}

func main() {
	//	go func() {
	// Listen on http: to raise an error and indicate that https: is required.
	//
	// This could also be achieved by passing the same `m` martini instance as
	// used by the https server, and by using a middleware that checks for https
	// and returns an error if it is not a secure connection. This would have the benefit
	// of handling only the defined routes. However, it is common practice to define
	// APIs on separate web servers from the web (html) pages, for maintenance and
	// scalability purposes, so it's not like it will block otherwise valid routes.
	//
	// It is also common practice to use a different subdomain so that cookies are
	// not transfered with every API request.
	// So with that in mind, it seems reasonable to refuse each and every request
	// on the non-https server, regardless of the route. This could of course be done
	// on a reverse-proxy in front of this web server.
	//
	//		if err := http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//			http.Error(w, "https scheme is required", http.StatusBadRequest)
	//		})); err != nil {
	//			log.Fatal(err)
	//		}
	//	}()

	// Listen on https: with the preconfigured martini instance. The certificate files
	// can be created using this command in this repository's root directory:
	//
	// go run /path/to/goroot/src/pkg/crypto/tls/generate_cert.go --host="localhost"
	//
	flag.Parse()
	defer dbmap.Db.Close()
	//	defer redis_pool.Close()
	//	defer redis_pool.Put(redis_resource)

	//	timer := time.NewTicker(24 * time.Hour)
	//	timer := time.NewTicker(10 * time.Second)
	//	logger.Println("==============================================================")
	//	for {
	//		select {
	//		case <-timer.C:
	//			go HotTimerList()
	//		}
	//	}

	go func() {
		log.Println("listening on 9000")
		if err := http.ListenAndServe(*addr, m); err != nil {
			log.Fatal(err)
		}
	}()
	for {
		select {
		case <-timer.C:
			go HotTimerList()
		}
	}
}
