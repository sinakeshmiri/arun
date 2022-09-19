package compiler

import (
	"archive/zip"

	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func builder(filepath string) error {

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
	buildEnv := "/tmp/arun-builder-" + rnd
	err := os.Mkdir(buildEnv, os.ModePerm)
	if err != nil {

		return "", err
	}
	err = unzip("./wrap.zip", buildEnv)
	if err != nil {

		return "", err
	}
	d1 := []byte(src)
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