package rest

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	httpmw "github.com/Jexim/HelloGo/internal/adapter/http/middleware"
	httprespond "github.com/Jexim/HelloGo/internal/adapter/http/respond"
	"github.com/Jexim/HelloGo/internal/modules/hello/model"
)

// @title Hello Service API
// @version 1.0
// @description This is a hello service API documentation
// @host localhost:8080
// @BasePath /api
type REST struct {
	helloUC model.Usecase
}

func New(mux *chi.Mux, prefix string, helloUC model.Usecase) model.REST {
	rest := &REST{helloUC: helloUC}

	mux.Route(prefix, func(r chi.Router) {
		r.Get("/", rest.ListHellos)
	})

	return rest
}

// @Summary Get Hello
// @Description Get Hello
// @Accept json
// @Produce json
// @Success 200 {object} model.Hello
// @Router /hello [get]
// ListHellos returns a paginated list of hellos
func (r *REST) ListHellos(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// parse pagination
	q := req.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	items, err := r.helloUC.GetAll(req.Context(), limit, offset)
	if err != nil {
		status, code, msg := httpmw.MapError(err)
		httprespond.JSON(w, status, map[string]any{
			"error": map[string]any{
				"code":     code,
				"message":  msg,
				"trace_id": httpmw.GetTraceID(req),
			},
		})
		return
	}
	httprespond.JSON(w, http.StatusOK, items)
}
