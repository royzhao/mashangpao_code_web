package main

import (
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	//	"github.com/codegangsta/martini-contrib/auth"
)

// 只有一个martini实例
var m *martini.Martini

func init() {
	// dbmap = nil
	// //init database
	dbmap := initDb()
	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
	log.Println("init db is successful!")
	// initrundb()
	//init code module
	//init_code()
	init_code(dbmap)
	init_codestep(dbmap)
	m = martini.New()
	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	//static
	m.Use(martini.Static("public"))
	//m.Use(auth.Basic(AuthToken, ""))
	m.Use(MapEncoder)

	r := martini.NewRouter()
	//得到所有的代码
	r.Get(`/api/code`, GetCodes)
	//得到某一个用户的所有代码
	r.Get(`/api/code/:userid`, GetCodesByUser)
	//一个用户增加一个代码
	r.Post(`/api/code/:userid`, AddCode)
	//查询一个指定的代码
	r.Get(`/api/code/:userid/:codeid`, GetCode)
	//修改指定的代码
	r.Put(`/api/code/:userid/:codeid`, UpdateCode)
	//删除指定的代码
	r.Delete(`/api/code/:userid/:codeid`, DeleteCode)

	//得到代码的具体步骤
	//得到全部的代码步骤元数据
	r.Get(`/api/code/:userid/:codeid/step`, GetCodeSteps)
	//得到某一个具体步骤
	r.Get(`/api/code/:userid/:codeid/step/:stepid`, GetCodeStep)
	//增加一个代码的具体步骤
	r.Post(`/api/code/:userid/:codeid/step`, AddCodeStep)
	//修改置顶的代码步骤的元数据
	r.Put(`/api/code/:userid/:codeid/step/:stepid`, UpdateCodeStep)
	//删除
	r.Delete(`/api/code/:userid/:codeid/step/:stepid`, DeleteCodeStep)

	//具体内容操作
	r.Get(`/api/code/:userid/:codeid/step/:stepid/detail`, GetCodeStepDetail)
	//修改具体内容
	r.Put(`/api/code/:userid/:codeid/step/:stepid/detail`, UpdateCodeStepDetail)

	//add command
	r.Put(`/api/code/:userid/:codeid/step/:stepid/cmd`, UpdateCodeStepCmd)
	r.Get(`/api/code/:userid/:codeid/step/:stepid/cmd`, GetCodeStepCmd)

	//code run
	r.Put(`/api/code/run/:imageid`, RunCodeStep)
	r.Get(`/api/code/run/:runid`, GetRunResult)

	// Inject database
	m.MapTo(code_db, (*codeDB_inter)(nil))
	m.MapTo(code_step_db, (*codeStepDB_inter)(nil))
	// Add the router action
	m.Action(r.Handle)

}

// slash).
var rxExt = regexp.MustCompile(`(\.(?:xml|text|json))\/?$`)

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
	log.Println("listening on 8001")
	if err := http.ListenAndServe(":8001", m); err != nil {
		log.Fatal(err)
	}

}
