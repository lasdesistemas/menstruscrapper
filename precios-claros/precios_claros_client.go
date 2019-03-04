package preciosclaros

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type PreciosClarosClient struct {
	restClient RestClient
}

func NewClient(rc RestClient) *PreciosClarosClient {
	return &PreciosClarosClient{rc}
}

const (
	host               = "https://d3e6htiiul5ek9.cloudfront.net/prueba"
	pathSucursales     = "/sucursales"
	pathProducto       = "/productos&id_categoria=%v"
	pathPrecioProducto = "/producto"
)

func (pc *PreciosClarosClient) ObtenerSucursales() ([]string, error) {

	sucursales := []string{}

	paginas, err := pc.obtenerSucursales("0", "30", &sucursales)

	if err != nil {
		return sucursales, err
	}

	if paginas > 1 {
		for i := 1; i <= paginas; i++ {
			offset := strconv.Itoa(i * 30)
			limit := "30"
			_, err := pc.obtenerSucursales(offset, limit, &sucursales)

			if err != nil {
				return sucursales, err
			}
		}
	}

	return sucursales, nil
}

func (pc *PreciosClarosClient) obtenerSucursales(offset, limit string, sucursales *[]string) (int, error) {
	response, err := pc.restClient.Get(host + pathSucursales + "?offset=" + offset + "&limit=" + limit)
	defer response.Body.Close()

	// Valida el resultado
	if err != nil {
		return 0, err
	}

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("el pedido dio status: %v", response.StatusCode)
	}

	// Obtengo la respuesta
	bodyBytes, errRead := ioutil.ReadAll(response.Body)

	if errRead != nil {
		return 0, fmt.Errorf("error al leer la respuesta: %v", errRead)
	}

	respuesta := struct {
		TotalPagina int `json:"totalPagina"`
		Total       int `json:"total"`
		Sucursales  []*Sucursal
	}{}

	json.Unmarshal(bodyBytes, &respuesta)
	paginas := int(math.Ceil(float64(respuesta.Total / respuesta.TotalPagina)))

	for _, sucursal := range respuesta.Sucursales {
		*sucursales = append(*sucursales, sucursal.Id)
	}

	return paginas, nil
}

func (pc *PreciosClarosClient) ObtenerListaDeTampones(sucursales []string) ([]int, error) {

	var sucursales50 []string
	tampones := []int{}

	for len(sucursales) > 0 {

		if len(sucursales) > 50 {
			sucursales50 = sucursales[0:50]
			sucursales = sucursales[50:]
		} else {
			sucursales50 = sucursales
			sucursales = nil
		}

		paginas, err := pc.obtenerProductos("090215", "0", "100", &tampones, sucursales50)

		if err != nil {
			return tampones, err
		}

		if paginas > 1 {
			for i := 1; i <= paginas; i++ {
				offset := strconv.Itoa(i * 100)
				limit := "100"
				_, err := pc.obtenerProductos("090215", offset, limit, &tampones, sucursales50)

				if err != nil {
					return tampones, err
				}
			}
		}
	}

	return tampones, nil
}

func (pc *PreciosClarosClient) obtenerProductos(categoria string, offset, limit string, productos *[]int, sucursales []string) (int, error) {
	sucursalesQueryString := "&array_sucursales=" + strings.Join(sucursales, ",")
	pathProductoConCategoria := fmt.Sprintf(pathProducto, categoria)
	response, err := pc.restClient.Get(host + pathProductoConCategoria + sucursalesQueryString + "&offset=" + offset + "&limit=" + limit)
	defer response.Body.Close()

	// Valida el resultado
	if err != nil {
		return 0, err
	}

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("el pedido dio status: %v", response.StatusCode)
	}

	// Obtengo la respuesta
	bodyBytes, errRead := ioutil.ReadAll(response.Body)

	if errRead != nil {
		return 0, fmt.Errorf("error al leer la respuesta: %v", errRead)
	}

	respuesta := struct {
		Total       int `json:"total"`
		TotalPagina int `json:"totalPagina"`
		Productos   []*Producto
	}{}

	json.Unmarshal(bodyBytes, &respuesta)

	// Convierto a lista de int
	for _, producto := range respuesta.Productos {
		id, err := strconv.Atoi(producto.Id)

		if err == nil {
			*productos = append(*productos, id)
		} else {
			fmt.Println("No se pudo convertir a int el id de producto: ", producto.Id)
		}
	}

	paginas := int(math.Ceil(float64(respuesta.Total / respuesta.TotalPagina)))

	return paginas, nil
}

func (pc *PreciosClarosClient) ObtenerListaDeToallitas(sucursales []string) ([]int, error) {

	var sucursales50 []string
	toallitas := []int{}

	for len(sucursales) > 0 {

		if len(sucursales) > 50 {
			sucursales50 = sucursales[0:50]
			sucursales = sucursales[50:]
		} else {
			sucursales50 = sucursales
			sucursales = nil
		}

		paginas, err := pc.obtenerProductos("090216", "0", "100", &toallitas, sucursales50)

		if err != nil {
			return toallitas, err
		}

		if paginas > 1 {
			for i := 1; i <= paginas; i++ {
				offset := strconv.Itoa(i * 100)
				limit := "100"
				_, err := pc.obtenerProductos("090216", offset, limit, &toallitas, sucursales50)

				if err != nil {
					return toallitas, err
				}
			}
		}
	}

	return toallitas, nil
}

func (pc *PreciosClarosClient) ObtenerListaDePrecios(sucursales []string, productos []int, categoria string) ([]*Producto, error) {

	precios := []*Producto{}

	for _, id := range productos {

		sucursalesQueryString := "&array_sucursales=" + strings.Join(sucursales, ",")
		response, err := pc.restClient.Get(host + fmt.Sprintf(pathPrecioProducto+"&id_producto=%v", id) + sucursalesQueryString)

		defer response.Body.Close()

		// Valida el resultado
		if err != nil {
			return nil, err
		}

		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("el pedido dio status: %v", response.StatusCode)
		}

		// Obtengo la respuesta
		bodyBytes, errRead := ioutil.ReadAll(response.Body)

		if errRead != nil {
			return nil, fmt.Errorf("error al leer la respuesta: %v", errRead)
		}

		respuesta := struct {
			TotalPagina int `json:"totalPagina"`
			Total       int `json:"total"`
			Producto    *Producto
			Sucursales  []*Sucursal
		}{}

		json.Unmarshal(bodyBytes, &respuesta)

		// Agrego cada producto con su precio y detalle de sucursal a la lista de precios
		for _, sucursal := range respuesta.Sucursales {
			if sucursal.PreciosProducto != nil {
				p := pc.generarRenglonProducto(sucursal, *respuesta.Producto, categoria)
				precios = append(precios, p)
			}
		}
	}

	return precios, nil
}

func (pc *PreciosClarosClient) generarRenglonProducto(sucursal *Sucursal, producto Producto, categoria string) *Producto {

	precioProducto := Producto{}

	precioProducto = producto

	precioProducto.Categoria = categoria
	precioProducto.Comercio = sucursal.Comercio
	precioProducto.Sucursal = sucursal.Nombre
	precioProducto.Direccion = sucursal.Direccion
	precioProducto.Localidad = sucursal.Localidad
	precioProducto.PrecioDeLista = sucursal.PreciosProducto.PrecioLista

	return &precioProducto
}
