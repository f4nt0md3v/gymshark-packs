package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/f4nt0md3v/gymshark-packs/pkg/httpx"
	"github.com/f4nt0md3v/gymshark-packs/server/api/services"
	"github.com/f4nt0md3v/gymshark-packs/server/models/api"
)

type PackController struct {
	packer services.Packer
}

func NewPackController(p services.Packer) *PackController {
	return &PackController{packer: p}
}

func (c *PackController) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", c.Pack)

	return r
}

// Pack returns request payload in response.
// @Tags Pack Package
// @Summary Pack Package for service
// @Description Implements Pack Package function for service
// @Accept  json
// @Produce  json
// @Param Items body api.PackRequest true "PackRequest"
// @Success 200 {object} api.PackResponse
// @Failure 400 {object} httpx.Response
// @Failure 500 {object} httpx.Response
// @Router /packer/pack [POST]
func (c *PackController) Pack(w http.ResponseWriter, r *http.Request) {
	var (
		req       api.PackRequest
		packSizes = []int{250, 500, 1000, 2000, 5000}
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = render.Render(w, r, httpx.BadRequest(err))
		return
	}

	render.JSON(w, r, NewPackResponse(c.packer.CalculateNumberOfPacks(req.Items, packSizes)))
}

// NewPackResponse returns a new PackResponse.
func NewPackResponse(packs map[int]int) *api.PackResponse {
	return &api.PackResponse{Packs: packs}
}
