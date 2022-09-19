package http

import (

	"github.com/sinakeshmiri/arun/internal/ports"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":80", r)
}
