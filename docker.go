package main

import (
	"encoding/json"
	//	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/fsouza/go-dockerclient"
)

var (
	//	endpoint        = "unix:///var/run/docker.sock"
	// endpoint = "http://imagehub.peilong.me:81"

	// dockerclient, _ = docker.NewClient(endpoint)
	// browserEndpoint = "http://vpn.peilong.me:8080"
	// dockerhub       = "docker2.peilong.me:5000"
	endpoint        string
	browserEndpoint string
	dockerhub       string
	dockerclient    *docker.Client
)

type DockerContainerID struct {
	ID string
}

func (c CRImage) dockerCommit() error {
	//func main() {
	//req, err := http.NewRequest("GET", c.getURL(path), params)
	logger.Warnln(c.UserId)
	resp, err := http.Get(browserEndpoint + "/containers/" + strconv.FormatInt(c.UserId, 10))
	if err != nil {
		logger.Warnf("error getting container id: %s", err)
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Warnf("error reading http body: %s", err)
		return err
	}
	var di DockerContainerID
	logger.Warnf(string(body))
	if err = json.Unmarshal(body, &di); err != nil {
		logger.Warnf("error decoding container id: %s", err)
		return err
	}
	//	if err = client.PauseContainer(di.ID); err != nil {
	//		logger.Warnf("error stopping container: %s", err)
	//		return err
	//	}
	//	if err = client.StopContainer(di.ID, 5); err != nil {
	//		logger.Warnf("error stopping container: %s", err)
	//		return err
	//	}

	commitOpts := docker.CommitContainerOptions{Container: di.ID, Repository: c.ImageName, Tag: strconv.Itoa(c.Tag)}
	if _, err := dockerclient.CommitContainer(commitOpts); err != nil {
		logger.Warnf("error committing container: %s", err)
		return err
	}
	if err = dockerclient.StopContainer(di.ID, 5); err != nil {
		logger.Warnf("error stopping container: %s", err)
		return err
	}
	//	if err = client.RemoveContainer(docker.RemoveContainerOptions{ID: di.ID, Force: true}); err != nil {
	//		logger.Warnf("error removing container: %s", err)
	//		return err
	//	}
	return nil
	//	err = client.RemoveContainer(docker.RemoveContainerOptions{ID: "ffc4dfc4827c"})
}

func (c CRImage) dockerPush() error {
	logger.Println(c.ImageName)
	name := c.ImageName + ":" + strconv.Itoa(c.Tag)
	if err := dockerclient.TagImage(name, docker.TagImageOptions{Repo: dockerhub + "/" + c.ImageName, Tag: strconv.Itoa(c.Tag), Force: true}); err != nil {
		logger.Warnf("error tagging container: %s", err)
		return err
	}
	logger.Println(strconv.Itoa(c.Tag))
	opts := docker.PushImageOptions{Name: dockerhub + "/" + c.ImageName, Tag: strconv.Itoa(c.Tag), Registry: dockerhub + "/"}
	var auth docker.AuthConfiguration
	if err := dockerclient.PushImage(opts, auth); err != nil {
		logger.Warnf("error pushing container: %s", err)
		return err
	}
	return nil
}

//func dockerDelete(id string) {
//	if err := client.CommitContainer(commitOpts); err != nil {
//		fmt.Println(err)
//		return err
//	}
//}
