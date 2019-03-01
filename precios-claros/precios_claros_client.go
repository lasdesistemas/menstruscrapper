package preciosclaros

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"strconv"
)

type PreciosClarosClient struct {
	restClient RestClient
}

func NewClient(rc RestClient) *PreciosClarosClient {
	return &PreciosClarosClient{rc}
}

const (
	host = "https://d3e6htiiul5ek9.cloudfront.net/prueba"
	sucursales = "/sucursales"
	tampones = "/productos&id_categoria=090215"
	toallitas = "/productos&id_categoria=090216"
	producto = "/producto"
)

func (pc *PreciosClarosClient) ObtenerSucursales() ([]string, error) {

	response, err := pc.restClient.Get(host + sucursales)
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
		TotalPagina int `json:"total_pagina"`
		Total int `json:"total"`
		Sucursales []*Sucursal
	}{0, 0, []*Sucursal{}}

	json.Unmarshal(bodyBytes, &respuesta)

	// Convierto a lista de string
	sucursales := []string{}
	for _, sucursal := range respuesta.Sucursales {
		sucursales = append(sucursales, sucursal.Id)
	}

	return sucursales, nil
}

func (pc *PreciosClarosClient) ObtenerListaDeTampones(sucursales []string) ([]int, error) {

	sucursalesQueryString := "&array_sucursales=" + strings.Join(sucursales, ",")

	response, err := pc.restClient.Get(host + tampones + sucursalesQueryString)
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
		Total int `json:"total"`
		Productos []*Producto
	}{0, []*Producto{}}

	json.Unmarshal(bodyBytes, &respuesta)

	// Convierto a lista de int
	tampones := []int{}
	for _, producto := range respuesta.Productos {
		id, err := strconv.Atoi(producto.Id)

		if err == nil {
			tampones = append(tampones, id)
		} else {
			fmt.Println("No se pudo convertir a int el id de producto: ", producto.Id)
		}

	}

	return tampones, nil
}

func (pc *PreciosClarosClient) ObtenerListaDeToallitas(sucursales []string) ([]int, error) {

	sucursalesQueryString := "&array_sucursales=" + strings.Join(sucursales, ",")

	response, err := pc.restClient.Get(host + toallitas + sucursalesQueryString)
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
		Total int `json:"total"`
		Productos []*Producto
	}{0, []*Producto{}}

	json.Unmarshal(bodyBytes, &respuesta)

	// Convierto a lista de int
	toallitas := []int{}
	for _, producto := range respuesta.Productos {
		id, err := strconv.Atoi(producto.Id)

		if err == nil {
			toallitas = append(toallitas, id)
		} else {
			fmt.Println("No se pudo convertir a int el id de producto: ", producto.Id)
		}

	}

	return toallitas, nil
}

func (pc *PreciosClarosClient) ObtenerListaDePrecios(sucursales []string, productos []int) ([]*Producto, error) {

	precios := []*Producto{}

	for _ , id := range productos {

		sucursalesQueryString := "&array_sucursales=" + strings.Join(sucursales, ",")
		response, err := pc.restClient.Get(host + fmt.Sprintf(producto + "&id_producto=%v", id) + sucursalesQueryString)

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
		Total int `json:"total"`
		Producto *Producto
		Sucursales []*Sucursal
		}{0, 0,&Producto{}, []*Sucursal{}}

		json.Unmarshal(bodyBytes, &respuesta)

		// Agrego cada producto con su precio y detalle de sucursal a la lista de precios
		for _, sucursal := range respuesta.Sucursales {
			if sucursal.PreciosProducto != nil {
				p := pc.generarRenglonProducto(sucursal, *respuesta.Producto)
				precios = append(precios, p)
			}
		}
	}

	return precios, nil
}

func (pc *PreciosClarosClient) generarRenglonProducto(sucursal *Sucursal, producto Producto) *Producto {

	precioProducto := Producto{}

	precioProducto = producto

	if strings.Contains(strings.ToLower(precioProducto.Nombre), "tampon") {
		precioProducto.Categoria = "Tampones"
	} else {
		precioProducto.Categoria = "Toallitas"
	}

	precioProducto.Comercio = sucursal.Comercio
	precioProducto.Sucursal = sucursal.Nombre
	precioProducto.Direccion = sucursal.Direccion
	precioProducto.Localidad = sucursal.Localidad
	precioProducto.PrecioDeLista = sucursal.PreciosProducto.PrecioLista

	return &precioProducto
}