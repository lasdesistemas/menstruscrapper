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
	host                    = "https://d3e6htiiul5ek9.cloudfront.net/prueba"
	sucursales              = "/sucursales?offset=0&limit=30"
	tampones                = "/productos&id_categoria=090215&array_sucursales=15-1-1803,15-1-8009&offset=0&limit=100"
	toallitas               = "/productos&id_categoria=090216&array_sucursales=15-1-1803,15-1-8009"
	producto                = "/producto&id_producto=%v&array_sucursales=15-1-1803,15-1-8009"
	pathTamponesConPaginado = "/productos&id_categoria=090215&array_sucursales=%v"
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

func TestObtenerSucursalesConPaginado(t *testing.T) {

	// Inicialización
	sucursalesEsperadas := []string{"11-2-1075", "16-1-1302", "15-1-8012", "15-1-8014",
		"9-1-140", "10-1-112", "15-1-8005", "15-1-8002", "10-1-171", "10-1-175", "15-1-8007",
		"15-1-806", "15-1-1802", "15-1-8013", "15-1-8001", "15-1-1801", "15-1-8006", "15-1-8003",
		"15-1-8015", "6-1-18", "15-1-8010", "15-1-8016", "15-1-804", "15-1-8011", "15-1-802",
		"15-1-8008", "15-1-1800", "15-1-803", "15-1-800", "6-1-9", "15-1-1803", "15-1-8009",
		"15-1-801", "15-1-8004", "9-3-5251", "9-1-655", "9-1-110", "9-1-657", "9-1-656",
		"11-2-1011", "19-1-03330", "9-1-64", "9-1-658", "9-1-731", "9-1-980", "9-1-40",
		"6-2-21", "11-2-1052", "11-2-1078", "36-3-32", "15-1-226", "10-3-521", "15-1-126",
		"2-1-260", "49-1-2", "13-1-111", "50-1-1", "50-1-2", "49-1-1", "12-1-67", "12-1-101",
		"19-1-03298"}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	sucursales := []string{"../archivos-test/sucursales-paginado-1.json",
		"../archivos-test/sucursales-paginado-2.json", "../archivos-test/sucursales-paginado-3.json"}
	urls := []string{"/sucursales?offset=0&limit=30", "/sucursales?offset=30&limit=30", "/sucursales?offset=60&limit=30"}

	mockRestClient := inicializarMockRestClient(mockCtrl, sucursales, urls)
	preciosClarosClient := preciosclaros.NewClient(mockRestClient)

	// Operación
	sucursalesObtenidas, err := preciosClarosClient.ObtenerSucursales()

	// Validación
	assert.Len(t, sucursalesObtenidas, 62)
	assert.Equal(t, sucursalesEsperadas, sucursalesObtenidas, "las sucursales no son iguales")
	assert.Nil(t, err)
}

func TestObtenerListaDeTamponesConPaginado(t *testing.T) {

	// Inicialización
	tamponesEsperados := []int{7702425540606, 7702026176341, 7702026176358, 7702026176334, 7798153712910,
		7798153712927, 7793620000811, 7793620000828, 662425026821, 662425026753, 7793620924414, 7798270242314,
		662425027002, 662425026746, 8480017134790, 7798270242307, 662425027040, 662425027019, 662425026845,
		7891010604806, 7702027435355, 7702027435447, 7792052055123, 7792052046367, 7792052055116, 7792052046350,
		7792052055130, 8480017184924, 7791720002179, 7891010604813, 7793620002129, 7793620002099, 7891010604905,
		7798024399981, 7798158709014, 7891010010560, 7891010886547, 7891010010546, 7891010991029, 7790010963848,
		662425026975, 7891010604943, 7891010604899, 7793620002112, 7793620002082, 7702027435966, 7891010010515,
		7891010952976, 7891010952952, 7891010952969, 7702027435300, 4015400151746, 7891010604929, 7891010010577,
		7798152221093, 7791720002186, 7793620002136, 7793620002105, 7891010604882, 7891010604912, 7798024399998,
		7798158709021, 7891010010591, 662425026760, 662425026982, 7891010604950, 7702027435409, 4015400151715,
		7793620001825, 7790147009952, 7790147009969, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702026176341, 7702026176358,
		7702026176334, 7798153712910, 7798153712927, 7793620000811, 7793620000828, 662425026821, 662425026753,
		7793620924414, 7798270242314, 662425027002, 662425026746, 8480017134790, 7798270242307, 662425027040,
		662425027019, 662425026845, 7891010604806, 7702027435355, 7702027435447, 7792052055123, 7792052046367,
		7792052055116, 7792052046350, 7792052055130, 8480017184924, 7791720002179, 7891010604813, 7793620002129,
		7793620002099, 7891010604905, 7798024399981, 7798158709014, 7891010010560, 7891010886547, 7891010010546,
		7891010991029, 7790010963848, 662425026975, 7891010604943, 7891010604899, 7793620002112, 7793620002082,
		7702027435966, 7891010010515, 7891010952976, 7891010952952, 7891010952969, 7702027435300, 4015400151746,
		7891010604929, 7891010010577, 7798152221093, 7791720002186, 7793620002136, 7793620002105, 7891010604882,
		7891010604912, 7798024399998, 7798158709021, 7891010010591, 662425026760, 662425026982, 7891010604950,
		7702027435409, 4015400151715, 7793620001825, 7790147009952, 7790147009969, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702026176341, 7702425540606, 7702026176341, 7702026176358, 7702026176334, 7798153712910,
		7798153712927, 7793620000811, 7793620000828, 662425026821, 662425026753, 7793620924414, 7798270242314,
		662425027002, 662425026746, 8480017134790, 7798270242307, 662425027040, 662425027019, 662425026845,
		7891010604806, 7702027435355, 7702027435447, 7792052055123, 7792052046367, 7792052055116, 7792052046350,
		7792052055130, 8480017184924, 7791720002179, 7891010604813, 7793620002129, 7793620002099, 7891010604905,
		7798024399981, 7798158709014, 7891010010560, 7891010886547, 7891010010546, 7891010991029, 7790010963848,
		662425026975, 7891010604943, 7891010604899, 7793620002112, 7793620002082, 7702027435966, 7891010010515,
		7891010952976, 7891010952952, 7891010952969, 7702027435300, 4015400151746, 7891010604929, 7891010010577,
		7798152221093, 7791720002186, 7793620002136, 7793620002105, 7891010604882, 7891010604912, 7798024399998,
		7798158709021, 7891010010591, 662425026760, 662425026982, 7891010604950, 7702027435409, 4015400151715,
		7793620001825, 7790147009952, 7790147009969, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702026176341, 7702026176358,
		7702026176334, 7798153712910, 7798153712927, 7793620000811, 7793620000828, 662425026821, 662425026753,
		7793620924414, 7798270242314, 662425027002, 662425026746, 8480017134790, 7798270242307, 662425027040,
		662425027019, 662425026845, 7891010604806, 7702027435355, 7702027435447, 7792052055123, 7792052046367,
		7792052055116, 7792052046350, 7792052055130, 8480017184924, 7791720002179, 7891010604813, 7793620002129,
		7793620002099, 7891010604905, 7798024399981, 7798158709014, 7891010010560, 7891010886547, 7891010010546,
		7891010991029, 7790010963848, 662425026975, 7891010604943, 7891010604899, 7793620002112, 7793620002082,
		7702027435966, 7891010010515, 7891010952976, 7891010952952, 7891010952969, 7702027435300, 4015400151746,
		7891010604929, 7891010010577, 7798152221093, 7791720002186, 7793620002136, 7793620002105, 7891010604882,
		7891010604912, 7798024399998, 7798158709021, 7891010010591, 662425026760, 662425026982, 7891010604950,
		7702027435409, 4015400151715, 7793620001825, 7790147009952, 7790147009969, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606, 7702425540606,
		7702026176341}

	sucursales := []string{"11-2-1075", "16-1-1302", "15-1-8012", "15-1-8014",
		"9-1-140", "10-1-112", "15-1-8005", "15-1-8002", "10-1-171", "10-1-175", "15-1-8007",
		"15-1-806", "15-1-1802", "15-1-8013", "15-1-8001", "15-1-1801", "15-1-8006", "15-1-8003",
		"15-1-8015", "6-1-18", "15-1-8010", "15-1-8016", "15-1-804", "15-1-8011", "15-1-802",
		"15-1-8008", "15-1-1800", "15-1-803", "15-1-800", "6-1-9", "15-1-1803", "15-1-8009",
		"15-1-801", "15-1-8004", "9-3-5251", "9-1-655", "9-1-110", "9-1-657", "9-1-656",
		"11-2-1011", "19-1-03330", "9-1-64", "9-1-658", "9-1-731", "9-1-980", "9-1-40",
		"6-2-21", "11-2-1052", "11-2-1078", "36-3-32", "15-1-226", "10-3-521", "15-1-126",
		"2-1-260", "49-1-2", "13-1-111", "50-1-1", "50-1-2", "49-1-1", "12-1-67", "12-1-101",
		"19-1-03298"}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tampones := []string{"../archivos-test/tampones-paginado-1.json",
		"../archivos-test/tampones-paginado-2.json", "../archivos-test/tampones-paginado-3.json",
		"../archivos-test/tampones-paginado-1.json", "../archivos-test/tampones-paginado-2.json",
		"../archivos-test/tampones-paginado-3.json"}
	pagina1 := fmt.Sprintf(pathTamponesConPaginado, "11-2-1075,16-1-1302,15-1-8012,15-1-8014,9-1-140,10-1-112,15-1-8005,15-1-8002,10-1-171,10-1-175,15-1-8007,15-1-806,15-1-1802,15-1-8013,15-1-8001,15-1-1801,15-1-8006,15-1-8003,15-1-8015,6-1-18,15-1-8010,15-1-8016,15-1-804,15-1-8011,15-1-802,15-1-8008,15-1-1800,15-1-803,15-1-800,6-1-9,15-1-1803,15-1-8009,15-1-801,15-1-8004,9-3-5251,9-1-655,9-1-110,9-1-657,9-1-656,11-2-1011,19-1-03330,9-1-64,9-1-658,9-1-731,9-1-980,9-1-40,6-2-21,11-2-1052,11-2-1078,36-3-32&offset=0&limit=100")
	pagina2 := fmt.Sprintf(pathTamponesConPaginado, "11-2-1075,16-1-1302,15-1-8012,15-1-8014,9-1-140,10-1-112,15-1-8005,15-1-8002,10-1-171,10-1-175,15-1-8007,15-1-806,15-1-1802,15-1-8013,15-1-8001,15-1-1801,15-1-8006,15-1-8003,15-1-8015,6-1-18,15-1-8010,15-1-8016,15-1-804,15-1-8011,15-1-802,15-1-8008,15-1-1800,15-1-803,15-1-800,6-1-9,15-1-1803,15-1-8009,15-1-801,15-1-8004,9-3-5251,9-1-655,9-1-110,9-1-657,9-1-656,11-2-1011,19-1-03330,9-1-64,9-1-658,9-1-731,9-1-980,9-1-40,6-2-21,11-2-1052,11-2-1078,36-3-32&offset=100&limit=100")
	pagina3 := fmt.Sprintf(pathTamponesConPaginado, "11-2-1075,16-1-1302,15-1-8012,15-1-8014,9-1-140,10-1-112,15-1-8005,15-1-8002,10-1-171,10-1-175,15-1-8007,15-1-806,15-1-1802,15-1-8013,15-1-8001,15-1-1801,15-1-8006,15-1-8003,15-1-8015,6-1-18,15-1-8010,15-1-8016,15-1-804,15-1-8011,15-1-802,15-1-8008,15-1-1800,15-1-803,15-1-800,6-1-9,15-1-1803,15-1-8009,15-1-801,15-1-8004,9-3-5251,9-1-655,9-1-110,9-1-657,9-1-656,11-2-1011,19-1-03330,9-1-64,9-1-658,9-1-731,9-1-980,9-1-40,6-2-21,11-2-1052,11-2-1078,36-3-32&offset=200&limit=100")
	pagina4 := fmt.Sprintf(pathTamponesConPaginado, "15-1-226,10-3-521,15-1-126,2-1-260,49-1-2,13-1-111,50-1-1,50-1-2,49-1-1,12-1-67,12-1-101,19-1-03298&offset=0&limit=100")
	pagina5 := fmt.Sprintf(pathTamponesConPaginado, "15-1-226,10-3-521,15-1-126,2-1-260,49-1-2,13-1-111,50-1-1,50-1-2,49-1-1,12-1-67,12-1-101,19-1-03298&offset=100&limit=100")
	pagina6 := fmt.Sprintf(pathTamponesConPaginado, "15-1-226,10-3-521,15-1-126,2-1-260,49-1-2,13-1-111,50-1-1,50-1-2,49-1-1,12-1-67,12-1-101,19-1-03298&offset=200&limit=100")
	urls := []string{pagina1, pagina2, pagina3, pagina4, pagina5, pagina6}
	mockRestClient := inicializarMockRestClient(mockCtrl, tampones, urls)
	preciosClarosClient := preciosclaros.NewClient(mockRestClient)

	// Operación
	tamponesObtenidos, err := preciosClarosClient.ObtenerListaDeTampones(sucursales)

	// Validación
	assert.Len(t, tamponesObtenidos, 404)
	assert.Equal(t, tamponesEsperados, tamponesObtenidos, "los tampones no son iguales")
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
