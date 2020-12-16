package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"io/ioutil"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	//authConfig := types.AuthConfig{
	//	Username: "username",
	//	Password: "password",
	//}
	//encodedJSON, err := json.Marshal(authConfig)
	//if err != nil {
	//	panic(err)
	//}
	//authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	out, err := cli.ImagePull(ctx, "harbor.sigsus.cn:8443/apulisedge/nginx_arm64:1.18.0", types.ImagePullOptions{ /*RegistryAuth: authStr*/ })
	if err != nil {
		panic(err)
	}

	defer out.Close()
	io.Copy(ioutil.Discard, out)
	fmt.Printf("pull image succ\n")
}
