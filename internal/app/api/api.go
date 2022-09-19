package api

import (
	"time"

	"github.com/sinakeshmiri/arun/internal/ports"
)

// Application implements the APIPort interface
type Application struct {
	fs ports.FsPort
	orc  ports.OrchestratorPort
	db   ports.DbPort
	f    Function
	cmp  ports.CompilerPort
}

// NewApplication creates a new Application
func NewApplication(
	db ports.DbPort, orc ports.OrchestratorPort, fs ports.FsPort, f Function) *Application {
	return &Application{db: db, orc: orc, fs: fs, f: f}
}

// Complie compiles the source code of the given fucntion
func (apia Application) Complie(name string, src string) (string, error) {
	err := apia.db.CheckName(name)
	if err != nil {
		return "", err
	}
	bin, err := apia.cmp.Compile(src)
	if err != nil {
		return "", err
	}

	binUrl, err := apia.fs.SaveBinary(bin)
	if err != nil {
		return "", err
	}

	return binUrl, nil
}

// Run starts the pod
func (apia Application) Run(binLocation string) (string, error) {
	podName, err := apia.orc.Run(binLocation)
	if err != nil {
		return "", err
	}

	return podName, nil
}

func (apia Application) Save(name string, binary string, cost time.Duration) error {
	 err := apia.db.SaveFunction(name,binary,cost)
	if err != nil {
		return  err
	}

	return nil
}

/*
// GetAddition gets the result of adding parameters a and b
func (apia Application) GetAddition(a, b int32) (int32, error) {
	answer, err := apia.arith.Addition(a, b)
	if err != nil {
		return 0, err
	}

	err = apia.db.AddToHistory(answer, "addition")
	if err != nil {
		return 0, err
	}

	return answer, nil
}*/