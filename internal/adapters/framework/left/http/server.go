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
	api ports.APIPort
}

// NewAdapter creates a new Adapter
func NewAdapter(api ports.APIPort) *Adapter {
	return &Adapter{api: api}
}

func (httpa Adapter) Run() {

	// Starts a new Gin instance with no middle-ware
	r := gin.New()

	// Define your handlers
	r.POST("/add", httpa.AddFunction)
	r.Any("/:path", httpa.RunFunction)

	if err := r.Run(); err != nil {
		log.Printf("Error: %v", err)
	}
}
