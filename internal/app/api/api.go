package api

import (
	"time"

	"github.com/sinakeshmiri/arun/internal/ports"
)

// Application implements the APIPort interface
type Application struct {
	fs  ports.FsPort
	orc ports.OrchestratorPort
	db  ports.DbPort
	f   Function
	cmp ports.CompilerPort
}

// NewApplication creates a new Application
func NewApplication(
	db ports.DbPort, orc ports.OrchestratorPort, fs ports.FsPort, comp ports.CompilerPort, f Function) *Application {
	return &Application{db: db, orc: orc, fs: fs, cmp: comp, f: f}
}

// Run starts the pod
func (apia Application) RunFunction(name string) (string, error) {
	binLocation,err := apia.db.GetFunction(name)
	if err != nil {
		return "", err
	}
	podName, err := apia.orc.Run(binLocation)
	if err != nil {
		return "", err
	}

	return podName, nil
}

func (apia Application) AddFunction(name string, src string) error {
	err := apia.db.CheckName(name)
	if err != nil {
		return err
	}
	bin, err := apia.cmp.Compile(src)
	if err != nil {
		return  err
	}
	binUrl, err := apia.fs.SaveBinary(bin)
	if err != nil {
		return  err
	}
	err = apia.db.SaveFunction(name, binUrl, time.Duration(0))
	if err != nil {
		return err
	}
	return nil
}
