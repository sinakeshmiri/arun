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
func (apia Application) RunFunction(name string) (string,int32, error) {
	binLocation, _, err := apia.db.GetFunction(name)
	if err != nil {
		return "",0,err
	}
	podName,lprot, err := apia.orc.Run(binLocation)
	if err != nil {
		return "",0,err
	}

	return podName,lprot, nil
}

func (apia Application) AddFunction(name string, src string) error {
	err := apia.db.CheckName(name)
	if err != nil {
		return err
	}
	bin, err := apia.cmp.Compile(src)
	if err != nil {
		return err
	}
	binUrl, err := apia.fs.SaveBinary(bin)
	if err != nil {
		return err
	}
	err = apia.db.SaveFunction(name, binUrl, time.Duration(0))
	if err != nil {
		return err
	}
	return nil
}

func (apia Application) UpdateFunction(name string,duration time.Duration) ( error) {
	src, oldDuration, err := apia.db.GetFunction(name)
	if err != nil {
		return  err
	}
	newDuration:=oldDuration+duration
	apia.db.SaveFunction(name,src,newDuration)
	if err != nil {
		return  err
	}
	return nil
}
