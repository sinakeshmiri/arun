package main

import (
	"fmt"
	"log"
	"time"

	dbms "github.com/sinakeshmiri/arun/internal/adapters/framework/right/db"
)

func main() {
	db, err := dbms.NewAdapter("mysql", "root:Admin123@tcp(192.168.110.253:3307)/arun")
	if err != nil {
		log.Fatal(err)
	}
	err=db.SaveFunction("GOz","http://37.32.24.125:1080/files/b7be83e7ad6c1fed37b29ddd3197b447",time.Duration(0))
	if err != nil {
		log.Fatal(err)
	}
	l,t,err:=db.GetFunction("GOz")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(l,t)
	err=db.CheckName("GOH")
	if err != nil {
		log.Fatal(err)
	}
}

/*
func main() {

	e, err := orc.NewAdapter("./config")
	if err != nil {
		log.Fatal(err)
	}
	x, err := fs.NewAdapter("http://37.32.24.125:1080/files/")
	if err != nil {
		log.Fatal(err)
	}
	t, err := x.SaveBinary("app")
	if err != nil {
		log.Fatal(err)
	}
	xx, err := e.Run(t)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(xx)
}
*/
/*
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
	dsourceName := "root:Admin123@tcp(db:3306)/hex_test"
	k8Configfile := ""
	tusdServer := ""

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
	HTTPAdapter := HTTP.NewAdapter(applicationAPI)
	HTTPAdapter.Run()
}
*/
