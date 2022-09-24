package main

import (
	"log"

	// application
	"github.com/sinakeshmiri/arun/internal/app/api"
	"github.com/sinakeshmiri/arun/internal/app/core/function"

	// adapters
	HTTP "github.com/sinakeshmiri/arun/internal/adapters/framework/left/http"
	"github.com/sinakeshmiri/arun/internal/adapters/framework/right/compiler"
	"github.com/sinakeshmiri/arun/internal/adapters/framework/right/db"
	"github.com/sinakeshmiri/arun/internal/adapters/framework/right/fs"
	"github.com/sinakeshmiri/arun/internal/adapters/framework/right/orchestrator"
)

func main() {
	var err error

	dbaseDriver := "mysql"
	dsourceName := "root:Admin123@tcp(72.30.52.64:3307)/arun"
	k8Configfile := "config"
	tusdServer := "http://172.30.52.64:1080/files/"
	k8surl:="192.168.188.129"

	dbAdapter, err := db.NewAdapter(dbaseDriver, dsourceName)
	if err != nil {
		log.Fatalf("failed to initiate dbase connection: %v", err)
	}
	defer dbAdapter.CloseDbConnection()
	cmpAdapter, err := compiler.NewAdapter()
	if err != nil {
		log.Fatalf("failed to initiate compiler: %v", err)
	}
	fsAdapter, err := fs.NewAdapter(tusdServer)
	if err != nil {
		log.Fatalf("failed to initiate file server: %v", err)
	}
	orcAdapter, err := orchestrator.NewAdapter(k8Configfile)
	if err != nil {
		log.Fatalf("failed to initiate orchestrator: %v", err)
	}
	// core
	core := function.New()

	// NOTE: The application's right side port for driven
	// adapters, in this case, a db adapter.
	// Therefore the type for the dbAdapter parameter
	// that is to be injected into the NewApplication will
	// be of type DbPort
	applicationAPI := api.NewApplication(dbAdapter, orcAdapter, fsAdapter, cmpAdapter, core)

	// NOTE: We use dependency injection to give the grpc
	// adapter access to the application, therefore
	// the location of the port is inverted. That is
	// the grpc adapter accesses the hexagon's driving port at the
	// application boundary via dependency injection,
	// therefore the type for the applicaitonAPI parameter
	// that is to be injected into the gRPC adapter will
	// be of type APIPort which is our hexagons left side
	// port for driving adapters
	HTTPAdapter := HTTP.NewAdapter(applicationAPI,k8surl)
	HTTPAdapter.Run()
}

