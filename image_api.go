package main

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"strconv"
	"time"
)

var logger = logrus.New()

func GetImageIssues(r *http.Request, enc Encoder, parms martini.Params) (int, string) {
	id, err := strconv.ParseInt(parms["imageid"], 10, 64)
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the image with id %s does not exist", parms["userid"]))))
	}
	// Get the query string arguments, if any
	qs := r.URL.Query()
	key := qs.Get("key")
	page := qs.Get("page")
	num := qs.Get("num")
	_num, err := strconv.Atoi(num)
	if err != nil {
		_num = 5
	}
	if _num == 0 {
		_num = 5
	}
	_page, err := strconv.Atoi(page)
	if err != nil {
		_page = 1
	}
	if _page <= 0 {
		_page = 1
	}
	// Otherwise, return all Codes
	return http.StatusOK, Must(enc.Encode(FindImageIssues(key, _page, _num, id)))
}

func AddImageIssue(r *http.Request, enc Encoder, parms martini.Params) (int, string) {
	imageid, err := strconv.ParseInt(parms["imageid"], 10, 64)
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the image with id %s does not exist", parms["imageid"]))))
	}
	al, err := getPostImageIssue(r)
	if err != nil {
		return http.StatusBadRequest, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the issue create failed no image %s", parms["imageid"]))))
	}
	al.Image_id = imageid
	id, err := AddOneImageIssue(al)
	if err != nil {
		return http.StatusBadRequest, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the issue create failed"))))
	}
	al.Id = int64(id)
	go func() {
		var obj CRImage
		image := obj.Querylog(imageid)
		if image.ImageId == 0 {
			// Invalid id, or does not exist
			return
		}
		_, err := NewMessage(image.UserId, al.Author, fmt.Sprintf("有人评论了您的讨论，<a href='/dashboard.html#/image/%d/issue/%d'>click to read</a>", imageid, id), 1)
		if err != nil {
			log.Println(err)
		}
	}()
	return http.StatusCreated, Must(enc.Encode(al))
}

func DeleteImageIssue(r *http.Request, enc Encoder, parms martini.Params) (int, string) {
	issue, err := strconv.ParseInt(parms["issueid"], 10, 64)
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the issue with id %s does not exist", parms["issueid"]))))
	}
	err = DeleteImageIssueById(issue)
	if err != nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the issue with id %s does not exist", parms["issueid"]))))
	}
	return http.StatusOK, fmt.Sprintf("delete issue=%d is ok", issue)
}

func UpdateImageIssue(r *http.Request, enc Encoder, parms martini.Params) (int, string) {
	imageid, err := strconv.ParseInt(parms["imageid"], 10, 64)
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the image with id %s does not exist", parms["imageid"]))))
	}
	issue, err := strconv.ParseInt(parms["issueid"], 10, 64)
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the issue with id %s does not exist", parms["issueid"]))))
	}
	al, err := getPostImageIssue(r)
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("Update failed"))))
	}
	al.Image_id = imageid
	al.Id = issue
	err = UpdateImageIssueById(al)
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("Update failed"))))
	}
	return http.StatusOK, Must(enc.Encode(al))
}

func GetImageIssueComments(r *http.Request, enc Encoder, parms martini.Params) (int, string) {
	id, err := strconv.ParseInt(parms["issueid"], 10, 64)
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the issue with id %s does not exist", parms["userid"]))))
	}
	// Get the query string arguments, if any
	qs := r.URL.Query()
	key := qs.Get("key")
	page := qs.Get("page")
	num := qs.Get("num")
	_num, err := strconv.Atoi(num)
	if err != nil {
		_num = 5
	}
	if _num == 0 {
		_num = 5
	}
	_page, err := strconv.Atoi(page)
	if err != nil {
		_page = 1
	}
	if _page <= 0 {
		_page = 1
	}
	// Otherwise, return all Codes
	return http.StatusOK, Must(enc.Encode(FindImageIssueComment(key, _page, _num, id)))
}

func AddImageIssueComment(r *http.Request, enc Encoder, parms martini.Params) (int, string) {
	issueid, err := strconv.ParseInt(parms["issueid"], 10, 64)
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("Add failed issueid is not accepted"))))
	}
	al, err := getPostImageIssueComment(r)
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("Add failed  post not accepted"))))
	}
	al.Issue_id = issueid
	id, err := AddOneImageIssueComment(al)
	if err != nil {
		return http.StatusBadRequest, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the issue create failed"))))
	}
	al.Id = id
	go func() {
		issue := GetImageIssueById(issueid)
		_, err := NewMessage(al.Reply_to, al.Author,
			fmt.Sprintf("someone attend <a href='/dashboard.html#/image/%d/issue/%d'>click to read</a>", issue.Image_id, issue.Id), 1)
		if err != nil {
			log.Println(err)
		}
	}()
	return http.StatusCreated, Must(enc.Encode(al))
}

func DeleteImageIssueComment(r *http.Request, enc Encoder, parms martini.Params) (int, string) {
	commentid, err := strconv.ParseInt(parms["commentid"], 10, 64)
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the commentid with id %s does not exist", parms["issue"]))))
	}
	err = DeleteDataImageIssueComment(commentid)
	if err != nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the commentid with id %s does not exist", parms["issue"]))))
	}
	return http.StatusOK, fmt.Sprintf("delete commentid=%d is ok", commentid)
}

func UpdateImageIssueComment(r *http.Request, enc Encoder, parms martini.Params) (int, string) {
	issueid, err := strconv.ParseInt(parms["issueid"], 10, 64)
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("Add failed"))))
	}
	commentid, err := strconv.ParseInt(parms["commentid"], 10, 64)
	if err != nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the commentid with id %s does not exist", parms["issue"]))))
	}
	al, err := getPostImageIssueComment(r)
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("Add failed"))))
	}
	al.Issue_id = issueid
	al.Id = commentid
	err = UpdateDataImageIssueComment(al)
	if err != nil {
		return http.StatusOK, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("Update failed"))))
	}
	return http.StatusOK, Must(enc.Encode(al))
}

//list all the images
func listImages(w http.ResponseWriter, r *http.Request, parms martini.Params) {
	//	val, err := redis_client.Get("hotimage")
	//	if err != nil {
	//		logger.Error(err)
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//	}
	//	var list HotImages
	//	if err = json.Unmarshal(val, &list); err != nil {
	//		logger.Error(err)
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//	}
	//	var images []CRImage
	//	if len(list.List) > 50 {
	//		images = list.List[0:50]
	//	} else {
	//		images = list.List[0:]
	//	}
	//	if err := json.NewEncoder(w).Encode(images); err != nil {
	//		logger.Error(err)
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//	}
	qs := r.URL.Query()
	key := qs.Get("key")
	page, err := strconv.Atoi(qs.Get("page"))
	if err != nil {
		page = 1
	}
	if page <= 0 {
		page = 1
	}
	num, err := strconv.Atoi(qs.Get("num"))
	if err != nil {
		num = 8
	}
	if num == 0 {
		num = 8
	}
	start := (page - 1) * num
	end := start + num
	var list ImageList
	var result HotImages
	if key == "" {
		result.Num = num
		result.Page = page
		conn := pool.Get()
		defer conn.Close()
		val, err := conn.Do("GET", "hotimage")
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// convert interface to []byte
		tmp, _ := val.([]byte)
		if err = json.Unmarshal(tmp, &list); err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if end > 50 {
			result.List = list.List[start:50]
		} else if end > len(list.List) {
			result.List = list.List[start:]
		} else {
			result.List = list.List[start:end]
		}
		if len(list.List) > 50 {
			result.Total = 50
		} else {
			result.Total = int64(len(list.List))
		}
	} else {
		result = QuerybyName(key, page, num)
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//list a user's images
func listMyImages(w http.ResponseWriter, r *http.Request, parms martini.Params) {
	uid, _ := strconv.ParseInt(parms["id"], 10, 64)
	var i CRImage
	logger.Println(uid)
	image := i.QuerybyUser(uid)
	logger.Println(image)
	if err := json.NewEncoder(w).Encode(image); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type imageFullName struct {
	fullname string
}

//get an image name from its id
func getImageName(w http.ResponseWriter, r *http.Request, parms martini.Params) {
	id, _ := strconv.ParseInt(parms["id"], 10, 64)
	var img CRImage
	image := img.Querylog(id)
	// name := image.ImageName + ":" + strconv.Itoa(image.Tag)
	// log.Println(name)
	// fullName := imageFullName{fullname: name}
	if err := json.NewEncoder(w).Encode(image); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//get an image's log
func imageLogs(w http.ResponseWriter, r *http.Request, parms martini.Params) {
	id, _ := strconv.ParseInt(parms["id"], 10, 64)
	var img CRImage
	image := img.Querylog(id)
	if err := json.NewEncoder(w).Encode(*image); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type unique struct {
	IsUnique bool
}

//verify if the image name exists
func imageVerify(w http.ResponseWriter, r *http.Request, parms martini.Params) {
	name := parms["name"]
	isUnique := QueryVerify(name)
	if err := json.NewEncoder(w).Encode(unique{IsUnique: isUnique}); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type newimage struct {
	UserId    int64
	ImageName string
	BaseImage string
	Tag       int
	Descrip   string
}

type baseImage struct {
	Bimage string
}

//create a new image from base image
func createImage(w http.ResponseWriter, r *http.Request) {
	//	vars := mux.Vars(r)
	//	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	var ni newimage
	if err := json.NewDecoder(r.Body).Decode(&ni); err != nil {
		logger.Warnf("error decoding image: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bi := baseImage{ni.BaseImage}
	cr := newImage(ni.UserId, ni.ImageName, ni.Tag, ni.Descrip)
	if err := cr.Add(); err != nil {
		logger.Warnf("error creating image: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(bi); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type myImageID struct {
	ID int64
}

//commit a new image
func commitImage(w http.ResponseWriter, r *http.Request) {
	//	var ni newimage
	var ci CRImage
	if err := json.NewDecoder(r.Body).Decode(&ci); err != nil {
		logger.Warnf("error decoding image: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := ci.dockerCommit(); err != nil {
		logger.Warnf("error committing image: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//response the image id
	//	mi := myImageID{ID: ci.ImageId}
	//	if err := json.NewEncoder(w).Encode(mi); err != nil {
	//		logger.Error(err)
	//	}
}

//edit an exist image
func editImage(w http.ResponseWriter, r *http.Request) {
	//	vars := mux.Vars(r)
	//	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	var ci CRImage
	if err := json.NewDecoder(r.Body).Decode(&ci); err != nil {
		logger.Warnf("error decoding image: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := ci.UpdateImage(); err != nil {
		logger.Warnf("error updating image: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

//push an a new image to the private registry
func pushImage(w http.ResponseWriter, r *http.Request) {
	var ci CRImage
	if err := json.NewDecoder(r.Body).Decode(&ci); err != nil {
		logger.Warnf("error decoding image: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := ci.dockerPush(); err != nil {
		logger.Warnf("error pushing image: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := ci.UpdateStatus(1); err != nil {
		logger.Warnf("error updating image status: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

//the parameter of fork image function
type starData struct {
	Uid   int64
	Image CRImage
}

//star or unstar a image
func starImage(w http.ResponseWriter, r *http.Request) {
	//	r.ParseForm()
	//	starStr := r.FormValue("sbool")
	//	star, _ := strconv.ParseBool(starStr)
	//	sid := r.FormValue("id")
	//	log.Println(sid)
	//	log.Println(star)
	//	var cr CRImage
	var data starData
	//	var cs CRStar
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logger.Warnf("error decoding image: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err := data.Image.UpdateStar(data.Uid)
	if err != nil {
		logger.Warnf("error staring image: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

type starID struct {
	ID int64
}

//query the star record
func queryStarid(w http.ResponseWriter, r *http.Request, parms martini.Params) {
	id, _ := strconv.ParseInt(parms["id"], 10, 64)
	uid, _ := strconv.ParseInt(parms["uid"], 10, 64)
	cs := CRStar{ImageId: id, UserId: uid}
	sid := cs.QueryStar()
	if err := json.NewEncoder(w).Encode(starID{ID: sid}); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//the parameter of fork image function
type forkData struct {
	Uid   int64
	Uname string
	Image CRImage
}

//fork an exist image
func forkImage(w http.ResponseWriter, r *http.Request) {
	//	uid, _ := strconv.ParseInt(parms["uid"], 10, 64)
	//	uname, _ := parms["uname"]
	var data forkData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logger.Warnf("error decoding image: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//can't fork one's own image
	if data.Uid == data.Image.UserId {
		http.Error(w, "Can not fork your own image", http.StatusInternalServerError)
		return
	}
	err := data.Image.UpdateFork(data.Uid, data.Uname)
	if err != nil {
		logger.Warnf("error forking image: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

type forked struct {
	Forked bool
}

func queryFork(w http.ResponseWriter, r *http.Request, parms martini.Params) {
	id, _ := strconv.ParseInt(parms["id"], 10, 64)
	uid, _ := strconv.ParseInt(parms["uid"], 10, 64)
	cf := CRFork{ImageId: id, UserId: uid}
	fork := cf.QueryFork()
	if err := json.NewEncoder(w).Encode(forked{Forked: fork}); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func searchImage(w http.ResponseWriter, r *http.Request, parms martini.Params) {
	//	name := parms["name"]
	//	result := QuerybyName(name)
	//	if err := json.NewEncoder(w).Encode(result); err != nil {
	//		logger.Error(err)
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//	}
}

// Parse the request body, load into an Code structure.
func getPostImageIssue(r *http.Request) (*Image_issue, error) {
	decoder := json.NewDecoder(r.Body)
	var t Image_issue
	err := decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	t.Create_date = time.Now().String()
	return &t, nil
}
func getPostImageIssueComment(r *http.Request) (*Image_issue_comment, error) {
	decoder := json.NewDecoder(r.Body)
	var t Image_issue_comment
	err := decoder.Decode(&t)
	if err != nil {
		return nil, err
	}
	t.Create_date = time.Now().String()
	return &t, nil
}
