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
	toallitas = "/productos&id_categoria=090216&array_sucursales=15-1-1803,15-1-8009"
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

func TestObtenerListaDeToallitas(t *testing.T) {

	// Inicialización
	toallitasEsperadas := []int{7501065922755, 7793620003386, 7790010614856, 7500435140560, 7506339393804,
	8480017180032, 7790010669061, 7790010599283, 7891010607135, 7506339391558, 7793620003249, 7506339356939,
	7790010599269, 7793620003263, 7790010596602, 7794626008559, 7790010669085, 7500435023306, 7500435023290,
	7793620003256, 7790010669078, 7793620003232, 7891010599294, 7500435023337, 8480017775313}
	sucursales:= []string{"15-1-1803", "15-1-8009"}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRestClient := inicializarMockRestClient(mockCtrl, "../archivos-test/toallitas.json", toallitas)
	preciosClarosClient := preciosclaros.NewClient(mockRestClient)

	// Operación
	toallitasObtenidas, err := preciosClarosClient.ObtenerListaDeToallitas(sucursales)

	// Validación
	assert.Equal(t, toallitasEsperadas, toallitasObtenidas, "las toallitas no son iguales")
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