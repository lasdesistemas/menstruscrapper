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
	tampones = "/productos&id_categoria=090215&array_sucursales=15-1-1803,15-1-8009"
)

func TestObtenerSucursales(t *testing.T) {

	// Inicialización
	sucursalesEsperadas := []string{"15-1-1803", "15-1-8009"}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRestClient := inicializarMockRestClient(mockCtrl, "../archivos-test/sucursales.json", sucursales)
	preciosClarosClient := preciosclaros.NewClient(mockRestClient)

	// Operación
	sucursalesObtenidas, err := preciosClarosClient.ObtenerSucursales()

	// Validación
	assert.Equal(t, sucursalesEsperadas, sucursalesObtenidas, "las sucursales no son iguales")
	assert.Nil(t, err)
}

func TestObtenerListaDeTampones(t *testing.T) {

	// Inicialización
	tamponesEsperados := []int{7891010604882, 7891010604943, 7891010604905, 7891010604912, 8480017134790,
	8480017184924, 7891010604813}
	sucursales:= []string{"15-1-1803", "15-1-8009"}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRestClient := inicializarMockRestClient(mockCtrl, "../archivos-test/tampones.json", tampones)
	preciosClarosClient := preciosclaros.NewClient(mockRestClient)

	// Operación
	tamponesObtenidos, err := preciosClarosClient.ObtenerListaDeTampones(sucursales)

	// Validación
	assert.Equal(t, tamponesEsperados, tamponesObtenidos, "los tampones no son iguales")
	assert.Nil(t, err)
}

func inicializarMockRestClient(mockCtrl *gomock.Controller, path string, url string) *mock_precios_claros.MockRestClient {
	mockRestClient := mock_precios_claros.NewMockRestClient(mockCtrl)
	json, _ := ioutil.ReadFile(path)
	body := ioutil.NopCloser(bytes.NewReader(json))
	respuesta := &http.Response{StatusCode: http.StatusOK, Body: body}
	mockRestClient.EXPECT().Get(host + url).Return(respuesta, nil)
	return mockRestClient
}