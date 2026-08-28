package httphandler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/6-clima-google-cloud-run/internal/cep"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/6-clima-google-cloud-run/internal/httphandler"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type fakeCityFinder struct {
	city string
	err  error
}

func (f fakeCityFinder) FindCity(ctx context.Context, zipcode string) (string, error) {
	return f.city, f.err
}

type fakeTemperatureFetcher struct {
	celsius float64
	err     error
}

func (f fakeTemperatureFetcher) CurrentCelsius(ctx context.Context, city string) (float64, error) {
	return f.celsius, f.err
}

func newRouter(h httphandler.Handler) http.Handler {
	router := chi.NewRouter()
	router.Get("/{cep}", h.ServeHTTP)
	return router
}

func TestHandler_InvalidZipcode(t *testing.T) {
	h := httphandler.Handler{Cities: fakeCityFinder{}, Temps: fakeTemperatureFetcher{}}
	req := httptest.NewRequest(http.MethodGet, "/123", nil)
	rec := httptest.NewRecorder()

	newRouter(h).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.JSONEq(t, `{"message":"invalid zipcode"}`, rec.Body.String())
}

func TestHandler_ZipcodeNotFound(t *testing.T) {
	h := httphandler.Handler{
		Cities: fakeCityFinder{err: cep.ErrNotFound},
		Temps:  fakeTemperatureFetcher{},
	}
	req := httptest.NewRequest(http.MethodGet, "/00000000", nil)
	rec := httptest.NewRecorder()

	newRouter(h).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `{"message":"can not find zipcode"}`, rec.Body.String())
}

func TestHandler_Success(t *testing.T) {
	h := httphandler.Handler{
		Cities: fakeCityFinder{city: "São Paulo"},
		Temps:  fakeTemperatureFetcher{celsius: 28.5},
	}
	req := httptest.NewRequest(http.MethodGet, "/01001000", nil)
	rec := httptest.NewRecorder()

	newRouter(h).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}`, rec.Body.String())
}
