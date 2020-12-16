// Copyright 2020 Apulis Technology Inc. All rights reserved.

package utils

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"github.com/docker/docker/client"
	"io"
	"os"
	"strings"
)

var (
	ErrCannotGetImageLoadResult = errors.New("Can`t get image load result")
)

type ImageLoadResult struct {
	Stream string
}

func NewDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func CloseDockerClient(cli *client.Client) error {
	return cli.Close()
}

func DockerImageLoad(ctx context.Context, cli *client.Client, srcFile string) (string, error) {
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
	}

	if len(loadedTag) == 0 {
		return loadedTag, ErrCannotGetImageLoadResult
	}

	return loadedTag, nil
}
