package http

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sinakeshmiri/arun/internal/ports"
)

/*type APIPort interface {
	AddFunction(name string ,source string) (error)
	RunFunction(name string) (time.Duration, error)
}
*/
// Adapter implements the GRPCPort interface
type Adapter struct {
	NodeUri string
	api ports.APIPort
}

// NewAdapter creates a new Adapter
func NewAdapter(api ports.APIPort ,nodeUri string) *Adapter {
	return &Adapter{api: api , NodeUri: nodeUri}
}

func (httpa Adapter) Run() {

	// Starts a new Gin instance with no middle-ware
	r := gin.New()

	// Define your handlers
	r.POST("/add", httpa.AddFunction)
	r.Any("/:id", httpa.RunFunction)
	r.Any("/", httpa.RunFunction)

	if err := r.Run(); err != nil {
		log.Printf("Error: %v", err)
	}
}
