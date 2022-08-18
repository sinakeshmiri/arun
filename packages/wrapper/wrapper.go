package wrapper

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func builder(filepath string) error {
	os.Chdir(filepath)
	cmd := exec.Command("go", "mod", "init", "arun.dev/builder/function")
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("go", "mod", "tidy")
	err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error)
		return err
	}
	cmd = exec.Command("go", "build", "-o", "bin.elf", "main.go")
	err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error)
		return err
	}
	return nil
}

func Make(zipPath string) error {
	buildEnv := uuid.New().String()

	buildEnv = "/tmp/arun-builder-" + buildEnv
	err := os.Mkdir(buildEnv, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error)
		return err
	}
	err = unzip(zipPath, buildEnv)
	if err != nil {
		fmt.Println(err.Error)
		return err
	}
	err = builder(buildEnv)
	//fmt.Println(err.Error)
	if err != nil {
		return err
	}
	err = imager(buildEnv)
	if err != nil {
		//fmt.Println(err.Error)
		return err
	}
	return nil
}

func unzip(src string, dest string) error {
	var filenames []string
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, f := range reader.File {
		fp := filepath.Join(dest, f.Name)
		strings.HasPrefix(fp, filepath.Clean(dest)+string(os.PathSeparator))
		filenames = append(filenames, fp)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(fp, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}
		err = os.MkdirAll(filepath.Dir(fp), os.ModePerm)
		if err != nil {
			return err
		}
		outFile, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()
	}
	return nil
}
