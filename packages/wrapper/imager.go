package wrapper

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/sinakeshmiri/arun/packages/config"
)

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func imageMacker(dockerClient *client.Client, fdt config.Lamda, buildEnv string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	tar, err := archive.TarWithOptions(fdt.Funcname, &archive.TarOptions{})
	if err != nil {
		return err
	}

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{fdt.Username + "/" + fdt.Funcname},
		Remove:     true,
	}
	res, err := dockerClient.ImageBuild(ctx, tar, opts)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = print(res.Body)
	if err != nil {
		return err
	}

	return nil
}
func print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func dockerfileMaker(buildEnv string) error {
	dockerFile := []byte(fmt.Sprintf("FROM scratch\nADD %s/bin.elf /\nCMD ['/bin.elf']", buildEnv))
	return os.WriteFile(buildEnv+"dockerFile", dockerFile, 0644)

}

func imager(buildEnv string) error {
	fmt.Println("EEE")
	var f config.Lamda
	f.Username = "sina"
	f.Funcname = "myfunc"
	err := dockerfileMaker(buildEnv)
	fmt.Println("EEE")
	if err != nil {
		return err
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = imageMacker(cli, f, buildEnv)
	if err != nil {
		return err
	}
	return nil
}
