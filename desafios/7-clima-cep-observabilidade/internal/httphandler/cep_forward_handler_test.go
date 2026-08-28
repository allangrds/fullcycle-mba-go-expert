package httphandler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/7-clima-cep-observabilidade/internal/serviceaclient"
	"github.com/stretchr/testify/assert"
)

type fakeServiceBForwarder struct {
	resp *serviceaclient.Response
	err  error
}

func (f fakeServiceBForwarder) ForwardCEP(ctx context.Context, cep string) (*serviceaclient.Response, error) {
	return f.resp, f.err
}

func TestCEPForwardHandler_ServeHTTP(t *testing.T) {
	t.Run("invalid cep does not call service b", func(t *testing.T) {
		handler := CEPForwardHandler{ServiceB: fakeServiceBForwarder{err: errors.New("should not be called")}}
		req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewReader([]byte(`{"cep":"123"}`)))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.JSONEq(t, `{"message":"invalid zipcode"}`, rec.Body.String())
	})

	t.Run("valid cep forwards to service b and relays response", func(t *testing.T) {
		handler := CEPForwardHandler{ServiceB: fakeServiceBForwarder{
			resp: &serviceaclient.Response{StatusCode: http.StatusOK, Body: []byte(`{"city":"São Paulo","temp_C":28.5,"temp_F":83.3,"temp_K":301.5}`)},
		}}
		req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewReader([]byte(`{"cep":"29902555"}`)))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"city":"São Paulo","temp_C":28.5,"temp_F":83.3,"temp_K":301.5}`, rec.Body.String())
	})

	t.Run("relays not found status from service b", func(t *testing.T) {
		handler := CEPForwardHandler{ServiceB: fakeServiceBForwarder{
			resp: &serviceaclient.Response{StatusCode: http.StatusNotFound, Body: []byte(`{"message":"can not find zipcode"}`)},
		}}
		req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewReader([]byte(`{"cep":"99999999"}`)))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.JSONEq(t, `{"message":"can not find zipcode"}`, rec.Body.String())
	})
}
