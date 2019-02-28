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

func TestGenerarListaDePrecios(t *testing.T) {

	// mock the client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mocks.NewMockClient(mockCtrl)
	preciosTampones := generarListaDePrecios()
	inicializarMockClient(mockClient, preciosTampones, nil)

	pathCsvEsperado := "archivos-test/esperados/precios.csv"

	scrapper := menstruscrapper.New(mockClient)

	pathCsvGenerado := scrapper.GenerarListaDePrecios()

	generaElCsvEsperado(t, pathCsvGenerado, pathCsvEsperado)
}

func inicializarMockClient(mockClient *mocks.MockClient, preciosTampones, preciosToallitas []*preciosclaros.Producto) {

	sucursales := []string{"15-1-1803", "15-1-8009"}
	tampones := []int{7891010604905}

	mockClient.EXPECT().ObtenerSucursales().Return(sucursales,nil)
	mockClient.EXPECT().ObtenerListaDeTampones(sucursales).Return(tampones, nil)
	mockClient.EXPECT().ObtenerListaDeToallitas(sucursales).Return(nil, nil)
	primero := mockClient.EXPECT().ObtenerListaDePrecios(sucursales, tampones).Return(preciosTampones, nil)
	segundo := mockClient.EXPECT().ObtenerListaDePrecios(sucursales, nil).Return(preciosToallitas, nil)

	gomock.InOrder(primero, segundo)
}

func generaElCsvEsperado(t *testing.T, pathCsvObtenido string, pathCsvEsperado string) {

	csvObtenido, errCsvObtenido := ioutil.ReadFile(pathCsvObtenido)
	csvEsperado, errCsvEsperado := ioutil.ReadFile(pathCsvEsperado)

	assert.Equal(t, string(csvEsperado), string(csvObtenido), "los archivos no son iguales")
	assert.Nil(t, errCsvObtenido)
	assert.Nil(t, errCsvEsperado)
}

func generarListaDePrecios() []*preciosclaros.Producto {

	var productos []*preciosclaros.Producto

	unTampon := preciosclaros.Producto{"Tampones","OB","Tampones Medio Helix Ob 20 Un",
	"20.0 un","DIA Argentina S.A","1803 - Salta","Radio Patagonia 0",
	"Salta",136.49}

	otroTampon := preciosclaros.Producto{"Tampones","OB","Tampones Medio Helix Ob 20 Un",
		"20.0 un","DIA Argentina S.A","8009 - Salta","Sarmiento 0",
		"Salta",136.49}

	return append(productos, &unTampon, &otroTampon)
}


