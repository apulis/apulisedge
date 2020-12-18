// Copyright 2020 Apulis Technology Inc. All rights reserved.

package cluster

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var (
	ErrCannotGetImageLoadResult = errors.New("Can`t get image load result")
	ErrImageTagFailed           = errors.New("Image tag failed")
)

type ImageLoadResult struct {
	Stream string
}

// init docker cli
func (c *Cluster) InitDockerCli(harborAddress string, harborProject string, harborUser string, harborPasswd string) {
	c.HarborAddress = harborAddress
	c.HarborProject = harborProject
	c.HarborUser = harborUser
	c.HarborPasswd = harborPasswd

	logger.Infof("HarborAddress = %s, HarborProject = %s, HarborUser = %s, HarborPasswd = %s",
		c.HarborAddress, c.HarborProject, c.HarborUser, c.HarborPasswd)
}

func (c *Cluster) GetHarborAddress() string {
	return c.HarborAddress
}

func (c *Cluster) GetHarborProject() string {
	return c.HarborProject
}

func (c *Cluster) GetImageNameAndVersion(imageTag string) (string, string, error) {
	var img string

	path := strings.Split(imageTag, "/")
	if len(path) > 0 {
		img = path[len(path)-1]
	} else {
		img = imageTag
	}

	tagVer := strings.Split(img, ":")
	if len(tagVer) != 2 {
		return "", "", ErrImageTagFailed
	}

	return tagVer[0], tagVer[1], nil
}

func (c *Cluster) NewDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func (c *Cluster) CloseDockerClient(cli *client.Client) error {
	return cli.Close()
}

func (c *Cluster) DockerImageLoad(ctx context.Context, cli *client.Client, srcFile string) (string, error) {
	var loadedTag string

	imgSrc, err := os.Open(srcFile)
	if err != nil {
		return loadedTag, err
	}
	defer imgSrc.Close()

	res, err := cli.ImageLoad(ctx, imgSrc, false)
	if err != nil {
		return loadedTag, err
	}
	defer res.Body.Close()

	// read load result and parse
	// ImageLoad will return like this: {"stream":"Loaded image: nginx_arm64:1.18.0\n"}
	reader := bufio.NewReader(res.Body)
	result := ImageLoadResult{}
	for {
		a, _, c := reader.ReadLine()
		if c == io.EOF {
			break
		}
		err = json.Unmarshal(a, &result)
		if err != nil || len(result.Stream) == 0 {
			continue
		}

		splitStream := strings.Split(result.Stream, ":")
		if len(splitStream) != 3 || !strings.Contains(splitStream[0], "Loaded image") {
			continue
		}

		loadedTag = splitStream[1] + ":" + splitStream[2]
		loadedTag = strings.Trim(loadedTag, " \n")
	}

	if len(loadedTag) == 0 {
		return loadedTag, ErrCannotGetImageLoadResult
	}

	return loadedTag, nil
}

func (c *Cluster) DockerImageTag(ctx context.Context, cli *client.Client, src string, target string) error {
	logger.Infof("DockerImageTag: src = %s, target = %s", src, target)
	return cli.ImageTag(ctx, src, target)
}

func (c *Cluster) DockerImagePush(ctx context.Context, cli *client.Client, image string) error {
	authConfig := types.AuthConfig{
		Username: c.HarborUser,
		Password: c.HarborPasswd,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	res, err := cli.ImagePush(ctx, image, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		return err
	}
	defer res.Close()

	io.Copy(ioutil.Discard, res)
	return nil
}
