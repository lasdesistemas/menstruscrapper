package preciosclaros_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	preciosclaros "github.com/lasdesistemas/menstruscrapper/precios-claros"
	mock_precios_claros "github.com/lasdesistemas/menstruscrapper/precios-claros/mocks"
	"github.com/stretchr/testify/assert"
)

const (
	host       = "https://d3e6htiiul5ek9.cloudfront.net/prueba"
	sucursales = "/sucursales"
	tampones   = "/productos&id_categoria=090215&array_sucursales=15-1-1803,15-1-8009"
	toallitas  = "/productos&id_categoria=090216&array_sucursales=15-1-1803,15-1-8009"
	producto   = "/producto&id_producto=%v&array_sucursales=15-1-1803,15-1-8009"
)

func TestObtenerSucursales(t *testing.T) {

	// Inicialización
	sucursalesEsperadas := []string{"15-1-1803", "15-1-8009"}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRestClient := inicializarMockRestClient(mockCtrl, []string{"../archivos-test/sucursales.json"}, []string{sucursales})
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
	sucursales := []string{"15-1-1803", "15-1-8009"}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRestClient := inicializarMockRestClient(mockCtrl, []string{"../archivos-test/tampones.json"}, []string{tampones})
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
	sucursales := []string{"15-1-1803", "15-1-8009"}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRestClient := inicializarMockRestClient(mockCtrl, []string{"../archivos-test/toallitas.json"}, []string{toallitas})
	preciosClarosClient := preciosclaros.NewClient(mockRestClient)

	// Operación
	toallitasObtenidas, err := preciosClarosClient.ObtenerListaDeToallitas(sucursales)

	// Validación
	assert.Equal(t, toallitasEsperadas, toallitasObtenidas, "las toallitas no son iguales")
	assert.Nil(t, err)
}

func TestObtenerListaDePreciosDeUnProducto(t *testing.T) {

	// Inicialización
	listaDePreciosEsperada := generarListaDePreciosTampones()
	sucursales := []string{"15-1-1803", "15-1-8009"}
	tampones := []int{7891010604905}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRestClient := inicializarMockRestClient(mockCtrl, []string{"../archivos-test/precios-tampones-7891010604905.json"},
		[]string{fmt.Sprintf(producto, 7891010604905)})
	preciosClarosClient := preciosclaros.NewClient(mockRestClient)

	// Operación
	listaDePreciosObtenida, err := preciosClarosClient.ObtenerListaDePrecios(sucursales, tampones, "Tampones")

	// Validación
	assert.ElementsMatch(t, listaDePreciosEsperada, listaDePreciosObtenida, "las listas de precios no son iguales")
	assert.Nil(t, err)
}

func TestObtenerListaDePreciosDeMasDeUnProducto(t *testing.T) {

	// Inicialización
	listaDePreciosEsperada := generarListaDePreciosDosTampones()
	sucursales := []string{"15-1-1803", "15-1-8009"}
	tampones := []int{7891010604905, 7891010604943}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRestClient := inicializarMockRestClient(mockCtrl, []string{"../archivos-test/precios-tampones-7891010604905.json",
		"../archivos-test/precios-tampones-7891010604943.json"},
		[]string{fmt.Sprintf(producto, 7891010604905), fmt.Sprintf(producto, 7891010604943)})
	preciosClarosClient := preciosclaros.NewClient(mockRestClient)

	// Operación
	listaDePreciosObtenida, err := preciosClarosClient.ObtenerListaDePrecios(sucursales, tampones, "Tampones")

	// Validación
	assert.ElementsMatch(t, listaDePreciosEsperada, listaDePreciosObtenida, "las listas de precios no son iguales")
	assert.Nil(t, err)
}

func inicializarMockRestClient(mockCtrl *gomock.Controller, paths []string, urls []string) *mock_precios_claros.MockRestClient {

	mockRestClient := mock_precios_claros.NewMockRestClient(mockCtrl)
	calls := []*gomock.Call{}

	for i, path := range paths {
		json, _ := ioutil.ReadFile(path)
		body := ioutil.NopCloser(bytes.NewReader(json))
		respuesta := &http.Response{StatusCode: http.StatusOK, Body: body}
		call := mockRestClient.EXPECT().Get(host+urls[i]).Return(respuesta, nil)
		calls = append(calls, call)
	}

	gomock.InOrder(calls...)

	return mockRestClient
}

func generarListaDePreciosTampones() []*preciosclaros.Producto {

	var productos []*preciosclaros.Producto

	unTampon := preciosclaros.Producto{"7891010604905", "Tampones", "OB", "Tampones Medio Helix Ob 20 Un",
		"20.0 un", "DIA Argentina S.A", "1803 - Salta", "Radio Patagonia 0",
		"Salta", 136.49}

	otroTampon := preciosclaros.Producto{"7891010604905", "Tampones", "OB", "Tampones Medio Helix Ob 20 Un",
		"20.0 un", "DIA Argentina S.A", "8009 - Salta", "Sarmiento 0",
		"Salta", 136.49}

	return append(productos, &unTampon, &otroTampon)
}

func generarListaDePreciosDosTampones() []*preciosclaros.Producto {

	productos := generarListaDePreciosTampones()

	unTampon := preciosclaros.Producto{"7891010604943", "Tampones", "OB", "Tampones Medio Pro Comfort Ob 10 Un",
		"10.0 un", "DIA Argentina S.A", "1803 - Salta", "Radio Patagonia 0",
		"Salta", 90.99}

	otroTampon := preciosclaros.Producto{"7891010604943", "Tampones", "OB", "Tampones Medio Pro Comfort Ob 10 Un",
		"10.0 un", "DIA Argentina S.A", "8009 - Salta", "Sarmiento 0",
		"Salta", 90.99}

	return append(productos, &unTampon, &otroTampon)
}
