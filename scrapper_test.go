package menstruscrapper_test

import (
	"io/ioutil"
	"testing"

	"github.com/lasdesistemas/menstruscrapper"
	"github.com/stretchr/testify/assert"
)

func TestGenerarListaDePrecios(t *testing.T) {

	pathCsvEsperado := "archivos-test/esperados/precios.csv"
	scrapper := menstruscrapper.New(NewMockClient())

	pathCsvGenerado := scrapper.GenerarListaDePrecios()

	generaElCsvEsperado(t, pathCsvGenerado, pathCsvEsperado)
}

func generaElCsvEsperado(t *testing.T, pathCsvObtenido string, pathCsvEsperado string) {

	csvObtenido, errCsvObtenido := ioutil.ReadFile(pathCsvObtenido)
	csvEsperado, errCsvEsperado := ioutil.ReadFile(pathCsvEsperado)

	assert.Equal(t, string(csvEsperado), string(csvObtenido), "los archivos no son iguales")
	assert.Nil(t, errCsvObtenido)
	assert.Nil(t, errCsvEsperado)

}

func NewMockClient() *menstruscrapper.MockClient {
	return &menstruscrapper.MockClient{}
}
