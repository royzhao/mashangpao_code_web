package main

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/martini"
	"net/http"
	"strconv"
)

var logger = logrus.New()

//list all the images
func listImages(w http.ResponseWriter, r *http.Request, parms martini.Params) {
	images := QueryImage()
	if err := json.NewEncoder(w).Encode(images); err != nil {
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
	name := image.ImageName + ":" + strconv.Itoa(image.Tag)
	fullName := imageFullName{fullname: name}
	if err := json.NewEncoder(w).Encode(fullName); err != nil {
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
	name := parms["name"]
	result := QuerybyName(name)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
