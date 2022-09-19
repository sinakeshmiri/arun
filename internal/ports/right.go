package ports

import (
	"time"
)

// DbPort is the port for a db adapter
type DbPort interface {
	CloseDbConnection()
	SaveFunction(name string, binary string, cost time.Duration) error
	CheckName(name string)error
}

// K8Port is the port for a kubernetese adapter
type OrchestratorPort interface {
	Run(binary string) (string,error)
}

// TusdPort is the port for a tusd adapter
type FsPort interface {
	SaveBinary(binary string) (string,error)
}

type CompilerPort interface{
Compile(src string) (string,error)
}