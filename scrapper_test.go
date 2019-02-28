package menstruscrapper_test

import (
	"io/ioutil"
	"testing"

	"github.com/lasdesistemas/menstruscrapper"
	"github.com/stretchr/testify/assert"
	"github.com/lasdesistemas/menstruscrapper/precios-claros"
	"github.com/golang/mock/gomock"
	"github.com/lasdesistemas/menstruscrapper/mocks"
)

func TestGenerarListaDePreciosCuandoSoloHayTampones(t *testing.T) {

	// Inicialización
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mocks.NewMockClient(mockCtrl)
	preciosTampones := generarListaDePreciosTampones()
	idsTampones := []int{7891010604905}
	inicializarMockClient(mockClient, preciosTampones, nil, idsTampones, nil)

	pathCsvEsperado := "archivos-test/esperados/precios-solo-tampones.csv"

	scrapper := menstruscrapper.New(mockClient)

	// Operación
	pathCsvGenerado := scrapper.GenerarListaDePrecios()

	// Validación
	generaElCsvEsperado(t, pathCsvGenerado, pathCsvEsperado)
}

func TestGenerarListaDePreciosCuandoHayTamponesYToallitas(t *testing.T) {

	// Inicialización
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mocks.NewMockClient(mockCtrl)
	preciosTampones := generarListaDePreciosTampones()
	preciosToallitas := generarListaDePreciosToallitas()
	idsTampones := []int{7891010604905}
	idsToallitas := []int{7501065922755}
	inicializarMockClient(mockClient, preciosTampones, preciosToallitas, idsTampones, idsToallitas)

	pathCsvEsperado := "archivos-test/esperados/precios-tampones-y-toallitas.csv"

	scrapper := menstruscrapper.New(mockClient)

	// Operación
	pathCsvGenerado := scrapper.GenerarListaDePrecios()

	// Validación
	generaElCsvEsperado(t, pathCsvGenerado, pathCsvEsperado)
}

func inicializarMockClient(mockClient *mocks.MockClient, preciosTampones, preciosToallitas []*preciosclaros.Producto,
	tampones []int, toallitas []int) {

	sucursales := []string{"15-1-1803", "15-1-8009"}

	mockClient.EXPECT().ObtenerSucursales().Return(sucursales,nil)
	mockClient.EXPECT().ObtenerListaDeTampones(sucursales).Return(tampones, nil)
	mockClient.EXPECT().ObtenerListaDeToallitas(sucursales).Return(toallitas, nil)
	primero := mockClient.EXPECT().ObtenerListaDePrecios(sucursales, tampones).Return(preciosTampones, nil)
	segundo := mockClient.EXPECT().ObtenerListaDePrecios(sucursales, toallitas).Return(preciosToallitas, nil)

	gomock.InOrder(primero, segundo)
}

func generaElCsvEsperado(t *testing.T, pathCsvObtenido string, pathCsvEsperado string) {

	csvObtenido, errCsvObtenido := ioutil.ReadFile(pathCsvObtenido)
	csvEsperado, errCsvEsperado := ioutil.ReadFile(pathCsvEsperado)

	assert.Equal(t, string(csvEsperado), string(csvObtenido), "los archivos no son iguales")
	assert.Nil(t, errCsvObtenido)
	assert.Nil(t, errCsvEsperado)
}

func generarListaDePreciosTampones() []*preciosclaros.Producto {

	var productos []*preciosclaros.Producto

	unTampon := preciosclaros.Producto{"Tampones","OB","Tampones Medio Helix Ob 20 Un",
	"20.0 un","DIA Argentina S.A","1803 - Salta","Radio Patagonia 0",
	"Salta",136.49}

	otroTampon := preciosclaros.Producto{"Tampones","OB","Tampones Medio Helix Ob 20 Un",
		"20.0 un","DIA Argentina S.A","8009 - Salta","Sarmiento 0",
		"Salta",136.49}

	return append(productos, &unTampon, &otroTampon)
}

func generarListaDePreciosToallitas() []*preciosclaros.Producto {

	var productos []*preciosclaros.Producto

	unaToallita := preciosclaros.Producto{"Toallitas","ALWAYS","Toallas Femeninas Ultrafinas flexialas Always 8 Un",
		"8.0 un","DIA Argentina S.A","1803 - Salta","Radio Patagonia 0",
		"Salta",71.79}

	otraToallita := preciosclaros.Producto{"Toallitas","ALWAYS","Toallas Femeninas Ultrafinas flexialas Always 8 Un",
		"8.0 un","DIA Argentina S.A","8009 - Salta","Sarmiento 0",
		"Salta",71.79}

	return append(productos, &unaToallita, &otraToallita)
}

