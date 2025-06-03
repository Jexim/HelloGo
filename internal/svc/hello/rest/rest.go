package rest

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/Jexim/HelloGo/internal/svc/hello/model"
)

// @title Hello Service API
// @version 1.0
// @description This is a hello service API documentation
// @host localhost:8080
// @BasePath /api/v1
type REST struct {
	helloUC model.Usecase
}

func New(mux *chi.Mux, prefix string, helloUC model.Usecase) model.REST {
	rest := &REST{helloUC: helloUC}

	mux.Route(prefix, func(r chi.Router) {
		r.Get("/", rest.GetHello)
	})

	return rest
}

// @Summary Get Hello
// @Description Get Hello
// @Accept json
// @Produce json
// @Success 200 {object} model.Hello
// @Router /hello [get]
func (r *REST) GetHello(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
