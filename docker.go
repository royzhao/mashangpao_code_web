package main

import (
	"encoding/json"
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
	endpoint        = "http://imagehub.learn4me.com"
	browserEndpoint string
	dockerhub       = "registry.learn4me.com:5000"
	dockerclient    *docker.Client
)

type DockerContainerID struct {
	ID string
}

func (c CRImage) dockerCommit() error {
	//func main() {
	//req, err := http.NewRequest("GET", c.getURL(path), params)
	logger.Warnln(c.UserId)
	resp, err := http.Get(conf.BrowserEndpoint + "/containers/" + strconv.FormatInt(c.UserId, 10))
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
	logger.Println(name)
	if err := dockerclient.TagImage(name, docker.TagImageOptions{Repo: conf.Dockerhub + "/" + c.ImageName, Tag: strconv.Itoa(c.Tag), Force: true}); err != nil {
		logger.Warnf("error tagging container: %s", err)
		return err
	}
	logger.Println(strconv.Itoa(c.Tag))
	opts := docker.PushImageOptions{Name: conf.Dockerhub + "/" + c.ImageName, Tag: strconv.Itoa(c.Tag), Registry: conf.Dockerhub + "/"}
	var auth docker.AuthConfiguration
	if err := dockerclient.PushImage(opts, auth); err != nil {
		logger.Warnf("error pushing container: %s", err)
		return err
	}
	return nil
}

func (c CRImage) dockerFork(oldName string) error {
	oldName = conf.Dockerhub + "/" + oldName
	bash := []string{"bash"}
	config := &docker.Config{AttachStdin: true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		Memory:       0,
		MemorySwap:   0,
		CPUShares:    512,
		CPUSet:       "0,1",
		StdinOnce:    false,
		Cmd:          bash,
		Image:        oldName}
	var hc *docker.HostConfig
	opt := docker.CreateContainerOptions{Config: config, HostConfig: hc}
	container, err := dockerclient.CreateContainer(opt)
	if err != nil {
		logger.Warnf("error creating new container: %s", err)
		return err
	}
	commitOpts := docker.CommitContainerOptions{Container: container.ID, Repository: c.ImageName, Tag: strconv.Itoa(c.Tag)}
	if _, err = dockerclient.CommitContainer(commitOpts); err != nil {
		logger.Warnf("error committing container: %s", err)
		return err
	}
	//	err = dockerclient.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true})
	//	if _, err = dockerclient.CommitContainer(commitOpts); err != nil {
	//		logger.Warnf("error removing container: %s", err)
	//		return err
	//	}
	if err = c.dockerPush(); err != nil {
		logger.Warnf("error pushing container: %s", err)
		return err
	}
	//	if err = dockerclient.StopContainer(di.ID, 5); err != nil {
	//		logger.Warnf("error stopping container: %s", err)
	//		return err
	//	}
	return nil
}

//func dockerDelete(id string) {
//	if err := client.CommitContainer(commitOpts); err != nil {
//		fmt.Println(err)
//		return err
//	}
//}
