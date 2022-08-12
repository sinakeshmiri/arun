package wrapper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func downloader(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

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
		return err
	}
	cmd = exec.Command("go", "build", "main.go")
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func Make(filepath string, url string) error {
	err := downloader(filepath+"/bin", url)
	if err != nil {
		return err
	}
	err = builder(filepath + "/bin")
	if err != nil {
		return err
	}
	return nil
}
