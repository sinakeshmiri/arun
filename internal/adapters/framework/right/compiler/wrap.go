package compiler

import (
	"archive/zip"
	b64 "encoding/base64"
	"fmt"

	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func builder(filepath string) error {
	os.Setenv("CGO_ENABLED","0" )
	os.Chdir(filepath)
	cmd := exec.Command("go", "mod", "tidy")
	err := cmd.Run()
	if err != nil {

		return err
	}
	cmd = exec.Command("go", "build", "-o", "app", "main.go")
	err = cmd.Run()
	if err != nil {

		return err
	}
	return nil
}

func Make(src string) (string, error) {
	rnd := uuid.New().String()
	buildEnv := "../arun-builder-" + rnd
	err := os.Mkdir(buildEnv, os.ModePerm)
	if err != nil {

		return "", err
	}
	////
	_,err=copy("wrap.zip",buildEnv+"/wrap.zip")
	if err != nil {
		return "", err
	}
	////
	err = unzip(buildEnv+"/wrap.zip", buildEnv)
	if err != nil {
		fmt.Println("wrapper zip not found")
		return "", err
	}
	d1, err := b64.StdEncoding.DecodeString(src)
	fmt.Println(string(d1),src)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(buildEnv+"/packages/function/fucntion.go", d1, 0644)
	if err != nil {

		return "", err
	}
	err = builder(buildEnv)

	if err != nil {
		return "", err
	}
	return buildEnv + "/app", nil
}

func unzip(wrp string, dest string) error {
	var filenames []string
	reader, err := zip.OpenReader(wrp)
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
func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
			return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
			return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
			return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
			return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}