package httphandler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/7-clima-cep-observabilidade/internal/cep"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/7-clima-cep-observabilidade/internal/serviceaclient"
)

type ServiceBForwarder interface {
	ForwardCEP(ctx context.Context, cep string) (*serviceaclient.Response, error)
}

type CEPForwardHandler struct {
	ServiceB ServiceBForwarder
}

func (h CEPForwardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req cepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !cep.IsValid(req.CEP) {
		respondError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	resp, err := h.ServiceB.ForwardCEP(r.Context(), req.CEP)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(resp.Body)
}
