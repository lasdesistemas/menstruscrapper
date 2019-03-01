package preciosclaros_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/golang/mock/gomock"
	"net/http"
	"io/ioutil"
	"bytes"
	"github.com/lasdesistemas/menstruscrapper/precios-claros/mocks"
	"github.com/lasdesistemas/menstruscrapper/precios-claros"
)

const (
	host = "https://d3e6htiiul5ek9.cloudfront.net/prueba"
	sucursales = "/sucursales"
)

func TestObtenerSucursales(t *testing.T) {

	// Inicialización
	sucursalesEsperadas := []string{"15-1-1803", "15-1-8009"}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRestClient := inicializarMockRestClient(mockCtrl)
	preciosClarosClient := preciosclaros.NewClient(mockRestClient)

	// Operación
	sucursalesObtenidas, err := preciosClarosClient.ObtenerSucursales()

	// Validación
	assert.Equal(t, sucursalesEsperadas, sucursalesObtenidas, "las sucursales no son iguales")
	assert.Nil(t, err)
}

func inicializarMockRestClient(mockCtrl *gomock.Controller) *mock_precios_claros.MockRestClient {
	mockRestClient := mock_precios_claros.NewMockRestClient(mockCtrl)
	sucursalesJson, _ := ioutil.ReadFile("../archivos-test/sucursales.json")
	body := ioutil.NopCloser(bytes.NewReader(sucursalesJson))
	respuesta := &http.Response{StatusCode: http.StatusOK, Body: body}
	mockRestClient.EXPECT().Get(host + sucursales).Return(respuesta, nil)
	return mockRestClient
}